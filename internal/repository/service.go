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
	Add(ctx context.Context, service model.Service, userIds []int) error
	GetUserService(ctx context.Context, serviceName string, userID, roleId int) (service model.Service, err error)
	GetUserServices(ctx context.Context, roleID int, userId int) (serviceRes []model.Service, err error)
	List(ctx context.Context) ([]model.Service, error)
	Update(ctx context.Context, service model.Service) error
	Delete(ctx context.Context, service model.Service) error
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

func (sr *ServicesRepository) Add(ctx context.Context, service model.Service, userIds []int) error {
	_, err := sr.DB.ExecContext(ctx, `
		INSERT INTO services (name, address, method, header, body,  access_level, execution_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, service.Name, service.Address, service.Method, service.Header, service.Body,
		service.AccessLevel, service.ExecutionTime)
	if err != nil {
		log.Fatal(err)
	}

	// insert user_services
	for _, userId := range userIds {
		_, err := sr.DB.ExecContext(ctx, `
			INSERT INTO user_services (user_id, service_id)
			VALUES ($1, (SELECT id FROM services WHERE name = $2))`, userId, service.Name)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (sr *ServicesRepository) GetUserService(ctx context.Context, serviceName string, userID, roleId int) (service model.Service, err error) {
	row, err := sr.DB.QueryContext(ctx, `
		select s.name, s.address, s.method, header, body ,s.access_level, s.execution_time, s.error_estimate	
		from services s
		join user_services us on s.id = us.service_id
		where s.name=$1 and us.user_id=$2 and (s.access_level <=$3 OR s.access_level = 1);
	`, serviceName, userID, roleId)
	if err != nil {
		return
	}
	defer row.Close()
	var header, body *string
	for row.Next() {
		err := row.Scan(
			&service.Name, &service.Address, &service.Method, &header, &body,
			&service.AccessLevel, &service.ExecutionTime, &service.ErrorEstimate,
		)
		if err != nil {
			log.Fatal(err)
		}

		if header == nil {
			header = new(string)
			*header = "{}"
		}
		if body == nil {
			body = new(string)
			*body = "{}"
		}

		h, b, err := bodyHeader_deserializer(*header, *body)
		if err != nil {
			return model.Service{}, err
		}

		service.Header = h
		service.Body = b
	}

	return
}

func (sr *ServicesRepository) GetUserServices(ctx context.Context, roleID int, userId int) (serviceRes []model.Service, err error) {

	q := `
		select s.name, s.address, s.method, header, body ,s.access_level, s.execution_time
		from services s
		join user_services us on s.id = us.service_id
		where us.user_id = $1 and (s.access_level <= $2 OR s.access_level = 1);
	`
	rows, err := sr.DB.QueryContext(ctx, q, userId, roleID)
	if err != nil {
		return []model.Service{}, err
	}

	var header, body *string
	defer rows.Close()

	for rows.Next() {
		var service model.Service
		err := rows.Scan(
			&service.Name, &service.Address, &service.Method, &header, &body,
			&service.AccessLevel, &service.ExecutionTime,
		)
		if err != nil {
			log.Fatal(err)
		}

		if header == nil {
			header = new(string)
			*header = "{}"
		}
		if body == nil {
			body = new(string)
			*body = "{}"
		}

		h, b, err := bodyHeader_deserializer(*header, *body)
		if err != nil {
			return []model.Service{}, err
		}
		service.Header = h
		service.Body = b
		serviceRes = append(serviceRes, service)
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
	if *serviceName != "" {
		fields = append(fields, "name")
	}
	if *serviceAddress != "" {
		fields = append(fields, "address")
	}
	if *serviceMethod != "" {
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
	if *serviceExecutionTime != 0 {
		fields = append(fields, "executiontime")
	}

	// write query based on fields
	switch {
	case *serviceName != "":
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
	q := `DELETE FROM user_services WHERE service_id = (SELECT id FROM services WHERE name = $1);`
	_, err := sr.DB.ExecContext(ctx, q, service.Name)
	if err != nil {
		return err
	}
	q = `DELETE FROM services WHERE name = $1;`
	_, err = sr.DB.ExecContext(ctx, q, service.Name)
	if err != nil {
		return err
	}
	return nil
}
