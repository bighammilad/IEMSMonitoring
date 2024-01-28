package userendpoint

import (
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase/useruc"
	"strconv"

	. "monitoring/internal/globals"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewRegisterUser() *UserEndpoint {
	var useruc useruc.UserUsecase = useruc.UserUsecase{DB: &userrepo.UserRepo{DB: GlobalPG}}
	return &UserEndpoint{
		UserUC: &useruc,
	}
}

func (lue *UserEndpoint) RegisterUser(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	if username == "" || password == "" || role == "" {
		return errors.New("bad request")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	roleId, _ := strconv.Atoi(role)

	if claims.Admin {
		_ = jwt.NewWithClaims(jwt.SigningMethodES256,
			jwt.MapClaims{
				"role": roleId,
				"sm":   claims.Name,
			})
		// s, err := t.SignedString([]byte("Secret"))
		// if err != nil {
		// 	return err
		// }
		lue.UserUC.DB.Register(c.Request().Context(), username, password, roleId)
	} else {
		return errors.New("Forbidden")
	}

	return nil
}
