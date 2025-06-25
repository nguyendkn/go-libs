package json

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	result, err := Marshal(data)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	if len(result) == 0 {
		t.Errorf("Marshal() returned empty result")
	}
}

func TestMarshalIndent(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	result, err := MarshalIndent(data, "", "  ")
	if err != nil {
		t.Errorf("MarshalIndent() error = %v", err)
		return
	}

	if len(result) == 0 {
		t.Errorf("MarshalIndent() returned empty result")
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`{"name": "John", "age": 30}`)
	var result map[string]interface{}

	err := Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if result["name"] != "John" {
		t.Errorf("Unmarshal() name = %v, want John", result["name"])
	}
}

func TestUnmarshalToValue(t *testing.T) {
	data := []byte(`{"name": "John", "age": 30}`)

	result, err := UnmarshalToValue(data)
	if err != nil {
		t.Errorf("UnmarshalToValue() error = %v", err)
		return
	}

	if !result.IsObject() {
		t.Errorf("UnmarshalToValue() should return an object")
	}
}

func TestEqual(t *testing.T) {
	v1, _ := Parse(`{"name": "John", "age": 30}`)
	v2, _ := Parse(`{"name": "John", "age": 30}`)
	v3, _ := Parse(`{"name": "Jane", "age": 25}`)

	if !v1.Equal(v2) {
		t.Errorf("Equal() should return true for identical values")
	}

	if v1.Equal(v3) {
		t.Errorf("Equal() should return false for different values")
	}

	// Test with nil
	if v1.Equal(nil) {
		t.Errorf("Equal() should return false when comparing with nil")
	}
}

func TestMerge(t *testing.T) {
	v1, _ := Parse(`{"name": "John", "age": 30}`)
	v2, _ := Parse(`{"city": "New York", "age": 31}`)

	err := v1.Merge(v2)
	if err != nil {
		t.Errorf("Merge() error = %v", err)
		return
	}

	// Check merged values
	city, err := v1.GetPath("city")
	if err != nil {
		t.Errorf("Merge() should have added city field")
		return
	}

	cityStr, _ := city.GetString()
	if cityStr != "New York" {
		t.Errorf("Merge() city = %v, want New York", cityStr)
	}

	// Check overwritten value
	age, _ := v1.GetPath("age")
	ageInt, _ := age.GetInt()
	if ageInt != 31 {
		t.Errorf("Merge() age = %v, want 31", ageInt)
	}
}

func TestKeys(t *testing.T) {
	v, _ := Parse(`{"name": "John", "age": 30, "city": "New York"}`)

	keys, err := v.Keys()
	if err != nil {
		t.Errorf("Keys() error = %v", err)
		return
	}

	if len(keys) != 3 {
		t.Errorf("Keys() length = %v, want 3", len(keys))
	}

	// Check that all expected keys are present
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	expectedKeys := []string{"name", "age", "city"}
	for _, expected := range expectedKeys {
		if !keyMap[expected] {
			t.Errorf("Keys() missing key: %s", expected)
		}
	}
}

func TestValues(t *testing.T) {
	// Test with object
	obj, _ := Parse(`{"name": "John", "age": 30}`)

	values, err := obj.Values()
	if err != nil {
		t.Errorf("Values() error = %v", err)
		return
	}

	if len(values) != 2 {
		t.Errorf("Values() length = %v, want 2", len(values))
	}

	// Test with array
	arr, _ := Parse(`["hello", 42, true]`)

	values, err = arr.Values()
	if err != nil {
		t.Errorf("Values() error = %v", err)
		return
	}

	if len(values) != 3 {
		t.Errorf("Values() length = %v, want 3", len(values))
	}
}

func TestToMap(t *testing.T) {
	v, _ := Parse(`{"name": "John", "age": 30}`)

	result, err := v.ToMap()
	if err != nil {
		t.Errorf("ToMap() error = %v", err)
		return
	}

	if result["name"] != "John" {
		t.Errorf("ToMap() name = %v, want John", result["name"])
	}

	if result["age"] != float64(30) {
		t.Errorf("ToMap() age = %v, want 30", result["age"])
	}
}

func TestToSlice(t *testing.T) {
	v, _ := Parse(`["hello", 42, true]`)

	result, err := v.ToSlice()
	if err != nil {
		t.Errorf("ToSlice() error = %v", err)
		return
	}

	if len(result) != 3 {
		t.Errorf("ToSlice() length = %v, want 3", len(result))
	}

	if result[0] != "hello" {
		t.Errorf("ToSlice() first element = %v, want hello", result[0])
	}
}

func TestUnmarshalTo(t *testing.T) {
	v, _ := Parse(`{"name": "John", "age": 30}`)

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var user User
	err := v.UnmarshalTo(&user)
	if err != nil {
		t.Errorf("UnmarshalTo() error = %v", err)
		return
	}

	if user.Name != "John" || user.Age != 30 {
		t.Errorf("UnmarshalTo() incorrect values: %+v", user)
	}
}

func TestFromStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	user := User{Name: "John", Age: 30}

	result, err := FromStruct(user)
	if err != nil {
		t.Errorf("FromStruct() error = %v", err)
		return
	}

	if !result.IsObject() {
		t.Errorf("FromStruct() should return an object")
	}

	name, _ := result.GetPath("name")
	nameStr, _ := name.GetString()
	if nameStr != "John" {
		t.Errorf("FromStruct() name = %v, want John", nameStr)
	}
}

func TestCompact(t *testing.T) {
	input := []byte(`{
		"name": "John",
		"age": 30
	}`)

	result := Compact(input)
	if len(result) >= len(input) {
		t.Errorf("Compact() should reduce size")
	}
}

func TestCompactString(t *testing.T) {
	input := `{
		"name": "John",
		"age": 30
	}`

	result := CompactString(input)
	if len(result) >= len(input) {
		t.Errorf("CompactString() should reduce size")
	}
}

func TestIndent(t *testing.T) {
	input := []byte(`{"name":"John","age":30}`)

	result := Indent(input, "", "  ")
	if len(result) <= len(input) {
		t.Errorf("Indent() should increase size")
	}
}

func TestIndentString(t *testing.T) {
	input := `{"name":"John","age":30}`

	result := IndentString(input, "", "  ")
	if len(result) <= len(input) {
		t.Errorf("IndentString() should increase size")
	}
}

func TestSize(t *testing.T) {
	v, _ := Parse(`{"name": "John", "age": 30}`)

	size := v.Size()
	if size <= 0 {
		t.Errorf("Size() should return positive value")
	}

	// Test with nil value
	nilValue := &Value{data: nil}
	nilSize := nilValue.Size()
	if nilSize != 4 { // "null"
		t.Errorf("Size() for nil should return 4, got %d", nilSize)
	}
}

func TestMergeEdgeCases(t *testing.T) {
	// Test merging with nil
	v1, _ := Parse(`{"name": "John"}`)
	err := v1.Merge(nil)
	if err != nil {
		t.Errorf("Merge() with nil should not error")
	}

	// Test merging with nil data
	v2 := &Value{data: nil}
	err = v1.Merge(v2)
	if err != nil {
		t.Errorf("Merge() with nil data should not error")
	}

	// Test merging non-objects
	v3, _ := Parse(`"string"`)
	v4, _ := Parse(`{"key": "value"}`)
	err = v3.Merge(v4)
	if err == nil {
		t.Errorf("Merge() should error when target is not object")
	}

	// Test merging with non-object source
	err = v1.Merge(v3)
	if err == nil {
		t.Errorf("Merge() should error when source is not object")
	}

	// Test recursive merge
	v5, _ := Parse(`{"user": {"name": "John", "age": 30}}`)
	v6, _ := Parse(`{"user": {"age": 31, "city": "NYC"}}`)
	err = v5.Merge(v6)
	if err != nil {
		t.Errorf("Merge() recursive error = %v", err)
		return
	}

	// Check merged result
	age, _ := v5.GetPath("user.age")
	ageInt, _ := age.GetInt()
	if ageInt != 31 {
		t.Errorf("Merge() recursive age = %v, want 31", ageInt)
	}

	city, _ := v5.GetPath("user.city")
	cityStr, _ := city.GetString()
	if cityStr != "NYC" {
		t.Errorf("Merge() recursive city = %v, want NYC", cityStr)
	}
}

func TestKeysEdgeCases(t *testing.T) {
	// Test with non-object
	v, _ := Parse(`"string"`)
	_, err := v.Keys()
	if err == nil {
		t.Errorf("Keys() should error with non-object")
	}
}

func TestValuesEdgeCases(t *testing.T) {
	// Test with non-object/array
	v, _ := Parse(`"string"`)
	_, err := v.Values()
	if err == nil {
		t.Errorf("Values() should error with non-object/array")
	}
}

func TestToMapEdgeCases(t *testing.T) {
	// Test with non-object
	v, _ := Parse(`"string"`)
	_, err := v.ToMap()
	if err == nil {
		t.Errorf("ToMap() should error with non-object")
	}
}

func TestToSliceEdgeCases(t *testing.T) {
	// Test with non-array
	v, _ := Parse(`"string"`)
	_, err := v.ToSlice()
	if err == nil {
		t.Errorf("ToSlice() should error with non-array")
	}
}

func TestUnmarshalToEdgeCases(t *testing.T) {
	// Test with marshal error (create invalid data)
	v := &Value{data: make(chan int)} // channels can't be marshaled
	var result map[string]interface{}
	err := v.UnmarshalTo(&result)
	if err == nil {
		t.Errorf("UnmarshalTo() should error with unmarshalable data")
	}
}

func TestFromStructEdgeCases(t *testing.T) {
	// Test with unmarshalable struct
	type BadStruct struct {
		Ch chan int
	}

	bad := BadStruct{Ch: make(chan int)}
	_, err := FromStruct(bad)
	if err == nil {
		t.Errorf("FromStruct() should error with unmarshalable struct")
	}
}

func TestCompactEdgeCases(t *testing.T) {
	// Test with invalid JSON
	invalid := []byte(`{invalid}`)
	result := Compact(invalid)
	if string(result) != string(invalid) {
		t.Errorf("Compact() should return original for invalid JSON")
	}
}

func TestIndentEdgeCases(t *testing.T) {
	// Test with invalid JSON
	invalid := []byte(`{invalid}`)
	result := Indent(invalid, "", "  ")
	if string(result) != string(invalid) {
		t.Errorf("Indent() should return original for invalid JSON")
	}
}

func TestCloneEdgeCases(t *testing.T) {
	// Test cloning value with marshal error
	v := &Value{data: make(chan int)}
	cloned := v.Clone()
	if !cloned.IsNull() {
		t.Errorf("Clone() should return null for unmarshalable data")
	}

	// Test cloning value with unmarshal error (this is harder to trigger)
	// We'll just test that Clone works with normal data
	normal := New("test")
	clonedNormal := normal.Clone()
	if !normal.Equal(clonedNormal) {
		t.Errorf("Clone() should create equal copy")
	}
}
