package services

import (
	"context"
	"errors"
	"search-country-system/internal/mocks"
	"testing"

	"search-country-system/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCountryData_CacheMiss_APISuccess_WithNoCountryData(t *testing.T) {
	mockCache := new(mocks.ICacheService)
	mockService := new(mocks.IServiceCall)

	apiResponse := `[]`

	// Cache miss
	mockCache.On("Get", "india").Return(nil, false)

	// Only ONE mock for API call
	mockService.
		On("Get", mock.Anything, "/v3.1/name/india?fullText=true", mock.Anything).
		Return([]byte(apiResponse), 200, nil).
		Once()

	service := InitCountryService(mockService, mockCache)

	result, err := service.GetCountryData(context.Background(), "india")

	// Expect no data scenario
	assert.Error(t, err)
	assert.Equal(t, dto.Country{}, result)

	mockCache.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetCountryData_CacheHit(t *testing.T) {
	mockCache := mocks.NewICacheService(t)
	mockService := mocks.NewIServiceCall(t)

	expected := dto.Country{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "INR",
		Population: 100,
	}

	mockCache.On("Get", "india").Return(expected, true)

	service := InitCountryService(mockService, mockCache)

	result, err := service.GetCountryData(context.Background(), "india")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockCache.AssertExpectations(t)
}

func TestGetCountryData_CacheMiss_APISuccess(t *testing.T) {
	mockCache := new(mocks.ICacheService)
	mockService := new(mocks.IServiceCall)

	apiResponse := `[{
		"name": {"common": "India"},
		"capital": ["New Delhi"],
		"population": 1400000000,
		"currencies": {"INR": {}}
	}]`

	mockCache.On("Get", "india").Return(nil, false)
	mockCache.On("Set", "india", mock.Anything).Return()

	mockService.
		On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return([]byte(apiResponse), 200, nil)

	service := InitCountryService(mockService, mockCache)

	result, err := service.GetCountryData(context.Background(), "india")

	assert.NoError(t, err)
	assert.Equal(t, "India", result.Name)
	assert.Equal(t, "New Delhi", result.Capital)
	assert.Equal(t, "INR", result.Currency)

	mockCache.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetCountryData_APIFailure(t *testing.T) {
	mockCache := new(mocks.ICacheService)
	mockService := new(mocks.IServiceCall)

	mockCache.On("Get", "india").Return(nil, false)

	mockService.
		On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, 500, errors.New("api failed"))

	service := InitCountryService(mockService, mockCache)

	_, err := service.GetCountryData(context.Background(), "india")

	assert.Error(t, err)

	mockService.AssertExpectations(t)
}

func TestGetCountryData_JSONError(t *testing.T) {
	mockCache := new(mocks.ICacheService)
	mockService := new(mocks.IServiceCall)

	mockCache.On("Get", "india").Return(nil, false)

	mockService.
		On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return([]byte("invalid-json"), 200, nil)

	service := InitCountryService(mockService, mockCache)

	_, err := service.GetCountryData(context.Background(), "india")

	assert.Error(t, err)
}
