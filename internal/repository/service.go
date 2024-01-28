package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"monitoring/internal/model"
	"monitoring/pkg/postgres"
)

type IServicesRepository interface {
	Add(ctx context.Context, service model.Service) error
	GetServiceByName(ctx context.Context, service model.Service, roldId int, userId int) (model.Service, error)
	GetServiceById(ctx context.Context, service model.Service, roldId int, userId int) (model.Service, error)
	List(ctx context.Context) ([]model.Service, error)
	Update(ctx context.Context, service model.Service) error
	Delete(ctx context.Context, service model.Service) error
	GetServicesForUser(ctx context.Context, userID int) ([]string, error)
}

type ServicesRepository struct {
	DB postgres.IPostgres
}

var bodyHeader_deserializer = func(header, body string) (map[string]string, map[string]interface{}, error) {
	var headerMap map[string]string
	err := json.Unmarshal([]byte(header), &headerMap)
	if err != nil {
		return nil, nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal([]byte(body), &bodyMap)
	if err != nil {
		return nil, nil, err
	}

	return headerMap, bodyMap, nil
}

func (sr *ServicesRepository) Add(ctx context.Context, service model.Service) error {
	_, err := sr.DB.ExecContext(ctx, `
		INSERT INTO services (name, address, method, header, body,  access_level, execution_time, allowed_users)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, service.Name, service.Address, service.Method, service.Header, service.Body,
		service.AccessLevel, service.ExecutionTime, service.AllowedUsers)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (sr *ServicesRepository) GetServicesForUser(ctx context.Context, userID int) ([]string, error) {
	var services []string
	rows, err := sr.DB.QueryContext(ctx, `
		SELECT s.name
		FROM services s
		JOIN users u ON u.access_level = s.access_level
		WHERE u.id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var serviceName string
		err := rows.Scan(&serviceName)
		if err != nil {
			return nil, err
		}
		services = append(services, serviceName)
	}

	return services, nil
}

func (sr *ServicesRepository) GetServiceByName(ctx context.Context, serviceName string, roleID int, userId int) (serviceRes model.Service, err error) {

	// q := `SELECT name,address,method,header,body,access_level,execution_time,allowed_users FROM services WHERE name = $1 order by id asc limit 1`

	q := `
	SELECT name,address,method,header,body,access_level,execution_time
	FROM services
	WHERE name = $1
	AND (access_level = $2 and $3 = ANY(allowed_users));
	`
	rows, err := sr.DB.QueryContext(ctx, q, serviceName, roleID, userId)
	if err != nil {
		return model.Service{}, err
	}

	var header, body string
	for rows.Next() {
		var service model.Service
		err := rows.Scan(
			&service.Name, &service.Address, &service.Method,
			&header, &body,
			&service.AccessLevel, &service.ExecutionTime,
		)
		if err != nil {
			log.Fatal(err)
		}
		serviceRes = service
	}

	h, b, err := bodyHeader_deserializer(header, body)
	if err != nil {
		return model.Service{}, err
	}
	serviceRes.Header = h
	serviceRes.Body = b

	return serviceRes, nil
}

func (sr *ServicesRepository) GetServiceById(ctx context.Context, serviceId int, userId int, roleId int) (serviceRes model.Service, err error) {

	q := `
	SELECT name,address,method,header,body,access_level,execution_time
	FROM services
	WHERE id = $1
	AND (access_level = $2 and $3 = ANY(allowed_users));
	`

	rows, err := sr.DB.QueryContext(ctx, q, serviceId)
	if err != nil {
		return model.Service{}, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&serviceRes.Name, &serviceRes.Address, &serviceRes.Method, &serviceRes.Header, &serviceRes.Body, &serviceRes.AccessLevel, &serviceRes.ExecutionTime)
		if err != nil {
			return model.Service{}, err
		}
	}

	return serviceRes, nil
}

func (sr *ServicesRepository) List(ctx context.Context) (services []model.Service, err error) {
	q := `SELECT * FROM services`

	rows, err := sr.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var service model.Service
		err := rows.Scan(&service.Name, &service.Address, &service.Method, &service.Header, &service.Body, &service.AccessLevel, &service.ExecutionTime)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (sr *ServicesRepository) Update(ctx context.Context, service model.Service) error {

	// check which field has been passed
	serviceId := service.ID
	serviceName := service.Name
	serviceAddress := service.Address
	serviceMethod := service.Method
	serviceHeader := service.Header
	serviceBody := service.Body
	serviceAccessLevel := service.AccessLevel
	serviceExecutionTime := service.ExecutionTime

	// check which fields have been filled
	var fields []string
	if serviceName != "" {
		fields = append(fields, "name")
	}
	if serviceAddress != "" {
		fields = append(fields, "address")
	}
	if serviceMethod != "" {
		fields = append(fields, "method")
	}
	if serviceHeader != nil {
		fields = append(fields, "header")
	}
	if serviceBody != nil {
		fields = append(fields, "body")
	}
	if serviceAccessLevel >= 0 {
		fields = append(fields, "accesslevel")
	}
	if serviceExecutionTime != 0 {
		fields = append(fields, "executiontime")
	}

	// write query based on fields
	switch {
	case serviceName != "":
		q := `UPDATE services SET `
		for i, field := range fields {
			if i == len(fields)-1 {
				q += field + " = $" + fmt.Sprint(i+1) + " WHERE name = $" + fmt.Sprint(i+2)
			} else {
				q += field + " = $" + fmt.Sprint(i+1) + ", "
			}
		}
		q += ";"

		// write values based on fields
		var values []interface{}
		for _, field := range fields {
			switch field {
			case "name":
				values = append(values, serviceName)
			case "address":
				values = append(values, serviceAddress)
			case "method":
				values = append(values, serviceMethod)
			case "header":
				values = append(values, serviceHeader)
			case "body":
				values = append(values, serviceBody)
			case "accesslevel":
				values = append(values, serviceAccessLevel)
			case "executiontime":
				values = append(values, serviceExecutionTime)
			}
		}
		values = append(values, serviceName)

		_, err := sr.DB.ExecContext(ctx, q, values...)
		if err != nil {
			return err
		}
	case serviceId > 0:
		q := `UPDATE services SET `
		for i, field := range fields {
			if i == len(fields)-1 {
				q += field + " = $" + fmt.Sprint(i+1) + " WHERE id = $" + fmt.Sprint(i+2)
			} else {
				q += field + " = $" + fmt.Sprint(i+1) + ", "
			}
		}
		q += ";"

		// write values based on fields
		var values []interface{}
		for _, field := range fields {
			switch field {
			case "name":
				values = append(values, serviceName)
			case "address":
				values = append(values, serviceAddress)
			case "method":
				values = append(values, serviceMethod)
			case "header":
				values = append(values, serviceHeader)
			case "body":
				values = append(values, serviceBody)
			case "accesslevel":
				values = append(values, serviceAccessLevel)
			case "executiontime":
				values = append(values, serviceExecutionTime)
			}

		}
		values = append(values, serviceId)

		_, err := sr.DB.ExecContext(ctx, q, values...)
		if err != nil {
			return err
		}
	default:
		return errors.New("id or name must be passed")
	}

	return nil
}

func (sr *ServicesRepository) Delete(ctx context.Context, service model.Service) error {

	// check id or name has been passed
	name := service.Name
	id := service.ID

	qByName := `DELETE FROM services WHERE name = $1`
	qById := `DELETE FROM services WHERE id = $1`

	switch {
	case name != "":
		_, err := sr.DB.ExecContext(ctx, qByName, service.Name)
		if err != nil {
			return err
		}
	case id != 0:
		_, err := sr.DB.ExecContext(ctx, qById, service.ID)
		if err != nil {
			return err
		}
	default:
		return errors.New("id or name must be passed")
	}

	return nil
}
