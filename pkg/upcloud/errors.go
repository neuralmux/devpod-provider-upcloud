package upcloud

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

// ErrorType represents the type of error for better handling
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeAuthentication
	ErrorTypeNotFound
	ErrorTypeQuotaExceeded
	ErrorTypeInvalidParameter
	ErrorTypeNetworkTimeout
	ErrorTypeServerBusy
	ErrorTypePermissionDenied
)

// ProviderError wraps UpCloud API errors with additional context
type ProviderError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *ProviderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// WrapError wraps an error with provider-specific context
func WrapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check if it's an UpCloud Problem error
	if problem, ok := err.(*upcloud.Problem); ok {
		return handleUpCloudProblem(problem, operation)
	}

	// Check for common error patterns
	errStr := err.Error()
	
	if strings.Contains(errStr, "401") || strings.Contains(errStr, "unauthorized") {
		return &ProviderError{
			Type:    ErrorTypeAuthentication,
			Message: "Authentication failed. Please check your UpCloud credentials",
			Err:     err,
		}
	}

	if strings.Contains(errStr, "404") || strings.Contains(errStr, "not found") {
		return &ProviderError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("Resource not found during %s", operation),
			Err:     err,
		}
	}

	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return &ProviderError{
			Type:    ErrorTypeNetworkTimeout,
			Message: fmt.Sprintf("Network timeout during %s. Please check your connection", operation),
			Err:     err,
		}
	}

	if strings.Contains(errStr, "quota") || strings.Contains(errStr, "limit exceeded") {
		return &ProviderError{
			Type:    ErrorTypeQuotaExceeded,
			Message: "UpCloud account quota exceeded. Please check your account limits",
			Err:     err,
		}
	}

	// Default error
	return &ProviderError{
		Type:    ErrorTypeUnknown,
		Message: fmt.Sprintf("Error during %s", operation),
		Err:     err,
	}
}

// handleUpCloudProblem handles specific UpCloud API problem responses
func handleUpCloudProblem(problem *upcloud.Problem, operation string) error {
	// Map UpCloud error codes to user-friendly messages
	switch problem.Status {
	case 401:
		return &ProviderError{
			Type:    ErrorTypeAuthentication,
			Message: "Invalid UpCloud credentials. Please check your username and password",
			Err:     problem,
		}
	case 402:
		return &ProviderError{
			Type:    ErrorTypeQuotaExceeded,
			Message: "Payment required. Please check your UpCloud account billing status",
			Err:     problem,
		}
	case 403:
		return &ProviderError{
			Type:    ErrorTypePermissionDenied,
			Message: fmt.Sprintf("Permission denied for %s. Please check your account permissions", operation),
			Err:     problem,
		}
	case 404:
		return &ProviderError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("Resource not found during %s", operation),
			Err:     problem,
		}
	case 409:
		return &ProviderError{
			Type:    ErrorTypeServerBusy,
			Message: fmt.Sprintf("Resource conflict during %s. The server may be in use or transitioning", operation),
			Err:     problem,
		}
	case 422:
		return &ProviderError{
			Type:    ErrorTypeInvalidParameter,
			Message: fmt.Sprintf("Invalid parameters for %s: %s", operation, problem.Title),
			Err:     problem,
		}
	case 429:
		return &ProviderError{
			Type:    ErrorTypeQuotaExceeded,
			Message: "Rate limit exceeded. Please wait a moment and try again",
			Err:     problem,
		}
	case 503:
		return &ProviderError{
			Type:    ErrorTypeServerBusy,
			Message: "UpCloud service temporarily unavailable. Please try again later",
			Err:     problem,
		}
	default:
		return &ProviderError{
			Type:    ErrorTypeUnknown,
			Message: fmt.Sprintf("UpCloud API error during %s: %s", operation, problem.Title),
			Err:     problem,
		}
	}
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	if perr, ok := err.(*ProviderError); ok {
		return perr.Type == ErrorTypeNotFound
	}
	if problem, ok := err.(*upcloud.Problem); ok {
		return problem.Status == 404
	}
	return strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404")
}

// IsAuthenticationError checks if the error is an authentication error
func IsAuthenticationError(err error) bool {
	if perr, ok := err.(*ProviderError); ok {
		return perr.Type == ErrorTypeAuthentication
	}
	if problem, ok := err.(*upcloud.Problem); ok {
		return problem.Status == 401
	}
	return false
}

// IsQuotaError checks if the error is a quota/limit error
func IsQuotaError(err error) bool {
	if perr, ok := err.(*ProviderError); ok {
		return perr.Type == ErrorTypeQuotaExceeded
	}
	if problem, ok := err.(*upcloud.Problem); ok {
		return problem.Status == 402 || problem.Status == 429
	}
	return false
}