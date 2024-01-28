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
	ServicesUC usecase.ServicesUsecase
}

func NewServicesEndpoints() *ServicesEndpoints {
	var serviceuc usecase.ServicesUsecase = usecase.ServicesUsecase{
		ServicesRepo: repository.ServicesRepository{DB: GlobalPG},
		Useruc:       useruc.UserUsecase{DB: &userrepo.UserRepo{DB: GlobalPG}}}
	return &ServicesEndpoints{
		ServicesUC: serviceuc,
	}
}

func (se *ServicesEndpoints) ListServices(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtCustomClaims)
	if !claims.Admin {
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

	service, err := se.checkPostParams(c)
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

	id := c.FormValue("id")
	name := c.FormValue("name")
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

	var err error

	tokenString := c.Get("user").(*jwt.Token)
	// claims := jwtuserToken.Claims.(*model.JwtCustomClaims)
	token, err := jwt.Parse(tokenString.Raw, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	jwtToken := make(map[string]interface{})
	jwtToken = token.Claims.(jwt.MapClaims)

	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*model.JwtCustomClaims)

	type RequestBody struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var requestBody RequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if requestBody.ID == "" && requestBody.Name == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}

	var idInt int
	service := model.Service{}
	if requestBody.ID != "" {
		idInt, err = strconv.Atoi(requestBody.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		if idInt <= 0 {
			return c.JSON(http.StatusBadRequest, "Bad request")
		}
		service.ID = idInt
	}

	if requestBody.Name != "" {
		service.Name = requestBody.Name
	}

	userId, err := se.getUsrId(c, jwtToken["name"].(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	svc, err := se.ServicesUC.Get(c.Request().Context(), service, int(jwtToken["role"].(float64)), userId)

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
		Name:          name,
		Address:       address,
		Method:        method,
		Header:        headerMap,
		Body:          bodyMap,
		AccessLevel:   accLevel,
		ExecutionTime: exeTimeInt64,
	}

	return service, nil
}

func (se *ServicesEndpoints) checkPostParams(c echo.Context) (service model.Service, err error) {
	name := c.FormValue("name")
	address := c.FormValue("address")
	method := c.FormValue("method")
	header := c.FormValue("header")
	body := c.FormValue("body")
	accesslevel := c.FormValue("access_level")
	executiontime := c.FormValue("execution_time")
	allowedusers := c.FormValue("allowed_users")

	// check for empty params
	if name == "" || accesslevel == "" {
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

	var exeTimeInt64 int64
	if executiontime != "" {
		exeTimeInt64, err = strconv.ParseInt(executiontime, 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	allowedusers = strings.Replace(allowedusers, " ", "", -1)
	allowU := strings.Split(allowedusers, ",")

	var userIds []int
	for _, username := range allowU {
		userId, err := se.getUsrId(c, username)
		if err != nil {
			return model.Service{}, errors.New("error while getting user id")
		}

		userIds = append(userIds, userId)
	}

	service = model.Service{
		Name:          name,
		Address:       address,
		Method:        method,
		Header:        headerMap,
		Body:          bodyMap,
		AccessLevel:   accLevel,
		ExecutionTime: exeTimeInt64,
		AllowedUsers:  userIds,
	}

	return service, nil
}

func (se *ServicesEndpoints) getUsrId(c echo.Context, username string) (userId int, err error) {
	userId, err = se.ServicesUC.Useruc.DB.GetUsrId(c.Request().Context(), username)
	if err != nil {
		return -1, errors.New("error while getting user id")
	}

	return userId, nil
}
