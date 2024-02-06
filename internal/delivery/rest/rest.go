package internal

import (
	"monitoring/config"
	// "monitoring/internal/delivery/rest/endpoints"
	"monitoring/internal/delivery/rest/endpoints"
	"monitoring/internal/delivery/rest/endpoints/userendpoint"
	"monitoring/internal/delivery/rest/middlewares"
	"monitoring/internal/model"

	// . "monitoring/internal/globals"

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
	restericted.Use(middlewares.IsAdminUser)

	user := userendpoint.NewUserEndpoint()
	restericted.Use()
	restericted.POST("/user/create", user.Create)
	restericted.POST("/user/read", user.Read)
	restericted.GET("/user/readall", user.ReadAll)
	restericted.POST("/user/update", user.Update)
	restericted.POST("/user/delete", user.Delete)

	service := endpoints.NewServicesEndpoints()
	restericted.GET("/service/getservices", service.GetUserServices)
	restericted.GET("/service/getservice", service.GetUserService)
	restericted.POST("/service/add", service.AddService)
	restericted.POST("/service/delete", service.DeleteService)

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
