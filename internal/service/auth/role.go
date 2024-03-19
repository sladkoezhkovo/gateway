package auth

import (
	"context"
	api "github.com/sladkoezhkovo/gateway/api/auth"
)

type roleService struct {
	client api.AuthServiceClient
}

func NewRoleService(client api.AuthServiceClient) (*roleService, error) {
	return &roleService{
		client: client,
	}, nil
}

func (s *roleService) List(ctx context.Context, limit, offset int32) (*api.ListRoleResponse, error) {
	req := &api.Bounds{
		Limit:  limit,
		Offset: offset,
	}

	res, err := s.client.ListRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *roleService) FindById(ctx context.Context, id int64) (*api.Role, error) {
	req := &api.FindRoleByIdRequest{
		Id: id,
	}

	role, err := s.client.FindByIdRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return role, nil
}
