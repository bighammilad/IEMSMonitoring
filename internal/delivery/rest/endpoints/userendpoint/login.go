package userendpoint

import (
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase/useruc"
	"net/http"
	"time"

	. "monitoring/internal/globals"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginUserEndpoint struct {
	Name  string `json:"name,omitempty"`
	Admin bool   `json:"admin,omitempty"`
	Role  int    `json:"role,omitempty"`
	jwt.RegisteredClaims
	LoginUC useruc.ILogin `json:"login_uc,omitempty"`
}

func NewLoginUserEndpoint() *LoginUserEndpoint {
	var loginuc useruc.LoginUC = useruc.LoginUC{LoginRepo: userrepo.LoginRepo{DB: GlobalPG}}
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
	var LUE *LoginUserEndpoint
	if loginuc.Role == 0 {
		LUE = &LoginUserEndpoint{
			loginuc.Username,
			true,
			loginuc.Role,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
			nil,
		}
	} else {
		LUE = &LoginUserEndpoint{
			loginuc.Username,
			false,
			loginuc.Role,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			},
			nil,
		}
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LUE)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}
