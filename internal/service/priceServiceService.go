// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"

	"github.com/eugenshima/PriceService/internal/model"
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
