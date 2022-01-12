package aternos_api

import (
	"encoding/json"
	"strings"
)

// Data represents an arbitrary payload that can either contain a string or a variable object.
type Data struct {
	// Content contains either a regular string or a serialized JSON object.
	Content      string
	ContentBytes []byte
}

func (d *Data) MarshalJSON() ([]byte, error) {
	// Content is an object stored as serialized JSON, so just return the bytes.
	if strings.HasPrefix(d.Content, "{") {
		if d.ContentBytes != nil {
			return d.ContentBytes, nil
		}
		return []byte(d.Content), nil
	}

	// Content is a regular string.
	return json.Marshal(d.Content)
}

func (d *Data) UnmarshalJSON(data []byte) error {
	var raw interface{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch raw.(type) {
	case string:
		str := raw.(string)
		*d = Data{
			Content:      str,
			ContentBytes: []byte(str),
		}
	case map[string]interface{}:
		b, err := json.Marshal(raw)
		if err != nil {
			return err
		}
		*d = Data{
			Content:      string(b),
			ContentBytes: b,
		}
	}

	return nil
}
