package scalars

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalDateTime marshals a time.Time to RFC3339 format
func MarshalDateTime(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(t.Format(time.RFC3339)))
	})
}

// UnmarshalDateTime unmarshals a datetime string with validation
func UnmarshalDateTime(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case string:
		return parseDateTime(v)
	case *string:
		if v == nil {
			return time.Time{}, fmt.Errorf("datetime cannot be null")
		}
		return parseDateTime(*v)
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case float64:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("datetime must be a string, int, or float64, got %T", v)
	}
}

// parseDateTime parses various datetime formats
func parseDateTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("datetime cannot be empty")
	}

	// Try different formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid datetime format: %s (expected RFC3339 or similar)", s)
}

// ValidateDateTime validates a datetime string
func ValidateDateTime(datetime string) error {
	if datetime == "" {
		return fmt.Errorf("datetime cannot be empty")
	}
	_, err := parseDateTime(datetime)
	return err
}

// ValidateDateTimeRange validates that a datetime is within a reasonable range
func ValidateDateTimeRange(t time.Time) error {
	now := time.Now()
	
	// Check if date is too far in the past (before 1900)
	if t.Year() < 1900 {
		return fmt.Errorf("datetime too far in the past (before 1900)")
	}
	
	// Check if date is too far in the future (more than 100 years from now)
	if t.After(now.AddDate(100, 0, 0)) {
		return fmt.Errorf("datetime too far in the future (more than 100 years)")
	}
	
	return nil
}