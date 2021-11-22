package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	bunny "github.com/simplesurance/bunny-go"
)

const (
	keySafeHopEnable                       = "enable"
	keySafeHopOriginConnectTimeout         = "origin_connect_timeout"
	keySafeHopOriginResponseTimeout        = "origin_response_timeout"
	keySafeHopOriginRetries                = "origin_retries"
	keySafeHopOriginRetry5xxResponses      = "origin_retry_5xx_response"
	keySafeHopOriginRetryConnectionTimeout = "origin_retry_connection_timeout"
	keySafeHopOriginRetryDelay             = "origin_retry_delay"
	keySafeHopOriginRetryResponseTimeout   = "origin_retry_response_timeout"
)

var resourcePullZoneSafeHop = &schema.Resource{
	Schema: map[string]*schema.Schema{
		keySafeHopEnable: {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If enabled, SafeHop will attempt to retry failed requests to the origin in case of errors or connection failures in a round-robin fashion.",
		},
		keySafeHopOriginConnectTimeout: {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "The amount of seconds to wait when connecting to the origin. Otherwise the request will fail or retry.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{3, 5, 10})),
			DiffSuppressFunc: diffSupressIntUnset,
		},
		keySafeHopOriginResponseTimeout: {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "The amount of seconds to wait when waiting for the origin reply. Otherwise the request will fail or retry.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{5, 15, 30, 45, 60})),
			DiffSuppressFunc: diffSupressIntUnset,
		},
		keySafeHopOriginRetries: {
			Type:     schema.TypeInt,
			Optional: true,
			Description: "Configure how many times bunny.net will re-attempt to connect to the origin before failing with a 502 or a 504 response.\n" +
				"If multiple IPs are set on the origin hostname, the CDN will automatically cycle between them on subsequent attempts.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{0, 1, 2})),
		},
		keySafeHopOriginRetry5xxResponses: {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Determines if we should retry the request in case of a 5XX response.",
		},
		keySafeHopOriginRetryConnectionTimeout: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Determines if we should retry the request in case of a connection timeout.",
		},
		keySafeHopOriginRetryDelay: {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "Determines the amount of time that the CDN should wait before retrying an origin request.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{0, 1, 3, 5, 10})),
			DiffSuppressFunc: diffSupressIntUnset,
		},
		keySafeHopOriginRetryResponseTimeout: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Determines if we should retry the request in case of a response timeout.",
		},
	},
}

func safeHopToResource(pz *bunny.PullZone, d *schema.ResourceData) error {
	safeHopSettings := map[string]interface{}{}

	safeHopSettings[keySafeHopEnable] = pz.EnableSafeHop
	safeHopSettings[keySafeHopOriginConnectTimeout] = pz.OriginConnectTimeout
	safeHopSettings[keySafeHopOriginResponseTimeout] = pz.OriginResponseTimeout
	safeHopSettings[keySafeHopOriginRetry5xxResponses] = pz.OriginRetry5xxResponses
	safeHopSettings[keySafeHopOriginRetryConnectionTimeout] = pz.OriginRetryConnectionTimeout
	safeHopSettings[keySafeHopOriginRetryDelay] = pz.OriginRetryDelay
	safeHopSettings[keySafeHopOriginRetryResponseTimeout] = pz.OriginRetryResponseTimeout
	safeHopSettings[keySafeHopOriginRetries] = pz.OriginRetries

	return d.Set(keySafeHop, []map[string]interface{}{safeHopSettings})
}

func safehopPullZoneUpdateOptionsFromResource(res *bunny.PullZoneUpdateOptions, d *schema.ResourceData) {
	m := structureFromResource(d, keySafeHop)
	if len(m) == 0 {
		return
	}

	res.EnableSafeHop = m.getBoolPtr(keySafeHopEnable)
	res.OriginConnectTimeout = m.getInt32Ptr(keySafeHopOriginConnectTimeout)
	res.OriginResponseTimeout = m.getInt32Ptr(keySafeHopOriginResponseTimeout)
	res.OriginRetries = m.getInt32Ptr(keySafeHopOriginRetries)
	res.OriginRetry5xxResponses = m.getBoolPtr(keySafeHopOriginRetry5xxResponses)
	res.OriginRetryConnectionTimeout = m.getBoolPtr(keySafeHopOriginRetryConnectionTimeout)
	res.OriginRetryDelay = m.getInt32Ptr(keySafeHopOriginRetryDelay)
	res.OriginRetryResponseTimeout = m.getBoolPtr(keySafeHopOriginRetryResponseTimeout)
}
