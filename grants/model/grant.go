package model

import "errors"

type Grant struct {
	key   string
	value string
}

// Key returns the key of the grant.
func (g *Grant) Key() string {
	return g.key
}

// Value returns the value of the grant.
func (g *Grant) Value() string {
	return g.value
}

// Marshal returns the grant as a map.
func (g *Grant) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"key":   g.key,
		"value": g.value,
	}
}

// Unmarshal sets the grant from a map.
func (g *Grant) Unmarshal(body map[string]interface{}) error {
	if k, ok := body["key"].(string); !ok {
		return errors.New("key is not a string")
	} else if v, ok := body["value"].(string); !ok {
		return errors.New("value is not a string")
	} else {
		g.value = v
		g.key = k
	}

	return nil
}
