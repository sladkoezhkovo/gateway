package auth

import (
	"context"
	"fmt"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct {
	client api.AuthServiceClient
}

func New(cfg config.RemoteConfig) (*service, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &service{
		client: api.NewAuthServiceClient(conn),
	}, nil
}

func (s *service) SignUp(ctx context.Context, user *entity.User) (*api.TokenResponse, error) {
	req := &api.SignUpRequest{
		Email:    user.Email,
		Password: user.Password,
		RoleId:   user.Role.Id,
	}

	res, err := s.client.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) SignIn(ctx context.Context, user *entity.User) (*api.TokenResponse, error) {
	req := &api.SignInRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	res, err := s.client.SignIn(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) Refresh(ctx context.Context, refresh string) (*api.TokenResponse, error) {
	req := &api.RefreshRequest{
		RefreshToken: refresh,
	}

	res, err := s.client.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) Auth(ctx context.Context, access string, roleId int64) (bool, error) {
	req := &api.AuthRequest{
		AccessToken: access,
		RoleId:      roleId,
	}

	res, err := s.client.Auth(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Approved, nil
}

func (s *service) Logout(ctx context.Context, access string) error {
	req := &api.LogoutRequest{
		AccessToken: access,
	}

	if _, err := s.client.Logout(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *service) List(ctx context.Context, limit, offset int32) (*api.ListUserResponse, error) {
	req := &api.Bounds{
		Limit:  limit,
		Offset: offset,
	}

	fmt.Println(req)

	res, err := s.client.ListUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
