package controllers

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	"github.com/google/uuid"
)

type ImageController struct {
	service app.ImageService
}

func NewImageController(s app.ImageService) ImageController {
	return ImageController{
		service: s,
	}
}

// AddImages godoc
// @Summary      Add a new image
// @Security     ApiKeyAuth
// @Description  load a new image
// @Tags         images
// @Accept       png
// @Accept       jpeg
// @Produce      json
// @Produce      xml
// @Param        userId   path      string  true  "User ID"
// @Param        formData   formData	file   true "The file to upload"	x-mimetype: image/png
// @Success      201  {object}  resources.ImageDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /users/{userId}/image [post]
func (c ImageController) AddImage(ctx echo.Context) error {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}

	buff, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}

	filetype := http.DetectContentType(buff)

	if filetype != "image/jpeg" && filetype != "image/png" {
		return FormatedResponse(ctx, http.StatusForbidden, err)
	}

	img := database.Image{
		UserId: userId,
		Name:   uuid.NewString() + "." + strings.TrimLeft(filetype, "image/"),
	}

	i, err := c.service.Save(img, buff)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}

	var imageDto resources.ImageDto
	return FormatedResponse(ctx, http.StatusCreated, imageDto.DatabaseToDto(i))
}

func (c ImageController) DeleteImage(ctx echo.Context) error {

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = c.service.Delete(id)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}

	return FormatedResponse(ctx, http.StatusOK, domain.OK)
}
