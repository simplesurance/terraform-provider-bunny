package provider

import (
	bunny "github.com/simplesurance/bunny-go"
)

var edgeRuleMatchingTypesStr = map[string]int{
	"any":  bunny.MatchingTypeAny,
	"all":  bunny.MatchingTypeAll,
	"none": bunny.MatchingTypeNone,
}

var edgeRuleMatchingTypesInt = reverseStrIntMap(edgeRuleMatchingTypesStr)

var edgeRuleMatchingTypeKeys = strIntMapKeysSorted(edgeRuleMatchingTypesStr)
