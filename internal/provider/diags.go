package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/diag"

func diagsErrFromErr(summary string, err error) diag.Diagnostics {
	return diagsFromErr(summary, err, diag.Error)
}

func diagsWarnFromErr(summary string, err error) diag.Diagnostics {
	return diagsFromErr(summary, err, diag.Warning)
}

func diagsFromErr(summary string, err error, severity diag.Severity) diag.Diagnostics {
	return diag.Diagnostics{{
		Severity: severity,
		Summary:  summary,
		Detail:   err.Error(),
	}}
}
