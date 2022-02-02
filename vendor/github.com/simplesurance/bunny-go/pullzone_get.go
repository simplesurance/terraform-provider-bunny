package bunny

import (
	"context"
	"fmt"
)

// Constants for the Type fields of a Pull Zone.
const (
	PullZoneTypeStandard int = 1
	PullZoneTypeVolume   int = 2
)

// Constants for the values of the PatternMatchingType of EdgeRuleTrigger and
// TriggerMatchingType of an EdgeRule.
const (
	MatchingTypeAny int = iota
	MatchingTypeAll
	MatchingTypeNone
)

// PullZone represents the response of the the List and Get Pull Zone API endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_index2 https://docs.bunny.net/reference/pullzonepublic_index
type PullZone struct {
	ID *int64 `json:"Id,omitempty"`

	AWSSigningEnabled                   *bool    `json:"AWSSigningEnabled,omitempty"`
	AWSSigningKey                       *string  `json:"AWSSigningKey,omitempty"`
	AWSSigningRegionName                *string  `json:"AWSSigningRegionName,omitempty"`
	AWSSigningSecret                    *string  `json:"AWSSigningSecret,omitempty"`
	AccessControlOriginHeaderExtensions []string `json:"AccessControlOriginHeaderExtensions,omitempty"`
	AddCanonicalHeader                  *bool    `json:"AddCanonicalHeader,omitempty"`
	AddHostHeader                       *bool    `json:"AddHostHeader,omitempty"`
	AllowedReferrers                    []string `json:"AllowedReferrers,omitempty"`
	BlockPostRequests                   *bool    `json:"BlockPostRequests,omitempty"`
	BlockRootPathAccess                 *bool    `json:"BlockRootPathAccess,omitempty"`
	BlockedCountries                    []string `json:"BlockedCountries,omitempty"`
	BlockedIPs                          []string `json:"BlockedIps,omitempty"`
	BlockedReferrers                    []string `json:"BlockedReferrers,omitempty"`
	BudgetRedirectedCountries           []string `json:"BudgetRedirectedCountries,omitempty"`
	BurstSize                           *int32   `json:"BurstSize,omitempty"`
	// CacheControlBrowserMaxAgeOverride is called
	// CacheControlPublicMaxAgeOverride in the API. Both names refer to the
	// same setting.
	CacheControlBrowserMaxAgeOverride     *int64      `json:"CacheControlPublicMaxAgeOverride,omitempty"`
	CacheControlMaxAgeOverride            *int64      `json:"CacheControlMaxAgeOverride,omitempty"`
	CacheErrorResponses                   *bool       `json:"CacheErrorResponses,omitempty"`
	CnameDomain                           *string     `json:"CnameDomain,omitempty"`
	ConnectionLimitPerIPCount             *int32      `json:"ConnectionLimitPerIPCount,omitempty"`
	CookieVaryParameters                  []string    `json:"CookieVaryParameters,omitempty"`
	DNSRecordID                           *int64      `json:"DnsRecordId,omitempty"`
	DNSRecordValue                        *string     `json:"DnsRecordValue,omitempty"`
	DNSZoneID                             *int64      `json:"DnsZoneId,omitempty"`
	DisableCookies                        *bool       `json:"DisableCookies,omitempty"`
	EdgeRules                             []*EdgeRule `json:"EdgeRules,omitempty"`
	EnableAccessControlOriginHeader       *bool       `json:"EnableAccessControlOriginHeader,omitempty"`
	EnableAvifVary                        *bool       `json:"EnableAvifVary,omitempty"`
	EnableCacheSlice                      *bool       `json:"EnableCacheSlice,omitempty"`
	EnableCookieVary                      *bool       `json:"EnableCookieVary,omitempty"`
	EnableCountryCodeVary                 *bool       `json:"EnableCountryCodeVary,omitempty"`
	EnableGeoZoneAF                       *bool       `json:"EnableGeoZoneAF,omitempty"`
	EnableGeoZoneAsia                     *bool       `json:"EnableGeoZoneASIA,omitempty"`
	EnableGeoZoneEU                       *bool       `json:"EnableGeoZoneEU,omitempty"`
	EnableGeoZoneSA                       *bool       `json:"EnableGeoZoneSA,omitempty"`
	EnableGeoZoneUS                       *bool       `json:"EnableGeoZoneUS,omitempty"`
	EnableHostnameVary                    *bool       `json:"EnableHostnameVary,omitempty"`
	EnableLogging                         *bool       `json:"EnableLogging,omitempty"`
	EnableMobileVary                      *bool       `json:"EnableMobileVary,omitempty"`
	EnableOriginShield                    *bool       `json:"EnableOriginShield,omitempty"`
	EnableSafeHop                         *bool       `json:"EnableSafeHop,omitempty"`
	EnableTLS1                            *bool       `json:"EnableTLS1,omitempty"`
	EnableTLS11                           *bool       `json:"EnableTLS1_1,omitempty"`
	EnableWebPVary                        *bool       `json:"EnableWebPVary,omitempty"`
	Enabled                               *bool       `json:"Enabled,omitempty"`
	ErrorPageCustomCode                   *string     `json:"ErrorPageCustomCode,omitempty"`
	ErrorPageEnableCustomCode             *bool       `json:"ErrorPageEnableCustomCode,omitempty"`
	ErrorPageEnableStatuspageWidget       *bool       `json:"ErrorPageEnableStatuspageWidget,omitempty"`
	ErrorPageStatuspageCode               *string     `json:"ErrorPageStatuspageCode,omitempty"`
	ErrorPageWhitelabel                   *bool       `json:"ErrorPageWhitelabel,omitempty"`
	FollowRedirects                       *bool       `json:"FollowRedirects,omitempty"`
	Hostnames                             []*Hostname `json:"Hostnames,omitempty"`
	IgnoreQueryStrings                    *bool       `json:"IgnoreQueryStrings,omitempty"`
	LimitRateAfter                        *float64    `json:"LimitRateAfter,omitempty"`
	LimitRatePerSecond                    *float64    `json:"LimitRatePerSecond,omitempty"`
	LogForwardingEnabled                  *bool       `json:"LogForwardingEnabled,omitempty"`
	LogForwardingHostname                 *string     `json:"LogForwardingHostname,omitempty"`
	LogForwardingPort                     *int32      `json:"LogForwardingPort,omitempty"`
	LogForwardingToken                    *string     `json:"LogForwardingToken,omitempty"`
	LoggingIPAnonymizationEnabled         *bool       `json:"LoggingIPAnonymizationEnabled,omitempty"`
	LoggingSaveToStorage                  *bool       `json:"LoggingSaveToStorage,omitempty"`
	LoggingStorageZoneID                  *int64      `json:"LoggingStorageZoneId,omitempty"`
	MonthlyBandwidthLimit                 *int64      `json:"MonthlyBandwidthLimit,omitempty"`
	MonthlyBandwidthUsed                  *int64      `json:"MonthlyBandwidthUsed,omitempty"`
	MonthlyCharges                        *float64    `json:"MonthlyCharges,omitempty"`
	Name                                  *string     `json:"Name,omitempty"`
	OptimizerAutomaticOptimizationEnabled *bool       `json:"OptimizerAutomaticOptimizationEnabled,omitempty"`
	OptimizerDesktopMaxWidth              *int32      `json:"OptimizerDesktopMaxWidth,omitempty"`
	OptimizerEnableManipulationEngine     *bool       `json:"OptimizerEnableManipulationEngine,omitempty"`
	OptimizerEnableWebP                   *bool       `json:"OptimizerEnableWebP,omitempty"`
	OptimizerEnabled                      *bool       `json:"OptimizerEnabled,omitempty"`
	OptimizerImageQuality                 *int32      `json:"OptimizerImageQuality,omitempty"`
	OptimizerMinifyCSS                    *bool       `json:"OptimizerMinifyCSS,omitempty"`
	OptimizerMinifyJavaScript             *bool       `json:"OptimizerMinifyJavaScript,omitempty"`
	OptimizerMobileImageQuality           *int32      `json:"OptimizerMobileImageQuality,omitempty"`
	OptimizerMobileMaxWidth               *int32      `json:"OptimizerMobileMaxWidth,omitempty"`
	OptimizerWatermarkEnabled             *bool       `json:"OptimizerWatermarkEnabled,omitempty"`
	OptimizerWatermarkMinImageSize        *int32      `json:"OptimizerWatermarkMinImageSize,omitempty"`
	OptimizerWatermarkOffset              *float64    `json:"OptimizerWatermarkOffset,omitempty"`
	OptimizerWatermarkPosition            *int        `json:"OptimizerWatermarkPosition,omitempty"`
	OptimizerWatermarkURL                 *string     `json:"OptimizerWatermarkUrl,omitempty"`
	OriginConnectTimeout                  *int32      `json:"OriginConnectTimeout,omitempty"`
	OriginResponseTimeout                 *int32      `json:"OriginResponseTimeout,omitempty"`
	OriginRetries                         *int32      `json:"OriginRetries,omitempty"`
	OriginRetry5xxResponses               *bool       `json:"OriginRetry5xxResponses,omitempty"`
	OriginRetryConnectionTimeout          *bool       `json:"OriginRetryConnectionTimeout,omitempty"`
	OriginRetryDelay                      *int32      `json:"OriginRetryDelay,omitempty"`
	OriginRetryResponseTimeout            *bool       `json:"OriginRetryResponseTimeout,omitempty"`
	OriginShieldEnableConcurrencyLimit    *bool       `json:"OriginShieldEnableConcurrencyLimit,omitempty"`
	OriginShieldMaxConcurrentRequests     *int32      `json:"OriginShieldMaxConcurrentRequests,omitempty"`
	OriginShieldMaxQueuedRequests         *int32      `json:"OriginShieldMaxQueuedRequests,omitempty"`
	OriginShieldQueueMaxWaitTime          *int32      `json:"OriginShieldQueueMaxWaitTime,omitempty"`
	OriginShieldZoneCode                  *string     `json:"OriginShieldZoneCode,omitempty"`
	OriginURL                             *string     `json:"OriginUrl,omitempty"`
	PermaCacheStorageZoneID               *int64      `json:"PermaCacheStorageZoneId,omitempty"`
	PriceOverride                         *float64    `json:"PriceOverride,omitempty"`
	QueryStringVaryParameters             []string    `json:"QueryStringVaryParameters,omitempty"`
	RequestLimit                          *int32      `json:"RequestLimit,omitempty"`
	StorageZoneID                         *int64      `json:"StorageZoneId,omitempty"`
	Type                                  *int        `json:"Type,omitempty"`
	UseStaleWhileOffline                  *bool       `json:"UseStaleWhileOffline,omitempty"`
	UseStaleWhileUpdating                 *bool       `json:"UseStaleWhileUpdating,omitempty"`
	VerifyOriginSSL                       *bool       `json:"VerifyOriginSSL,omitempty"`
	VideoLibraryID                        *int64      `json:"VideoLibraryId,omitempty"`
	ZoneSecurityEnabled                   *bool       `json:"ZoneSecurityEnabled,omitempty"`
	ZoneSecurityIncludeHashRemoteIP       *bool       `json:"ZoneSecurityIncludeHashRemoteIP,omitempty"`
	ZoneSecurityKey                       *string     `json:"ZoneSecurityKey,omitempty"`
}

// Hostname represents a Hostname returned from the Get and List Pull Zone API Endpoints.
type Hostname struct {
	ID               *int64  `json:"Id,omitempty"`
	Value            *string `json:"Value,omitempty"`
	ForceSSL         *bool   `json:"ForceSSL,omitempty"`
	IsSystemHostname *bool   `json:"IsSystemHostname,omitempty"`
	HasCertificate   *bool   `json:"HasCertificate,omitempty"`
}

// EdgeRule represents an EdgeRule.
// It is returned from the Get and List Pull Zone and passed to the AddorUpdateEdgeRule API Endpoints.
type EdgeRule struct {
	GUID                *string            `json:"Guid,omitempty"`
	ActionType          *int               `json:"ActionType,omitempty"`
	ActionParameter1    *string            `json:"ActionParameter1,omitempty"`
	ActionParameter2    *string            `json:"ActionParameter2,omitempty"`
	Triggers            []*EdgeRuleTrigger `json:"Triggers,omitempty"`
	TriggerMatchingType *int               `json:"TriggerMatchingType,omitempty"`
	Description         *string            `json:"Description,omitempty"`
	Enabled             *bool              `json:"Enabled,omitempty"`
}

// Get retrieves the Pull Zone with the given id.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_index2
func (s *PullZoneService) Get(ctx context.Context, id int64) (*PullZone, error) {
	var res PullZone

	req, err := s.client.newGetRequest(fmt.Sprintf("pullzone/%d", id), nil)
	if err != nil {
		return nil, err
	}

	if err := s.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, err
}
