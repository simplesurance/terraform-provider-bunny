package provider

import (
	"fmt"

	ptr "github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// structure represents a nested Terraform block.
// The Terraform sdk has no type for structs/objects. They are commonly
// represented as a TypeList with with a single map[string]interface{} element.
type structure map[string]interface{}

// structureFromResource returns a new structure from the field with the passed
// key from ResourceData.
// d.Get() must return a value of type []interface{} with 0 or 1
// map[string]interface{} elements. Otherwise the functon will panic.
func structureFromResource(d *schema.ResourceData, key string) structure {
	list := d.Get(key).([]interface{})
	if len(list) == 0 {
		return nil
	}

	if len(list) != 1 {
		panic(fmt.Sprintf("expected list with length 0 or 1, got length: %d", len(list)))
	}

	return list[0].(map[string]interface{})
}

// getBoolPtr returns the value of the passed key as *bool.
func (m structure) getBoolPtr(key string) *bool {
	return ptr.ToBool(m[key].(bool))
}

// getStr returns the value of the passed key as string.
func (m structure) getStr(key string) string {
	return m[key].(string)
}

func (m structure) getInt32Ptr(key string) *int32 {
	return ptr.ToInt32(int32(m[key].(int)))
}

func (m structure) getInt64Ptr(key string) *int64 {
	return ptr.ToInt64(int64(m[key].(int)))
}
