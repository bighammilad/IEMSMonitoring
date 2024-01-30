package useruc

import (
	"context"
	"monitoring/internal/repository/userrepo"
)

type ILogin interface {
	Login(ctx context.Context, username, password string) (*LoginUC, error)
}

type LoginUC struct {
	Username   string
	Role       int
	ILoginRepo userrepo.ILoginRepo
}

func (l *LoginUC) Login(ctx context.Context, username, password string) (*LoginUC, error) {

	user, err := l.ILoginRepo.Auth(ctx, username, password)
	if err != nil {
		return &LoginUC{}, err
	}

	return &LoginUC{
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
