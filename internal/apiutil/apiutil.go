package apiutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

type ErrorResponse struct {
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// Validate performs struct validation and writes error response if needed.
func Validate(w http.ResponseWriter, input any) bool {
	if err := validate.Struct(input); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			RespondValidationErrors(w, ve)
			return false
		}
		RespondJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Validation error",
		})
		return false
	}
	return true
}

// RespondValidationErrors formats and sends validation errors as JSON.
func RespondValidationErrors(w http.ResponseWriter, ve validator.ValidationErrors) {
	fields := make(map[string]string)

	for _, fe := range ve {
		fields[fe.Field()] = validationMessage(fe)
	}

	RespondJSON(w, http.StatusBadRequest, ErrorResponse{
		Message: "Validation error",
		Fields:  fields,
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	}
	return "is invalid"
}

// RespondError sends a generic JSON error response.
func RespondError(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, ErrorResponse{
		Message: msg,
	})
}

// RespondJSON sends a JSON response with custom payload and status code.
func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func ParseJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid JSON")
		return false
	}

	return true
}

// ParseDate parses a date string in YYYY-MM-DD format and writes an error response if invalid.
func ParseDate(w http.ResponseWriter, date string) (time.Time, bool) {

	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid date: %s (expected YYYY-MM-DD)", date))
		return time.Time{}, false
	}

	return t, true
}

func ParseDatePtr(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}

	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %s (expected YYYY-MM-DD)", date)
	}

	return &d, nil
}

func ParseIntPtr(value string) (*int, error) {
	if value == "" {
		return nil, nil
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid integer: %s", value)
	}

	return &i, nil
}

type contextKey string

const UserIDKey contextKey = "userId"

// UserID retrieves the user ID from the request context and writes an error response if not found.
func UserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "missing user ID")
		return "", false
	}

	return userID, true
}
