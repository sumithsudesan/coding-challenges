package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	// for storage
	storage persistence.Storage
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, store persistence.Storage) *Server {
	return &Server{
		listenAddress: listenAddress,
		storage:       store,
	}
}

func (s *Server) GetStorage() persistence.Storage {
	return s.storage
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	// handler
	mux.HandleFunc("/api/v0/health", s.withCORS(s.Health))
	mux.HandleFunc("/api/v0/device", s.withCORS(s.CreateSignatureDevice))
	mux.HandleFunc("/api/v0/sign", s.withCORS(s.SignTransaction))
	mux.HandleFunc("/api/v0/devices", s.withCORS(s.ListDevices))

	return http.ListenAndServe(s.listenAddress, mux)
}

// withCORS wraps a handler to add CORS headers.
// To allow or restrict resources based on the origin of the request,
// ensuring secure interactions between different origins.
func (s *Server) withCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins, methods, and headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to the handler function
		next.ServeHTTP(w, r)
	})
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
