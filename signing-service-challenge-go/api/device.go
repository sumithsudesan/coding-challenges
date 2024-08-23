package api

// TODO: REST endpoints ...

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
)

// Respone for creating an Signature device
type CreateSignatureDeviceResponse struct {
	DeviceID string `json:"device_id"`
	Label    string `json:"label"`
}

// Response for signing a transaction.
type SignatureResponse struct {
	Signature  string `json:"Signature"`
	SignedData string `json:"signed_data"`
}

// Requests
type CreateDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type SigninRequest struct {
	DeviceID string `json:"device_id"`
	Data     string `json:"data"`
}

// CreateSignatureDevice handles the creation of a new Signature device.
func (s *Server) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	// Check the POST method
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	var err error

	// handling create Device Req
	var createDeviceReq CreateDeviceRequest
	if err = json.NewDecoder(r.Body).Decode(&createDeviceReq); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	var keyPair domain.KeyDataPair

	// Generate key based on the algoirthm
	keyPair, err = s.generateKeyPair(createDeviceReq.Algorithm)
	if err != nil {
		WriteInternalError(w)
		return
	}

	// device UUID
	deviceID := uuid.New().String()

	// New device
	device := domain.SignatureDevice{
		ID:               deviceID,
		Label:            createDeviceReq.Label,
		SignatureCounter: 0,
		KeyPair:          keyPair,
	}

	// Save the device
	err = s.storage.SaveDevice(device)
	if err != nil {
		WriteInternalError(w)
		return
	}

	// Response
	resp := CreateSignatureDeviceResponse{
		DeviceID: deviceID,
		Label:    createDeviceReq.Label,
	}
	WriteAPIResponse(w, http.StatusOK, resp)
}

// SignTransaction handles the signing of transaction data.
func (s *Server) SignTransaction(w http.ResponseWriter, r *http.Request) {
	// Check the POST method
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	var signinReq SigninRequest
	var err error

	/// handling Signin Request
	if err = json.NewDecoder(r.Body).Decode(&signinReq); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	if signinReq.DeviceID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Device ID cannot be empty"})
		return
	}

	// Get the devide details
	device, err := s.storage.GetDevice(signinReq.DeviceID)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, []string{"Device not found"})
		return
	}

	counter := device.SignatureCounter
	lastAuth := base64.StdEncoding.EncodeToString([]byte(device.ID))
	securedData := fmt.Sprintf("%d_%s_%s", counter, signinReq.Data, lastAuth)
	signature, err := device.KeyPair.Sign(securedData)
	if err != nil {
		WriteInternalError(w)
		return
	}

	// increment the counter
	device.SignatureCounter++

	// Save the device details
	err = s.storage.SaveDevice(device)
	if err != nil {
		WriteInternalError(w)
		return
	}

	// Return the response
	resp := SignatureResponse{
		Signature:  base64.StdEncoding.EncodeToString(signature),
		SignedData: securedData,
	}
	WriteAPIResponse(w, http.StatusOK, resp)
}

// ListDevices handles listing all devices.
func (s *Server) ListDevices(w http.ResponseWriter, r *http.Request) {
	// Check the get method
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	// Get the list of devices
	devices, err := s.storage.ListDevices()
	if err != nil {
		WriteInternalError(w)
		return
	}

	WriteAPIResponse(w, http.StatusOK, devices)
}

// Generate KeyPair
func (s *Server) generateKeyPair(algorithm string) (domain.KeyDataPair, error) {
	switch algorithm {
	case "ECC":
		return domain.CreateECCKeyPair()
	case "RSA":
		return domain.CreateRSAKeyPair()
	default:
		return nil, errors.New("unsupported algorithm")
	}
}
