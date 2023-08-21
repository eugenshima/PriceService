// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"

	"github.com/eugenshima/PriceService/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PriceServiceService struct ....
type PriceServiceService struct {
	rps PriceServiceRepository
}

// NewPriceServiceService creates a new PriceServiceService
func NewPriceServiceService(rps PriceServiceRepository) *PriceServiceService {
	return &PriceServiceService{rps: rps}
}

// PriceServiceRepository interface represents a repository methods
type PriceServiceRepository interface {
	RedisConsumer(context.Context) ([]*model.Share, error)
}

// GetLatestPrice return latest price from redis stream
func (s *PriceServiceService) GetLatestPrice(ctx context.Context) ([]*model.Share, error) {
	return s.rps.RedisConsumer(ctx)
}

// AddSubscription function adds a subscription to concrete price
func (s *PriceServiceService) AddSubscription(ctx context.Context) (map[string]string, error) {
	shares, err := s.rps.RedisConsumer(ctx)
	if err != nil {
		logrus.Errorf("RedisConsumer: %v", err)
	}
	sharesMap := make(map[string]string)
	for _, result := range shares {
		sharesMap[result.ShareName] = result.SharePrice.(string)
	}
	return sharesMap, nil
}

// DeleteSubscription function deletes a subscription from concrete price
func (s *PriceServiceService) DeleteSubscription(ctx context.Context, ID uuid.UUID) error {
	return nil
}
