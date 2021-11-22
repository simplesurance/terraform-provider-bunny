package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyLimitsConnectionLimitPerIPCount = "connection_limit_per_ip_count"
	keyLimitsMonthlyBandwidthLimit     = "monthly_bandwidth_limit"
	keyLimitsRequestLimit              = "request_limit"
)

var resourcePullZoneLimits = &schema.Resource{
	Schema: map[string]*schema.Schema{
		keyLimitsConnectionLimitPerIPCount: {
			Type:             schema.TypeInt,
			Description:      "Limit the maximum number of allowed connections to the zone per IP.Set to 0 for unlimited.",
			Optional:         true,
			ValidateDiagFunc: validateIsInt32,
		},
		keyLimitsRequestLimit: {
			Type:             schema.TypeInt,
			Description:      "Limit the maximum number of requests per second coming from a single IP. Set to 0 for unlimited.",
			Optional:         true,
			ValidateDiagFunc: validateIsInt32,
		},
		keyLimitsMonthlyBandwidthLimit: {
			Type:        schema.TypeInt,
			Description: "Limits the allowed bandwidth used in a month, in Bytes. If the limit is reached the zone will be disabled.",
			Optional:    true,
		},
	},
}

func limitsToResource(pz *bunny.PullZone, d *schema.ResourceData) error {
	m := map[string]interface{}{}

	m[keyLimitsRequestLimit] = pz.RequestLimit
	m[keyLimitsMonthlyBandwidthLimit] = pz.MonthlyBandwidthLimit
	m[keyLimitsConnectionLimitPerIPCount] = pz.ConnectionLimitPerIPCount

	logger.Infof("limitsToResource: setting to :%+v", m)
	return d.Set(keyLimits, []map[string]interface{}{m})
}

func limitsFromResource(res *bunny.PullZoneUpdateOptions, d *schema.ResourceData) {
	m := structureFromResource(d, keyLimits)
	if len(m) == 0 {
		return
	}

	res.RequestLimit = m.getInt32Ptr(keyLimitsRequestLimit)
	res.MonthlyBandwidthLimit = m.getInt64Ptr(keyLimitsMonthlyBandwidthLimit)
	res.ConnectionLimitPerIPCount = m.getInt32Ptr(keyLimitsConnectionLimitPerIPCount)
}
