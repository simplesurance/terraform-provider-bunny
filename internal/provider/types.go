package provider

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// strSliceEqual returns true if all elements in the slices are equal.
func strSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// strSliceToLower converts all elements in strs to lowercase.
func strSliceToLower(strs []string) {
	for i := range strs {
		strs[i] = strings.ToLower(strs[i])
	}
}

// strSliceAsInterfaceSlice returns an []interface{} containing all elements of strs.
func strSliceAsInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))

	for i, elem := range strs {
		res[i] = elem
	}

	return res
}

type strSetOpts struct {
	caseInsensitive bool
	ignoreOrder     bool
}

type setStrSetOpt func(*strSetOpts)

// ignoreOrderOpt is an option for setStrSet. When passed, d.Set() is not
// called if the current *schema.Set stored in d contains the same strings then
// in strSlice but the order differs.
func ignoreOrderOpt(opts *strSetOpts) {
	opts.ignoreOrder = true
}

// caseInsensitiveOpt is an option for setStrSet. When passed, d.Set() is not
// called if the elements of the current *schema.Set stored in d contain the
// same strings then in strSlice but the upper/lower case of the strings differ.
func caseInsensitiveOpt(opts *strSetOpts) {
	opts.caseInsensitive = true
}

// setStrSet converts strSlice to a *schema.Set with strings elements and sets the field key in d to the set.
// If the slices are unset or empty, d.Set() will not be called, to prevent
// that terraform shows differences between empty and nil slices.
// opts can be passed to further customize suppressing differences between the
// current set and the new value.
func setStrSet(d *schema.ResourceData, key string, strSlice []string, opts ...setStrSetOpt) error {
	var options strSetOpts

	for _, opt := range opts {
		opt(&options)
	}

	curSet := d.Get(key).(*schema.Set)
	curStrSlice := interfaceSlicetoStrSlice(curSet.List())

	newSet := schema.NewSet(schema.HashString, strSliceAsInterfaceSlice(strSlice))
	newStrSlice := interfaceSlicetoStrSlice(newSet.List())

	if options.caseInsensitive {
		strSliceToLower(curStrSlice)
		strSliceToLower(newStrSlice)
	}

	if options.ignoreOrder {
		sort.Strings(curStrSlice)
		sort.Strings(newStrSlice)
	}

	if strSliceEqual(curStrSlice, newStrSlice) {
		logger.Debugf("setStrSet: %s %+v and +%v are equal, not setting value", key, curSet.GoString(), newSet.GoString())
		return nil
	}

	logger.Debugf("setStrSet: %s %+v and +%v are not equal, setting value", key, curSet.GoString(), newSet.GoString())

	return d.Set(key, newSet)
}

// getStrSetAsSlice converts a TypeSet with string elements to a []string
func getStrSetAsSlice(d *schema.ResourceData, key string) []string {
	return strSetAsSlice(d.Get(key))
}

func strSetAsSlice(val interface{}) []string {
	if val == nil {
		return []string{}
	}

	set := val.(*schema.Set)
	return interfaceSlicetoStrSlice(set.List())
}

func interfaceSlicetoStrSlice(in []interface{}) []string {
	res := make([]string, len(in))
	for i, elem := range in {
		res[i] = elem.(string)
	}

	return res
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
func getFloat64Ptr(d *schema.ResourceData, keyName string) *float64 {
	val := d.Get(keyName)
	if val == nil {
		return nil
	}

	v := val.(float64)
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

func idAsInt64(d *schema.ResourceData) (int64, error) {
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

// strIntMapKeys returns a []string containing the keys of m.
func strIntMapKeys(m map[string]int) []string {
	res := make([]string, 0, len(m))

	for key := range m {
		res = append(res, key)
	}

	return res
}

func reverseStrIntMap(m map[string]int) map[int]string {
	res := make(map[int]string, len(m))

	for k, v := range m {
		res[v] = k
	}

	return res
}

func strIntMapGet(m map[string]int, key string) (int, error) {
	if v, exists := m[key]; exists {
		return v, nil
	}

	return -1, fmt.Errorf("key %q not found", key)
}

func intStrMapGet(m map[int]string, key *int) (string, error) {
	if key == nil {
		return "", errors.New("key is nil")
	}

	if v, exists := m[*key]; exists {
		return v, nil
	}

	return "", fmt.Errorf("key '%d' not found", key)
}
