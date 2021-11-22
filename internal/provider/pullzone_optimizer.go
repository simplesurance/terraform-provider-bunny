package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyOptimizerEnabled                  = "enabled"
	keyOptimizerEnableManipulationEngine = "enable_manipulation_engine"
	keyOptimizerEnableWebP               = "enable_webp"
	keyOptimizerMinifyCSS                = "minify_css"
	keyOptimizerMinifyJavaScript         = "minify_javascript"

	keyOptimizerSmartImageOptimizationBlock  = "smart_image_optimization"
	keyOptimizerAutomaticOptimizationEnabled = "enabled"
	keyOptimizerImageQuality                 = "image_quality"
	keyOptimizerDesktopMaxWidth              = "desktop_max_width"
	keyOptimizerMobileImageQuality           = "mobile_image_quality"
	keyOptimizerMobileMaxWidth               = "mobile_max_width"

	keyOptimizerWatermarkBlock        = "watermark"
	keyOptimizerWatermarkEnabled      = "enabled"
	keyOptimizerWatermarkMinImageSize = "min_image_size"
	keyOptimizerWatermarkOffset       = "offset"
	keyOptimizerWatermarkPosition     = "position"
	keyOptimizerWatermarkURL          = "url"
)

var resourcePullZoneOptimizer = &schema.Resource{
	Schema: map[string]*schema.Schema{
		keyOptimizerEnabled: {
			Type:        schema.TypeBool,
			Description: "Determines if the optimizer should be enabled for this zone.",
			Default:     false,
			Optional:    true,
		},
		keyOptimizerEnableWebP: {
			Type:        schema.TypeBool,
			Description: "If enabled, images will be automatically converted into an efficient WebP format when supported by the client to greatly reduce file size and improve load times.",
			Default:     true,
			Optional:    true,
		},
		keyOptimizerMinifyCSS: {
			Type:        schema.TypeBool,
			Description: "If enabled, CSS files will be automatically minified to reduce their file size without modifying the functionality.",
			Default:     true,
			Optional:    true,
		},

		keyOptimizerMinifyJavaScript: {
			Type:        schema.TypeBool,
			Description: "Determines if the JavaScript minifcation should be enabled.",
			Default:     true,
			Optional:    true,
		},
		keyOptimizerEnableManipulationEngine: {
			Type:        schema.TypeBool,
			Description: "Enable on the fly image manipulation engine for dynamic URL based image manipulation.",
			Optional:    true,
			Default:     true,
		},

		keyOptimizerSmartImageOptimizationBlock: {
			Type:             schema.TypeList,
			MaxItems:         1,
			Optional:         true,
			DiffSuppressFunc: diffSupressMissingOptionalBlock,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					keyOptimizerAutomaticOptimizationEnabled: {
						Type:        schema.TypeBool,
						Description: "If enabled, Bunny Optimizer will automatically resize and compress images for desktop and mobile devices.",
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

					keyOptimizerImageQuality: {
						Type:             schema.TypeInt,
						Description:      "Determines the image quality for desktop clients.",
						Optional:         true,
						Default:          85,
						ValidateDiagFunc: validateIsInt32,
					},
					keyOptimizerMobileMaxWidth: {
						Type:             schema.TypeInt,
						Description:      "Determines the maximum automatic image size for mobile clients.",
						Optional:         true,
						Default:          800,
						ValidateDiagFunc: validateIsInt32,
					},
					keyOptimizerMobileImageQuality: {
						Type:             schema.TypeInt,
						Description:      "Determines the image quality for mobile clients.",
						Optional:         true,
						Default:          70,
						ValidateDiagFunc: validateIsInt32,
					},
				},
			},
		},

		keyOptimizerWatermarkBlock: {
			Type:             schema.TypeList,
			MaxItems:         1,
			Optional:         true,
			DiffSuppressFunc: diffSupressMissingOptionalBlock,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					keyOptimizerWatermarkEnabled: {
						Type:        schema.TypeBool,
						Description: "Determines if image watermarking should be enabled.",
						Default:     true,
						Optional:    true,
					},
					keyOptimizerWatermarkURL: {
						Type:        schema.TypeString,
						Description: "Sets the URL of the watermark image.",
						Optional:    true,
					},
					keyOptimizerWatermarkOffset: {
						Type:        schema.TypeFloat,
						Description: "Sets the offset of the watermark image.",
						Optional:    true,
						Default:     3,
					},
					keyOptimizerWatermarkMinImageSize: {
						Type:             schema.TypeInt,
						Description:      "Sets the minimum image size to which the watermark will be added.",
						Optional:         true,
						Default:          300,
						ValidateDiagFunc: validateIsInt32,
					},
					keyOptimizerWatermarkPosition: {
						Type:        schema.TypeInt,
						Description: "Sets the position of the watermark image.",
						Optional:    true,
						Default:     0,
					},
				},
			},
		},
	},
}

func optimizerFlatten(pz *bunny.PullZone, d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{{
		keyOptimizerEnabled:                     pz.OptimizerEnabled,
		keyOptimizerEnableManipulationEngine:    pz.OptimizerEnableManipulationEngine,
		keyOptimizerEnableWebP:                  pz.OptimizerEnableWebP,
		keyOptimizerMinifyCSS:                   pz.OptimizerMinifyCSS,
		keyOptimizerMinifyJavaScript:            pz.OptimizerMinifyJavaScript,
		keyOptimizerSmartImageOptimizationBlock: optimizerSmartImageOptimizationFlatten(pz),
		keyOptimizerWatermarkBlock:              optimizerWatermarkFlatten(pz),
	}}
}

func optimizerSmartImageOptimizationFlatten(pz *bunny.PullZone) []map[string]interface{} {
	return []map[string]interface{}{{
		keyOptimizerAutomaticOptimizationEnabled: pz.OptimizerAutomaticOptimizationEnabled,
		keyOptimizerDesktopMaxWidth:              pz.OptimizerDesktopMaxWidth,
		keyOptimizerImageQuality:                 pz.OptimizerImageQuality,
		keyOptimizerMobileMaxWidth:               pz.OptimizerMobileMaxWidth,
		keyOptimizerMobileImageQuality:           pz.OptimizerMobileImageQuality,
	}}
}

func optimizerSmartImageOptimizationExpand(res *bunny.PullZoneUpdateOptions, m structure) {
	if len(m) == 0 {
		return
	}

	res.OptimizerAutomaticOptimizationEnabled = m.getBoolPtr(keyOptimizerAutomaticOptimizationEnabled)
	res.OptimizerImageQuality = m.getInt32Ptr(keyOptimizerImageQuality)
	res.OptimizerDesktopMaxWidth = m.getInt32Ptr(keyOptimizerDesktopMaxWidth)
	res.OptimizerMobileImageQuality = m.getInt32Ptr(keyOptimizerMobileImageQuality)
	res.OptimizerMobileMaxWidth = m.getInt32Ptr(keyOptimizerMobileMaxWidth)
}

func optimizerWatermarkFlatten(pz *bunny.PullZone) []map[string]interface{} {
	return []map[string]interface{}{{
		keyOptimizerWatermarkEnabled:      pz.OptimizerWatermarkEnabled,
		keyOptimizerWatermarkURL:          pz.OptimizerWatermarkURL,
		keyOptimizerWatermarkOffset:       pz.OptimizerWatermarkOffset,
		keyOptimizerWatermarkMinImageSize: pz.OptimizerWatermarkMinImageSize,
		keyOptimizerWatermarkPosition:     pz.OptimizerWatermarkPosition,
	}}
}

func optimizerWatermarkExpand(res *bunny.PullZoneUpdateOptions, m structure) {
	if len(m) == 0 {
		return
	}

	res.OptimizerWatermarkEnabled = m.getBoolPtr(keyOptimizerWatermarkEnabled)
	res.OptimizerWatermarkURL = m.getStrPtr(keyOptimizerWatermarkURL)
	res.OptimizerWatermarkOffset = m.getFloat64Ptr(keyOptimizerWatermarkOffset)
	res.OptimizerWatermarkMinImageSize = m.getInt32Ptr(keyOptimizerWatermarkMinImageSize)
	res.OptimizerWatermarkPosition = m.getIntPtr(keyOptimizerWatermarkPosition)
}

func optimizerFromResource(res *bunny.PullZoneUpdateOptions, d *schema.ResourceData) {
	m := structureFromResource(d, keyOptimizer)
	if len(m) == 0 {
		return
	}

	res.OptimizerEnabled = m.getBoolPtr(keyOptimizerEnabled)
	res.OptimizerEnableManipulationEngine = m.getBoolPtr(keyOptimizerEnableManipulationEngine)
	res.OptimizerEnableWebP = m.getBoolPtr(keyOptimizerEnableWebP)
	res.OptimizerMinifyCSS = m.getBoolPtr(keyOptimizerMinifyCSS)
	res.OptimizerMinifyJavaScript = m.getBoolPtr(keyOptimizerMinifyJavaScript)

	smartImageBlock := m[keyOptimizerSmartImageOptimizationBlock].([]interface{})
	optimizerSmartImageOptimizationExpand(res, structureFromElem(smartImageBlock))

	watermarkBlock := m[keyOptimizerWatermarkBlock].([]interface{})
	optimizerWatermarkExpand(res, structureFromElem(watermarkBlock))
}
