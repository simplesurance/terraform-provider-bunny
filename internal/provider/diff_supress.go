package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func diffSupressIntUnset(k, old, new string, d *schema.ResourceData) bool {
	return new == "0"
}

func diffSupressMissingOptionalBlock(k, old, new string, d *schema.ResourceData) bool {
	return old == "1" && new == "0"
}
