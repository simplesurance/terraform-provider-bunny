package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyAWSSigningEnabled                     = "aws_signing_enabled"
	keyAWSSigningKey                         = "aws_signing_key"
	keyAWSSigningRegionName                  = "aws_signing_region_name"
	keyAWSSigningSecret                      = "aws_signing_secret"
	keyAllowedReferrers                      = "allowed_referrers"
	keyBlockPostRequests                     = "block_post_requests"
	keyBlockRootPathAccess                   = "block_root_path_access"
	keyBlockedCountries                      = "blocked_countries"
	keyBlockedIPs                            = "blocked_ips"
	keyBudgetRedirectedCountries             = "budget_redirected_countries"
	keyCacheControlBrowserMaxAgeOverride     = "cache_control_browser_max_age_override"
	keyCacheControlMaxAgeOverride            = "cache_control_max_age_override"
	keyCacheErrorResponses                   = "cache_error_responses"
	keyConnectionLimitPerIPCount             = "connection_limit_per_ip_count"
	keyDisableCookies                        = "disable_cookies"
	keyEnableAvifVary                        = "enable_avif_vary"
	keyEnableCacheSlice                      = "enable_cache_slice"
	keyEnableCountryCodeVary                 = "enable_country_code_vary"
	keyEnableGeoZoneAF                       = "enable_geo_zone_af"
	keyEnableGeoZoneAsia                     = "enable_geo_zone_asia"
	keyEnableGeoZoneEU                       = "enable_geo_zone_eu"
	keyEnableGeoZoneSA                       = "enable_geo_zone_sa"
	keyEnableGeoZoneUS                       = "enable_geo_zone_us"
	keyEnableHostnameVary                    = "enable_hostname_vary"
	keyCnameDomain                           = "cname_domain"
	keyEnableLogging                         = "enable_logging"
	keyEnableMobileVary                      = "enable_mobile_vary"
	keyEnableOriginShield                    = "enable_origin_shield"
	keyEnableTLS1                            = "enable_tlsv1"
	keyEnableTLS11                           = "enable_tls1_1"
	keyEnableWebPVary                        = "enable_webp_vary"
	keyErrorPageCustomCode                   = "error_page_custom_code"
	keyErrorPageEnableCustomCode             = "error_page_enable_custom_code"
	keyErrorPageEnableStatuspageWidget       = "error_page_enable_statuspage_widget"
	keyErrorPageStatuspageCode               = "error_page_statuspage_code"
	keyErrorPageWhitelabel                   = "error_page_whitelabel"
	keyFollowRedirects                       = "follow_redirects"
	keyVideoLibraryID                        = "video_library_id"
	keyIgnoreQueryStrings                    = "ignore_query_strings"
	keyLogForwardingEnabled                  = "log_forwarding_enabled"
	keyLogForwardingHostname                 = "log_forwarding_hostname"
	keyLogForwardingPort                     = "log_forwarding_port"
	keyLogForwardingToken                    = "log_forwarding_token"
	keyLoggingIPAnonymizationEnabled         = "logging_ip_anonymization_enabled"
	keyLoggingSaveToStorage                  = "logging_save_to_storage"
	keyLoggingStorageZoneID                  = "logging_storage_zone_id"
	keyMonthlyBandwidthLimit                 = "monthly_bandwidth_limit"
	keyOptimizerAutomaticOptimizationEnabled = "optimizer_automatic_optimization_enabled"
	keyOptimizerDesktopMaxWidth              = "optimizer_desktop_max_width"
	keyOptimizerEnableManipulationEngine     = "optimizier_enable_manipulation_engine"
	keyOptimizerEnableWebP                   = "optimizer_enable_webp"
	keyOptimizerEnabled                      = "optimizer_enabled"
	keyOptimizerImageQuality                 = "optimizer_image_quality"
	keyOptimizerMinifyCSS                    = "optimizer_minify_css"
	keyOptimizerMinifyJavaScript             = "optimizer_minify_javascript"
	keyOptimizerMobileImageQuality           = "optimizer_mobile_image_quality"
	keyOptimizerMobileMaxWidth               = "optimizer_mobile_max_width"
	keyOptimizerWatermarkEnabled             = "optimizer_watermark_enabled"
	keyOptimizerWatermarkMinImageSize        = "optimizer_watermark_min_image_size"
	keyOptimizerWatermarkOffset              = "optimizer_watermark_offset"
	keyOptimizerWatermarkPosition            = "optimizer_watermark_position"
	keyOptimizerWatermarkURL                 = "optimizer_watermark_url"
	keyOriginShieldZoneCode                  = "origin_shield_zone_code"
	keyOriginURL                             = "origin_url"
	keyEnabled                               = "enabled"
	keyPermaCacheStorageZoneID               = "perma_cache_storage_zone_id"
	keyRequestLimit                          = "request_limit"
	keyType                                  = "type"
	keyVerifyOriginSSL                       = "verify_origin_ssl"
	keyZoneSecurityEnabled                   = "zone_security_enabled"
	keyZoneSecurityIncludeHashRemoteIP       = "zone_security_include_hash_remote_ip"

	keyBlockedReferrers = "blocked_referrers" // uses different API
	keyName             = "name"
	keyStorageZoneID    = "storage_zone_id"
	keyZoneSecurityKey  = "zone_security_key"

	keyLastUpdated = "last_updated"

	keySafeHop = "safehop"
	keyHeaders = "headers"
)

func resourcePullZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePullZoneCreate,
		ReadContext:   resourcePullZoneRead,
		UpdateContext: resourcePullZoneUpdate,
		DeleteContext: resourcePullZoneDelete,

		Schema: map[string]*schema.Schema{
			keyAWSSigningEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if the AWS signing should be enabled or not.",
				Default:     false,
				Optional:    true,
			},
			keyAWSSigningKey: {
				Type:        schema.TypeString,
				Description: "AWS Signing Key",
				Optional:    true,
			},
			keyAWSSigningRegionName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			keyAWSSigningSecret: {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			keyAllowedReferrers: {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Sets the list of referrer hostnames that are allowed to access the Pull Zone. Requests containing the header Referer: hostname that is not on the list will be rejected. If empty, all the referrers are allowed.",
				Optional:    true,
			},
			keyBlockPostRequests: {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			keyBlockRootPathAccess: {
				Type:        schema.TypeBool,
				Default:     false,
				Description: "Determines if the zone should block requests to the root of the zone.",
				Optional:    true,
			},
			keyBlockedCountries: {
				Type:        schema.TypeSet,
				Description: "Sets the list of two letter Alpha2 country codes that will be blocked from accessing the zone.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			keyBlockedIPs: {
				Type:        schema.TypeSet,
				Description: "Sets the list of IPs that are blocked from accessing the Pull Zone. Requests coming from the following IPs will be rejected. If empty, all the IPs will be allowed.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			keyBudgetRedirectedCountries: {
				Type:        schema.TypeSet,
				Description: "Sets the list of two letter Alpha2 country codes that will be redirected to the cheapest possible region.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			keyCacheControlBrowserMaxAgeOverride: {
				Type:        schema.TypeInt,
				Description: "Sets the browser cache control override setting for this zone.",
				Optional:    true,
				Default:     -1,
			},
			keyCacheControlMaxAgeOverride: {
				Type:        schema.TypeInt,
				Description: "Sets the cache control override setting for this zone.",
				Optional:    true,
				Default:     -1,
			},
			keyCacheErrorResponses: {
				Type:        schema.TypeBool,
				Description: "If enabled, bunny.net will temporarily cache error responses (304+ HTTP status codes) from your servers for 5 seconds to prevent DDoS attacks on your origin.\nIf disabled, error responses will be set to no-cache.",
				Optional:    true,
				Default:     false,
			},
			keyConnectionLimitPerIPCount: {
				Type:             schema.TypeInt,
				Description:      "Determines the maximum number of connections per IP that will be allowed to connect to this Pull Zone.",
				Optional:         true,
				ValidateDiagFunc: validateIsInt32,
			},
			keyDisableCookies: {
				Type:        schema.TypeBool,
				Description: "Determines if the Pull Zone should automatically remove cookies from the responses.",
				Optional:    true,
				Default:     true,
			},
			keyEnableAvifVary: {
				Type:        schema.TypeBool,
				Description: "Determines if the AVIF Vary feature should be enabled..",
				Default:     false,
				Optional:    true,
			},
			keyEnableCacheSlice: {
				Type:        schema.TypeBool,
				Description: "Determines if cache slicing (Optimize for video) should be enabled for this zone.",
				Default:     false,
				Optional:    true,
			},
			keyEnableCountryCodeVary: {
				Type:        schema.TypeBool,
				Description: "Determines if the Country Code Vary feature should be enabled.",
				Default:     false,
				Optional:    true,
			},
			keyEnableGeoZoneAF: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Serve data from the Middle East & Africa Zone.",
			},
			keyEnableGeoZoneAsia: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Serve data from the Asia & Oceania Zone.",
			},
			keyEnableGeoZoneEU: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Serve data from the Europe Zone.",
			},
			keyEnableGeoZoneSA: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Serve data from the South America Zone.",
			},
			keyEnableGeoZoneUS: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Serve data from the US Zone.",
			},
			keyEnableHostnameVary: {
				Type:        schema.TypeBool,
				Description: "Determines if the Hostname Vary feature should be enabled.",
				Default:     false,
				Optional:    true,
			},
			keyCnameDomain: {
				Type:        schema.TypeString,
				Description: "The CNAME domain of the Pull Zone for setting up custom hostnames.",
				Computed:    true,
			},

			keyEnableLogging: {
				Type:        schema.TypeBool,
				Description: "Determines if the logging should be enabled for this zone.",
				Default:     true,
				Optional:    true,
			},
			keyEnableMobileVary: {
				Type:        schema.TypeBool,
				Description: "Determines if the Mobile Vary feature is enabled.",
				Default:     false,
				Optional:    true,
			},
			keyEnableOriginShield: {
				Type:        schema.TypeBool,
				Description: "Determines if the origin shield should be enabled.",
				Default:     false,
				Optional:    true,
			},
			keyEnableTLS1: {
				Type:        schema.TypeBool,
				Description: "Determines if the TLS 1 should be enabled on this zone.",
				Default:     true,
				Optional:    true,
			},
			keyEnableTLS11: {
				Type:        schema.TypeBool,
				Description: "Determines if the TLS 1.1 should be enabled on this zone.",
				Default:     true,
				Optional:    true,
			},
			keyEnableWebPVary: {
				Type:        schema.TypeBool,
				Description: "Determines if the WebP Vary feature should be enabled.",
				Default:     false,
				Optional:    true,
			},
			keyErrorPageCustomCode: {
				Type:        schema.TypeString,
				Description: "Contains the custom error page code that will be returned",
				Optional:    true,
			},
			keyErrorPageEnableCustomCode: {
				Type:        schema.TypeBool,
				Description: "Determines if custom error page code should be enabled.",
				Default:     false,
				Optional:    true,
			},
			keyErrorPageEnableStatuspageWidget: {
				Type:        schema.TypeBool,
				Description: "Determines if the statuspage widget should be displayed on the error pages.",
				Default:     false,
				Optional:    true,
			},
			keyErrorPageStatuspageCode: {
				Type:        schema.TypeString,
				Description: "The statuspage code that will be used to build the status widget.",
				Optional:    true,
			},
			keyErrorPageWhitelabel: {
				Type:        schema.TypeBool,
				Description: "Determines if the error pages should be whitelabel or not.",
				Default:     false,
				Optional:    true,
			},
			keyFollowRedirects: {
				Type:        schema.TypeBool,
				Description: "Determines if the zone should follow redirects return by the oprigin and cache the response.",
				Default:     false,
				Optional:    true,
			},
			keyVideoLibraryID: {
				Type:        schema.TypeInt,
				Description: "The ID of the video library that the zone is linked to.",
				Computed:    true,
			},
			keyIgnoreQueryStrings: {
				Type:        schema.TypeBool,
				Description: "Determines if the Pull Zone should ignore query strings when serving cached objects (Vary by Query String).",
				Default:     true,
				Optional:    true,
			},
			keyLogForwardingEnabled: {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			keyLogForwardingHostname: {
				Type:        schema.TypeString,
				Description: "Sets the log forwarding destination hostname for the zone.",
				Optional:    true,
			},
			keyLogForwardingPort: {
				Type:             schema.TypeInt,
				Description:      "Sets the log forwarding port for the zone.",
				Default:          0,
				Optional:         true,
				ValidateDiagFunc: validateIsInt32,
			},
			keyLogForwardingToken: {
				Type:        schema.TypeString,
				Description: "Sets the log forwarding token for the zone.",
				Sensitive:   true,
				Optional:    true,
			},
			keyLoggingIPAnonymizationEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if the log anonoymization should be enabled. The field can only be set if the DPA agreement was set in the webinterface.",
				Default:     true,
				Optional:    true,
			},
			keyLoggingSaveToStorage: {
				Type:         schema.TypeBool,
				Description:  "Determines if the logging permanent storage should be enabled.",
				Default:      false,
				Optional:     true,
				RequiredWith: []string{keyLoggingStorageZoneID},
			},
			keyLoggingStorageZoneID: {
				Type:        schema.TypeInt,
				Description: "Sets the Storage Zone id that should contain the logs from this Pull Zone.",
				Default:     0,
				Optional:    true,
			},
			keyMonthlyBandwidthLimit: {
				Type:        schema.TypeInt,
				Description: "Sets the monthly limit of bandwidth in bytes that the pullzone is allowed to use.",
				Default:     0,
				Optional:    true,
			},
			keyOptimizerAutomaticOptimizationEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if the automatic image optimization should be enabled.",
				Optional:    true,
				Default:     true,
			},
			keyOptimizerDesktopMaxWidth: {
				Type:             schema.TypeInt,
				Description:      "Determines if the automatic image optimization should be enabled.",
				Optional:         true,
				Default:          1600,
				ValidateDiagFunc: validateIsInt32,
			},
			keyOptimizerEnableManipulationEngine: {
				Type:        schema.TypeBool,
				Description: "Determines if the image manipulation should be enabled.",
				Optional:    true,
				Default:     true,
			},
			keyOptimizerEnableWebP: {
				Type:        schema.TypeBool,
				Description: "Determines if the WebP optimization should be enabled.",
				Default:     true,
				Optional:    true,
			},
			keyOptimizerEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if the optimizer should be enabled for this zone.",
				Default:     false,
				Optional:    true,
			},
			keyOptimizerImageQuality: {
				Type:             schema.TypeInt,
				Description:      "Determines the image quality for desktop clients.",
				Optional:         true,
				Default:          85,
				ValidateDiagFunc: validateIsInt32,
			},
			keyOptimizerMinifyCSS: {
				Type:        schema.TypeBool,
				Description: "Determines if the CSS minifcation should be enabled.",
				Default:     true,
				Optional:    true,
			},
			keyOptimizerMinifyJavaScript: {
				Type:        schema.TypeBool,
				Description: "Determines if the JavaScript minifcation should be enabled.",
				Default:     true,
				Optional:    true,
			},
			keyOptimizerMobileImageQuality: {
				Type:             schema.TypeInt,
				Description:      "Determines the image quality for mobile clients.",
				Optional:         true,
				Default:          70,
				ValidateDiagFunc: validateIsInt32,
			},
			keyOptimizerMobileMaxWidth: {
				Type:             schema.TypeInt,
				Description:      "Determines the maximum automatic image size for mobile clients.",
				Optional:         true,
				Default:          800,
				ValidateDiagFunc: validateIsInt32,
			},
			keyOptimizerWatermarkEnabled: {
				Type:        schema.TypeBool,
				Description: "Determines if image watermarking should be enabled.",
				Default:     false,
				Optional:    true,
			},
			keyOptimizerWatermarkMinImageSize: {
				Type:             schema.TypeInt,
				Description:      "Sets the minimum image size to which the watermark will be added.",
				Optional:         true,
				Default:          300,
				ValidateDiagFunc: validateIsInt32,
			},
			keyOptimizerWatermarkOffset: {
				Type:        schema.TypeFloat,
				Description: "Sets the offset of the watermark image.",
				Optional:    true,
				Default:     3,
			},
			keyOptimizerWatermarkPosition: {
				Type:        schema.TypeInt,
				Description: "Sets the position of the watermark image.",
				Optional:    true,
				Default:     0,
			},
			keyOptimizerWatermarkURL: {
				Type:        schema.TypeString,
				Description: "Sets the URL of the watermark image.",
				Optional:    true,
			},
			keyOriginShieldZoneCode: {
				Type:        schema.TypeString,
				Description: "Determines the zone code where the origin shield should be set up.",
				Optional:    true,
				Default:     "FR",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice([]string{"FR", "IL"}, false),
				),
			},
			keyOriginURL: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The origin URL of the Pull Zone where the files are fetched from.",
			},
			keyPermaCacheStorageZoneID: {
				Type:        schema.TypeInt,
				Description: "The ID of the storage zone that should be used as the Perma-Cache.",
				Default:     0,
				Optional:    true,
			},
			keyRequestLimit: {
				Type:        schema.TypeInt,
				Description: "Determines the maximum number of requests per second that will be allowed to connect to this Pull Zone.",
				Default:     0,
				Optional:    true,
			},
			keySafeHop: {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Elem:             resourcePullZoneSafeHop,
				DiffSuppressFunc: diffSupressMissingOptionalBlock,
			},
			keyHeaders: {
				Type:             schema.TypeList,
				MaxItems:         1,
				Optional:         true,
				Elem:             resourcePullZoneHeaders,
				DiffSuppressFunc: diffSupressMissingOptionalBlock,
			},
			keyType: {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 1)),
				Description:      "The type of the Pull Zone. Standard = 0, Volume = 1.",
			},
			keyVerifyOriginSSL: {
				Type:        schema.TypeBool,
				Description: "Determines if the SSL certificate should be verified when connecting to the origin.",
				Default:     false,
				Optional:    true,
			},
			keyZoneSecurityEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			keyZoneSecurityIncludeHashRemoteIP: {
				Type:     schema.TypeBool,
				Optional: true,
			},

			keyEnabled: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			keyName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Pull Zone.",
			},
			keyStorageZoneID: {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the storage zone that the Pull Zone is linked to.",
			},
			keyZoneSecurityKey: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			keyBlockedReferrers: {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true, // must be updated via Add/Remove Blocked Referrer API Endpoint, not implemented
				ForceNew:    true,
				Description: "The list of hostnames that will be blocked from accessing the Pull Zone.",
			},

			keyLastUpdated: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePullZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	pz, err := clt.PullZone.Add(ctx, &bunny.PullZoneAddOptions{
		Name:          d.Get(keyName).(string),
		OriginURL:     d.Get(keyOriginURL).(string),
		StorageZoneID: getInt64Ptr(d, keyStorageZoneID),
		Type:          d.Get(keyType).(int),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating pull zone failed: %w", err))
	}

	d.SetId(strconv.FormatInt(*pz.ID, 10))
	if err := d.Set(keyLastUpdated, time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}

	// PullZone.Add() only supports to set a subset of a Pull Zone object,
	// call Update to set the remaining ones.
	if diags := resourcePullZoneUpdate(ctx, d, meta); diags.HasError() {
		// if updating fails the pz was still created, initialize with the PZ
		// returned from the Add operation
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "setting pull zone attributes via update failed",
		})

		if err := pullZoneToResourceData(pz, d); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "converting api-type to resource data failed: " + err.Error(),
			})

		}

		return diags
	}

	return nil
}

func resourcePullZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	pullZone, err := resourceDataToPullZoneUpdate(d)
	if err != nil {
		return diagsErrFromErr("converting resource to API type failed", err)
	}

	id, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedPullZone, err := clt.PullZone.Update(ctx, id, pullZone)
	if err != nil {
		return diagsErrFromErr("updating pull zone via API failed", err)
	}

	if err := pullZoneToResourceData(updatedPullZone, d); err != nil {
		return diagsErrFromErr("converting api type to resource data after successful update failed: %w", err)
	}

	if err := d.Set(keyLastUpdated, time.Now().Format(time.RFC850)); err != nil {
		return diagsWarnFromErr(
			fmt.Sprintf("could not set %s", keyLastUpdated),
			err,
		)
	}

	return nil
}

func resourcePullZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	id, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	pz, err := clt.PullZone.Get(ctx, id)
	if err != nil {
		return diagsErrFromErr("could not retrieve pull zone", err)
	}

	if err := pullZoneToResourceData(pz, d); err != nil {
		return diagsErrFromErr("converting api type to resource data after successful read failed: %w", err)
	}

	return nil
}

func resourcePullZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	id, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = clt.PullZone.Delete(ctx, id)
	if err != nil {
		return diagsErrFromErr("could not delete pull zone", err)
	}

	d.SetId("")

	return nil
}

// pullZoneToResourceData sets fields in d to the values in pz.
func pullZoneToResourceData(pz *bunny.PullZone, d *schema.ResourceData) error {
	if pz.ID != nil {
		d.SetId(strconv.FormatInt(*pz.ID, 10))
	}

	if err := d.Set(keyAWSSigningEnabled, pz.AWSSigningEnabled); err != nil {
		return err
	}
	if err := d.Set(keyAWSSigningKey, pz.AWSSigningKey); err != nil {
		return err
	}
	if err := d.Set(keyAWSSigningRegionName, pz.AWSSigningRegionName); err != nil {
		return err
	}
	if err := d.Set(keyAWSSigningSecret, pz.AWSSigningSecret); err != nil {
		return err
	}
	if err := setStrSet(d, keyAllowedReferrers, pz.AllowedReferrers, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := d.Set(keyBlockPostRequests, pz.BlockPostRequests); err != nil {
		return err
	}
	if err := d.Set(keyBlockRootPathAccess, pz.BlockRootPathAccess); err != nil {
		return err
	}
	if err := setStrSet(d, keyBlockedCountries, pz.BlockedCountries, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := setStrSet(d, keyBlockedIPs, pz.BlockedIPs, ignoreOrderOpt); err != nil {
		return err
	}
	if err := setStrSet(d, keyBudgetRedirectedCountries, pz.BudgetRedirectedCountries, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := d.Set(keyCacheControlBrowserMaxAgeOverride, pz.CacheControlBrowserMaxAgeOverride); err != nil {
		return err
	}
	if err := d.Set(keyCacheControlMaxAgeOverride, pz.CacheControlMaxAgeOverride); err != nil {
		return err
	}
	if err := d.Set(keyCacheErrorResponses, pz.CacheErrorResponses); err != nil {
		return err
	}
	if err := d.Set(keyConnectionLimitPerIPCount, pz.ConnectionLimitPerIPCount); err != nil {
		return err
	}
	if err := d.Set(keyDisableCookies, pz.DisableCookies); err != nil {
		return err
	}
	if err := d.Set(keyEnableAvifVary, pz.EnableAvifVary); err != nil {
		return err
	}
	if err := d.Set(keyEnableCacheSlice, pz.EnableCacheSlice); err != nil {
		return err
	}
	if err := d.Set(keyEnableCountryCodeVary, pz.EnableCountryCodeVary); err != nil {
		return err
	}
	if err := d.Set(keyEnableGeoZoneAF, pz.EnableGeoZoneAF); err != nil {
		return err
	}
	if err := d.Set(keyEnableGeoZoneAsia, pz.EnableGeoZoneAsia); err != nil {
		return err
	}
	if err := d.Set(keyEnableGeoZoneEU, pz.EnableGeoZoneEU); err != nil {
		return err
	}
	if err := d.Set(keyEnableGeoZoneSA, pz.EnableGeoZoneSA); err != nil {
		return err
	}
	if err := d.Set(keyEnableGeoZoneUS, pz.EnableGeoZoneUS); err != nil {
		return err
	}
	if err := d.Set(keyEnableHostnameVary, pz.EnableHostnameVary); err != nil {
		return err
	}
	if err := d.Set(keyCnameDomain, pz.CnameDomain); err != nil {
		return err
	}
	if err := d.Set(keyEnableLogging, pz.EnableLogging); err != nil {
		return err
	}
	if err := d.Set(keyEnableMobileVary, pz.EnableMobileVary); err != nil {
		return err
	}
	if err := d.Set(keyEnableOriginShield, pz.EnableOriginShield); err != nil {
		return err
	}
	if err := d.Set(keyEnableTLS1, pz.EnableTLS1); err != nil {
		return err
	}
	if err := d.Set(keyEnableTLS11, pz.EnableTLS11); err != nil {
		return err
	}
	if err := d.Set(keyEnableWebPVary, pz.EnableWebPVary); err != nil {
		return err
	}
	if err := d.Set(keyErrorPageCustomCode, pz.ErrorPageCustomCode); err != nil {
		return err
	}
	if err := d.Set(keyErrorPageEnableCustomCode, pz.ErrorPageEnableCustomCode); err != nil {
		return err
	}
	if err := d.Set(keyErrorPageEnableStatuspageWidget, pz.ErrorPageEnableStatuspageWidget); err != nil {
		return err
	}
	if err := d.Set(keyErrorPageStatuspageCode, pz.ErrorPageStatuspageCode); err != nil {
		return err
	}
	if err := d.Set(keyErrorPageWhitelabel, pz.ErrorPageWhitelabel); err != nil {
		return err
	}
	if err := d.Set(keyFollowRedirects, pz.FollowRedirects); err != nil {
		return err
	}
	if err := d.Set(keyVideoLibraryID, pz.VideoLibraryID); err != nil {
		return err
	}
	if err := d.Set(keyIgnoreQueryStrings, pz.IgnoreQueryStrings); err != nil {
		return err
	}
	if err := d.Set(keyLogForwardingEnabled, pz.LogForwardingEnabled); err != nil {
		return err
	}
	if err := d.Set(keyLogForwardingHostname, pz.LogForwardingHostname); err != nil {
		return err
	}
	if err := d.Set(keyLogForwardingPort, pz.LogForwardingPort); err != nil {
		return err
	}
	if err := d.Set(keyLogForwardingPort, pz.LogForwardingPort); err != nil {
		return err
	}
	if err := d.Set(keyLogForwardingToken, pz.LogForwardingToken); err != nil {
		return err
	}
	if err := d.Set(keyLoggingIPAnonymizationEnabled, pz.LoggingIPAnonymizationEnabled); err != nil {
		return err
	}
	if err := d.Set(keyLoggingSaveToStorage, pz.LoggingSaveToStorage); err != nil {
		return err
	}
	if err := d.Set(keyLoggingStorageZoneID, pz.LoggingStorageZoneID); err != nil {
		return err
	}
	if err := d.Set(keyMonthlyBandwidthLimit, pz.MonthlyBandwidthLimit); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerAutomaticOptimizationEnabled, pz.OptimizerAutomaticOptimizationEnabled); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerDesktopMaxWidth, pz.OptimizerDesktopMaxWidth); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerEnableManipulationEngine, pz.OptimizerEnableManipulationEngine); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerEnableWebP, pz.OptimizerEnableWebP); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerEnabled, pz.OptimizerEnabled); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerImageQuality, pz.OptimizerImageQuality); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerMinifyCSS, pz.OptimizerMinifyCSS); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerMinifyJavaScript, pz.OptimizerMinifyJavaScript); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerMobileImageQuality, pz.OptimizerMobileImageQuality); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerMobileMaxWidth, pz.OptimizerMobileMaxWidth); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerWatermarkEnabled, pz.OptimizerWatermarkEnabled); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerWatermarkMinImageSize, pz.OptimizerWatermarkMinImageSize); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerWatermarkOffset, pz.OptimizerWatermarkOffset); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerWatermarkPosition, pz.OptimizerWatermarkPosition); err != nil {
		return err
	}
	if err := d.Set(keyOptimizerWatermarkURL, pz.OptimizerWatermarkURL); err != nil {
		return err
	}
	if err := d.Set(keyOriginShieldZoneCode, pz.OriginShieldZoneCode); err != nil {
		return err
	}
	if err := d.Set(keyOriginURL, pz.OriginURL); err != nil {
		return err
	}
	if err := d.Set(keyPermaCacheStorageZoneID, pz.PermaCacheStorageZoneID); err != nil {
		return err
	}
	if err := d.Set(keyRequestLimit, pz.RequestLimit); err != nil {
		return err
	}
	if err := d.Set(keyType, pz.Type); err != nil {
		return err
	}
	if err := d.Set(keyVerifyOriginSSL, pz.VerifyOriginSSL); err != nil {
		return err
	}
	if err := d.Set(keyZoneSecurityEnabled, pz.ZoneSecurityEnabled); err != nil {
		return err
	}
	if err := d.Set(keyZoneSecurityIncludeHashRemoteIP, pz.ZoneSecurityIncludeHashRemoteIP); err != nil {
		return err
	}
	if err := setStrSet(d, keyBlockedReferrers, pz.BlockedReferrers, ignoreOrderOpt, caseInsensitiveOpt); err != nil {
		return err
	}
	if err := d.Set(keyEnabled, pz.Enabled); err != nil {
		return err
	}
	if err := d.Set(keyName, pz.Name); err != nil {
		return err
	}
	if err := d.Set(keyZoneSecurityKey, pz.ZoneSecurityKey); err != nil {
		return err
	}

	if err := safeHopToResource(pz, d); err != nil {
		return err
	}

	if err := headersToResource(pz, d); err != nil {
		return err
	}

	return nil
}

// resourceDataToPullZoneUpdate returns a PullZoneUpdateOptions API type that
// has fields set to the values in d.
func resourceDataToPullZoneUpdate(d *schema.ResourceData) (*bunny.PullZoneUpdateOptions, error) {
	var res bunny.PullZoneUpdateOptions

	res.AWSSigningEnabled = getBoolPtr(d, keyAWSSigningEnabled)
	res.AWSSigningKey = getStrPtr(d, keyAWSSigningKey)
	res.AWSSigningRegionName = getStrPtr(d, keyAWSSigningRegionName)
	res.AWSSigningSecret = getStrPtr(d, keyAWSSigningSecret)
	res.AllowedReferrers = getStrSetAsSlice(d, keyAllowedReferrers)
	res.BlockPostRequests = getBoolPtr(d, keyBlockPostRequests)
	res.BlockRootPathAccess = getBoolPtr(d, keyBlockRootPathAccess)
	res.BlockedCountries = getStrSetAsSlice(d, keyBlockedCountries)
	res.BlockedIPs = getStrSetAsSlice(d, keyBlockedIPs)
	res.BudgetRedirectedCountries = getStrSetAsSlice(d, keyBudgetRedirectedCountries)
	res.CacheControlBrowserMaxAgeOverride = getInt64Ptr(d, keyCacheControlBrowserMaxAgeOverride)
	res.CacheControlMaxAgeOverride = getInt64Ptr(d, keyCacheControlMaxAgeOverride)
	res.CacheErrorResponses = getBoolPtr(d, keyCacheErrorResponses)
	res.ConnectionLimitPerIPCount = getInt32Ptr(d, keyConnectionLimitPerIPCount)
	res.DisableCookies = getBoolPtr(d, keyDisableCookies)
	res.EnableAvifVary = getBoolPtr(d, keyEnableAvifVary)
	res.EnableCacheSlice = getBoolPtr(d, keyEnableCacheSlice)
	res.EnableCountryCodeVary = getBoolPtr(d, keyEnableCountryCodeVary)
	res.EnableHostnameVary = getBoolPtr(d, keyEnableHostnameVary)
	res.EnableLogging = getBoolPtr(d, keyEnableLogging)
	res.EnableMobileVary = getBoolPtr(d, keyEnableMobileVary)
	res.EnableOriginShield = getBoolPtr(d, keyEnableOriginShield)
	res.EnableTLS1 = getBoolPtr(d, keyEnableTLS1)
	res.EnableTLS11 = getBoolPtr(d, keyEnableTLS11)
	res.EnableWebPVary = getBoolPtr(d, keyEnableWebPVary)
	res.ErrorPageCustomCode = getStrPtr(d, keyErrorPageCustomCode)
	res.ErrorPageEnableCustomCode = getBoolPtr(d, keyErrorPageEnableCustomCode)
	res.ErrorPageEnableStatuspageWidget = getBoolPtr(d, keyErrorPageEnableStatuspageWidget)
	res.ErrorPageStatuspageCode = getStrPtr(d, keyErrorPageStatuspageCode)
	res.ErrorPageWhitelabel = getBoolPtr(d, keyErrorPageWhitelabel)
	res.FollowRedirects = getBoolPtr(d, keyFollowRedirects)
	res.IgnoreQueryStrings = getBoolPtr(d, keyIgnoreQueryStrings)
	res.LogForwardingEnabled = getBoolPtr(d, keyLogForwardingEnabled)
	res.LogForwardingHostname = getStrPtr(d, keyLogForwardingHostname)
	res.LogForwardingPort = getInt32Ptr(d, keyLogForwardingPort)
	res.LogForwardingToken = getStrPtr(d, keyLogForwardingToken)
	res.LoggingIPAnonymizationEnabled = getBoolPtr(d, keyLoggingIPAnonymizationEnabled)
	res.LoggingSaveToStorage = getBoolPtr(d, keyLoggingSaveToStorage)
	res.LoggingStorageZoneID = getInt64Ptr(d, keyLoggingStorageZoneID)
	res.MonthlyBandwidthLimit = getInt64Ptr(d, keyMonthlyBandwidthLimit)
	res.OptimizerAutomaticOptimizationEnabled = getBoolPtr(d, keyOptimizerAutomaticOptimizationEnabled)
	res.OptimizerDesktopMaxWidth = getInt32Ptr(d, keyOptimizerDesktopMaxWidth)
	res.OptimizerEnableManipulationEngine = getBoolPtr(d, keyOptimizerEnableManipulationEngine)
	res.OptimizerEnableWebP = getBoolPtr(d, keyOptimizerEnableWebP)
	res.OptimizerEnabled = getBoolPtr(d, keyOptimizerEnabled)
	res.OptimizerImageQuality = getInt32Ptr(d, keyOptimizerImageQuality)
	res.OptimizerMinifyCSS = getBoolPtr(d, keyOptimizerMinifyCSS)
	res.OptimizerMinifyJavaScript = getBoolPtr(d, keyOptimizerMinifyJavaScript)
	res.OptimizerMobileImageQuality = getInt32Ptr(d, keyOptimizerMobileImageQuality)
	res.OptimizerMobileMaxWidth = getInt32Ptr(d, keyOptimizerMobileMaxWidth)
	res.OptimizerWatermarkEnabled = getBoolPtr(d, keyOptimizerWatermarkEnabled)
	res.OptimizerWatermarkMinImageSize = getInt32Ptr(d, keyOptimizerWatermarkMinImageSize)
	res.OptimizerWatermarkOffset = getFloat64Ptr(d, keyOptimizerWatermarkOffset)
	res.OptimizerWatermarkPosition = getIntPtr(d, keyOptimizerWatermarkPosition)
	res.OptimizerWatermarkURL = getStrPtr(d, keyOptimizerWatermarkURL)
	res.OriginShieldZoneCode = getStrPtr(d, keyOriginShieldZoneCode)
	res.OriginURL = getStrPtr(d, keyOriginURL)
	res.PermaCacheStorageZoneID = getInt64Ptr(d, keyPermaCacheStorageZoneID)
	res.RequestLimit = getInt32Ptr(d, keyRequestLimit)
	res.Type = getIntPtr(d, keyType)
	res.VerifyOriginSSL = getBoolPtr(d, keyVerifyOriginSSL)
	res.ZoneSecurityEnabled = getBoolPtr(d, keyZoneSecurityEnabled)
	res.ZoneSecurityIncludeHashRemoteIP = getBoolPtr(d, keyZoneSecurityIncludeHashRemoteIP)

	safehopPullZoneUpdateOptionsFromResource(&res, d)
	headersFromResource(&res, d)

	return &res, nil
}
