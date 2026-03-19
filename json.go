package orderedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

var (
	ErrInvalidJSON        = &jsonError{"invalid JSON format"}
	ErrUnsupportedKeyType = &jsonError{"only string keys are supported"}
)

type jsonError struct {
	msg string
}

func (e *jsonError) Error() string {
	return e.msg
}

func (o *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	if o == nil {
		return []byte("{}"), nil
	}

	var sb strings.Builder
	encoder := json.NewEncoder(&sb)
	encoder.SetEscapeHTML(false)

	sb.WriteString("{")

	first := true
	var marshalErr error
	o.Iterate(func(pair *Pair[K, V]) bool {
		if !first {
			sb.WriteString(",")
		}
		first = false

		keyJSON, e := json.Marshal(pair.Key)
		if e != nil {
			marshalErr = e
			return false
		}

		valJSON, e := json.Marshal(pair.Value)
		if e != nil {
			marshalErr = e
			return false
		}

		sb.Write(keyJSON)
		sb.WriteString(":")
		sb.Write(valJSON)
		return true
	})

	sb.WriteString("}")
	return []byte(sb.String()), marshalErr
}

func (o *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.keys == nil {
		o.keys = make([]K, 0)
	}
	if o.values == nil {
		o.values = make(map[K]V)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil {
		return err
	}

	if token != json.Delim('{') {
		return ErrInvalidJSON
	}

	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		key, ok := token.(string)
		if !ok {
			return ErrInvalidJSON
		}

		var val V
		if err := decoder.Decode(&val); err != nil {
			return err
		}

		var k K
		switch any(k).(type) {
		case string:
			o.keys = append(o.keys, any(key).(K))
			o.values[any(key).(K)] = val
		default:
			return ErrUnsupportedKeyType
		}
	}

	token, err = decoder.Token()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
