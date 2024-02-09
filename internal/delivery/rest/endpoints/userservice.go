package endpoints

import (
	. "monitoring/internal/globals"
	"monitoring/internal/model"
	"monitoring/internal/repository"
	"monitoring/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserServiceEndpoints struct {
	UserServiceUsecase usecase.IUserService
}

func NewUserServiceEndpoints() *UserServiceEndpoints {
	return &UserServiceEndpoints{
		UserServiceUsecase: &usecase.UserService{
			UserServiceRepo: &repository.UserServiceRepository{
				DB: GlobalPG,
			},
		},
	}

}

func (us *UserServiceEndpoints) Add(c echo.Context) error {
	var userServices model.UserService
	if err := c.Bind(&userServices); err != nil {
		return err
	}
	if userServices.ServiceID == 0 || userServices.UserID == 0 {
		return c.JSON(echo.ErrBadRequest.Code, "ServiceID and UserID are required")
	}

	err := us.UserServiceUsecase.Add(c.Request().Context(), userServices)
	if err != nil {
		return c.JSON(echo.ErrNotFound.Code, err.Error())
	}
	return c.JSON(http.StatusOK, "User Service added successfully")
}

func (us *UserServiceEndpoints) GetUserService(c echo.Context) error {
	var userservice model.UserService
	if err := c.Bind(&userservice); err != nil {
		return err
	}
	service, err := us.UserServiceUsecase.GetUserServices(c.Request().Context(), userservice)
	if err != nil {
		return c.JSON(echo.ErrNotFound.Code, err.Error())
	}
	return c.JSON(http.StatusOK, service)
}

func (us *UserServiceEndpoints) GetUserServices(c echo.Context) error {
	var userservice model.UserService
	if err := c.Bind(&userservice); err != nil {
		return err
	}
	services, err := us.UserServiceUsecase.GetUserServices(c.Request().Context(), userservice)
	if err != nil {
		return c.JSON(echo.ErrNotFound.Code, err.Error())
	}
	return c.JSON(http.StatusOK, services)
}

func (us *UserServiceEndpoints) List(c echo.Context) error {
	services, err := us.UserServiceUsecase.List(c.Request().Context())
	if err != nil {
		return c.JSON(echo.ErrNotFound.Code, err.Error())
	}
	return c.JSON(http.StatusOK, services)
}

func (us *UserServiceEndpoints) Update(c echo.Context) error {
	var userServices model.UserService
	if err := c.Bind(&userServices); err != nil {
		return err
	}
	err := us.UserServiceUsecase.Update(c.Request().Context(), userServices)
	if err != nil {
		return c.JSON(echo.ErrNotFound.Code, err.Error())
	}
	return c.JSON(http.StatusOK, "user service updated successfully")
}

func (us *UserServiceEndpoints) Delete(c echo.Context) error {
	var userservice model.UserService
	if err := c.Bind(&userservice); err != nil {
		return err
	}
	err := us.UserServiceUsecase.DeleteUserService(c.Request().Context(), userservice)
	if err != nil {
		return c.JSON(echo.ErrNotFound.Code, err.Error())
	}
	return c.JSON(http.StatusOK, "user service deleted successfully")
}
