// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"fmt"

	"github.com/eugenshima/PriceService/internal/model"
	proto "github.com/eugenshima/PriceService/proto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate /home/yauhenishymanski/work/bin/mockery --name=PriceServiceRepository --case=underscore --output=./mocks

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
func (ps *PriceServiceService) Subscribe(ID uuid.UUID, selectedShares []string) error {
	ps.pubSub.Mu.Lock()
	defer ps.pubSub.Mu.Unlock()
	if _, ok := ps.pubSub.Subs[ID]; !ok {
		ps.pubSub.Subs[ID] = selectedShares
		ps.pubSub.SubsShares[ID] = make(chan []*model.Share, 1)
		return nil
	}

	return fmt.Errorf("error subscribing to ID %v", ID)
}

// CloseSubscription function deletes a subscription from concrete price
func (ps *PriceServiceService) CloseSubscription(ID uuid.UUID) error {
	ps.pubSub.Mu.Lock()
	defer ps.pubSub.Mu.Unlock()

	if !ps.pubSub.Closed[ID] {
		for closedID, closedSub := range ps.pubSub.SubsShares {
			if closedID == ID {
				logrus.Info("Subscription closed")
				delete(ps.pubSub.SubsShares, ID)
				delete(ps.pubSub.Subs, ID)
				delete(ps.pubSub.Closed, ID)
				close(closedSub)
			}

		}
	}

	return nil
}

// PublishToAllSubscribers function ....
func (ps *PriceServiceService) PublishToAllSubscribers(ctx context.Context) {
	for {
		if len(ps.pubSub.Subs) == 0 {
			continue
		}
		allShares, err := ps.rps.RedisConsumer(ctx)
		if err != nil {
			logrus.Errorf("RedisConsumer: %v", err)
		}
		ps.pubSub.Mu.Lock()
		for ID, selectedShares := range ps.pubSub.Subs {
			shares := make([]*model.Share, 0)
			for _, share := range allShares {
				for _, selectedShare := range selectedShares {
					if share.ShareName == selectedShare {
						shares = append(shares, share)
						break
					}
				}
			}

			select {
			case <-ctx.Done():
				logrus.Info("stream ended (ctx done)")
				return
			default:
				if !ps.pubSub.Closed[ID] {
					ps.pubSub.SubsShares[ID] <- shares
				} else {
					logrus.Info("stream ended (subscription closed)")
					return
				}
			}
		}
		ps.pubSub.Mu.Unlock()
	}
}

// Publish function ....
func (ps *PriceServiceService) Publish(ctx context.Context, ID uuid.UUID) ([]*proto.Shares, error) {
	select {
	case <-ctx.Done():
		logrus.Info("context closed")
		return nil, nil
	case shares := <-ps.pubSub.SubsShares[ID]:
		res := []*proto.Shares{}
		for _, share := range shares {
			protoShare := &proto.Shares{
				ShareName:  share.ShareName,
				SharePrice: share.SharePrice.(float64),
			}
			res = append(res, protoShare)
		}
		return res, nil
	}
}
