package userendpoint

import (
	"monitoring/internal/model"
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase/useruc"
	"net/http"
	"time"

	. "monitoring/internal/globals"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginUserEndpoint struct {
	Name string `json:"name,omitempty"`
	Role int    `json:"role,omitempty"`
	jwt.RegisteredClaims
	LoginUC useruc.ILogin `json:"login_uc,omitempty"`
}

func NewLoginUserEndpoint() *LoginUserEndpoint {
	var loginuc useruc.LoginUC = useruc.LoginUC{ILoginRepo: &userrepo.LoginRepo{DB: GlobalPG}}
	return &LoginUserEndpoint{
		LoginUC: &loginuc,
	}
}

func (le *LoginUserEndpoint) Login(c echo.Context) error {
	var usr model.UserAuth
	if err := c.Bind(&usr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	if usr.Username == "" || usr.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	loginuc, err := le.LoginUC.Login(c.Request().Context(), usr.Username, usr.Password)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		} else if err.Error() == "incorrect password" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "incorrect password"})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}
	var mapclaim jwt.MapClaims
	if loginuc.Role == 1 {
		mapclaim = jwt.MapClaims{
			"name": loginuc.Username,
			"role": loginuc.Role,
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
		}
	} else {
		mapclaim = jwt.MapClaims{
			"name": loginuc.Username,
			"role": loginuc.Role,
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapclaim)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
