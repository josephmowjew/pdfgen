package errors

import "fmt"

// ErrorCode represents different types of PDF generation errors
type ErrorCode int

const (
	ErrInvalidTemplate ErrorCode = iota + 1
	ErrInvalidData
	ErrRenderFailed
	ErrLayoutFailed
	ErrGenerationFailed
)

// PDFError represents a custom error type for PDF operations
type PDFError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *PDFError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewPDFError creates a new PDFError
func NewPDFError(code ErrorCode, message string, cause error) *PDFError {
	return &PDFError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Is implements error interface for error comparison
func (e *PDFError) Is(target error) bool {
	t, ok := target.(*PDFError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Unwrap returns the underlying error
func (e *PDFError) Unwrap() error {
	return e.Cause
}
