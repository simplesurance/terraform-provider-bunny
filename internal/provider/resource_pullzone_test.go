package provider

import (
	"context"
	"regexp"

	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	ptr "github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	bunny "github.com/simplesurance/bunny-go"
)

func randHostname() string {
	return resource.PrefixedUniqueId(resourcePrefix) + ".test"
}

type pullZoneWanted struct {
	TerraformResourceName string
	bunny.PullZone
	Name                string
	OriginURL           string
	EnableGeoZoneAsia   bool
	EnableGeoZoneEU     bool
	ZoneSecurityEnabled bool
	DisableCookies      bool
	EnableTLSV1         bool
	FollowRedirects     bool
}

func newAPIClient() *bunny.Client {
	return bunny.NewClient(
		os.Getenv(envVarAPIKey),
		bunny.WithUserAgent(userAgent+"-test"),
	)
}

func stringsAreEqual(a string, b *string) error {
	if b != nil && a == *b {
		return nil
	}

	if b == nil {
		return fmt.Errorf("%q != %v", a, b)
	}

	return fmt.Errorf("%q != %q", a, *b)
}

func boolsAreEqual(a bool, b *bool) error {
	if b != nil && a == *b {
		return nil
	}

	if b == nil {
		return fmt.Errorf("%t != %v", a, b)
	}

	return fmt.Errorf("%t != %t", a, *b)
}

func checkBasicPullZoneAPIState(wanted *pullZoneWanted) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		strID, err := idFromState(s, wanted.TerraformResourceName)
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

		if err := stringsAreEqual(wanted.Name, pz.Name); err != nil {
			return fmt.Errorf("name of created pullzone differs: %w", err)
		}

		if err := stringsAreEqual(wanted.OriginURL, pz.OriginURL); err != nil {
			return fmt.Errorf("originURL of created pullzone differs: %w", err)
		}

		if err := boolsAreEqual(wanted.EnableGeoZoneAsia, pz.EnableGeoZoneAsia); err != nil {
			return fmt.Errorf("EnableGeoZoneAsia of created pullzone differs: %w", err)
		}

		if err := boolsAreEqual(wanted.EnableGeoZoneEU, pz.EnableGeoZoneEU); err != nil {
			return fmt.Errorf("EnableGeoZoneEU of created pullzone differs: %w", err)
		}

		if err := boolsAreEqual(wanted.ZoneSecurityEnabled, pz.ZoneSecurityEnabled); err != nil {
			return fmt.Errorf("ZoneSecurityEnabled of created pullzone differs: %w", err)
		}

		if err := boolsAreEqual(wanted.DisableCookies, pz.DisableCookies); err != nil {
			return fmt.Errorf("DisableCookies of created pullzone differs: %w", err)
		}

		if err := boolsAreEqual(wanted.EnableTLSV1, pz.EnableTLS1); err != nil {
			return fmt.Errorf("EnableTLS1 of created pullzone differs: %w", err)
		}

		if err := boolsAreEqual(wanted.FollowRedirects, pz.FollowRedirects); err != nil {
			return fmt.Errorf("FollowRedirects of created pullzone differs: %w", err)
		}

		return nil
	}
}

func checkPullZoneNotExists(pullZoneName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		var page int32

		for {
			pullzones, err := clt.PullZone.List(context.Background(), &bunny.PaginationOptions{
				Page:    page,
				PerPage: 1000,
			})
			if err != nil {
				return fmt.Errorf("listing pullzones failed: %w", err)
			}

			for _, pz := range pullzones.Items {
				if pz.Name == nil {
					return fmt.Errorf("got pullzone from api with empty Name: %+v", pz)
				}

				if pullZoneName == *pz.Name {
					return &resource.UnexpectedStateError{
						State:         "exists",
						ExpectedState: []string{"not exists"},
					}

				}

				if !*pullzones.HasMoreItems {
					return nil
				}

				page++
			}
		}
	}
}

func TestAccPullZone_basic(t *testing.T) {
	/*
	   TODO:
	   - set a TypeList field and check if it was set correctly
	   - test updating keyBlockedReferrers, should recreate the pz
	   - only set required values in this test (https://github.com/hashicorp/terraform-provider-google/wiki/Developer-Best-Practices#acceptance-tests)

	*/
	attrs := pullZoneWanted{
		TerraformResourceName: "bunny_pullzone.mytest1",
		Name:                  randResourceName(),
		OriginURL:             "https://tabletennismap.de",
		EnableGeoZoneAsia:     true,
		EnableGeoZoneEU:       true,
		DisableCookies:        true,
		EnableTLSV1:           true,
		FollowRedirects:       false,
	}

	tf := fmt.Sprintf(`
resource "bunny_pullzone" "mytest1" {
	name = "%s"
	origin_url ="%s"
}
`,
		attrs.Name,
		attrs.OriginURL,
	)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check:  checkBasicPullZoneAPIState(&attrs),
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
		CheckDestroy: checkPullZoneNotExists(attrs.Name),
	})
}

func tfStrList(sl []string) string {
	var res string

	if len(sl) == 0 {
		return "[]"
	}

	res = "[ "
	for i, elem := range sl {
		res += fmt.Sprintf("%q", elem)
		if i < len(sl)-1 {
			res += ", "
		}
	}
	res += " ]"

	return res
}

func TestAccPullZone_full(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_pullzone." + resourceName

	// set fields to different values then their defaults, to be able to test if the settings are applied
	attrs := bunny.PullZone{
		AWSSigningEnabled:                 ptr.ToBool(true),
		AWSSigningKey:                     ptr.ToString("12345"),
		AWSSigningRegionName:              ptr.ToString("eu"),
		AWSSigningSecret:                  ptr.ToString("456"),
		AddCanonicalHeader:                ptr.ToBool(true),
		AddHostHeader:                     ptr.ToBool(true),
		AllowedReferrers:                  []string{"myhost", "localhost"},
		BlockPostRequests:                 ptr.ToBool(true),
		BlockRootPathAccess:               ptr.ToBool(true),
		BlockedCountries:                  []string{"KP", "US"},
		BlockedIPs:                        []string{"1.1.1.1", "127.0.0.1", "::1"},
		BudgetRedirectedCountries:         []string{"DE", "GB"},
		CacheControlBrowserMaxAgeOverride: ptr.ToInt64(100),
		CacheControlMaxAgeOverride:        ptr.ToInt64(3),
		CacheErrorResponses:               ptr.ToBool(true),
		ConnectionLimitPerIPCount:         ptr.ToInt32(23),
		DisableCookies:                    ptr.ToBool(false),
		EnableAccessControlOriginHeader:   ptr.ToBool(false),
		EnableAvifVary:                    ptr.ToBool(true),
		EnableCacheSlice:                  ptr.ToBool(true),
		EnableCountryCodeVary:             ptr.ToBool(true),
		EnableHostnameVary:                ptr.ToBool(true),
		EnableLogging:                     ptr.ToBool(false),
		EnableMobileVary:                  ptr.ToBool(true),
		EnableOriginShield:                ptr.ToBool(true),
		EnableTLS1:                        ptr.ToBool(false),
		EnableTLS11:                       ptr.ToBool(false),
		EnableWebPVary:                    ptr.ToBool(true),
		ErrorPageCustomCode:               ptr.ToString("error"),
		ErrorPageEnableCustomCode:         ptr.ToBool(true),
		ErrorPageEnableStatuspageWidget:   ptr.ToBool(true),
		ErrorPageStatuspageCode:           ptr.ToString("statuspage-error"),
		ErrorPageWhitelabel:               ptr.ToBool(true),
		FollowRedirects:                   ptr.ToBool(true),
		IgnoreQueryStrings:                ptr.ToBool(false),
		LogForwardingEnabled:              ptr.ToBool(true),
		LogForwardingHostname:             ptr.ToString("localhost"),
		LogForwardingPort:                 ptr.ToInt32(22),
		LogForwardingToken:                ptr.ToString("abcd"),
		LoggingIPAnonymizationEnabled:     ptr.ToBool(false),
		// TODO: can only be set if LoggingStorageZoneId is set to an existing storagezone
		//LoggingSaveToStorage:             ptr.ToBool(true),
		// TODO: Test LoggingStorageZoneId
		MonthlyBandwidthLimit: ptr.ToInt64(10240),
		OriginShieldZoneCode:  ptr.ToString("IL"),
		OriginURL:             ptr.ToString("http://terraform.io"),
		// TODO: Test PermaCacheStorageZoneID
		RequestLimit:                    ptr.ToInt32(3),
		Type:                            ptr.ToInt(1),
		VerifyOriginSSL:                 ptr.ToBool(true),
		ZoneSecurityEnabled:             ptr.ToBool(true),
		ZoneSecurityIncludeHashRemoteIP: ptr.ToBool(false),
		Name:                            ptr.ToString(randResourceName()),
		// TODO: Test StorageZoneID
		ZoneSecurityKey: ptr.ToString("xyz"),

		EnableSafeHop:                       ptr.ToBool(true),
		AccessControlOriginHeaderExtensions: []string{"txt", "exe", "json"},
		OriginConnectTimeout:                ptr.ToInt32(3),
		OriginResponseTimeout:               ptr.ToInt32(45),
		OriginRetries:                       ptr.ToInt32(2),
		OriginRetry5xxResponses:             ptr.ToBool(true),
		OriginRetryConnectionTimeout:        ptr.ToBool(false),
		OriginRetryDelay:                    ptr.ToInt32(3),
		OriginRetryResponseTimeout:          ptr.ToBool(false),

		OptimizerAutomaticOptimizationEnabled: ptr.ToBool(true),
		OptimizerDesktopMaxWidth:              ptr.ToInt32(1024),
		OptimizerEnableManipulationEngine:     ptr.ToBool(false),
		OptimizerEnableWebP:                   ptr.ToBool(false),
		OptimizerEnabled:                      ptr.ToBool(true),
		OptimizerImageQuality:                 ptr.ToInt32(81),
		OptimizerMinifyCSS:                    ptr.ToBool(false),
		OptimizerMinifyJavaScript:             ptr.ToBool(false),
		OptimizerMobileImageQuality:           ptr.ToInt32(10),
		OptimizerMobileMaxWidth:               ptr.ToInt32(200),
		OptimizerWatermarkEnabled:             ptr.ToBool(true),
		OptimizerWatermarkMinImageSize:        ptr.ToInt32(150),
		OptimizerWatermarkOffset:              ptr.ToFloat64(1),
		OptimizerWatermarkPosition:            ptr.ToInt(150),
		OptimizerWatermarkURL:                 ptr.ToString("https://via.placeholder.com/150"),
	}

	tf := fmt.Sprintf(`
resource "bunny_pullzone" "%s" {
	aws_signing_enabled = %t
	aws_signing_key = "%s"
	aws_signing_region_name = "%s"
	aws_signing_secret = "%s"
	allowed_referrers = %s
	block_post_requests = %t
	block_root_path_access = %t
	blocked_countries = %s
	blocked_ips = %s
	budget_redirected_countries = %s
	cache_control_browser_max_age_override  = %d
	cache_control_max_age_override = %d
	cache_error_responses = %t
	disable_cookies = %t
	enable_avif_vary = %t
	enable_cache_slice = %t
	enable_country_code_vary = %t
	#enable_geo_zone_af
	#enable_geo_zone_asia
	#enable_geo_zone_eu
	#enable_geo_zone_sa
	#enable_geo_zone_us
	enable_hostname_vary = %t
	enable_logging = %t
	enable_mobile_vary = %t
	enable_origin_shield = %t
	enable_tlsv1 = %t
	enable_tls1_1 = %t
	enable_webp_vary = %t
	error_page_custom_code = "%s"
	error_page_enable_custom_code = "%t"
	error_page_enable_statuspage_widget = %t
	error_page_statuspage_code = "%s"
	error_page_whitelabel = "%t"
	follow_redirects = %t
	ignore_query_strings = %t
	log_forwarding_enabled = %t
	log_forwarding_hostname = "%s"
	log_forwarding_port = %d
	log_forwarding_token = "%s"
	# logging_ip_anonymization_enabled // the field can only bet set after signing the dpa-agreement in the webinterface
	# logging_save_to_storage
	# logging_storage_zone_id
	origin_shield_zone_code = "%s"
	origin_url = "%s"
	# perma_cache_storage_zone_id
	type = %d
	verify_origin_ssl = %t
	zone_security_enabled = %t
	zone_security_include_hash_remote_ip = %t
	# blocked_referrers // uses different API
	name = "%s"
	# storage_zone_id
	# zone_security_key

	safehop {
		enable = %t
	    	origin_connect_timeout = %d
	    	origin_response_timeout  =  %d
	    	origin_retries = %d
	    	origin_retry_5xx_response = %t
	    	origin_retry_connection_timeout = %t
	    	origin_retry_delay = %d
	    	origin_retry_response_timeout = %t
	}

	headers {
		enable_access_control_origin_header = %t
		access_control_origin_header_extensions = "%s"
		add_canonical_header = %t
		add_host_header = %t
	}

	limits {
		request_limit = %d
		monthly_bandwidth_limit = %d
		connection_limit_per_ip_count = %d
	}

	optimizer {
		enabled = %t
		enable_webp = %t
		minify_css = %t
		minify_javascript = %t
		enable_manipulation_engine = %t

		smart_image_optimization {
			enabled = %t
			desktop_max_width = %d
			image_quality = %d
			mobile_max_width = %d
			mobile_image_quality = %d
		}

		watermark {
			enabled = %t
			url = "%s"
			offset = %f
			min_image_size = %d
			position = %d
		}
	}
}
`,
		resourceName,

		ptr.GetBool(attrs.AWSSigningEnabled),
		ptr.GetString(attrs.AWSSigningKey),
		ptr.GetString(attrs.AWSSigningRegionName),
		ptr.GetString(attrs.AWSSigningSecret),
		tfStrList(attrs.AllowedReferrers),
		ptr.GetBool(attrs.BlockPostRequests),
		ptr.GetBool(attrs.BlockRootPathAccess),
		tfStrList(attrs.BlockedCountries),
		tfStrList(attrs.BlockedIPs),
		tfStrList(attrs.BudgetRedirectedCountries),
		ptr.GetInt64(attrs.CacheControlBrowserMaxAgeOverride),
		ptr.GetInt64(attrs.CacheControlMaxAgeOverride),
		ptr.GetBool(attrs.CacheErrorResponses),
		ptr.GetBool(attrs.DisableCookies),
		ptr.GetBool(attrs.EnableAvifVary),
		ptr.GetBool(attrs.EnableCacheSlice),
		ptr.GetBool(attrs.EnableCountryCodeVary),
		ptr.GetBool(attrs.EnableHostnameVary),
		ptr.GetBool(attrs.EnableLogging),
		ptr.GetBool(attrs.EnableMobileVary),
		ptr.GetBool(attrs.EnableOriginShield),
		ptr.GetBool(attrs.EnableTLS1),
		ptr.GetBool(attrs.EnableTLS11),
		ptr.GetBool(attrs.EnableWebPVary),
		ptr.GetString(attrs.ErrorPageCustomCode),
		ptr.GetBool(attrs.ErrorPageEnableCustomCode),
		ptr.GetBool(attrs.ErrorPageEnableStatuspageWidget),
		ptr.GetString(attrs.ErrorPageStatuspageCode),
		ptr.GetBool(attrs.ErrorPageWhitelabel),
		ptr.GetBool(attrs.FollowRedirects),
		ptr.GetBool(attrs.IgnoreQueryStrings),
		ptr.GetBool(attrs.LogForwardingEnabled),
		ptr.GetString(attrs.LogForwardingHostname),
		ptr.GetInt32(attrs.LogForwardingPort),
		ptr.GetString(attrs.LogForwardingToken),
		//ptr.GetBool(attrs.LoggingIPAnonymizationEnabled),
		//ptr.GetBool(attrs.LoggingSaveToStorage),
		ptr.GetString(attrs.OriginShieldZoneCode),
		ptr.GetString(attrs.OriginURL),
		ptr.GetInt(attrs.Type),
		ptr.GetBool(attrs.VerifyOriginSSL),
		ptr.GetBool(attrs.ZoneSecurityEnabled),
		ptr.GetBool(attrs.ZoneSecurityIncludeHashRemoteIP),
		ptr.GetString(attrs.Name),

		ptr.GetBool(attrs.EnableSafeHop),
		ptr.GetInt32(attrs.OriginConnectTimeout),
		ptr.GetInt32(attrs.OriginResponseTimeout),
		ptr.GetInt32(attrs.OriginRetries),
		ptr.GetBool(attrs.OriginRetry5xxResponses),
		ptr.GetBool(attrs.OriginRetryConnectionTimeout),
		ptr.GetInt32(attrs.OriginRetryDelay),
		ptr.GetBool(attrs.OriginRetryResponseTimeout),

		ptr.GetBool(attrs.EnableAccessControlOriginHeader),
		strings.Join(attrs.AccessControlOriginHeaderExtensions, ","),
		ptr.GetBool(attrs.AddCanonicalHeader),
		ptr.GetBool(attrs.AddHostHeader),

		ptr.GetInt32(attrs.RequestLimit),
		ptr.GetInt64(attrs.MonthlyBandwidthLimit),
		ptr.GetInt32(attrs.ConnectionLimitPerIPCount),

		ptr.GetBool(attrs.OptimizerEnabled),
		ptr.GetBool(attrs.OptimizerEnableWebP),
		ptr.GetBool(attrs.OptimizerMinifyCSS),
		ptr.GetBool(attrs.OptimizerMinifyJavaScript),
		ptr.GetBool(attrs.OptimizerEnableManipulationEngine),

		ptr.GetBool(attrs.OptimizerAutomaticOptimizationEnabled),
		ptr.GetInt32(attrs.OptimizerDesktopMaxWidth),
		ptr.GetInt32(attrs.OptimizerImageQuality),
		ptr.GetInt32(attrs.OptimizerMobileMaxWidth),
		ptr.GetInt32(attrs.OptimizerMobileImageQuality),

		ptr.GetBool(attrs.OptimizerWatermarkEnabled),
		ptr.GetString(attrs.OptimizerWatermarkURL),
		ptr.GetFloat64(attrs.OptimizerWatermarkOffset),
		ptr.GetInt32(attrs.OptimizerWatermarkMinImageSize),
		ptr.GetInt(attrs.OptimizerWatermarkPosition),
	)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check:  checkPzState(t, fullResourceName, &attrs),
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
		CheckDestroy: checkPullZoneNotExists(fullResourceName),
	})
}

func TestAccPullZone_CaseInsensitiveOrderIndependentFields(t *testing.T) {
	pzName := randResourceName()

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "testpz" {
	name = "%s"
	origin_url ="http://bunny.net"
	budget_redirected_countries = ["DE","us","aL","Dz"]
	blocked_countries = ["ER","ee","fK","Fj"]
	allowed_referrers = ["ID","ir","iQ","Ie"]
}`, pzName),
			},
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "testpz" {
	name = "%s"
	origin_url ="http://bunny.net"
	budget_redirected_countries = ["us","DE","al","dz"]
	blocked_countries = ["ee","ER","fk","fj"]
	allowed_referrers = ["ie","ID","ir","iq",]
}`, pzName),
			},

			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "testpz" {
	name = "%s"
	origin_url ="http://bunny.net"
}`, pzName),
				Destroy: true,
			},
		},
	})
}

func checkPzState(t *testing.T, resourceName string, wanted *bunny.PullZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		resourceState := s.Modules[0].Resources[resourceName]
		if resourceState == nil {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		insState := resourceState.Primary
		if insState == nil {
			return fmt.Errorf("resource %s has no primary state", resourceName)
		}

		if insState.ID == "" {
			return errors.New("ID is empty")
		}

		id, err := strconv.Atoi(insState.ID)
		if err != nil {
			return fmt.Errorf("could not convert resource ID %q to int64: %w", id, err)
		}

		pz, err := clt.PullZone.Get(context.Background(), int64(id))
		if err != nil {
			return fmt.Errorf("fetching pull-zone with id %d from api client failed: %w", id, err)
		}

		diff := pzDiff(t, wanted, pz)
		if len(diff) != 0 {
			return fmt.Errorf("wanted and actual state differs:\n%s", strings.Join(diff, "\n"))
		}

		return nil

	}
}

func strDiff(a, b *string) string {
	if a == b {
		return ""
	}

	if a == nil || b == nil {
		return fmt.Sprintf("%q (%v) != %q (%v)", ptr.GetString(a), a, ptr.GetString(b), b)
	}

	if ptr.GetString(a) != ptr.GetString(b) {
		return fmt.Sprintf("%q (%v) != %q (%v)", ptr.GetString(a), a, ptr.GetString(b), b)
	}

	return ""
}

func boolDiff(a, b *bool) string {
	if a == b {
		return ""
	}

	if a == nil || b == nil {
		return fmt.Sprintf("%t (%v) != %t (%v)", ptr.GetBool(a), a, ptr.GetBool(b), b)
	}

	if ptr.GetBool(a) != ptr.GetBool(b) {
		return fmt.Sprintf("%t (%v) != %t (%v)", ptr.GetBool(a), a, ptr.GetBool(b), b)
	}

	return ""
}

func intDiff(a, b *int) string {
	if a == b {
		return ""
	}

	if a == nil || b == nil {
		return fmt.Sprintf("%d (%v) != %d (%v)", ptr.GetInt(a), a, ptr.GetInt(b), b)
	}

	if ptr.GetInt(a) != ptr.GetInt(b) {
		return fmt.Sprintf("%d (%v) != %d (%v)", ptr.GetInt(a), a, ptr.GetInt(b), b)
	}

	return ""
}

func int64Diff(a, b *int64) string {
	if a == b {
		return ""
	}

	if a == nil || b == nil {
		return fmt.Sprintf("%d (%v) != %d (%v)", ptr.GetInt64(a), a, ptr.GetInt64(b), b)
	}

	if ptr.GetInt64(a) != ptr.GetInt64(b) {
		return fmt.Sprintf("%d (%v) != %d (%v)", ptr.GetInt64(a), a, ptr.GetInt64(b), b)
	}

	return ""
}

func int32Diff(a, b *int32) string {
	if a == b {
		return ""
	}

	if a == nil || b == nil {
		return fmt.Sprintf("%d (%v) != %d (%v)", ptr.GetInt32(a), a, ptr.GetInt32(b), b)
	}

	if ptr.GetInt32(a) != ptr.GetInt32(b) {
		return fmt.Sprintf("%d (%v) != %d (%v)", ptr.GetInt32(a), a, ptr.GetInt32(b), b)
	}

	return ""
}

func float64Diff(a, b *float64) string {
	if a == b {
		return ""
	}

	if a == nil || b == nil {
		return fmt.Sprintf("%f (%v) != %f (%v)", ptr.GetFloat64(a), a, ptr.GetFloat64(b), b)
	}

	if ptr.GetFloat64(a) != ptr.GetFloat64(b) {
		return fmt.Sprintf("%f (%v) != %f (%v)", ptr.GetFloat64(a), a, ptr.GetFloat64(b), b)
	}

	return ""
}

func strSliceDiff(a, b []string) string {
	if !strSliceEqual(a, b) {
		return fmt.Sprintf("%+v != %+v ", a, b)
	}

	return ""
}

// pullZoneDiffIgnoredFields contains a list of fieldsnames in a bunny.PullZone struct that are ignored by pzDiff.
var pullZoneDiffIgnoredFields = map[string]struct{}{
	"AccessControlOriginHeaderExtensions": {}, // computed field
	"BlockedReferrers":                    {}, // computed field
	"CnameDomain":                         {}, // computed field
	"EnableGeoZoneAF":                     {}, // computed field
	"EnableGeoZoneAsia":                   {}, // computed field
	"EnableGeoZoneEU":                     {}, // computed field
	"EnableGeoZoneSA":                     {}, // computed field
	"EnableGeoZoneUS":                     {}, // computed field
	"Enabled":                             {}, // computed field
	"ID":                                  {}, // computed field
	"LoggingIPAnonymizationEnabled":       {}, // can only bet set if DPA agreement was signed in the webinterface
	"VideoLibraryID":                      {}, // computed field
	"ZoneSecurityKey":                     {}, // computed field

	// the following fields are ignored because they are not implemented in the provider
	"BurstSize":                          {},
	"CacheErrorResponses":                {},
	"DNSRecordID":                        {},
	"DNSZoneID":                          {},
	"EnableAutoSSL":                      {},
	"EnableCookieVary":                   {},
	"EnableSmartCache":                   {},
	"LimitRateAfter":                     {},
	"LimitRatePerSecond":                 {},
	"LogAnonymizationType":               {},
	"LogFormat":                          {},
	"LogForwardingFormat":                {},
	"LogForwardingProtocol":              {},
	"OptimizerForceClasses":              {},
	"OriginHostHeader":                   {},
	"OriginShieldEnableConcurrencyLimit": {},
	"OriginShieldMaxConcurrentRequests":  {},
	"OriginShieldMaxQueuedRequests":      {},
	"OriginShieldQueueMaxWaitTime":       {},
	"OriginType":                         {},
	"ShieldDDosProtectionEnabled":        {},
	"ShieldDDosProtectionType":           {},
	"UseBackgroundUpdate":                {},
	"UseStaleWhileOffline":               {},
	"UseStaleWhileUpdating":              {},

	// The following fields are tested by separate testcases and ignored in
	// pull zone testcases.
	"Hostnames": {},
	"EdgeRules": {},

	// the following fields are ignored because they export accounting data
	// and aren't configuration settings:
	"MonthlyBandwidthUsed": {},
	"PriceOverride":        {},
	"MonthlyCharges":       {},

	// the following fields can not be tested because they require a
	// storage zone, which currently can not be created via the provider
	"LoggingSaveToStorage":       {},
	"PermaCacheStorageZoneID":    {},
	"LoggingSaveToStorageZoneID": {},
	"StorageZoneID":              {},
	"LoggingStorageZoneID":       {},
}

func pzDiff(t *testing.T, a, b interface{}) []string {
	t.Helper()
	return diffStructs(t, a, b, pullZoneDiffIgnoredFields)
}

func TestAccPullZone_OriginURLAndStorageZoneIDAreExclusive(t *testing.T) {
	pzName := randResourceName()

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "testpz" {
	name = "%s"
	origin_url ="http://bunny.net"
	storage_zone_id = 300
}`, pzName),
				ExpectError: regexp.MustCompile("only one of.*can be specified"),
			},
		},
	})
}
