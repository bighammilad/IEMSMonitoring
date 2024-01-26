package endpoints

import (
	"encoding/json"
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/usecase"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ServicesEndpoints struct {
	ServicesUC usecase.ServicesUsecase
}

func NewServicesEndpoints(srvUC usecase.ServicesUsecase) *ServicesEndpoints {
	return &ServicesEndpoints{
		ServicesUC: srvUC,
	}
}

func (se *ServicesEndpoints) ListServices(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if !(claims.RoleId > 1) {
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
	if !claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	service, err := se.checkPostParams(c)
	if err != nil {
		var t any
		return c.JSON(http.StatusBadRequest, t)
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
	if !claims.Admin {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	service, err := se.checkGetParams(c)
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
	if !claims.Admin {
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
	if claims.RoleId > 1 {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	type RequestBody struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var requestBody RequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if requestBody.ID == "" && requestBody.Name == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	var err error
	var idInt int
	service := model.Service{}
	if requestBody.ID != "" {
		idInt, err = strconv.Atoi(requestBody.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		if idInt < 0 {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}
		service.ID = idInt
	}

	if requestBody.Name != "" {
		service.Name = requestBody.Name
	}

	svc, err := se.ServicesUC.Get(c.Request().Context(), service)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, svc)

}

func (se *ServicesEndpoints) checkGetParams(c echo.Context) (service model.Service, err error) {
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
	headerMap := make(map[string]string, 0)
	err = json.Unmarshal([]byte(header), &headerMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	bodyMap := make(map[string]interface{}, 0)
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

func (se *ServicesEndpoints) checkPostParams(c echo.Context) (service model.Service, err error) {
	name := c.FormValue("name")
	address := c.FormValue("address")
	method := c.FormValue("method")
	header := c.FormValue("header")
	body := c.FormValue("body")
	accesslevel := c.FormValue("accesslevel")
	executiontime := c.FormValue("executiontime")
	allowedusers := c.FormValue("allowedusers")

	// check for empty params
	if name == "" && accesslevel == "" {
		return model.Service{}, errors.New("service name & accesslevel must be filled")
	}
	if address == "" && method == "" && header == "" && body == "" && executiontime == "" {
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

	var accessLevelInt int
	var accLevel model.AccessLevel
	if accesslevel != "" {
		accessLevelInt, err = strconv.Atoi(accesslevel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		if accessLevelInt < 0 && accessLevelInt > 2 {
			return model.Service{}, errors.New("access level value isn't valid")
		}
		accLevel = model.AccessLevel(accessLevelInt)
	} else {
		return model.Service{}, errors.New("access level must be filled")
	}

	exeTime, err := time.ParseDuration(executiontime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	allowU := strings.Split(allowedusers, ", ")
	for i, value := range allowU {
		allowU[i] = strings.TrimSpace(value)
	}

	service = model.Service{
		Name:          name,
		Address:       address,
		Method:        method,
		Header:        headerMap,
		Body:          bodyMap,
		AccessLevel:   accLevel,
		ExecutionTime: exeTime,
		AllowedUsers:  allowU,
	}

	return service, nil
}
