package provider

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
// opts can be passed to further customize suppressing differences between the
// current set and the new value.
func setStrSet(d *schema.ResourceData, key string, strSlice []string, opts ...setStrSetOpt) error {
	var options strSetOpts

	for _, opt := range opts {
		opt(&options)
	}

	curISet, isSet := d.GetOk(key)
	curSet := curISet.(*schema.Set)
	newSet := schema.NewSet(schema.HashString, strSliceAsInterfaceSlice(strSlice))

	if isSet {
		curStrSlice := interfaceSlicetoStrSlice(curSet.List())
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
	}

	logger.Debugf("setStrSet: %s %+v and +%v are not equal, setting value", key, curSet.GoString(), newSet.GoString())

	return d.Set(key, newSet)
}

func strSetAsSlice(val interface{}) []string {
	if val == nil {
		return []string{}
	}

	set := val.(*schema.Set)
	return interfaceSlicetoStrSlice(set.List())
}
