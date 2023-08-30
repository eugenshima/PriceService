// Package handlers handles gRPC requests and returns stream of current prices
package handlers

import (
	"context"
	"fmt"

	"github.com/eugenshima/PriceService/internal/model"
	proto "github.com/eugenshima/PriceService/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate /home/yauhenishymanski/work/bin/mockery --name=PriceServiceService --case=underscore --output=./mocks

// PriceServiceHandler struct ....
type PriceServiceHandler struct {
	srv PriceServiceService
	proto.UnimplementedPriceServiceServer
}

// NewPriceServiceHandler creates a new PriceServiceHandler
func NewPriceServiceHandler(srv PriceServiceService) *PriceServiceHandler {
	return &PriceServiceHandler{
		srv: srv,
	}
}

// PriceServiceService is an interface for accessing PriceService
type PriceServiceService interface {
	GetLatestPrice(ctx context.Context) (map[string]float64, error)
	Subscribe(context.Context, uuid.UUID) <-chan *model.Share
	Publish(context.Context, uuid.UUID) error
	CloseSubscription(uuid.UUID) error
}

// GetLatestPrice function receives request to get current prices
func (s *PriceServiceHandler) GetLatestPrices(req *proto.LatestPriceRequest, stream proto.PriceService_GetLatestPricesServer) error {
	ID := uuid.New()
	response := s.srv.Subscribe(stream.Context(), ID)

	for {
		select {
		case <-stream.Context().Done():
			err := s.srv.CloseSubscription(ID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"ID": ID}).Errorf("CloseSubscription: %v", err)
				return fmt.Errorf("CloseSubscription: %w", err)
			}
			logrus.Info("Stream is probably ended :D")
			return stream.Context().Err()
		case <-response:
			sharesMap, err := s.srv.GetLatestPrice(stream.Context())
			if err != nil {
				logrus.Errorf("GetLatestPrice: %v", err)
				return fmt.Errorf("GetLatestPrice: %w", err)
			}
			res := []*proto.Shares{}
			for key, value := range sharesMap {
				for i := 0; i < len(req.ShareName); i++ {
					if key == req.ShareName[i] {
						share := &proto.Shares{
							ShareName:  key,
							SharePrice: value,
						}
						res = append(res, share)
					}
				}

			}
			stream.Send(&proto.LatestPriceResponse{
				Shares: res,
				ID:     ID.String(),
			})
		default:
			err := s.srv.Publish(stream.Context(), ID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"ID": ID}).Errorf("Publish: %v", err)
				return fmt.Errorf("publish: %w", err)
			}
		}
	}
}
