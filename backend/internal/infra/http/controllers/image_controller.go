package controllers

import (
	"fmt"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DmytroKha/nix-chat/internal/app"
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

func (c ImageController) AddImage(ctx echo.Context) error {

	userCtxValue := ctx.Request().Context().Value("user")
	if userCtxValue == nil {
		err := fmt.Errorf("Not authenticated")
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

	ctx.Response().Write([]byte(i.Name))

	return nil
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
