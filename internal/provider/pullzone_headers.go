package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyEnableAccessControlOriginHeader     = "enable_access_control_origin_header"
	keyAccessControlOriginHeaderExtensions = "access_control_origin_header_extensions"
	keyAddCanonicalHeader                  = "add_canonical_header"
	keyAddHostHeader                       = "add_host_header"
)

var resourcePullZoneHeaders = &schema.Resource{
	Schema: map[string]*schema.Schema{
		keyEnableAccessControlOriginHeader: {
			Type:        schema.TypeBool,
			Description: fmt.Sprintf("Determines if the CORS headers listed in the %s attribute are applied", keyAccessControlOriginHeaderExtensions),
			Default:     true,
			Optional:    true,
		},
		keyAccessControlOriginHeaderExtensions: {
			Type:        schema.TypeString,
			Description: "CORS Headers will be added to all requests of files with the listed extensions.",
			Default:     "eot, ttf, woff, woff2, css",
			Optional:    true,
			DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
				oldSl := normalizeStrList(old, ',')
				newSl := normalizeStrList(new, ',')

				return strSliceEqual(oldSl, newSl)
			},
		},
		keyAddCanonicalHeader: {
			Type:        schema.TypeBool,
			Description: "Determines if the canonical header should be added by this zone.",
			Optional:    true,
		},
		keyAddHostHeader: {
			Type:        schema.TypeBool,
			Description: "If enabled, the original host header of the request will be forwarded to the origin server.",
			Optional:    true,
		},
	},
}

func headersToResource(pz *bunny.PullZone, d *schema.ResourceData) error {
	m := map[string]interface{}{}

	m[keyEnableAccessControlOriginHeader] = pz.EnableAccessControlOriginHeader
	m[keyAccessControlOriginHeaderExtensions] = strSliceAsNormalizedStr(
		pz.AccessControlOriginHeaderExtensions, ",",
	)
	m[keyAddCanonicalHeader] = pz.AddCanonicalHeader
	m[keyAddHostHeader] = pz.AddHostHeader

	return d.Set(keyHeaders, []map[string]interface{}{m})
}

func headersFromResource(res *bunny.PullZoneUpdateOptions, d *schema.ResourceData) {
	m := structureFromResource(d, keyHeaders)
	if len(m) == 0 {
		return
	}

	res.EnableAccessControlOriginHeader = m.getBoolPtr(keyEnableAccessControlOriginHeader)
	res.AccessControlOriginHeaderExtensions = normalizeStrList(
		m.getStr(keyAccessControlOriginHeaderExtensions), ',',
	)
	res.AddCanonicalHeader = m.getBoolPtr(keyAddCanonicalHeader)
	res.AddHostHeader = m.getBoolPtr(keyAddHostHeader)
}
