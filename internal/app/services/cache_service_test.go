package services

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"search-country-system/internal/dto"
	"search-country-system/internal/mocks"
)

// =========================
// Mock-based cache tests
// =========================

func TestCacheService_SetAndGet_Mock(t *testing.T) {
	mockCache := mocks.NewICacheService(t)

	// Setup expected behavior
	mockCache.On("Set", "india", "data").Return()
	mockCache.On("Get", "india").Return("data", true)

	// Call Set
	mockCache.Set("india", "data")

	// Call Get
	val, ok := mockCache.Get("india")

	assert.True(t, ok, "expected key to be found")
	assert.Equal(t, "data", val)

	// Assert all expectations were met
	mockCache.AssertExpectations(t)
}

func TestCacheService_KeyNotFound_Mock(t *testing.T) {
	mockCache := mocks.NewICacheService(t)

	mockCache.On("Get", "unknown").Return(nil, false)

	val, ok := mockCache.Get("unknown")

	assert.False(t, ok, "expected key not to be found")
	assert.Nil(t, val)
}

// =========================
// Example: Cache used in CountryService
// =========================
func TestCountryService_UsesCache(t *testing.T) {
	mockCache := mocks.NewICacheService(t)
	mockServiceCall := mocks.NewIServiceCall(t)

	countryService := InitCountryService(mockServiceCall, mockCache)

	countryName := "India"
	expectedCountry := dto.Country{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "INR",
		Population: 1390000000,
	}

	// Setup cache hit
	mockCache.On("Get", countryName).Return(expectedCountry, true)

	// Call the service
	country, err := countryService.GetCountryData(context.Background(), countryName)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCountry, country)

	// Assert all expectations
	mockCache.AssertExpectations(t)
}

// =========================
// Real cache behavior tests
// =========================

func TestCacheService_Overwrite(t *testing.T) {
	cache := InitCacheService()

	cache.Set("india", "old")
	cache.Set("india", "new")

	val, _ := cache.Get("india")

	assert.Equal(t, "new", val, "expected value to be overwritten")
}

func TestCacheService_ConcurrentAccess(t *testing.T) {
	cache := InitCacheService()
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			cache.Set(key, i)
			val, ok := cache.Get(key)
			assert.True(t, ok)
			assert.Equal(t, i, val)
		}(i)
	}

	wg.Wait()
}

func TestCacheService_NilValue(t *testing.T) {
	cache := InitCacheService()

	cache.Set("nil", nil)

	val, ok := cache.Get("nil")

	assert.True(t, ok, "expected key to be found")
	assert.Nil(t, val)
}
