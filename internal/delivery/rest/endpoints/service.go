package endpoints

import (
	"encoding/json"
	"monitoring/internal/model"
	"monitoring/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ServicesEndpoints struct {
	ServicesUC usecase.ServicesUsecase
}

func NewServicesEndpoints() *ServicesEndpoints {
	return &ServicesEndpoints{}
}

func (se *ServicesEndpoints) ListServices(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	services, err := se.ServicesUC.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}

func (se *ServicesEndpoints) AddService(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	service, err := se.checkParams(c)
	if err != nil {
		return err
	}

	// add service
	err = se.ServicesUC.Add(c.Request().Context(), service)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, service)

}

func (se *ServicesEndpoints) UpdateService(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	service, err := se.checkParams(c)
	if err != nil {
		return err
	}

	// update service
	err = se.ServicesUC.Update(c.Request().Context(), service)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Service updated")
}

func (se *ServicesEndpoints) DeleteService(c echo.Context) error {

	// first check if the user is admin
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	id := c.Param("id")
	name := c.Param("name")
	service := model.Service{}

	// check for empty params
	if id == "" || name == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	service.Name = name
	service.ID = idInt

	if idInt < 0 {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	// delete service
	err = se.ServicesUC.Delete(c.Request().Context(), service)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Service deleted")

}

func (se *ServicesEndpoints) GetService(c echo.Context) error {

	// first check if the user is admin
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	id := c.Param("id")
	name := c.Param("name")
	service := model.Service{}

	// check for empty params
	if id == "" && name == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	service.Name = name

	if idInt < 0 {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	// get service by id, name
	svc, err := se.ServicesUC.Get(c.Request().Context(), service)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, svc)

}

func (se *ServicesEndpoints) checkParams(c echo.Context) (service model.Service, err error) {
	name := c.Param("name")
	address := c.Param("address")
	method := c.Param("method")
	header := c.Param("header")
	body := c.Param("body")
	accessLevel := c.Param("accessLevel")
	executionTime := c.Param("executionTime")

	// check for empty params
	if (name == "") && (address == "" || method == "" || header == "" || body == "" || accessLevel == "" || executionTime == "") {
		return model.Service{}, c.JSON(http.StatusBadRequest, "Bad request")
	}

	// put params in service
	headerMap := make(map[string]string)
	err = json.Unmarshal([]byte(header), &headerMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	bodyMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &bodyMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	accessLevelInt, err := strconv.Atoi(accessLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	accLevel := model.AccessLevel(accessLevelInt)
	exeTime, err := time.ParseDuration(executionTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	service = model.Service{
		Name:          name,
		Address:       address,
		Method:        method,
		Header:        headerMap,
		Body:          bodyMap,
		AccessLevel:   accLevel,
		ExecutionTime: exeTime,
	}

	return service, nil
}
