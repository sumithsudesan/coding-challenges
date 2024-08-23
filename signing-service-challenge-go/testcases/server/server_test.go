package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

// Test storage

// Setup server
func setupServer() *api.Server {
	store := persistence.NewInMemoryStorage()
	return api.NewServer(":8080", store)
}

// Test create device
func TestCreateDevice(t *testing.T) {
	server := setupServer()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid request",
			requestBody:    `{"id": "device123", "algorithm": "ECC", "label": "My Device"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"device_id":"device123","label":"My Device"}`,
		},
		{
			name:           "Missing ID",
			requestBody:    `{"algorithm": "ECC", "label": "My Device"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":["Device ID cannot be empty"]}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":["Invalid request"]}`,
		},
		{
			name:           "Unsupported Algorithm",
			requestBody:    `{"id": "device123", "algorithm": "UnsupportedAlgo", "label": "My Device"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":["Unsupported algorithm"]}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/v0/device", bytes.NewBufferString(tt.requestBody.(string)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.CreateSignatureDevice(w, req)

		resp := w.Result()
		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("[ERROR] %s: expected  %d, recived %d", tt.name, tt.expectedStatus, resp.StatusCode)
		}

		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			t.Errorf("[ERROR] %s: error decoding response body: %v", tt.name, err)
			continue
		}

		bodyBytes, _ := json.Marshal(responseBody)
		if string(bodyBytes) != tt.expectedBody {
			t.Errorf("[ERROR] %s: expected %s, recived %s", tt.name, tt.expectedBody, string(bodyBytes))
		}
	}
}

// Test create sign
func TestSignTransaction(t *testing.T) {
	server := setupServer()

	// Create a device first
	device := domain.SignatureDevice{
		ID:               "device123",
		SignatureCounter: 1,
		Label:            "My Device",
	}
	server.GetStorage().SaveDevice(device)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid request",
			requestBody:    `{"device_id": "device123", "data": "transaction data"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Signature":"<generated_signature>","signed_data":"1_transaction data_<base64_encoded_device_id>"}`, // Adjust based on actual response
		},
		{
			name:           "Device not found",
			requestBody:    `{"device_id": "nonexistent", "data": "transaction data"}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"errors":["Device not found"]}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":["Invalid request"]}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/v0/sign", bytes.NewBufferString(tt.requestBody.(string)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.SignTransaction(w, req)

		resp := w.Result()
		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("[ERROR] %s: expected  %d, recived %d", tt.name, tt.expectedStatus, resp.StatusCode)
		}

		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			t.Errorf("[ERROR] %s: error decoding response body: %v", tt.name, err)
			continue
		}

		bodyBytes, _ := json.Marshal(responseBody)
		if string(bodyBytes) != tt.expectedBody {
			t.Errorf("[ERROR] %s: expected body %s, recived %s", tt.name, tt.expectedBody, string(bodyBytes))
		}
	}
}

// Test create list devices
func TestListDevices(t *testing.T) {
	server := setupServer()

	// Create a couple of devices for testing
	server.GetStorage().SaveDevice(domain.SignatureDevice{
		ID:               "device123",
		SignatureCounter: 1,
		Label:            "Device 123",
	})
	server.GetStorage().SaveDevice(domain.SignatureDevice{
		ID:               "device456",
		SignatureCounter: 2,
		Label:            "Device 456",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v0/devices", nil)
	w := httptest.NewRecorder()
	server.ListDevices(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("[ERROR] ListDevices: expected %d, recived %d", http.StatusOK, resp.StatusCode)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Errorf("[ERROR] ListDevices: error decoding response body: %v", err)
		return
	}

	devices, ok := responseBody["devices"].([]interface{})
	if !ok {
		t.Errorf("[ERROR] ListDevices: expected 'devices' to be a list")
		return
	}
	if len(devices) != 2 {
		t.Errorf("[ERROR]  ListDevices: expected 2 devices, recived %d", len(devices))
	}
}
