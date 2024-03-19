package auth

import (
	"context"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/entity"
)

type userService struct {
	client api.AuthServiceClient
}

func NewUserService(client api.AuthServiceClient) (*userService, error) {
	return &userService{
		client: client,
	}, nil
}

func (s *userService) SignUp(ctx context.Context, user *entity.User) (*api.TokenResponse, error) {
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

func (s *userService) SignIn(ctx context.Context, user *entity.User) (*api.TokenResponse, error) {
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

func (s *userService) Refresh(ctx context.Context, refresh string) (*api.TokenResponse, error) {
	req := &api.RefreshRequest{
		RefreshToken: refresh,
	}

	res, err := s.client.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userService) Auth(ctx context.Context, access string, roleId int64) (bool, error) {
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

func (s *userService) Logout(ctx context.Context, access string) error {
	req := &api.LogoutRequest{
		AccessToken: access,
	}

	if _, err := s.client.Logout(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *userService) List(ctx context.Context, limit, offset int32) (*api.ListUserResponse, error) {
	req := &api.Bounds{
		Limit:  limit,
		Offset: offset,
	}

	res, err := s.client.ListUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *userService) ListByRole(ctx context.Context, roleId int64, limit, offset int32) (*api.ListUserResponse, error) {
	bounds := &api.Bounds{
		Limit:  limit,
		Offset: offset,
	}

	req := &api.ListUserByRoleRequest{
		RoleId: roleId,
		Bounds: bounds,
	}

	res, err := s.client.ListUserByRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (s *userService) FindById(ctx context.Context, id int64) (*api.UserDetails, error) {
	req := &api.FindUserByIdRequest{
		Id: id,
	}

	user, err := s.client.FindUserById(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
