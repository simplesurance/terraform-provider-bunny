package provider

import (
	"fmt"

	ptr "github.com/AlekSi/pointer"
)

// structure represents a nested Terraform block.
// The Terraform sdk has no type for structs/objects. They are commonly
// represented as a TypeList with with a single map[string]interface{} element.
type structure map[string]interface{}

type resourceDataGetter interface {
	Get(string) interface{}
}

// structureFromResource returns a new structure from the field with the passed
// key from ResourceData.
// If the key does not exist in d, the type of the value is not []interface{}
// with elements of type map[string]interface{} or has not a size of 0 or 1,
// the function panics.
func structureFromResource(d resourceDataGetter, key string) structure {
	v := d.Get(key)
	if v == nil {
		panic(fmt.Sprintf("structureFromResource: key %q is nil in ResourceData", key))
	}

	return structureFromElem(v.([]interface{}))
}

func structureFromElem(e []interface{}) structure {
	if len(e) == 0 {
		logger.Debugf("structureFromResource: slice is empty")
		return nil
	}

	if len(e) != 1 {
		panic(fmt.Sprintf("expected list with length 0 or 1, got length: %d", len(e)))
	}

	return e[0].(map[string]interface{})
}

func (m structure) isEmpty() bool {
	return len(m) == 0
}

// getBoolPtr returns the value of the passed key as *bool.
func (m structure) getBoolPtr(key string) *bool {
	return ptr.ToBool(m[key].(bool))
}

// getStr returns the value of the passed key as string.
func (m structure) getStr(key string) string {
	return m[key].(string)
}

func (m structure) getStrPtr(key string) *string {
	return ptr.ToString(m[key].(string))
}

func (m structure) getIntPtr(key string) *int {
	return ptr.ToInt(m[key].(int))
}

func (m structure) getInt32Ptr(key string) *int32 {
	return ptr.ToInt32(int32(m[key].(int)))
}

func (m structure) getInt64Ptr(key string) *int64 {
	return ptr.ToInt64(int64(m[key].(int)))
}

func (m structure) getFloat64Ptr(key string) *float64 {
	return ptr.ToFloat64(m[key].(float64))
}
