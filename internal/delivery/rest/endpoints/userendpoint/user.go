package userendpoint

import (
	. "monitoring/internal/globals"
	"monitoring/internal/model"
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase/useruc"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	midware "monitoring/internal/delivery/rest/middlewares"
)

type UserEndpoint struct {
	Name string `json:"name,omitempty"`
	jwt.RegisteredClaims
	UserUC useruc.IUserUsecase `json:"repo,omitempty"`
}

func NewUserEndpoint() *UserEndpoint {
	db := userrepo.UserRepo{DB: GlobalPG}
	var useruc useruc.UserUsecase = useruc.UserUsecase{IUser: &db}
	return &UserEndpoint{
		UserUC: &useruc,
	}
}

func (ue *UserEndpoint) Create(c echo.Context) error {
	admin, err := midware.IsAdmin(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !admin {
		return echo.ErrForbidden
	}
	var usr model.UserAuth
	if err := c.Bind(&usr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}
	if usr.Username == "" || usr.Password == "" || usr.Role == 0 {
		return echo.ErrBadRequest
	}
	ok, err := ue.UserUC.Create(c.Request().Context(), usr.Username, usr.Password, usr.Role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !ok {
		return echo.ErrBadRequest
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return nil
}

func (ue *UserEndpoint) Read(c echo.Context) error {
	admin, err := midware.IsAdmin(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !admin {
		return echo.ErrForbidden
	}
	var usr model.UserAuth
	if err := c.Bind(&usr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}
	if usr.Username == "" {
		return echo.ErrBadRequest
	}
	user, err := ue.UserUC.Read(c.Request().Context(), usr.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (ue *UserEndpoint) ReadAll(c echo.Context) error {
	admin, err := midware.IsAdmin(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !admin {
		return echo.ErrForbidden
	}
	users, err := ue.UserUC.ReadAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"users": users,
	})
}

func (ue *UserEndpoint) Update(c echo.Context) error {
	admin, err := midware.IsAdmin(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !admin {
		return echo.ErrForbidden
	}
	var usr model.UserAuth
	if err := c.Bind(&usr); err != nil {
		return err
	}
	ok, err := ue.UserUC.Update(c.Request().Context(), usr.Username, usr.Password, usr.Role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !ok {
		return echo.ErrBadRequest
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "User Updated",
	})
}

func (ue *UserEndpoint) Delete(c echo.Context) error {
	admin, err := midware.IsAdmin(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if !admin {
		return echo.ErrForbidden
	}
	var usr model.UserAuth
	if err := c.Bind(&usr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}
	if usr.Username == "" {
		return echo.ErrBadRequest
	}
	err = ue.UserUC.Delete(c.Request().Context(), usr.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "User Deleted",
	})
}
