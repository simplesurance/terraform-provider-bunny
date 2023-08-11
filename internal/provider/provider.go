package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bunny "github.com/simplesurance/bunny-go"
)

const userAgent = "terraform-provider-bunny"
const envVarAPIKey = "BUNNY_API_KEY"
const keyAPIKey = "api_key"

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

// New instantiates a bunny terraform provider.
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			keyAPIKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envVarAPIKey, ""),
				Description: "The bunny.net API Key.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bunny_pullzone": dataSourcePullZone(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"bunny_pullzone":    resourcePullZone(),
			"bunny_edgerule":    resourceEdgeRule(),
			"bunny_hostname":    resourceHostname(),
			"bunny_storagezone": resourceStorageZone(),
		},
		ConfigureContextFunc: newProvider,
	}
}

func newProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get(keyAPIKey).(string)
	if apiKey == "" {
		return nil, diag.FromErr(
			fmt.Errorf("credentials not configured, either %s must be set in the provider config or the environment variable %s",
				keyAPIKey, envVarAPIKey,
			))
	}

	ua := userAgent
	if Version != "" {
		ua += "-" + Version
	}

	log.SetFlags(0)
	return bunny.NewClient(
		apiKey,
		bunny.WithUserAgent(ua),
		bunny.WithHTTPRequestLogger(logger.Debugf),
		bunny.WithHTTPResponseLogger(logger.Debugf),
	), nil
}
