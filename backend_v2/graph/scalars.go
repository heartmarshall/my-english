package graph

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// Time представляет время в GraphQL.
type Time time.Time

// MarshalTime сериализует Time в строку.
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.MarshalTime(t)
}

// UnmarshalTime десериализует строку в Time.
func UnmarshalTime(v interface{}) (time.Time, error) {
	return graphql.UnmarshalTime(v)
}

// JSON представляет JSON объект в GraphQL.
type JSON map[string]interface{}

// MarshalJSON сериализует JSON.
func MarshalJSON(j map[string]interface{}) graphql.Marshaler {
	return graphql.MarshalAny(j)
}

// UnmarshalJSON десериализует JSON.
func UnmarshalJSON(v interface{}) (map[string]interface{}, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		return val, nil
	case []byte:
		var result map[string]interface{}
		if err := json.Unmarshal(val, &result); err != nil {
			return nil, fmt.Errorf("cannot unmarshal JSON: %w", err)
		}
		return result, nil
	case string:
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			return nil, fmt.Errorf("cannot unmarshal JSON: %w", err)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot unmarshal %T into JSON", v)
	}
}
