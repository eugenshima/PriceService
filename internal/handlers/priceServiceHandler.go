// Package handlers handles gRPC requests and returns current prices
package handlers

import (
	"context"

	"github.com/eugenshima/PriceService/internal/model"
	proto "github.com/eugenshima/PriceService/proto"

	"github.com/sirupsen/logrus"
)

// PriceServiceHandler struct ....
type PriceServiceHandler struct {
	srv PriceServiceService
	proto.UnimplementedPriceServiceServer
}

// NewPriceServiceHandler creates a new PriceServiceHandler
func NewPriceServiceHandler(srv PriceServiceService) *PriceServiceHandler {
	return &PriceServiceHandler{srv: srv}
}

// PriceServiceService is an interface for accessing PriceService
type PriceServiceService interface {
	GetLatestPrice(context.Context) ([]*model.Share, error)
}

// GetLatestPrice function receives request to Get user from database by ID
func (s *PriceServiceHandler) GetLatestPrice(ctx context.Context, _ *proto.LatestPriceRequest) (*proto.LatestPriceResponse, error) {
	results, err := s.srv.GetLatestPrice(ctx)
	if err != nil {
		logrus.Errorf("GetAll: %v", err)
		return nil, err
	}
	res := []*proto.Shares{}
	for _, result := range results {
		share := &proto.Shares{
			ShareName:  result.ShareName,
			SharePrice: result.SharePrice.(string),
		}
		res = append(res, share)
	}
	return &proto.LatestPriceResponse{Shares: res}, nil
}
