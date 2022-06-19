package provider

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getOkStrPtr returns the value of the string field keyName in d.
// If the field is not it returns nil.
func getOkStrPtr(d *schema.ResourceData, keyName string) *string {
	val, isSet := d.GetOk(keyName)
	if !isSet || val == nil {
		return nil
	}

	v := val.(string)
	return &v
}

func getStrPtr(d *schema.ResourceData, keyName string) *string {
	val := d.Get(keyName)
	if val == nil {
		return nil
	}

	v := val.(string)
	return &v
}

func getBoolPtr(d *schema.ResourceData, keyName string) *bool {
	val := d.Get(keyName)
	if val == nil {
		return nil
	}

	v := val.(bool)
	return &v
}

func getInt32Ptr(d *schema.ResourceData, keyName string) *int32 {
	val := d.Get(keyName)
	if val == nil {
		return nil
	}

	v := int32(val.(int))
	return &v
}

func getInt64Ptr(d *schema.ResourceData, keyName string) *int64 {
	val := d.Get(keyName)
	if val == nil {
		return nil
	}

	v := int64(val.(int))
	return &v
}

func getIntPtr(d *schema.ResourceData, keyName string) *int {
	val := d.Get(keyName)
	if val == nil {
		return nil
	}

	v := val.(int)
	return &v
}

func getIDAsInt64(d *schema.ResourceData) (int64, error) {
	strID := d.Id()
	if strID == "" {
		return -1, errors.New("id is empty")
	}

	id, err := strconv.Atoi(strID)
	if err != nil {
		return -1, fmt.Errorf("converting id to integer failed: %w", err)
	}

	return int64(id), nil
}

// getStrSetAsSlice converts a TypeSet with string elements to a []string
func getStrSetAsSlice(d *schema.ResourceData, key string) []string {
	return strSetAsSlice(d.Get(key))
}
