// Package handlers handles gRPC requests and returns stream of current prices
package handlers

import (
	"context"
	"fmt"

	proto "github.com/eugenshima/PriceService/proto"

	vld "github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate /home/yauhenishymanski/work/bin/mockery --name=PriceServiceService --case=underscore --output=./mocks

// PriceServiceHandler struct ....
type PriceServiceHandler struct {
	srv PriceServiceService
	vl  *vld.Validate
	proto.UnimplementedPriceServiceServer
}

// NewPriceServiceHandler creates a new PriceServiceHandler
func NewPriceServiceHandler(srv PriceServiceService, vl *vld.Validate) *PriceServiceHandler {
	return &PriceServiceHandler{
		srv: srv,
		vl:  vl,
	}
}

// PriceServiceService is an interface for accessing PriceService
type PriceServiceService interface {
	GetLatestPrice(ctx context.Context) (map[string]float64, error)
	Subscribe(uuid.UUID, []string) error
	PublishToAllSubscribers(ctx context.Context)
	Publish(context.Context, uuid.UUID) ([]*proto.Shares, error)
	CloseSubscription(uuid.UUID) error
}

// CustomValidation function for custom validation
func (ph *PriceServiceHandler) CustomValidation(ctx context.Context, i interface{}) error {
	return nil
}

// GetLatestPrices function receives request to get current prices
func (ph *PriceServiceHandler) GetLatestPrices(req *proto.LatestPriceRequest, stream proto.PriceService_GetLatestPricesServer) error {
	ID := uuid.New()

	err := ph.srv.Subscribe(ID, req.ShareName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ID": ID, "req.ShareName": req.ShareName}).Errorf("Subscribe: %v", err)
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			protoShares, err := ph.srv.Publish(stream.Context(), ID)
			if err != nil {
				deletionErr := ph.srv.CloseSubscription(ID)
				if deletionErr != nil {
					logrus.Errorf("CloseSubscription: %v", err)
					return fmt.Errorf("CloseSubscription: %w", err)
				}
				logrus.WithFields(logrus.Fields{"ID": ID}).Errorf("Publish: %v", err)
				return fmt.Errorf("publish: %w", err)
			}
			err = stream.Send(&proto.LatestPriceResponse{
				Shares: protoShares,
				ID:     ID.String(),
			})
			if err != nil {
				deletionErr := ph.srv.CloseSubscription(ID)
				if deletionErr != nil {
					logrus.Errorf("CloseSubscription: %v", err)
					return fmt.Errorf("CloseSubscription: %w", err)
				}
			}

		}

	}

}
