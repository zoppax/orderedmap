package orderedmap

import (
	"encoding/json"
	"testing"
)

func Test_JSON_Marshal(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	expected := `{"a":1,"b":2}`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func Test_JSON_MarshalJSONEmpty(t *testing.T) {
	om := New[string, int]()

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON empty failed: %v", err)
	}
	expected := `{}`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func Test_JSON_MarshalJSONNil(t *testing.T) {
	var om *OrderedMap[string, int]

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON nil failed: %v", err)
	}
	expected := `{}`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func Test_JSON_MarshalJSONSingle(t *testing.T) {
	om := New[string, int]()
	om.Set("only", 1)

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON single failed: %v", err)
	}
	expected := `{"only":1}`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func Test_JSON_MarshalJSONComplex(t *testing.T) {
	om := New[string, interface{}]()
	om.Set("name", "John")
	om.Set("active", true)

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON complex failed: %v", err)
	}

	om2 := New[string, interface{}]()
	if err := json.Unmarshal(data, om2); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if name, _ := om2.Get("name"); name != "John" {
		t.Errorf("name = %v, want John", name)
	}
}

type testMarshalErrorType struct{}

func (m testMarshalErrorType) MarshalJSON() ([]byte, error) {
	return nil, ErrInvalidJSON
}

func Test_JSON_MarshalJSONKeyMarshalError(t *testing.T) {
	om := New[testMarshalErrorType, int]()
	om.Set(testMarshalErrorType{}, 1)

	_, err := om.MarshalJSON()
	if err == nil {
		t.Error("MarshalJSON with key marshal error should fail")
	}
}

func Test_JSON_MarshalJSONValueMarshalError(t *testing.T) {
	om := New[string, testMarshalErrorType]()
	om.Set("key", testMarshalErrorType{})

	_, err := om.MarshalJSON()
	if err == nil {
		t.Error("MarshalJSON with value marshal error should fail")
	}
}

func Test_JSON_MarshalJSONOrder(t *testing.T) {
	om := New[string, int]()
	om.Set("z", 26)
	om.Set("a", 1)
	om.Set("m", 13)

	data, _ := om.MarshalJSON()
	expected := `{"z":26,"a":1,"m":13}`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func Test_JSON_Unmarshal(t *testing.T) {
	var om OrderedMap[string, int]
	err := om.UnmarshalJSON([]byte(`{"a":1,"b":2}`))
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if om.Len() != 2 {
		t.Errorf("Len() = %d, want 2", om.Len())
	}
	keys := om.Keys()
	if keys[0] != "a" || keys[1] != "b" {
		t.Errorf("Keys() = %v, want [a, b]", keys)
	}
}

func Test_JSON_UnmarshalJSONEmpty(t *testing.T) {
	var om OrderedMap[string, int]
	err := om.UnmarshalJSON([]byte(`{}`))
	if err != nil {
		t.Fatalf("UnmarshalJSON empty failed: %v", err)
	}
	if om.Len() != 0 {
		t.Errorf("Len() = %d, want 0", om.Len())
	}
}

func Test_JSON_UnmarshalJSONInvalid(t *testing.T) {
	var om OrderedMap[string, int]
	err := om.UnmarshalJSON([]byte(`invalid`))
	if err == nil {
		t.Error("UnmarshalJSON invalid should fail")
	}
}

func Test_JSON_UnmarshalJSONNotObject(t *testing.T) {
	var om OrderedMap[string, int]
	err := om.UnmarshalJSON([]byte(`[1,2,3]`))
	if err == nil {
		t.Error("UnmarshalJSON array should fail")
	}
}

func Test_JSON_UnmarshalJSONUnsupportedKeyType(t *testing.T) {
	var om OrderedMap[bool, int]
	err := om.UnmarshalJSON([]byte(`{"key":1}`))
	if err == nil {
		t.Error("UnmarshalJSON with bool key should fail")
	}
}

func Test_JSON_UnmarshalJSONTokenError(t *testing.T) {
	var om OrderedMap[string, int]
	err := om.UnmarshalJSON([]byte(`{a`))
	if err == nil {
		t.Error("UnmarshalJSON with token error should fail")
	}
}

func Test_JSON_UnmarshalJSONNonStringKey(t *testing.T) {
	var om OrderedMap[int, int]
	err := om.UnmarshalJSON([]byte(`{1:2}`))
	if err == nil {
		t.Error("UnmarshalJSON with non-string key should fail")
	}
}

func Test_JSON_UnmarshalJSONDecodeError(t *testing.T) {
	type customType struct {
		X int
	}
	var om OrderedMap[string, customType]
	err := om.UnmarshalJSON([]byte(`{"a":{"invalid"}}`))
	if err == nil {
		t.Error("UnmarshalJSON with decode error should fail")
	}
}

func Test_JSON_UnmarshalJSONNonNilKeysAndValues(t *testing.T) {
	om := New[string, int]()
	om.Set("existing", 999)

	var om2 OrderedMap[string, int]
	om2.keys = []string{"old"}
	om2.values = map[string]int{"old": 100}

	err := om2.UnmarshalJSON([]byte(`{"a":1}`))
	if err != nil {
		t.Fatalf("UnmarshalJSON with non-nil fields failed: %v", err)
	}

	if om2.Len() != 2 {
		t.Errorf("Len() = %d, want 2", om2.Len())
	}
}

func Test_JSON_RoundTrip(t *testing.T) {
	om := New[string, int]()
	om.Set("z", 26)
	om.Set("a", 1)
	om.Set("m", 13)

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var om2 OrderedMap[string, int]
	if err := om2.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	keys := om2.Keys()
	if len(keys) != 3 || keys[0] != "z" || keys[1] != "a" || keys[2] != "m" {
		t.Errorf("Keys() = %v, want [z, a, m]", keys)
	}

	data2, _ := om2.MarshalJSON()
	if string(data2) != string(data) {
		t.Errorf("Round trip mismatch: %s vs %s", string(data2), string(data))
	}
}

func Test_JSON_StandardLibrary(t *testing.T) {
	om := New[string, int]()
	om.Set("b", 2)
	om.Set("a", 1)

	data, err := json.Marshal(om)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	expected := `{"b":2,"a":1}`
	if string(data) != expected {
		t.Errorf("json.Marshal = %s, want %s", string(data), expected)
	}

	var om2 OrderedMap[string, int]
	if err := json.Unmarshal(data, &om2); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if om2.Len() != 2 {
		t.Errorf("Len() = %d, want 2", om2.Len())
	}
}

func Test_JSON_OrderPreservation(t *testing.T) {
	var om OrderedMap[string, int]
	om.UnmarshalJSON([]byte(`{"third":3,"first":1,"second":2}`))

	keys := om.Keys()
	if keys[0] != "third" || keys[1] != "first" || keys[2] != "second" {
		t.Errorf("Order not preserved: %v", keys)
	}
}

func Test_JSON_Error(t *testing.T) {
	err := ErrInvalidJSON
	if err.Error() != "invalid JSON format" {
		t.Errorf("ErrInvalidJSON.Error() = %s, want 'invalid JSON format'", err.Error())
	}

	err = ErrUnsupportedKeyType
	if err.Error() != "only string keys are supported" {
		t.Errorf("ErrUnsupportedKeyType.Error() = %s, want 'only string keys are supported'", err.Error())
	}
}
