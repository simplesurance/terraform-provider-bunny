package provider

import (
	bunny "github.com/simplesurance/bunny-go"
)

var edgeRuleActionTypesStr = map[string]int{
	"force_ssl":                  bunny.EdgeRuleActionTypeForceSSL,
	"redirect":                   bunny.EdgeRuleActionTypeRedirect,
	"origin_url":                 bunny.EdgeRuleActionTypeOriginURL,
	"override_cache_time":        bunny.EdgeRuleActionTypeOverrideCacheTime,
	"block_request":              bunny.EdgeRuleActionTypeBlockRequest,
	"set_response_header":        bunny.EdgeRuleActionTypeSetResponseHeader,
	"set_request_header":         bunny.EdgeRuleActionTypeSetRequestHeader,
	"force_download":             bunny.EdgeRuleActionTypeForceDownload,
	"disable_token_auth":         bunny.EdgeRuleActionTypeDisableTokenAuthentication,
	"enable_token_auth":          bunny.EdgeRuleActionTypeEnableTokenAuthentication,
	"override_cache_time_public": bunny.EdgeRuleActionTypeOverrideCacheTimePublic,
	"ignore_quiery_string":       bunny.EdgeRuleActionTypeIgnoreQueryString,
	"disable_optimizer":          bunny.EdgeRuleActionTypeDisableOptimizer,
	"force_compression":          bunny.EdgeRuleActionTypeForceCompression,
}

var edgeRuleActionTypesInt = reverseStrIntMap(edgeRuleActionTypesStr)

var edgeRuleActionTypeKeys = strIntMapKeys(edgeRuleActionTypesStr)