package userendpoint

import (
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase/useruc"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewRegisterUser() *UserEndpoint {
	var useruc useruc.UserUsecase = useruc.UserUsecase{Repo: &userrepo.UserRepo{}}
	return &UserEndpoint{
		UserUC: &useruc,
	}
}

func (lue *UserEndpoint) RegisterUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		accessLevel, _ := strconv.Atoi(role)
		lue.UserUC.Repo.Register(c.Request().Context(), username, password, accessLevel)
	} else {
		return errors.New("Forbidden")
	}

	return nil
}
