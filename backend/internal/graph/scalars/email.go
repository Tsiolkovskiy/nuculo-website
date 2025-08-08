package scalars

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// EmailRegex is a simple email validation regex
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// MarshalEmail marshals an email string
func MarshalEmail(email string) graphql.Marshaler {
	if !EmailRegex.MatchString(email) {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(email))
	})
}

// UnmarshalEmail unmarshals an email string with validation
func UnmarshalEmail(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		if !EmailRegex.MatchString(v) {
			return "", fmt.Errorf("invalid email format: %s", v)
		}
		return v, nil
	case *string:
		if v == nil {
			return "", fmt.Errorf("email cannot be null")
		}
		if !EmailRegex.MatchString(*v) {
			return "", fmt.Errorf("invalid email format: %s", *v)
		}
		return *v, nil
	default:
		return "", fmt.Errorf("email must be a string, got %T", v)
	}
}

// ValidateEmail validates an email string
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if len(email) > 254 {
		return fmt.Errorf("email too long (max 254 characters)")
	}
	if !EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}