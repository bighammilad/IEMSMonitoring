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

func NewUpdateUser() *UserEndpoint {
	var userUC useruc.UserUsecase = useruc.UserUsecase{DB: &userrepo.UserRepo{}}
	return &UserEndpoint{
		UserUC: &userUC,
	}
}

func (lue *UserEndpoint) UpdateUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		accessLevel, _ := strconv.Atoi(role)
		lue.UserUC.DB.Update(c.Request().Context(), username, password, accessLevel)
	} else {
		return errors.New("Forbidden")
	}

	return nil

}
