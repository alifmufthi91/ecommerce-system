package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReserveStocks_ConcurrentHTTP(t *testing.T) {
	// Setup test environment
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	warehouseID := "a0ebb46d-6482-405c-a340-c4a144591fce"
	// Seed test data
	productID, warehouseID := seedTestDataViaAPI(t, env, warehouseID)

	t.Run("concurrent HTTP requests should not exceed available stock", func(t *testing.T) {
		concurrentRequests := 10
		requestQuantity := 15 // Each request wants 15 items
		// Total demand: 10 * 15 = 150, but only 100 available

		var wg sync.WaitGroup
		results := make([][]payload.ReserveStocksResp, concurrentRequests)
		var mu sync.Mutex
		successCount := 0

		// Launch concurrent HTTP requests
		for i := 0; i < concurrentRequests; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				reqBody := payload.ReserveStocksReq{
					Stocks: []payload.ReserveStocksData{
						{
							ProductID: productID,
							Quantity:  requestQuantity,
						},
					},
				}

				resp := makeReserveStocksRequest(t, env.APIBaseURL, reqBody)

				mu.Lock()
				results[index] = resp
				fmt.Println(results[index])
				if len(resp) > 0 {
					successCount++
				}
				mu.Unlock()
			}(i)
		}

		wg.Wait()

		// Verify results
		mu.Lock()
		defer mu.Unlock()

		// Should have some successful reservations but not all
		assert.Equal(t, successCount, 6, "6 * 15 = 90 should be reserved, but not all requests can succeed")
		assert.Less(t, successCount, concurrentRequests, "Not all reservations should succeed due to stock limit")

		// Verify final stock state via API
		finalStock := getStockViaAPI(t, env.APIBaseURL, productID, warehouseID)
		fmt.Println(finalStock)
		assert.Equal(t, finalStock.Quantity, 100, "Initial stock quantity should remain unchanged")
		assert.Equal(t, finalStock.Reserved, 90, "Total reserved stock should be 100 or less")

		t.Logf("Success count: %d/%d, Final reserved: %d, Final available: %d",
			successCount, concurrentRequests, finalStock.Reserved, finalStock.Quantity-finalStock.Reserved)
	})
}

// Helper functions for HTTP requests

func makeAuthenticatedRequest(method, url string, body []byte, token string) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	// Add authentication header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	return client.Do(req)
}

func makeReserveStocksRequest(t *testing.T, baseURL string, req payload.ReserveStocksReq) []payload.ReserveStocksResp {
	jsonData, err := json.Marshal(req)
	require.NoError(t, err)

	testToken := "test-static-token"

	resp, err := makeAuthenticatedRequest(
		"POST",
		fmt.Sprintf("%s/api/stocks/reserve", baseURL),
		jsonData,
		testToken,
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Data []payload.ReserveStocksResp `json:"data"`
		Meta httpresp.Meta               `json:"meta"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)

	require.NoError(t, err)

	return result.Data
}

func getStockViaAPI(t *testing.T, baseURL, productID, warehouseID string) model.WarehouseStock {
	url := fmt.Sprintf("%s/api/stocks?product_id_in=%s&warehouse_id_in=%s", baseURL, productID, warehouseID)

	testToken := "test-static-token"

	resp, err := makeAuthenticatedRequest("GET", url, nil, testToken)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result struct {
		Data []model.WarehouseStock `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	require.Len(t, result.Data, 1)

	return result.Data[0]
}

func seedTestDataViaAPI(t *testing.T, env *TestEnvironment, whID string) (productID, warehouseID string) {
	productID = uuid.New().String()
	warehouseID = whID

	// Create initial stock via API call
	createStockReq := map[string]interface{}{
		"warehouse_id": warehouseID,
		"product_id":   productID,
		"quantity":     100,
	}

	jsonData, err := json.Marshal(createStockReq)
	require.NoError(t, err)

	testToken := "test-static-token"

	resp, err := makeAuthenticatedRequest(
		"POST",
		fmt.Sprintf("%s/api/stocks", env.APIBaseURL),
		jsonData,
		testToken,
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	return productID, warehouseID
}

func resetStockViaAPI(t *testing.T, baseURL, productID, warehouseID string, quantity int) {
	updateReq := map[string]interface{}{
		"quantity": quantity,
		"reserved": 0,
	}

	jsonData, err := json.Marshal(updateReq)
	require.NoError(t, err)

	testToken := "test-static-token"

	url := fmt.Sprintf("%s/api/stocks/%s/%s", baseURL, warehouseID, productID)
	resp, err := makeAuthenticatedRequest("PUT", url, jsonData, testToken)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
