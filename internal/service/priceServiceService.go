// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"fmt"

	"github.com/eugenshima/PriceService/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PriceServiceService struct ....
type PriceServiceService struct {
	rps    PriceServiceRepository
	pubSub *model.PubSub
}

// NewPriceServiceService creates a new PriceServiceService
func NewPriceServiceService(rps PriceServiceRepository, pubSub *model.PubSub) *PriceServiceService {
	return &PriceServiceService{
		rps:    rps,
		pubSub: pubSub,
	}
}

// PriceServiceRepository interface represents a repository methods
type PriceServiceRepository interface {
	RedisConsumer(context.Context) ([]*model.Share, error)
}

// GetLatestPrice return latest price from redis stream
func (ps *PriceServiceService) GetLatestPrice(ctx context.Context) (map[string]float64, error) {
	shares, err := ps.rps.RedisConsumer(ctx)
	if err != nil {
		logrus.Errorf("RedisConsumer: %v", err)
	}
	sharesMap := make(map[string]float64)
	for _, result := range shares {
		sharesMap[result.ShareName] = result.SharePrice.(float64)
	}
	return sharesMap, nil
}

// Subscribe function adds a subscription to ID
func (ps *PriceServiceService) Subscribe(ctx context.Context, ID uuid.UUID) <-chan map[string]float64 {
	ps.pubSub.Mu.Lock()
	defer ps.pubSub.Mu.Unlock()

	responseChan := make(chan map[string]float64, 1)
	ps.pubSub.Subs[ID] = append(ps.pubSub.Subs[ID], responseChan)
	return responseChan
}

// CloseSubscription function deletes a subscription from concrete price
func (ps *PriceServiceService) CloseSubscription(ID uuid.UUID) error {
	ps.pubSub.Mu.Lock()
	defer ps.pubSub.Mu.Unlock()

	if !ps.pubSub.Closed {
		ps.pubSub.Closed = true
		for _, subs := range ps.pubSub.Subs {
			for _, responseChan := range subs {
				close(responseChan)
			}
		}
	}

	return nil
}

// Publish function publishes info to channel
func (ps *PriceServiceService) Publish(ctx context.Context, ID uuid.UUID) error {
	ps.pubSub.Mu.RLock()
	defer ps.pubSub.Mu.RUnlock()

	shares, err := ps.rps.RedisConsumer(ctx)
	if err != nil {
		logrus.Errorf("RedisConsumer: %v", err)
		return fmt.Errorf("RedisConsumer: %w", err)
	}
	sharesMap := make(map[string]float64)
	for _, result := range shares {
		sharesMap[result.ShareName] = result.SharePrice.(float64)
	}

	for _, responseChan := range ps.pubSub.Subs[ID] {
		responseChan <- sharesMap
	}

	return nil
}
