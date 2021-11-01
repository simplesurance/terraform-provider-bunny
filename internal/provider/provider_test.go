package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePrefix is the prefix that should be used when creating resources at
// the provider in integration tests.
const resourcePrefix = "tf-test-"

var testProvider *schema.Provider
var testProviders map[string]*schema.Provider

func init() {
	testProvider = New()
	testProviders = map[string]*schema.Provider{
		"bunny": testProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
