package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyActionParameter1           = "action_parameter_1"
	keyActionParameter2           = "action_parameter_2"
	keyActionType                 = "action_type"
	keyDescription                = "description"
	keyEdgeRuleEnabled            = "enabled"
	keyPullZoneID                 = "pull_zone_id"
	keyTriggerMatchingType        = "trigger_matching_type"
	keyTriggerParameter1          = "parameter_1"
	keyTriggerPatternMatches      = "pattern_matches"
	keyTriggerPatternMatchingType = "pattern_matching_type"
	keyTriggerType                = "type"
	keyTriggers                   = "trigger"
)

func resourceEdgeRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeRuleCreate,
		ReadContext:   resourceEdgeRuleRead,
		DeleteContext: resourceEdgeRuleDelete,
		UpdateContext: resourceEdgeRuleUpdate,
		Schema: map[string]*schema.Schema{
			keyPullZoneID: {
				Type:        schema.TypeInt,
				Description: "The ID of the Pull Zone to that Edge Rule belongs.",
				Required:    true,
			},
			keyActionType: {
				Type: schema.TypeString,
				Description: "The action type of the Edge Rule.\nValid values: " +
					strings.Join(edgeRuleActionTypeKeys, ", "),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(edgeRuleActionTypeKeys, false),
				),
				Required: true,
			},
			keyActionParameter1: {
				Type:        schema.TypeString,
				Description: "The Action parameter 1. The value depends on other parameters of the edge rule.",
				Optional:    true,
			},
			keyActionParameter2: {
				Type:        schema.TypeString,
				Description: "The Action parameter 2. The value depends on other parameters of the edge rule.",
				Optional:    true,
			},
			keyTriggers: {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 5, // otherwise the API returns the error: Maximum 5 condition are allowed per rule.
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						keyTriggerType: {
							Type: schema.TypeString,
							Description: "The type of the Trigger.\nValid values: " +
								strings.Join(edgeRuleTriggerTypeKeys, ", "),
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(edgeRuleTriggerTypeKeys, false),
							),
						},
						keyTriggerPatternMatches: {
							Type:        schema.TypeSet,
							Description: "The list of pattern matches that will trigger the edge rule.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						keyTriggerPatternMatchingType: {
							Type: schema.TypeString,
							Description: "The type of pattern matching.\nValid values: " +
								strings.Join(edgeRuleMatchingTypeKeys, ", "),
							Required: true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice(edgeRuleMatchingTypeKeys, false),
							),
						},

						keyTriggerParameter1: {
							Type:        schema.TypeString,
							Description: "The trigger parameter 1. The value depends on the type of trigger.",
							Optional:    true,
						},
					},
				},
			},
			keyTriggerMatchingType: {
				Type:        schema.TypeString,
				Description: "The trigger matching type.\nValid values: " + strings.Join(edgeRuleMatchingTypeKeys, ", "),
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(edgeRuleMatchingTypeKeys, false),
				),
				Optional: true,
			},
			keyDescription: {
				Type:        schema.TypeString,
				Description: "The description of the Edge Rule. This field is used internally by Terraform bunny-provider.",
				Computed:    true,
			},
			keyEdgeRuleEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if the edge rule is currently enabled or not.",
				Optional:    true,
				Default:     true,
			},
		},
	}
}

// findEdgeRuleGUID retrieves the Pull Zone from the bunny API and returns the guid of the first found edge rule that matches the Description.
func findEdgeRuleGUID(ctx context.Context, clt *bunny.Client, pullZoneID int64, description string) (string, error) {
	pz, err := clt.PullZone.Get(ctx, pullZoneID)
	if err != nil {
		return "", fmt.Errorf("retrieving pull zone failed: %w", err)
	}

	for _, edgeRule := range pz.EdgeRules {
		if edgeRule.Description != nil && *edgeRule.Description == description {
			if edgeRule.GUID == nil {
				return "", errors.New("found edge rule with matching description but guid is nil")
			}

			return *edgeRule.GUID, nil
		}
	}

	return "", fmt.Errorf("pull zone has no edge rule rule with internal identifier in description")
}

func resourceEdgeRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	// The bunny API endpoint does not return the ID of a newly created
	// Edge Rule.  To be able to identify the created edge rule uniquely
	// (including distinguishing them from edge-rules with the same
	// settings created in parallel via e.g. the UI), we store an own ID in
	// the description field.  The ID is only used once during creation to
	// initially find the created Edge Rule. After it was found the GUID is
	// used for identification. The description could be safely overwritten
	internalEdgeRuleID := "terraform-provider-bunny id: " + uuid.New().String()
	if err := d.Set(keyDescription, internalEdgeRuleID); err != nil {
		return diag.FromErr(err)
	}

	opts, err := resourceDataToAddOrUpdateEdgeRuleOptions(d)
	if err != nil {
		return diagsErrFromErr("setting description failed", err)
	}

	pullZoneID := int64(d.Get(keyPullZoneID).(int))

	err = clt.PullZone.AddOrUpdateEdgeRule(ctx, pullZoneID, opts)
	if err != nil {
		return diagsErrFromErr("creating edge rule failed", err)
	}

	guid, err := findEdgeRuleGUID(ctx, clt, pullZoneID, internalEdgeRuleID)
	if err != nil {
		return diagsErrFromErr(
			fmt.Sprintf("edge rule (description: %q) created successfully, looking up its guid failed", internalEdgeRuleID), err)
	}

	d.SetId(guid)

	return nil
}

func resourceEdgeRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	opts, err := resourceDataToAddOrUpdateEdgeRuleOptions(d)
	if err != nil {
		return diag.FromErr(err)
	}

	pullZoneID := int64(d.Get(keyPullZoneID).(int))

	err = clt.PullZone.AddOrUpdateEdgeRule(ctx, pullZoneID, opts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating edge rule failed: %w", err))
	}

	return nil
}

func edgeRuleTriggerTypeToInt(triggerType string) (int, error) {
	if k, exists := edgeRuleTriggerTypesStr[triggerType]; exists {
		return k, nil
	}

	return -1, fmt.Errorf("unsupported trigger type type: %q", triggerType)
}

func resourceDataToEdgeRuleTriggers(d *schema.ResourceData) ([]*bunny.EdgeRuleTrigger, error) {
	triggerSet := d.Get(keyTriggers).(*schema.Set)
	if triggerSet.Len() == 0 {
		return nil, nil
	}

	res := make([]*bunny.EdgeRuleTrigger, 0, triggerSet.Len())

	for _, item := range triggerSet.List() {
		i := item.(map[string]interface{})

		triggerType, err := edgeRuleTriggerTypeToInt(i[keyTriggerType].(string))
		if err != nil {
			return nil, err
		}

		var patternMatches []string
		if val := i[keyTriggerPatternMatches]; val != nil {
			patternMatches = strSetAsSlice(val)
		}

		patternMatchingType, err := strIntMapGet(
			edgeRuleMatchingTypesStr, i[keyTriggerPatternMatchingType].(string),
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", keyTriggerMatchingType, err)
		}

		triggerParam1 := i[keyTriggerParameter1].(string)

		res = append(res, &bunny.EdgeRuleTrigger{
			Type:                &triggerType,
			PatternMatches:      patternMatches,
			PatternMatchingType: &patternMatchingType,
			Parameter1:          &triggerParam1,
		})
	}

	return res, nil
}

func resourceDataToAddOrUpdateEdgeRuleOptions(d *schema.ResourceData) (*bunny.AddOrUpdateEdgeRuleOptions, error) {
	var guid *string
	if id := d.Id(); id != "" {
		guid = &id
	}

	actionType, err := strIntMapGet(
		edgeRuleActionTypesStr,
		d.Get(keyActionType).(string),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", keyActionType, err)
	}

	matchingType, err := strIntMapGet(
		edgeRuleMatchingTypesStr,
		d.Get(keyTriggerMatchingType).(string),
	)
	if err != nil {
		return nil, err
	}

	triggers, err := resourceDataToEdgeRuleTriggers(d)
	if err != nil {
		return nil, fmt.Errorf("converting edge rule triggers failed: %w", err)
	}

	return &bunny.AddOrUpdateEdgeRuleOptions{
		GUID:                guid,
		Enabled:             getBoolPtr(d, keyEnabled),
		ActionType:          &actionType,
		ActionParameter1:    getStrPtr(d, keyActionParameter1),
		ActionParameter2:    getStrPtr(d, keyActionParameter2),
		Triggers:            triggers,
		TriggerMatchingType: &matchingType,
		Description:         getStrPtr(d, keyDescription),
	}, nil
}

func resourceEdgeRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	edgeRuleGUID := d.Id()
	pullZoneID := int64(d.Get(keyPullZoneID).(int))

	err := clt.PullZone.DeleteEdgeRule(ctx, pullZoneID, edgeRuleGUID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceEdgeRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	edgeRuleGUID := d.Id()
	pullZoneID := int64(d.Get(keyPullZoneID).(int))

	pz, err := clt.PullZone.Get(ctx, pullZoneID)
	if err != nil {
		return diagsErrFromErr("retrieving pull zone failed", err)
	}

	if len(pz.EdgeRules) == 0 {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "pull zone has no edge rules",
		}}
	}

	var edgeRule *bunny.EdgeRule
	for _, er := range pz.EdgeRules {
		if er.GUID != nil && *er.GUID == edgeRuleGUID {
			edgeRule = er
			break
		}
	}

	if edgeRule == nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "edge rule not found",
			Detail:   fmt.Sprintf("pull zone with id %d, has no edge rule with guid: %q", pullZoneID, edgeRuleGUID),
		}}
	}

	if err := edgeRuleToResourceData(edgeRule, d); err != nil {
		return diagsErrFromErr("converting edge rule api type to terraform ResourceData failed", err)
	}

	return nil

}

func edgeRuleToResourceData(edgeRule *bunny.EdgeRule, d *schema.ResourceData) error {
	if edgeRule.GUID == nil || *edgeRule.GUID == "" {
		return errors.New("guid is empty")
	}

	d.SetId(*edgeRule.GUID)

	actionType, err := intStrMapGet(edgeRuleActionTypesInt, edgeRule.ActionType)
	if err != nil {
		return fmt.Errorf("%s: %w", keyActionType, err)
	}

	if err := d.Set(keyActionType, actionType); err != nil {
		return err
	}

	if err := d.Set(keyActionParameter1, edgeRule.ActionParameter1); err != nil {
		return err
	}

	if err := d.Set(keyActionParameter2, edgeRule.ActionParameter2); err != nil {
		return err
	}

	err = edgeRuleTriggerToResourceData(edgeRule.Triggers, d)
	if err != nil {
		return fmt.Errorf("converting triggers to resource data failed: %w", err)
	}

	matchingType, err := intStrMapGet(edgeRuleMatchingTypesInt, edgeRule.TriggerMatchingType)
	if err != nil {
		return fmt.Errorf("%s: %w", keyTriggerMatchingType, err)
	}

	if err := d.Set(keyTriggerMatchingType, matchingType); err != nil {
		return err
	}

	if err := d.Set(keyDescription, edgeRule.Description); err != nil {
		return err
	}

	if err := d.Set(keyEnabled, edgeRule.Enabled); err != nil {
		return err
	}

	return nil
}

func edgeRuleTriggerToResourceData(triggers []*bunny.EdgeRuleTrigger, d *schema.ResourceData) error {
	res := make([]map[string]interface{}, 0, len(triggers))

	for _, trigger := range triggers {
		triggerType, err := intStrMapGet(edgeRuleTriggerTypesInt, trigger.Type)
		if err != nil {
			return fmt.Errorf("%s: %w", triggerType, err)
		}

		patternMatchingType, err := intStrMapGet(edgeRuleMatchingTypesInt, trigger.PatternMatchingType)
		if err != nil {
			return fmt.Errorf("%s: %w", triggerType, err)
		}

		entry := make(map[string]interface{}, 4)
		entry[keyTriggerType] = triggerType
		entry[keyTriggerPatternMatches] = trigger.PatternMatches
		entry[keyTriggerPatternMatchingType] = patternMatchingType
		entry[keyTriggerParameter1] = trigger.Parameter1

		res = append(res, entry)
	}

	return d.Set(keyTriggers, res)
}
