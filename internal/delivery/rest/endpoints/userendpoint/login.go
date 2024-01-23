package userendpoint

import (
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase/useruc"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type LoginUserEndpoint struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
	LoginUC useruc.ILogin
}

func NewLoginUserEndpoint() *LoginUserEndpoint {
	var loginuc useruc.LoginUC = useruc.LoginUC{LoginRepo: userrepo.LoginRepo{}}
	return &LoginUserEndpoint{
		LoginUC: &loginuc,
	}
}

func (le *LoginUserEndpoint) Login(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	loginuc, err := le.LoginUC.Login(c.Request().Context(), username, password)
	if err != nil {
		return err
	}

	// Set custom claims
	var claims *LoginUserEndpoint
	if loginuc.Role == 0 {
		claims = &LoginUserEndpoint{
			loginuc.Username,
			true,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
			nil,
		}
	} else {
		claims = &LoginUserEndpoint{
			loginuc.Username,
			false,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
			nil,
		}
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}
