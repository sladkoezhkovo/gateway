package auth

import (
	"context"
	"github.com/sladkoezhkovo/gateway/internal/entity"
)

type Service interface {
	SignUp(ctx context.Context, user *entity.User) (*entity.Tokens, error)
	SignIn(ctx context.Context, user *entity.User) (*entity.Tokens, error)
	Refresh(ctx context.Context, refresh string) (*entity.Tokens, error)
	Auth(ctx context.Context, access string) (*entity.Tokens, error)
	Logout(ctx context.Context, access string) error
}

type handler struct {
	service Service
}

func New(service Service) *handler {
	return &handler{
		service: service,
	}
}
