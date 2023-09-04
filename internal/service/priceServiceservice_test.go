package service

import (
	"context"
	"os"
	"testing"

	"github.com/eugenshima/PriceService/internal/model"
	"github.com/eugenshima/PriceService/internal/service/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockPriceServiceRepository *mocks.PriceServiceRepository

// TestMain execute all tests
func TestMain(m *testing.M) {
	mockPriceServiceRepository = new(mocks.PriceServiceRepository)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestPublishToAllSubscribers(t *testing.T) {
	mockPriceServiceRepository.On("RedisConsumer", mock.AnythingOfType("*context.emptyCtx")).Return([]*model.Share{}, nil).Once()

	shares, err := mockPriceServiceRepository.RedisConsumer(context.Background())
	require.NoError(t, err)
	require.NotNil(t, shares)

	assertion := mockPriceServiceRepository.AssertExpectations(t)
	require.True(t, assertion)
}
