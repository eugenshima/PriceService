package handlers

import (
	"context"
	"os"
	"testing"

	"github.com/eugenshima/PriceService/internal/handlers/mocks"
	"github.com/eugenshima/PriceService/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var mockPriceService *mocks.PriceServiceService

// TestMain execute all tests
func TestMain(m *testing.M) {
	mockPriceService = new(mocks.PriceServiceService)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetLatestPrice(t *testing.T) {
	mockPriceService.On("GetLatestPrice", mock.Anything).Return([]*model.Share{}, nil).Twice()
	handler := NewPriceServiceHandler(mockPriceService)
	res, err := mockPriceService.GetLatestPrice(context.Background())
	require.NoError(t, err)
	results, err := handler.srv.GetLatestPrice(context.Background())
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Equal(t, len(res), len(results))
}
