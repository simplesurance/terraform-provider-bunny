package provider

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func randResourceName() string {
	return resource.PrefixedUniqueId(resourcePrefix)
}

func idFromState(s *terraform.State, resourceName string) (string, error) {
	resourceState := s.Modules[0].Resources[resourceName]
	if resourceState == nil {
		return "", fmt.Errorf("resource %s not found in state", resourceName)
	}

	insState := resourceState.Primary
	if insState == nil {
		return "", fmt.Errorf("resource %s has no primary state", resourceName)
	}

	if insState.ID == "" {
		return "", errors.New("ID is empty")
	}

	return insState.ID, nil
}

func diffStructs(t *testing.T, a, b interface{}, ignoredFields map[string]struct{}) []string {
	t.Helper()
	var res []string

	valStructA := reflect.Indirect(reflect.ValueOf(a))
	valStructB := reflect.Indirect(reflect.ValueOf(b))

	for i := 0; i < valStructA.NumField(); i++ {
		typeFieldA := valStructA.Type().Field(i)
		if ignoredFields != nil {
			if _, exists := ignoredFields[typeFieldA.Name]; exists {
				continue
			}
		}

		fieldA := valStructA.Field(i)
		fieldB := valStructB.Field(i)

		switch typeFieldA.Type.Kind() {
		case reflect.Slice:
			aLen := fieldA.Len()
			bLen := fieldB.Len()

			if aLen != bLen {
				res = append(res,
					fmt.Sprintf("slice %s differs, has different length, %d and %d",
						typeFieldA.Name, aLen, bLen),
				)
				break
			}

			for j := 0; j < aLen; j++ {
				diffs := diffStructs(t, fieldA.Index(j), fieldB.Index(j), nil)
				for _, d := range diffs {
					res = append(res, fmt.Sprintf("%s[%d]: %s", typeFieldA.Name, j, d))
				}
			}

		default:
			// skip comparing unexported values
			if !fieldA.CanInterface() {
				continue
			}

			diff, err := diffSimpleVal(fieldA, fieldB)
			if err != nil {
				t.Errorf("%s: %s", typeFieldA.Name, err)
				continue
			}

			if diff != "" {
				res = append(res, fmt.Sprintf("%s: %s", typeFieldA.Name, diff))
			}
		}

	}
	return res
}

func diffSimpleVal(a, b reflect.Value) (string, error) {
	valFieldA := a.Interface()
	valFieldB := b.Interface()

	switch valA := valFieldA.(type) {
	case *string:
		valB := valFieldB.(*string)

		return strDiff(valA, valB), nil

	case *bool:
		valB := valFieldB.(*bool)

		return boolDiff(valA, valB), nil

	case *int:
		valB := valFieldB.(*int)

		return intDiff(valA, valB), nil

	case *int64:
		valB := valFieldB.(*int64)

		return int64Diff(valA, valB), nil

	case *int32:
		valB := valFieldB.(*int32)

		return int32Diff(valA, valB), nil

	case *float64:
		valB := valFieldB.(*float64)

		return float64Diff(valA, valB), nil
	case []string:
		valB := valFieldB.([]string)

		return strSliceDiff(valA, valB), nil

	default:
		return "", fmt.Errorf("can not compare field, unsupported diff type: %T", valFieldA)
	}
}
