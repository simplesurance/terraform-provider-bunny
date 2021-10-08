package provider

import (
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validateIsInt32 = validation.ToDiagFunc(validation.IntBetween(math.MinInt32, math.MaxInt32))
