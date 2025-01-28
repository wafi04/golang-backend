package common

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Status     bool        `json:"status"`
	StatusCode int32       `json:"statusCode"`
	Message    *string     `json:"message,omitempty"`
	Error      *string     `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Status:     true,
		StatusCode: 200,
	}
}

func Success(data interface{}, message string) *Response {
	return &Response{
		Status:     true,
		StatusCode: 200,
		Message:    &message,
		Data:       data,
	}
}

func Error(statusCode int32, errorMessage string) *Response {
	return &Response{
		Status:     false,
		StatusCode: statusCode,
		Error:      &errorMessage,
	}
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithMessage(message string) *Response {
	r.Message = &message
	return r
}

func (r *Response) WithError(err string) *Response {
	r.Error = &err
	r.Status = false
	return r
}

func (r *Response) WithStatusCode(code int32) *Response {
	r.StatusCode = code
	return r
}

type Responses struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Status:  statusCode,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SendErrorResponseWithDetails(w http.ResponseWriter, statusCode int, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Status:  statusCode,
		Message: message,
		Details: details,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func SendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding success response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Middleware to measure response time
func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Start time of the request

		// Call the next handler
		next.ServeHTTP(w, r)

		// Calculate response time
		duration := time.Since(start)
		log.Printf("Response time for %s %s: %v\n", r.Method, r.URL.Path, duration)
	})
}
