package provider

import bunny "github.com/simplesurance/bunny-go"

var edgeRuleTriggerTypesStr = map[string]int{
	"url":             bunny.EdgeRuleTriggerTypeURL,
	"request_header":  bunny.EdgeRuleTriggerTypeRequestHeader,
	"response_header": bunny.EdgeRuleTriggerTypeResponseHeader,
	"url_extensions":  bunny.EdgeRuleTriggerTypeURLExtension,
	"country_code":    bunny.EdgeRuleTriggerTypeCountryCode,
	"remote_ip":       bunny.EdgeRuleTriggerTypeRemoteIP,
	"query_string":    bunny.EdgeRuleTriggerTypeURLQueryString,
	"random_chance":   bunny.EdgeRuleTriggerTypeRandomChance,
}

var edgeRuleTriggerTypesInt = reverseStrIntMap(edgeRuleTriggerTypesStr)

var edgeRuleTriggerTypeKeys = strIntMapKeysSorted(edgeRuleTriggerTypesStr)
