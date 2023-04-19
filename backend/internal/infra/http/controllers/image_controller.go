package controllers

import (
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/DmytroKha/nix-chat/internal/app"
)

type ImageController struct {
	service app.ImageService
}

func NewImageController(s app.ImageService) ImageController {
	return ImageController{
		service: s,
	}
}

// AddImage godoc
// @Summary Add an image
// @Description Adds an image to the user's account
// @Tags images
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param image formData file true "Image file to upload"
// @Success 200 {string} string "Returns the URL of the uploaded image"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Router /images [post]
func (c ImageController) AddImage(ctx echo.Context) error {

	userCtxValue := ctx.Request().Context().Value("user")
	if userCtxValue == nil {
		err := AuthErrNotAuth
		log.Println(err)
		return err
	}
	user := userCtxValue.(domain.User)
	userId := user.GetId()

	buff, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	filetype := http.DetectContentType(buff)

	if filetype != "image/jpeg" && filetype != "image/png" {
		returnErrorResponse(ctx.Response().Writer, http.StatusForbidden)
		return err
	}

	img := database.Image{
		UserId: userId,
		Name:   uuid.NewString() + "." + strings.TrimLeft(filetype, "image/"),
	}

	i, err := c.service.Save(img, buff)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	ctx.Response().Write([]byte("../././file_storage/" + i.Name))

	return nil
}
