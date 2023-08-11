package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)
const (
	keyPullZoneID  = "pull_zone_id"
)

func dataSourcePullZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePullZoneRead,

		Schema: map[string]*schema.Schema{
			keyPullZoneID: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			keyAWSSigningKey: {
				Type:        schema.TypeString,
				Description: "AWS Signing Key",
				Computed:    true,
			},
			keyAWSSigningRegionName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The AWS Signing region name.",
			},
			keyAWSSigningSecret: {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "The AWS Signing region secret.",
			},
			keyCnameDomain: {
				Type:        schema.TypeString,
				Description: "The CNAME domain of the Pull Zone for setting up custom hostnames.",
				Computed:    true,
			},
			keyName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Pull Zone.",
			},
			keyStorageZoneID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the storage zone that the Pull Zone is linked to.",
			},
			keyZoneSecurityEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the URL secure token authentication security is enabled.",
			},
			keyZoneSecurityIncludeHashRemoteIP: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the zone security hash should include the remote IP.",
			},
			keyZoneSecurityKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The security key used for secure URL token authentication.",
			},
		},
	}
}

func dataSourcePullZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get(keyPullZoneID).(int)
	d.SetId(fmt.Sprint(id))

	pz, err := readPullZone(ctx, d, meta)
	if err != nil {
		return err
	}

	if err := pullZoneToResourceShared(pz, d); err != nil {
		return diagsErrFromErr("converting api type to datasource after successful read failed", err)
	}

	return nil
}
