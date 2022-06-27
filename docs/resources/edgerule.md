---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bunny_edgerule Resource - bunny"
subcategory: ""
description: |-
  
---

# bunny_edgerule (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action_type` (String) The action type of the Edge Rule.
Valid values: block_request, bypass_perma_cache, disable_optimizer, disable_token_auth, enable_token_auth, force_compression, force_download, force_ssl, ignore_query_string, origin_url, override_cache_time, override_cache_time_public, redirect, set_request_header, set_response_header, set_status_code
- `pull_zone_id` (Number) The ID of the Pull Zone to that Edge Rule belongs.
- `trigger` (Block Set, Min: 1, Max: 5) (see [below for nested schema](#nestedblock--trigger))

### Optional

- `action_parameter_1` (String) The Action parameter 1. The value depends on other parameters of the edge rule.
- `action_parameter_2` (String) The Action parameter 2. The value depends on other parameters of the edge rule.
- `enabled` (Boolean) Determines if the edge rule is currently enabled or not.
- `trigger_matching_type` (String) The trigger matching type.
Valid values: all, any, none

### Read-Only

- `description` (String) The description of the Edge Rule. This field is used internally by Terraform bunny-provider.
- `id` (String) The ID of this resource.

<a id="nestedblock--trigger"></a>
### Nested Schema for `trigger`

Required:

- `pattern_matching_type` (String) The type of pattern matching.
Valid values: all, any, none
- `type` (String) The type of the Trigger.
Valid values: country_code, query_string, random_chance, remote_ip, request_header, request_method, response_header, status_code, url, url_extensions

Optional:

- `parameter_1` (String) The trigger parameter 1. The value depends on the type of trigger.
- `pattern_matches` (Set of String) The list of pattern matches that will trigger the edge rule.

