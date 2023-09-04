package handlers

import (
	"context"
	"os"
	"testing"

	"github.com/eugenshima/PriceService/internal/handlers/mocks"
	"github.com/eugenshima/PriceService/internal/model"
	PriceService "github.com/eugenshima/PriceService/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mockPriceService       *mocks.PriceServiceService
	mockPriceServiceEntity = model.Share{
		ShareName:  "testShare",
		SharePrice: 22.85,
	}
)

// TestMain execute all tests
func TestMain(m *testing.M) {
	mockPriceService = new(mocks.PriceServiceService)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetLatestPricesSubscribe(t *testing.T) {
	mockPriceService.On("Subscribe", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("[]string")).Return(nil).Once()

	mockSelectedShares := make([]string, 0)
	err := mockPriceService.Subscribe(uuid.New(), mockSelectedShares)
	require.NoError(t, err)

	assertion := mockPriceService.AssertExpectations(t)
	require.True(t, assertion)
}

func TestGetLatestPricePublish(t *testing.T) {
	mockPriceService.On("Publish", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uuid.UUID")).Return([]*PriceService.Shares{}, nil).Once()

	result, err := mockPriceService.Publish(context.Background(), uuid.New())
	require.NoError(t, err)
	require.NotNil(t, result)

	assertion := mockPriceService.AssertExpectations(t)
	require.True(t, assertion)
}

func TestGetLatestPriceSend(t *testing.T) {
	mockPriceService.On("CloseSubscription", mock.AnythingOfType("uuid.UUID")).Return(nil).Once()

	err := mockPriceService.CloseSubscription(uuid.New())
	require.NoError(t, err)

	assertion := mockPriceService.AssertExpectations(t)
	require.True(t, assertion)
}
