package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyVideoCount                       = "video_count"
	keyTrafficUsage                     = "traffic_usage"
	keyStorageUsage                     = "storage_usage"
	keyDateCreated                      = "date_created"
	keyReadOnlyAPIKey                   = "read_only_api_key"
	keyHasWatermark                     = "has_watermark"
	keyWatermarkPositionLeft            = "watermark_position_left"
	keyWatermarkPositionTop             = "watermark_position_top"
	keyWatermarkWidth                   = "watermark_width"
	keyWatermarkHeight                  = "watermark_height"
	keyEnabledResolutions               = "enabled_resolutions"
	keyViAiPublisherID                  = "vi_ai_publisher_id"
	keyVastTagURL                       = "vast_tag_url"
	keyWebhookURL                       = "webhook_url"
	keyCaptionsFontSize                 = "captions_font_size"
	keyCaptionsFontColor                = "captions_font_color"
	keyCaptionsBackground               = "captions_background"
	keyUILanguage                       = "ui_language"
	keyAllowEarlyPlay                   = "allow_early_play"
	keyPlayerTokenAuthenticationEnabled = "player_token_authentication_enabled"
	keyBlockNoneReferrer                = "block_none_referrer"
	keyEnableMP4Fallback                = "enable_mp4_fallback"
	keyKeepOriginalFiles                = "keep_original_files"
	keyAllowDirectPlay                  = "allow_direct_play"
	keyEnableDRM                        = "enable_drm"
	keyBitrate240p                      = "bitrate_240p"
	keyBitrate360p                      = "bitrate_360p"
	keyBitrate480p                      = "bitrate_480p"
	keyBitrate720p                      = "bitrate_720p"
	keyBitrate1080p                     = "bitrate_1080p"
	keyBitrate1440p                     = "bitrate_1440p"
	keyBitrate2160p                     = "bitrate_2160p"
	keyAPIAccessKey                     = "api_access_key"
	keyShowHeatmap                      = "show_heatmap"
	keyEnableContentTagging             = "enable_content_tagging"
	keyPullZoneType                     = "pull_zone_type"
	keyCustomHTML                       = "custom_html"
	keyControls                         = "controls"
	keyPlayerKeyColor                   = "player_key_color"
	keyFontFamily                       = "font_family"
	keyEnableTokenAuthentication        = "enable_token_authentication"
	keyEnableTokenIPVerification        = "enable_token_ip_verification"
	keyPullZoneID                       = "pull_zone_id"
	keyResetToken                       = "reset_token"
)

func resourceVideoLibrary() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVideoLibraryCreate,
		ReadContext:   resourceVideoLibraryRead,
		UpdateContext: resourceVideoLibraryUpdate,
		DeleteContext: resourceVideoLibraryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// TODO: missing update props
			// resettoken

			// immutable properties
			// NOTE: immutable properties are made immutable via
			// validation in the `CustomizeDiff` function. There
			// should be a validation function for each immutable
			// property.
			keyReplicationRegions: {
				Type:        schema.TypeSet,
				Description: "The geo-replication regions of the underlying storage zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(
						validation.StringInSlice([]string{"UK", "SE", "NY", "LA", "SG", "SY", "BR", "JH"}, false),
					),
				},
				Optional: true,
			},

			// mutable properties
			//  - required
			keyName: {
				Type:        schema.TypeString,
				Description: "The name of the Video Library.",
				Required:    true,
			},

			//  - mutable
			keyWatermarkPositionLeft: {
				Type:        schema.TypeInt,
				Description: "The left offset of the watermark position (in %).",
				Optional:    true,
			},
			keyWatermarkPositionTop: {
				Type:        schema.TypeInt,
				Description: "The top offset of the watermark position (in %).",
				Optional:    true,
			},
			keyWatermarkWidth: {
				Type:        schema.TypeInt,
				Description: "The width of the watermark (in %).",
				Optional:    true,
			},
			keyWatermarkHeight: {
				Type:        schema.TypeInt,
				Description: "The height of the watermark (in %).",
				Optional:    true,
			},
			keyEnabledResolutions: {
				Type:        schema.TypeString,
				Description: "The comma separated list of enabled resolutions. Possible values: 240p, 360p, 480p, 720p, 1080p, 1440p, 2160p.",
				Optional:    true,
				Default:     "240p,360p,480p,720p,1080p",
			},
			keyViAiPublisherID: {
				Type:        schema.TypeString,
				Description: "The vi.ai publisher id for advertising configuration.",
				Optional:    true,
			},
			keyVastTagURL: {
				Type:        schema.TypeString,
				Description: "The URL of the VAST tag endpoint for advertising configuration.",
				Optional:    true,
			},
			keyWebhookURL: {
				Type:        schema.TypeString,
				Description: "The webhook URL of the video library.",
				Optional:    true,
			},
			keyCaptionsFontSize: {
				Type:        schema.TypeInt,
				Description: "The captions display font size.",
				Optional:    true,
			},
			keyCaptionsFontColor: {
				Type:        schema.TypeString,
				Description: "The captions display font color.",
				Optional:    true,
			},
			keyCaptionsBackground: {
				Type:        schema.TypeString,
				Description: "The captions display background color.",
				Optional:    true,
			},
			keyUILanguage: {
				Type:        schema.TypeString,
				Description: "The UI language of the player.",
				Optional:    true,
			},
			keyAllowEarlyPlay: {
				Type:        schema.TypeBool,
				Description: "Determines if the Early-Play feature is enabled.",
				Optional:    true,
			},
			keyPlayerTokenAuthenticationEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if the player token authentication is enabled.",
				Optional:    true,
			},
			keyBlockNoneReferrer: {
				Type:        schema.TypeBool,
				Description: "Determines if the requests without a referrer are blocked.",
				Optional:    true,
				Default:     true,
			},
			keyEnableMP4Fallback: {
				Type:        schema.TypeBool,
				Description: "Determines if the MP4 fallback feature is enabled.",
				Optional:    true,
				Default:     true,
			},
			keyKeepOriginalFiles: {
				Type:        schema.TypeBool,
				Description: "Determines if the original video files should be stored after encoding.",
				Optional:    true,
				Default:     true,
			},
			keyAllowDirectPlay: {
				Type:        schema.TypeBool,
				Description: "Determines direct play URLs are enabled for the library.",
				Optional:    true,
				Default:     true,
			},
			keyEnableDRM: {
				Type:        schema.TypeBool,
				Description: "Determines if the MediaCage basic DRM is enabled.",
				Optional:    true,
			},
			keyBitrate240p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 240p videos.",
				Optional:    true,
				Default:     600,
			},
			keyBitrate360p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 360p videos.",
				Optional:    true,
				Default:     800,
			},
			keyBitrate480p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 480p videos.",
				Optional:    true,
				Default:     1400,
			},
			keyBitrate720p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 720p videos.",
				Optional:    true,
				Default:     2800,
			},
			keyBitrate1080p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 1080p videos.",
				Optional:    true,
				Default:     5000,
			},
			keyBitrate1440p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 1440p videos.",
				Optional:    true,
				Default:     8000,
			},
			keyBitrate2160p: {
				Type:        schema.TypeInt,
				Description: "The bitrate used for encoding 2160p videos.",
				Optional:    true,
				Default:     25000,
			},
			keyShowHeatmap: {
				Type:        schema.TypeBool,
				Description: "Determines if the video watch heatmap should be displayed in the player.",
				Optional:    true,
			},
			keyEnableContentTagging: {
				Type:        schema.TypeBool,
				Description: "Determines if content tagging should be enabled for this library.",
				Optional:    true,
				Default:     true,
			},
			keyCustomHTML: {
				Type:        schema.TypeString,
				Description: "The custom HTMl that is added into the head of the HTML player.",
				Optional:    true,
			},
			keyControls: {
				Type:        schema.TypeString,
				Description: "The comma separated list of controls that will be displayed in the video player. Possible values: play-large, play, progress, current-time, mute, volume, captions, settings, pip, airplay, fullscreen.",
				Optional:    true,
			},
			keyPlayerKeyColor: {
				Type:        schema.TypeString,
				Description: "The key color of the player.",
				Optional:    true,
			},
			keyFontFamily: {
				Type:        schema.TypeString,
				Description: "The captions font family.",
				Optional:    true,
				Default:     "Rubik",
			},
			// note: this is not available on the read
			keyEnableTokenAuthentication: {
				Type:        schema.TypeBool,
				Description: "Determines if the token authentication should be enabled.",
				Optional:    true,
			},
			// note: this is not available on the read
			keyEnableTokenIPVerification: {
				Type:        schema.TypeBool,
				Description: "Determines if the token IP verification should be enabled.",
				Optional:    true,
			},

			// computed properties
			keyVideoCount: {
				Type:        schema.TypeInt,
				Description: "The number of videos in the video library.",
				Computed:    true,
			},
			keyTrafficUsage: {
				Type:        schema.TypeInt,
				Description: "The amount of traffic usage this month.",
				Computed:    true,
			},
			keyStorageUsage: {
				Type:        schema.TypeInt,
				Description: "The total amount of storage used by the library.",
				Computed:    true,
			},
			keyDateCreated: {
				Type:        schema.TypeString,
				Description: "The date when the video library was created.",
				Computed:    true,
			},
			keyAPIKey: {
				Type:        schema.TypeString,
				Description: "The API key used for authenticating with the video library.",
				Computed:    true,
				Sensitive:   true,
			},
			keyReadOnlyAPIKey: {
				Type:        schema.TypeString,
				Description: "The read-only API key used for authenticating with the video library.",
				Computed:    true,
				Sensitive:   true,
			},
			keyHasWatermark: {
				Type:        schema.TypeBool,
				Description: "Determines if the video library has a watermark configured.",
				Computed:    true,
			},
			keyPullZoneID: {
				Type:        schema.TypeInt,
				Description: "The ID of the connected underlying pull zone.",
				Computed:    true,
			},
			keyStorageZoneID: {
				Type:        schema.TypeInt,
				Description: "The ID of the connected underlying storage zone.",
				Computed:    true,
			},
			keyAPIAccessKey: {
				Type:        schema.TypeString,
				Description: "The API access key for the library. Only added when the includeAccessKey parameter is set.",
				Computed:    true,
				Sensitive:   true,
			},
			keyPullZoneType: {
				Type:        schema.TypeInt,
				Description: "The type of the pull zone attached. Premium = 0, Volume = 1.",
				Computed:    true,
			},
			keyAllowedReferrers: {
				Type:        schema.TypeSet,
				Description: "The list of allowed referrer domains allowed to access the library.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			keyBlockedReferrers: {
				Type:        schema.TypeSet,
				Description: "The list of blocked referrer domains blocked from accessing the library.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},

		CustomizeDiff: customdiff.All(
			func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
				// if the resource hasn't been created yet, skip validation
				if d.Id() == "" {
					return nil
				}

				old, new := d.GetChange(keyReplicationRegions)

				// verify that none of the previous replication regions have been removed.
				var oldRep *schema.Set = old.(*schema.Set)
				var newRep *schema.Set = new.(*schema.Set)
				areEqual := cmp.Equal(oldRep.List(), newRep.List())
				if !areEqual {
					return immutableVideoLibraryReplicationRegionError(
						keyReplicationRegions,
						oldRep.List(),
						newRep.List(),
					)
				}

				return nil
			},
		),
	}
}

func immutableVideoLibraryReplicationRegionError(key string, from []interface{}, to []interface{}) error {
	const message = "'%s' can not be mutated.\n" +
		"This error occurred when attempting to updates values %+q from '%s' to '%s'.\n" +
		"To modify an existing '%s' the 'bunny_videolibrary' must be deleted and recreated.\n" +
		"WARNING: deleting a 'bunny_videolibrary' will also delete all the data it contains!"
	return fmt.Errorf(
		message,
		key,
		from,
		to,
		key,
		key,
	)
}

func resourceVideoLibraryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	vl, err := clt.VideoLibrary.Add(ctx, &bunny.VideoLibraryAddOptions{
		Name:               getStrPtr(d, keyName),
		ReplicationRegions: getStrSetAsSlice(d, keyReplicationRegions),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating video library failed: %w", err))
	}

	d.SetId(strconv.FormatInt(*vl.ID, 10))
	if err := d.Set(keyDateCreated, time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}

	// VideoLibrary.Add() only supports to set a subset of a Video Library object,
	// call Update to set the remaining ones.
	if diags := resourceVideoLibraryUpdate(ctx, d, meta); diags.HasError() {
		// if updating fails the vl was still created, initialize with the vl
		// returned from the Add operation
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "setting video library attributes via update failed",
		})

		if err := videoLibraryToResource(vl, d); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "converting api-type to resource data failed: " + err.Error(),
			})

		}

		return diags
	}

	return nil
}

func resourceVideoLibraryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	videoLibrary, err := videoLibraryFromResource(d)
	if err != nil {
		return diagsErrFromErr("converting resource to API type failed", err)
	}

	id, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedVideoLibrary, err := clt.VideoLibrary.Update(ctx, id, videoLibrary)
	if err != nil {
		// The video library contains fields (EnableTokenAuthentication,
		// EnableTokenIPVerification, & ResetToken) that are only updated by the
		// provider and not retrieved via ReadContext(). This causes that we run into the bug
		// https://github.com/hashicorp/terraform-plugin-sdk/issues/476.
		// As workaround d.Partial(true) is called.
		d.Partial(true)
		return diagsErrFromErr("updating video library via API failed", err)
	}

	if err := videoLibraryToResource(updatedVideoLibrary, d); err != nil {
		return diagsErrFromErr("converting api type to resource data after successful update failed", err)
	}

	return nil
}

func resourceVideoLibraryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	id, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vl, err := clt.VideoLibrary.Get(ctx, id, &bunny.VideoLibraryGetOpts{IncludeAccessKey: true})
	if err != nil {
		return diagsErrFromErr("could not retrieve video library", err)
	}

	if err := videoLibraryToResource(vl, d); err != nil {
		return diagsErrFromErr("converting api type to resource data after successful read failed", err)
	}

	return nil
}

func resourceVideoLibraryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	id, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = clt.VideoLibrary.Delete(ctx, id)
	if err != nil {
		return diagsErrFromErr("could not delete video library", err)
	}

	d.SetId("")

	return nil
}

// videoLibraryToResource sets fields in d to the values in vl.
func videoLibraryToResource(vl *bunny.VideoLibrary, d *schema.ResourceData) error {
	if vl.ID != nil {
		d.SetId(strconv.FormatInt(*vl.ID, 10))
	}

	if err := d.Set(keyName, vl.Name); err != nil {
		return err
	}
	if err := d.Set(keyVideoCount, vl.VideoCount); err != nil {
		return err
	}
	if err := d.Set(keyTrafficUsage, vl.TrafficUsage); err != nil {
		return err
	}
	if err := d.Set(keyStorageUsage, vl.StorageUsage); err != nil {
		return err
	}
	if err := d.Set(keyDateCreated, vl.DateCreated); err != nil {
		return err
	}
	if err := setStrSet(d, keyReplicationRegions, vl.ReplicationRegions, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := d.Set(keyAPIKey, vl.APIKey); err != nil {
		return err
	}
	if err := d.Set(keyReadOnlyAPIKey, vl.ReadOnlyAPIKey); err != nil {
		return err
	}
	if err := d.Set(keyHasWatermark, vl.HasWatermark); err != nil {
		return err
	}
	if err := d.Set(keyWatermarkPositionLeft, vl.WatermarkPositionLeft); err != nil {
		return err
	}
	if err := d.Set(keyWatermarkPositionTop, vl.WatermarkPositionTop); err != nil {
		return err
	}
	if err := d.Set(keyWatermarkWidth, vl.WatermarkWidth); err != nil {
		return err
	}
	if err := d.Set(keyPullZoneID, vl.PullZoneID); err != nil {
		return err
	}
	if err := d.Set(keyStorageZoneID, vl.StorageZoneID); err != nil {
		return err
	}
	if err := d.Set(keyWatermarkHeight, vl.WatermarkHeight); err != nil {
		return err
	}
	if err := d.Set(keyEnabledResolutions, vl.EnabledResolutions); err != nil {
		return err
	}
	if err := d.Set(keyViAiPublisherID, vl.ViAiPublisherID); err != nil {
		return err
	}
	if err := d.Set(keyVastTagURL, vl.VastTagURL); err != nil {
		return err
	}
	if err := d.Set(keyWebhookURL, vl.WebhookURL); err != nil {
		return err
	}
	if err := d.Set(keyCaptionsFontSize, vl.CaptionsFontSize); err != nil {
		return err
	}
	if err := d.Set(keyCaptionsFontColor, vl.CaptionsFontColor); err != nil {
		return err
	}
	if err := d.Set(keyCaptionsBackground, vl.CaptionsBackground); err != nil {
		return err
	}
	if err := d.Set(keyUILanguage, vl.UILanguage); err != nil {
		return err
	}
	if err := d.Set(keyAllowEarlyPlay, vl.AllowEarlyPlay); err != nil {
		return err
	}
	if err := d.Set(keyPlayerTokenAuthenticationEnabled, vl.PlayerTokenAuthenticationEnabled); err != nil {
		return err
	}
	if err := setStrSet(d, keyAllowedReferrers, vl.AllowedReferrers, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := setStrSet(d, keyBlockedReferrers, vl.BlockedReferrers, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := d.Set(keyBlockNoneReferrer, vl.BlockNoneReferrer); err != nil {
		return err
	}
	if err := d.Set(keyEnableMP4Fallback, vl.EnableMP4Fallback); err != nil {
		return err
	}
	if err := d.Set(keyKeepOriginalFiles, vl.KeepOriginalFiles); err != nil {
		return err
	}
	if err := d.Set(keyAllowDirectPlay, vl.AllowDirectPlay); err != nil {
		return err
	}
	if err := d.Set(keyEnableDRM, vl.EnableDRM); err != nil {
		return err
	}
	if err := d.Set(keyBitrate240p, vl.Bitrate240p); err != nil {
		return err
	}
	if err := d.Set(keyBitrate360p, vl.Bitrate360p); err != nil {
		return err
	}
	if err := d.Set(keyBitrate480p, vl.Bitrate480p); err != nil {
		return err
	}
	if err := d.Set(keyBitrate720p, vl.Bitrate720p); err != nil {
		return err
	}
	if err := d.Set(keyBitrate1080p, vl.Bitrate1080p); err != nil {
		return err
	}
	if err := d.Set(keyBitrate1440p, vl.Bitrate1440p); err != nil {
		return err
	}
	if err := d.Set(keyBitrate2160p, vl.Bitrate2160p); err != nil {
		return err
	}
	if err := d.Set(keyAPIAccessKey, vl.APIAccessKey); err != nil {
		return err
	}
	if err := d.Set(keyShowHeatmap, vl.ShowHeatmap); err != nil {
		return err
	}
	if err := d.Set(keyEnableContentTagging, vl.EnableContentTagging); err != nil {
		return err
	}
	if err := d.Set(keyPullZoneType, vl.PullZoneType); err != nil {
		return err
	}
	if err := d.Set(keyCustomHTML, vl.CustomHTML); err != nil {
		return err
	}
	if err := d.Set(keyControls, vl.Controls); err != nil {
		return err
	}
	if err := d.Set(keyPlayerKeyColor, vl.PlayerKeyColor); err != nil {
		return err
	}
	if err := d.Set(keyFontFamily, vl.FontFamily); err != nil {
		return err
	}

	return nil
}

// videoLibraryFromResource returns a VideoLibraryUpdateOptions API type that
// has fields set to the values in d.
func videoLibraryFromResource(d *schema.ResourceData) (*bunny.VideoLibraryUpdateOptions, error) {
	var res bunny.VideoLibraryUpdateOptions
	res.Name = getStrPtr(d, keyName)
	res.CustomHTML = getStrPtr(d, keyCustomHTML)
	res.PlayerKeyColor = getStrPtr(d, keyPlayerKeyColor)
	res.EnableTokenAuthentication = getBoolPtr(d, keyEnableTokenAuthentication)
	res.EnableTokenIPVerification = getBoolPtr(d, keyEnableTokenIPVerification)
	res.ResetToken = getBoolPtr(d, keyResetToken)
	res.WatermarkPositionLeft = getInt32Ptr(d, keyWatermarkPositionLeft)
	res.WatermarkPositionTop = getInt32Ptr(d, keyWatermarkPositionTop)
	res.WatermarkWidth = getInt32Ptr(d, keyWatermarkWidth)
	res.WatermarkHeight = getInt32Ptr(d, keyWatermarkHeight)
	res.EnabledResolutions = getStrPtr(d, keyEnabledResolutions)
	res.ViAiPublisherID = getStrPtr(d, keyViAiPublisherID)
	res.VastTagURL = getStrPtr(d, keyVastTagURL)
	res.WebhookURL = getStrPtr(d, keyWebhookURL)
	res.CaptionsFontSize = getInt32Ptr(d, keyCaptionsFontSize)
	res.CaptionsFontColor = getStrPtr(d, keyCaptionsFontColor)
	res.CaptionsBackground = getStrPtr(d, keyCaptionsBackground)
	res.UILanguage = getStrPtr(d, keyUILanguage)
	res.AllowEarlyPlay = getBoolPtr(d, keyAllowEarlyPlay)
	res.PlayerTokenAuthenticationEnabled = getBoolPtr(d, keyPlayerTokenAuthenticationEnabled)
	res.BlockNoneReferrer = getBoolPtr(d, keyBlockNoneReferrer)
	res.EnableMP4Fallback = getBoolPtr(d, keyEnableMP4Fallback)
	res.KeepOriginalFiles = getBoolPtr(d, keyKeepOriginalFiles)
	res.AllowDirectPlay = getBoolPtr(d, keyAllowDirectPlay)
	res.EnableDRM = getBoolPtr(d, keyEnableDRM)
	res.Controls = getStrPtr(d, keyControls)
	res.Bitrate240p = getInt32Ptr(d, keyBitrate240p)
	res.Bitrate360p = getInt32Ptr(d, keyBitrate360p)
	res.Bitrate480p = getInt32Ptr(d, keyBitrate480p)
	res.Bitrate720p = getInt32Ptr(d, keyBitrate720p)
	res.Bitrate1080p = getInt32Ptr(d, keyBitrate1080p)
	res.Bitrate1440p = getInt32Ptr(d, keyBitrate1440p)
	res.Bitrate2160p = getInt32Ptr(d, keyBitrate2160p)
	res.ShowHeatmap = getBoolPtr(d, keyShowHeatmap)
	res.EnableContentTagging = getBoolPtr(d, keyEnableContentTagging)
	res.FontFamily = getStrPtr(d, keyFontFamily)

	return &res, nil
}
