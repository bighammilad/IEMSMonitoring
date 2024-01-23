package internal

import (
	"monitoring/config"
	// "monitoring/internal/delivery/rest/endpoints"
	"monitoring/internal/delivery/rest/endpoints"
	"monitoring/internal/delivery/rest/endpoints/userendpoint"
	"monitoring/internal/model"

	"monitoring/pkg/postgres"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Rest struct {
	cfg        *config.Config
	monitoring postgres.IPostgres
	e          *echo.Echo
}

func New() (r *Rest, err error) {
	e := echo.New()
	r = &Rest{
		// cfg: cfg,
		e: e,
	}

	loginEndpoint := userendpoint.NewLoginUserEndpoint()
	e.POST("/login", loginEndpoint.Login).Name = "login"

	restericted := e.Group("/panel")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	restericted.Use(echojwt.WithConfig(config))

	restericted.GET("", restricted)
	registerUsr := userendpoint.NewRegisterUser()
	restericted.POST("/user/register", registerUsr.RegisterUser)
	updateUsr := userendpoint.NewRegisterUser()
	restericted.POST("/user/update", updateUsr.UpdateUser)

	service := endpoints.NewServicesEndpoints()
	restericted.POST("/service/add", service.AddService)

	e.GET("/demo", demo)
	e.GET("/test", test, echojwt.WithConfig(config))

	return
}

func (r Rest) Start(address string) (err error) {
	err = r.e.Start(address)

	return
}

func test(c echo.Context) error {
	return c.String(http.StatusOK, "test")
}

func demo(c echo.Context) error {
	return c.String(http.StatusOK, "demo")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	name := claims.Name
	if claims.Admin {
		return c.String(http.StatusOK, "Welcome "+name+"! Admin")
	} else {
		return c.String(http.StatusOK, "Welcome "+name+"!")
	}
}
