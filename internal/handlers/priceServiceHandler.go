// Package handlers handles gRPC requests and returns stream of current prices
package handlers

import (
	"context"
	"fmt"

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
	Subscribe(uuid.UUID) <-chan map[string]float64
	Publish(context.Context, uuid.UUID) error
	CloseSubscription() error
}

// GetLatestPrices function receives request to get current prices
func (ph *PriceServiceHandler) GetLatestPrices(req *proto.LatestPriceRequest, stream proto.PriceService_GetLatestPricesServer) error {
	ID := uuid.New()
	responseChan := ph.srv.Subscribe(ID)

	for {
		select {
		case <-stream.Context().Done():
			err := ph.srv.CloseSubscription()
			if err != nil {
				logrus.WithFields(logrus.Fields{"ID": ID}).Errorf("CloseSubscription: %v", err)
				return fmt.Errorf("CloseSubscription: %w", err)
			}
			logrus.Info("Stream is probably ended :D")
			return stream.Context().Err()
		case sharesMap := <-responseChan:
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
			err := stream.Send(&proto.LatestPriceResponse{
				Shares: res,
				ID:     ID.String(),
			})
			if err != nil {
				logrus.Errorf("Send: %v", err)
				return fmt.Errorf("send: %w", err)
			}
		default:
			err := ph.srv.Publish(stream.Context(), ID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"ID": ID}).Errorf("Publish: %v", err)
				return fmt.Errorf("publish: %w", err)
			}
		}
	}
}
