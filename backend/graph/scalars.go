package graph

import (
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// Time is a custom scalar type for GraphQL Time
type Time time.Time

// MarshalGQL implements the graphql.Marshaler interface
func (t Time) MarshalGQL(w io.Writer) {
	w.Write([]byte(`"` + time.Time(t).Format(time.RFC3339) + `"`))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *Time) UnmarshalGQL(v interface{}) error {
	if str, ok := v.(string); ok {
		parsed, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
		*t = Time(parsed)
		return nil
	}
	return nil
}

// MarshalTime serializes a Time to graphql.Marshaler
func MarshalTime(t Time) graphql.Marshaler {
	return t
}

// UnmarshalTime deserializes a string to a Time
func UnmarshalTime(v interface{}) (Time, error) {
	if str, ok := v.(string); ok {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return Time{}, err
		}
		return Time(t), nil
	}
	return Time{}, nil
}
