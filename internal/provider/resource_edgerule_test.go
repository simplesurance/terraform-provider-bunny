package provider

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	ptr "github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	bunny "github.com/simplesurance/bunny-go"
)

type edgeRulesWanted struct {
	// pullZoneName is the name of the pull-zone to that the edgeRules belong
	TerraformPullZoneResourceName string
	PullZoneName                  string
	EdgeRules                     []*bunny.EdgeRule
}

func checkEdgeRulesState(t *testing.T, wanted *edgeRulesWanted) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		strID, err := idFromState(s, wanted.TerraformPullZoneResourceName)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			return fmt.Errorf("could not convert resource ID %q to int64: %w", id, err)
		}

		pz, err := clt.PullZone.Get(context.Background(), int64(id))
		if err != nil {
			return fmt.Errorf("fetching pull-zone with id %d from api client failed: %w", id, err)
		}

		if err := stringsAreEqual(wanted.PullZoneName, pz.Name); err != nil {
			return fmt.Errorf("name of created pullzone differs: %w", err)
		}

		if len(pz.EdgeRules) != len(wanted.EdgeRules) {
			return fmt.Errorf("api returned pull request with %d edge rules, expected %d",
				len(pz.EdgeRules), len(wanted.EdgeRules),
			)
		}

		sortEdgeRules(t, wanted.EdgeRules)
		sortEdgeRules(t, pz.EdgeRules)

		for i := range pz.EdgeRules {
			diff := edgeRuleDiff(t, wanted.EdgeRules[i], pz.EdgeRules[i])
			if len(diff) != 0 {
				return fmt.Errorf("wanted and actual edge rule with idx %d differs:\n%s", i, strings.Join(diff, "\n"))
			}
		}

		return nil
	}
}

func sortEdgeRules(t *testing.T, rules []*bunny.EdgeRule) {
	sort.Slice(rules, func(i, j int) bool {
		a := rules[i]
		b := rules[j]

		if a.ActionType != nil && b.ActionType != nil && ptr.GetInt(a.ActionType) != ptr.GetInt(b.ActionType) {
			return ptr.GetInt(a.ActionType) < ptr.GetInt(b.ActionType)
		}

		if a.Description != nil && b.Description != nil && ptr.GetString(a.Description) != ptr.GetString(b.Description) {
			return ptr.GetString(a.Description) < ptr.GetString(b.Description)
		}

		t.Logf("WARN: edge rules slice elements not sorted, found no comparable attributes were equal")
		return false
	})
}

var edgeRuleDiffIgnoredFields = map[string]struct{}{
	"GUID":        {}, // is set as ID in resourceData, GUID does not exist in resourceData
	"Description": {}, // computed field, used internally by our provider for initial identification
}

func edgeRuleDiff(t *testing.T, a, b interface{}) []string {
	t.Helper()
	return diffStructs(t, a, b, edgeRuleDiffIgnoredFields)
}

func defPullZoneHostname(pullzoneName string) string {
	return fmt.Sprintf("%s.b-cdn.net", pullzoneName)
}

func TestAccEdgeRule_full(t *testing.T) {
	pzName := randPullZoneName()

	tfPz := fmt.Sprintf(`
resource "bunny_pullzone" "mypz" {
	name = "%s"
	origin_url ="https://bunny.net"
}`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tfPz + `
resource "bunny_edgerule" "myer" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "block_request"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
}`,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					PullZoneName:                  pzName,
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					EdgeRules: []*bunny.EdgeRule{
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeBlockRequest),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
					},
				}),
			},

			// change the trigger and add a second trigger
			{
				Config: tfPz + `
resource "bunny_edgerule" "myer" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "block_request"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["40"]
	}
	trigger {
		pattern_matching_type = "none"
		type = "remote_ip"
		pattern_matches = ["127.0.0.1"]
	}
}`,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					PullZoneName:                  pzName,
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					EdgeRules: []*bunny.EdgeRule{
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeBlockRequest),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"40"},
								},
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeNone),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRemoteIP),
									PatternMatches:      []string{"127.0.0.1"},
								},
							},
						},
					},
				}),
			},

			// replace the edge rule, add other combinations
			{
				Config: tfPz + `
resource "bunny_edgerule" "er1" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "override_cache_time"
	action_parameter_1 = "10"
	trigger_matching_type = "any"
	trigger {
		type = "request_header"
		parameter_1 = "user"
		pattern_matching_type = "any"
		pattern_matches = ["hans", "franz"]
	}
	trigger {
		pattern_matching_type = "none"
		type = "country_code"
		pattern_matches = ["de","dk"]
	}
}

resource "bunny_edgerule" "er2" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "force_download"
	trigger_matching_type = "any"
	trigger {
		type = "response_header"
		parameter_1 = "force_dl"
		pattern_matching_type = "any"
		pattern_matches = ["yes", "true"]
	}
	trigger {
		pattern_matching_type = "all"
		type = "url"
		pattern_matches = ["https://localhost", "https://bunny.net", "https://google.com"]
	}
}

resource "bunny_edgerule" "er3" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "set_request_header"
	action_parameter_1 = "hostname"
	action_parameter_2 = "{{hostname}}"
	trigger_matching_type = "any"
	trigger {
		type = "query_string"
		pattern_matching_type = "any"
		pattern_matches = ["set_hostname"]
	}
}
`,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					PullZoneName:                  pzName,
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					EdgeRules: []*bunny.EdgeRule{
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeSetRequestHeader),
							ActionParameter1:    ptr.ToString("hostname"),
							ActionParameter2:    ptr.ToString("{{hostname}}"),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeURLQueryString),
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									PatternMatches:      []string{"set_hostname"},
								},
							},
						},
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeOverrideCacheTime),
							ActionParameter1:    ptr.ToString("10"),
							ActionParameter2:    ptr.ToString(""),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRequestHeader),
									Parameter1:          ptr.ToString("user"),
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									PatternMatches:      []string{"hans", "franz"},
								},
								{
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeCountryCode),
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeNone),
									PatternMatches:      []string{"de", "dk"},
								},
							},
						},
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeForceDownload),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeResponseHeader),
									Parameter1:          ptr.ToString("force_dl"),
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
									PatternMatches:      []string{"yes", "true"},
								},
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeURL),
									PatternMatches:      []string{"https://localhost", "https://bunny.net", "https://google.com"},
								},
							},
						},
					},
				}),
			},

			{
				Config:  tfPz,
				Destroy: true,
			},
		},
		CheckDestroy: checkPullZoneNotExists(pzName),
	})
}

func TestAccEdgeRule_basic(t *testing.T) {
	pzName := randPullZoneName()
	tf := fmt.Sprintf(`
resource "bunny_pullzone" "mypz" {
	name = "%s"
	origin_url ="https://bunny.net"
}
resource "bunny_edgerule" "myer" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "block_request"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
} `, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					PullZoneName:                  pzName,
					EdgeRules: []*bunny.EdgeRule{
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeBlockRequest),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
					},
				}),
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
	})
}

func TestAccEdgeRule_delete(t *testing.T) {
	pzName := randPullZoneName()

	tfPz := fmt.Sprintf(`
resource "bunny_pullzone" "mypz" {
	name = "%s"
	origin_url ="https://bunny.net"
}
`, pzName)

	tfEdgeRule := `resource "bunny_edgerule" "myer" {
	pull_zone_id = bunny_pullzone.mypz.id
	action_type = "block_request"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
}`

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tfPz + tfEdgeRule,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					PullZoneName:                  pzName,
					EdgeRules: []*bunny.EdgeRule{
						{
							Enabled:             ptr.ToBool(true),
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeBlockRequest),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
					},
				}),
			},
			{
				Config: tfPz,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					PullZoneName:                  pzName,
				}),
			},
		},
	})
}

func TestAccEdgeRule_enable_disable(t *testing.T) {
	pzName := randPullZoneName()

	tfPz := fmt.Sprintf(`
resource "bunny_pullzone" "mypz" {
	name = "%s"
	origin_url ="https://bunny.net"
}`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tfPz + `
resource "bunny_edgerule" "er1" {
	pull_zone_id = bunny_pullzone.mypz.id
	enabled = false
	action_type = "block_request"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
}
resource "bunny_edgerule" "er2" {
	pull_zone_id = bunny_pullzone.mypz.id
	enabled = true
	action_type = "force_download"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
}
`,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					PullZoneName:                  pzName,
					EdgeRules: []*bunny.EdgeRule{
						{
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeBlockRequest),
							Enabled:             ptr.ToBool(false),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
						{
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeForceDownload),
							Enabled:             ptr.ToBool(true),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
					},
				}),
			},
			// enable the previously disabled edge rule, disable the previously enabled one
			{
				Config: tfPz + `
resource "bunny_edgerule" "er1" {
	pull_zone_id = bunny_pullzone.mypz.id
	enabled = true
	action_type = "block_request"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
}
resource "bunny_edgerule" "er2" {
	pull_zone_id = bunny_pullzone.mypz.id
	enabled = false
	action_type = "force_download"
	trigger_matching_type = "all"
	trigger {
		pattern_matching_type = "any"
		type = "random_chance"
		pattern_matches = ["30"]
	}
}
`,
				Check: checkEdgeRulesState(t, &edgeRulesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.mypz",
					PullZoneName:                  pzName,
					EdgeRules: []*bunny.EdgeRule{
						{
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeBlockRequest),
							Enabled:             ptr.ToBool(true),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
						{
							ActionType:          ptr.ToInt(bunny.EdgeRuleActionTypeForceDownload),
							Enabled:             ptr.ToBool(false),
							TriggerMatchingType: ptr.ToInt(bunny.MatchingTypeAll),
							ActionParameter1:    ptr.ToString(""),
							ActionParameter2:    ptr.ToString(""),
							Triggers: []*bunny.EdgeRuleTrigger{
								{
									PatternMatchingType: ptr.ToInt(bunny.MatchingTypeAny),
									Type:                ptr.ToInt(bunny.EdgeRuleTriggerTypeRandomChance),
									PatternMatches:      []string{"30"},
								},
							},
						},
					},
				}),
			},
		},
	})
}
