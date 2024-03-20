package admin

import (
	"context"
	"github.com/sladkoezhkovo/gateway/api/admin"
)

type cityService struct {
	client admin.AdminServiceClient
}

func NewCityService(client admin.AdminServiceClient) (*cityService, error) {
	return &cityService{
		client: client,
	}, nil
}

func (s *cityService) List(ctx context.Context, limit, offset int32) (*admin.ListCityResponse, error) {
	req := &admin.ListRequest{
		Limit:  limit,
		Offset: offset,
	}

	res, err := s.client.ListCity(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *cityService) FindById(ctx context.Context, id int64) (*admin.City, error) {
	req := &admin.FindByIdRequest{
		Id: id,
	}

	role, err := s.client.FindByIdCity(ctx, req)
	if err != nil {
		return nil, err
	}

	return role, nil

}
