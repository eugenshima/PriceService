// Package handlers handles gRPC requests and returns current prices
package handlers

import (
	"context"
	"fmt"
	"sync"

	"github.com/eugenshima/PriceService/internal/model"
	proto "github.com/eugenshima/PriceService/proto"

	"github.com/sirupsen/logrus"
)

//go:generate /home/yauhenishymanski/work/bin/mockery --name=PriceServiceService --case=underscore --output=./mocks

// PriceServiceHandler struct ....
type PriceServiceHandler struct {
	srv PriceServiceService
	mu  *sync.RWMutex
	proto.UnimplementedPriceServiceServer
}

// NewPriceServiceHandler creates a new PriceServiceHandler
func NewPriceServiceHandler(srv PriceServiceService, mu *sync.RWMutex) *PriceServiceHandler {
	return &PriceServiceHandler{
		srv: srv,
		mu:  mu,
	}
}

// PriceServiceService is an interface for accessing PriceService
type PriceServiceService interface {
	GetLatestPrice(context.Context) ([]*model.Share, error)
	AddSubscription(context.Context) (map[string]string, error)
}

// GetLatestPrice function receives request to get current prices
func (s *PriceServiceHandler) GetLatestPrices(req *proto.LatestPriceRequest, stream proto.PriceService_GetLatestPricesServer) error {
	for {
		select {
		case <-stream.Context().Done():
			logrus.Info("Stream is probably ended :D")
			return stream.Context().Err()
		default:
			results, err := s.srv.AddSubscription(stream.Context())
			if err != nil {
				logrus.Errorf("GetLatestPrice: %v", err)
				return fmt.Errorf("GetLatestPrice: %w", err)
			}
			res := []*proto.Shares{}
			for key, value := range results {
				if key == req.ShareName {
					share := &proto.Shares{
						ShareName:  key,
						SharePrice: value,
					}
					res = append(res, share)
				}
			}
			stream.Send(&proto.LatestPriceResponse{Shares: res})
		}
	}
}
