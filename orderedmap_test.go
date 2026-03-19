package orderedmap

import "testing"

func TestNew(t *testing.T) {
	om := New[string, int]()
	if om == nil {
		t.Fatal("New returned nil")
	}
	if om.Len() != 0 {
		t.Errorf("Len() = %d, want 0", om.Len())
	}
}

func TestSet(t *testing.T) {
	om := New[string, int]()

	om.Set("a", 1)
	if om.Len() != 1 {
		t.Errorf("Len() = %d, want 1", om.Len())
	}

	om.Set("b", 2)
	if om.Len() != 2 {
		t.Errorf("Len() = %d, want 2", om.Len())
	}

	om.Set("a", 10)
	if om.Len() != 2 {
		t.Errorf("Len() after update = %d, want 2", om.Len())
	}
}

func TestGet(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)

	val, ok := om.Get("a")
	if !ok {
		t.Error("Get(a) returned ok=false, want ok=true")
	}
	if val != 1 {
		t.Errorf("Get(a) = %d, want 1", val)
	}

	val, ok = om.Get("z")
	if ok {
		t.Error("Get(z) returned ok=true, want ok=false")
	}
}

func TestDelete(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	om.Delete("b")
	if om.Len() != 2 {
		t.Errorf("Len() after delete = %d, want 2", om.Len())
	}

	keys := om.Keys()
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "c" {
		t.Errorf("Keys() after delete = %v, want [a, c]", keys)
	}

	om.Delete("z")
	if om.Len() != 2 {
		t.Errorf("Len() after deleting non-existent key = %d, want 2", om.Len())
	}
}

func TestLen(t *testing.T) {
	om := New[string, int]()

	if om.Len() != 0 {
		t.Errorf("Len() = %d, want 0", om.Len())
	}

	om.Set("a", 1)
	if om.Len() != 1 {
		t.Errorf("Len() = %d, want 1", om.Len())
	}
}

func TestKeys(t *testing.T) {
	om := New[string, int]()
	om.Set("c", 3)
	om.Set("a", 1)
	om.Set("b", 2)

	keys := om.Keys()
	if len(keys) != 3 {
		t.Errorf("Keys() len = %d, want 3", len(keys))
	}
	if keys[0] != "c" || keys[1] != "a" || keys[2] != "b" {
		t.Errorf("Keys() = %v, want [c, a, b]", keys)
	}

	keysCopy := om.Keys()
	keysCopy[0] = "z"
	originalKeys := om.Keys()
	if originalKeys[0] != "c" {
		t.Errorf("Keys() returned a copy, modifying it should not affect original")
	}
}

func TestIndex(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	if om.Index("a") != 0 {
		t.Errorf("Index(a) = %d, want 0", om.Index("a"))
	}
	if om.Index("c") != 2 {
		t.Errorf("Index(c) = %d, want 2", om.Index("c"))
	}
	if om.Index("z") != -1 {
		t.Errorf("Index(z) = %d, want -1", om.Index("z"))
	}
}

func TestValues(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	values := om.Values()
	if len(values) != 3 {
		t.Errorf("Values() len = %d, want 3", len(values))
	}
	if values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf("Values() = %v, want [1, 2, 3]", values)
	}
}

func TestFirst(t *testing.T) {
	om := New[string, int]()

	if om.First() != nil {
		t.Error("First() on empty map should return nil")
	}

	om.Set("b", 2)
	om.Set("a", 1)
	om.Set("c", 3)

	first := om.First()
	if first == nil {
		t.Fatal("First() returned nil")
	}
	if first.Key != "b" {
		t.Errorf("First().Key = %s, want b", first.Key)
	}
	if first.Value != 2 {
		t.Errorf("First().Value = %d, want 2", first.Value)
	}
}

func TestLast(t *testing.T) {
	om := New[string, int]()

	if om.Last() != nil {
		t.Error("Last() on empty map should return nil")
	}

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	last := om.Last()
	if last == nil {
		t.Fatal("Last() returned nil")
	}
	if last.Key != "c" {
		t.Errorf("Last().Key = %s, want c", last.Key)
	}
	if last.Value != 3 {
		t.Errorf("Last().Value = %d, want 3", last.Value)
	}
}

func TestIterate(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	var keys []string
	var values []int

	om.Iterate(func(pair *Pair[string, int]) bool {
		keys = append(keys, pair.Key)
		values = append(values, pair.Value)
		return true
	})

	if len(keys) != 3 {
		t.Errorf("Iterate collected %d items, want 3", len(keys))
	}
	if keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Errorf("Iterate keys = %v, want [a, b, c]", keys)
	}
	if values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf("Iterate values = %v, want [1, 2, 3]", values)
	}
}

func TestIterateStop(t *testing.T) {
	om := New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	var count int

	om.Iterate(func(pair *Pair[string, int]) bool {
		count++
		if pair.Key == "b" {
			return false
		}
		return true
	})

	if count != 2 {
		t.Errorf("Iterate stopped early, count = %d, want 2", count)
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]int{"x": 10, "y": 20, "z": 30, "a": 40}
	keys := []string{"z", "x", "y"}

	om := FromMap(m, keys)

	if om.Len() != 3 {
		t.Errorf("Len() = %d, want 3", om.Len())
	}

	expectedKeys := []string{"z", "x", "y"}
	expectedValues := []int{30, 10, 20}

	omKeys := om.Keys()
	for i, k := range expectedKeys {
		if omKeys[i] != k {
			t.Errorf("Keys()[%d] = %s, want %s", i, omKeys[i], k)
		}
	}

	omValues := om.Values()
	for i, v := range expectedValues {
		if omValues[i] != v {
			t.Errorf("Values()[%d] = %d, want %d", i, omValues[i], v)
		}
	}
}

func TestFromMapMissingKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	keys := []string{"a", "c", "b"}

	om := FromMap(m, keys)

	if om.Len() != 2 {
		t.Errorf("Len() = %d, want 2", om.Len())
	}

	val, ok := om.Get("c")
	if ok {
		t.Errorf("Get(c) should return ok=false for missing key")
	}
	_ = val
}

func TestFromMapEmptyKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	keys := []string{}

	om := FromMap(m, keys)

	if om.Len() != 0 {
		t.Errorf("Len() = %d, want 0", om.Len())
	}
}

func TestConcurrentAccess(t *testing.T) {
	om := New[string, int]()

	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			om.Set("key", i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			om.Get("key")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			om.Len()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			om.Keys()
		}
		done <- true
	}()

	<-done
	<-done
	<-done
	<-done
}

func TestPairStruct(t *testing.T) {
	pair := &Pair[string, int]{Key: "test", Value: 42}

	if pair.Key != "test" {
		t.Errorf("Pair.Key = %s, want test", pair.Key)
	}
	if pair.Value != 42 {
		t.Errorf("Pair.Value = %d, want 42", pair.Value)
	}
}

func TestOrderPreserved(t *testing.T) {
	om := New[string, int]()

	order := []string{"first", "second", "third", "fourth", "fifth"}
	for i, k := range order {
		om.Set(k, i)
	}

	keys := om.Keys()
	for i, k := range order {
		if keys[i] != k {
			t.Errorf("Key at index %d = %s, want %s (order not preserved)", i, keys[i], k)
		}
	}

	om.Delete("third")
	order = []string{"first", "second", "fourth", "fifth"}
	keys = om.Keys()
	for i, k := range order {
		if keys[i] != k {
			t.Errorf("After delete, key at index %d = %s, want %s", i, keys[i], k)
		}
	}
}
