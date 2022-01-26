package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	keyCertificateCertificateData = "certificate_data"
	keyCertificatePrivateKeyData  = "private_key_data"
)

var resourceHostnameCertificate = &schema.Resource{
	Schema: map[string]*schema.Schema{
		keyCertificateCertificateData: {
			Type:        schema.TypeString,
			Description: "The public key.",
			Required:    true,
			ForceNew:    true,
		},
		keyCertificatePrivateKeyData: {
			Type:        schema.TypeString,
			Description: "The private key.",
			Required:    true,
			ForceNew:    true,
		},
	},
}
