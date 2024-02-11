package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	. "monitoring/internal/globals"
	"monitoring/internal/model"
	"monitoring/internal/repository"
	"monitoring/internal/repository/userrepo"
	"monitoring/internal/usecase"
	"monitoring/internal/usecase/useruc"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ServicesEndpoints struct {
	IServicesUC usecase.IServicesUsecase
	IUseruc     useruc.IUserUsecase
}

func NewServicesEndpoints() *ServicesEndpoints {

	return &ServicesEndpoints{
		IServicesUC: &usecase.ServicesUsecase{
			IServicesRepo: &repository.ServicesRepository{DB: GlobalPG},
		},
		IUseruc: &useruc.UserUsecase{
			IUserRepo: &userrepo.UserRepo{DB: GlobalPG},
		},
	}

}

func (se *ServicesEndpoints) ListServices(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.RoleId != 1 {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	services, err := se.IServicesUC.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, services)
}

func (se *ServicesEndpoints) AddService(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.RoleId != 1 {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	service, userIds, err := se.checkPostParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// add service
	// if len(userIds) == 0 {
	// 	return c.JSON(http.StatusBadRequest, "No users added")
	// }
	err = se.IServicesUC.Add(c.Request().Context(), service, userIds)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, service)

}

func (se *ServicesEndpoints) UpdateService(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.RoleId != 1 {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	service, _, err := se.checkPostParams(c)
	if err != nil {
		return err
	}

	// update service
	err = se.IServicesUC.Update(c.Request().Context(), service)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Service updated")
}

func (se *ServicesEndpoints) DeleteService(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if claims.RoleId != 1 {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	type RequestBody struct {
		Name string `json:"name"`
	}

	var req RequestBody
	if err := c.Bind(&req); err != nil {
		return errors.New("invalid JSON")
	}

	service := model.Service{}

	// check for empty params
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	service.Name = &req.Name

	// delete service
	err := se.IServicesUC.Delete(c.Request().Context(), service)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Service deleted")

}

func (se *ServicesEndpoints) GetUserService(c echo.Context) error {

	var err error
	tokenString := c.Get("user").(*jwt.Token)
	token, err := jwt.Parse(tokenString.Raw, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	jwtToken := token.Claims.(jwt.MapClaims)
	type RequestBody struct {
		Name string `json:"name"`
	}
	var requestBody RequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if requestBody.Name == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	service := model.Service{}

	if requestBody.Name != "" {
		service.Name = &requestBody.Name
	}

	userId, err := se.getUsrId(c, jwtToken["name"].(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	svc, err := se.IServicesUC.GetUserService(c.Request().Context(), requestBody.Name, int(jwtToken["role"].(float64)), userId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, svc)

}

func (se *ServicesEndpoints) GetUserServices(c echo.Context) error {

	var err error
	tokenString := c.Get("user").(*jwt.Token)
	token, err := jwt.Parse(tokenString.Raw, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	jwtToken := token.Claims.(jwt.MapClaims)
	userId, err := se.getUsrId(c, jwtToken["name"].(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	svc, err := se.IServicesUC.GetUserServices(c.Request().Context(), int(jwtToken["role"].(float64)), userId)

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

	var exeTimeInt64 int64
	if executionTime != "" {
		exeTimeInt64, err = strconv.ParseInt(executionTime, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	service = model.Service{
		Name:          &name,
		Address:       &address,
		Method:        &method,
		Header:        headerMap,
		Body:          bodyMap,
		AccessLevel:   accLevel,
		ExecutionTime: &exeTimeInt64,
	}

	return service, nil
}

func (se *ServicesEndpoints) checkPostParams(c echo.Context) (service model.Service, userIds []int, err error) {

	type RequestBody struct {
		Name          string `json:"name,omitempty"`
		Address       string `json:"address,omitempty"`
		Method        string `json:"method,omitempty"`
		Header        string `json:"header,omitempty"`
		Body          string `json:"body,omitempty"`
		AccessLevel   string `json:"accesslevel,omitempty"`
		ExecutionTime string `json:"execution_time,omitempty"`
		AllowedUsers  string `json:"users,omitempty"`
	}

	var req RequestBody
	if err := c.Bind(&req); err != nil {
		return model.Service{}, userIds, errors.New("invalid JSON")
	}

	// check for empty params
	if req.Name == "" || req.AccessLevel == "" {
		return model.Service{}, userIds, errors.New("service name & accesslevel must be filled")
	}
	if req.Address == "" && req.Method == "" && req.Header == "" && req.Body == "" && req.ExecutionTime == "" {
		return model.Service{}, userIds, c.JSON(http.StatusBadRequest, "Bad request")
	}

	// put params in service
	headerMap := make(map[string]string)
	if req.Header == "" {
		req.Header = "{}"
	}
	if req.Body == "" {
		req.Body = "{}"
	}

	bodyMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(req.Body), &bodyMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	var accessLevelInt int
	var accLevel model.AccessLevel
	if req.AccessLevel != "" {
		accessLevelInt, err = strconv.Atoi(req.AccessLevel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		if accessLevelInt < 0 && accessLevelInt > 2 {
			return model.Service{}, userIds, errors.New("access level value isn't valid")
		}
		accLevel = model.AccessLevel(accessLevelInt)
	} else {
		return model.Service{}, userIds, errors.New("access level must be filled")
	}

	var exeTimeInt64 int64
	if req.ExecutionTime != "" {
		exeTimeInt64, err = strconv.ParseInt(req.ExecutionTime, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	req.AllowedUsers = strings.Replace(req.AllowedUsers, " ", "", -1)
	allowU := strings.Split(req.AllowedUsers, ",")

	for _, username := range allowU {
		userId, err := se.getUsrId(c, username)
		if err != nil {
			return model.Service{}, userIds, errors.New("error while getting user id")
		}

		userIds = append(userIds, userId)
	}

	service = model.Service{
		Name:          &req.Name,
		Address:       &req.Address,
		Method:        &req.Method,
		Header:        headerMap,
		Body:          bodyMap,
		AccessLevel:   accLevel,
		ExecutionTime: &exeTimeInt64,
	}

	return service, userIds, nil
}

func (se *ServicesEndpoints) getUsrId(c echo.Context, username string) (userId int, err error) {
	userId, err = se.IUseruc.GetUsrId(c.Request().Context(), username)
	if err != nil {
		return -1, errors.New("error while getting user id")
	}

	return userId, nil
}
