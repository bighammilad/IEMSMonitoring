package repository

import (
	"context"
	"monitoring/internal/model"
	"monitoring/pkg/postgres"
)

type IUserServiceRepo interface {
	Add(ctx context.Context, userservice model.UserService) error
	GetUserService(ctx context.Context, userservice model.UserService) (model.UserService, error)
	DeleteUserService(ctx context.Context, userservice model.UserService) error
	GetServiceUsers(ctx context.Context, userservice model.UserService) ([]model.UserService, error)
	GetUserServices(ctx context.Context, userservice model.UserService) ([]model.UserService, error)
	List(ctx context.Context) ([]model.UserService, error)
	Update(ctx context.Context, newUserService model.UserService) error
	DeleteUserServices(ctx context.Context, usrservices []model.UserService) error
	AddUserServices(ctx context.Context, usrservices []model.UserService) error
}

type UserServiceRepository struct {
	DB postgres.IPostgres
}

func (Us *UserServiceRepository) Add(ctx context.Context, userservice model.UserService) error {
	_, err := Us.DB.ExecContext(ctx, `
		INSERT INTO user_services (service_id, user_id)
		VALUES ($1, $2)`, userservice.ServiceID, userservice.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (Us *UserServiceRepository) GetUserService(ctx context.Context, userservice model.UserService) (model.UserService, error) {
	rows, err := Us.DB.QueryContext(ctx, `
		SELECT service_id, user_id
		FROM user_services 
		WHERE service_id = $1 AND user_id = $2`, userservice.ServiceID, userservice.UserID)
	if err != nil {
		return userservice, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userservice.ServiceID, &userservice.UserID)
		if err != nil {
			return userservice, err
		}
	}
	return userservice, nil
}

func (Us *UserServiceRepository) DeleteUserService(ctx context.Context, userservice model.UserService) error {
	_, err := Us.DB.ExecContext(ctx, `
		DELETE FROM user_services
		WHERE service_id = $1 AND user_id = $2`, userservice.ServiceID, userservice.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (Us *UserServiceRepository) GetServiceUsers(ctx context.Context, userservice model.UserService) ([]model.UserService, error) {
	var userServices []model.UserService
	rows, err := Us.DB.QueryContext(ctx, `
		SELECT service_id, user_id
		FROM user_services 
		WHERE service_id = $1`, userservice.ServiceID)
	if err != nil {
		return userServices, err
	}
	defer rows.Close()
	for rows.Next() {
		var userService model.UserService
		err := rows.Scan(&userService.ServiceID, &userService.UserID)
		if err != nil {
			return userServices, err
		}
		userServices = append(userServices, userService)
	}
	return userServices, nil
}

func (Us *UserServiceRepository) GetUserServices(ctx context.Context, userservice model.UserService) ([]model.UserService, error) {
	var userServices []model.UserService
	rows, err := Us.DB.QueryContext(ctx, `
		SELECT service_id, user_id
		FROM user_services 
		WHERE user_id = $1`, userservice.UserID)
	if err != nil {
		return userServices, err
	}
	defer rows.Close()
	for rows.Next() {
		var userService model.UserService
		err := rows.Scan(&userService.ServiceID, &userService.UserID)
		if err != nil {
			return userServices, err
		}
		userServices = append(userServices, userService)
	}
	return userServices, nil
}

func (Us *UserServiceRepository) List(ctx context.Context) ([]model.UserService, error) {
	var userServices []model.UserService
	rows, err := Us.DB.QueryContext(ctx, `
		SELECT service_id, user_id
		FROM user_services`)
	if err != nil {
		return userServices, err
	}
	defer rows.Close()
	for rows.Next() {
		var userService model.UserService
		err := rows.Scan(&userService.ServiceID, &userService.UserID)
		if err != nil {
			return userServices, err
		}
		userServices = append(userServices, userService)
	}
	return userServices, nil
}

func (Us *UserServiceRepository) Update(ctx context.Context, newUserService model.UserService) (err error) {
	_, err = Us.DB.ExecContext(ctx, `
		UPDATE user_services SET user_id = $1 WHERE service_id = $2;`,
		newUserService.UserID, newUserService.ServiceID)
	if err != nil {
		return
	}
	return nil

}

func (Us *UserServiceRepository) DeleteUserServices(ctx context.Context, usrservices []model.UserService) (err error) {
	for _, us := range usrservices {
		err = Us.DeleteUserService(ctx, us)
		if err != nil {
			return err
		}
	}
	return nil
}

func (Us *UserServiceRepository) AddUserServices(ctx context.Context, usrservices []model.UserService) (err error) {
	for _, us := range usrservices {
		err = Us.Add(ctx, us)
		if err != nil {
			return err
		}
	}
	return nil
}
