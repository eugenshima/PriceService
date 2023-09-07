// Package handlers handles gRPC requests and returns stream of current prices
package handlers

import (
	"context"
	"fmt"

	proto "github.com/eugenshima/price-service/proto"

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
	AddSubscriber(uuid.UUID, []string) error
	CloseSubscription(uuid.UUID) error
	Publish(context.Context, uuid.UUID) ([]*proto.Shares, error)
	PublishToAllSubscribers(ctx context.Context)
}

// CustomValidation function for custom validation
func (ph *PriceServiceHandler) CustomValidation(ctx context.Context, i interface{}) error {
	err := ph.vl.VarCtx(ctx, i, "required")
	if err != nil {
		return fmt.Errorf("VarCtx: %w", err)
	}
	return nil
}

// Subscribe function adds subscriber
func (ph *PriceServiceHandler) Subscribe(req *proto.SubscribeRequest, stream proto.PriceService_SubscribeServer) error {
	err := ph.CustomValidation(stream.Context(), req.ShareName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"req.ShareName": req.ShareName}).Errorf("CustomValidate: %v", err)
		return fmt.Errorf("validate: %w", err)
	}
	ID := uuid.New()
	err = ph.srv.AddSubscriber(ID, req.ShareName)
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
				return err
			}
			err = stream.Send(&proto.SubscribeResponse{
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

// GetLatestPrices function receives request to get current prices
func (ph *PriceServiceHandler) SubscribeOld(req *proto.SubscribeRequest, stream proto.PriceService_SubscribeServer) error {
	err := ph.CustomValidation(stream.Context(), req.ShareName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"req.ShareName": req.ShareName}).Errorf("CustomValidate: %v", err)
		return fmt.Errorf("validate: %w", err)
	}
	fmt.Println("gg")
	ID := uuid.New()
	err = ph.srv.AddSubscriber(ID, req.ShareName)
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
			err = stream.Send(&proto.SubscribeResponse{
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
