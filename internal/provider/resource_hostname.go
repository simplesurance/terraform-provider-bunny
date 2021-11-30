package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	bunny "github.com/simplesurance/bunny-go"
)

const (
	keyHostnamePullZoneID          = "pull_zone_id"
	keyHostnameHostname            = "hostname" // yes, that variable name is intentional :-)
	keyHostnameForceSSL            = "force_ssl"
	keyHostnameIsSystemHostname    = "is_system_hostname"
	keyHostnameHasCertificate      = "has_certificate"
	keyHostnameLoadFreeCertificate = "load_free_certificate"
)

const (
	loadFreeCertMinDelay = 5 * time.Second
)

func resourceHostname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostnameCreate,
		UpdateContext: resourceHostnameUpdate,
		ReadContext:   resourceHostnameRead,
		DeleteContext: resourceHostnameDelete,

		Schema: map[string]*schema.Schema{
			keyHostnamePullZoneID: {
				Type:        schema.TypeInt,
				Description: "The ID of the pull zone to that the hostname belongs.",
				Required:    true,
				ForceNew:    true,
			},
			keyHostnameHostname: {
				Type:        schema.TypeString,
				Description: "The hostname value for the domain name.",
				Required:    true,
				ForceNew:    true,
			},
			keyHostnameForceSSL: {
				Type:        schema.TypeBool,
				Description: "Determines if the Force SSL feature is enabled.",
				Optional:    true,
				Default:     false,
			},
			keyHostnameIsSystemHostname: {
				Type:        schema.TypeBool,
				Description: "Determines if this is a system hostname controlled by bunny.net.",
				Computed:    true,
			},
			keyHostnameHasCertificate: {
				Type:        schema.TypeBool,
				Description: "Determines if the hostname has an SSL certificate configured.",
				Computed:    true,
			},
			keyHostnameLoadFreeCertificate: {
				Type:        schema.TypeBool,
				Description: "Determines if a free SSL certificate should be generated and loaded for the hostname",
				ForceNew:    true,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceHostnameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	pullZoneID := int64(d.Get(keyHostnamePullZoneID).(int))
	hostnameOpt := resourceDataToAddCustomHostnameOption(d)

	err := clt.PullZone.AddCustomHostname(ctx, pullZoneID, hostnameOpt)
	if err != nil {
		return diagsErrFromErr("could not add hostname", err)
	}

	var diag diag.Diagnostics

	if d.Get(keyHostnameLoadFreeCertificate).(bool) {
		if err := loadFreeCertRetry(ctx, clt, d.Timeout(schema.TimeoutCreate), *hostnameOpt.Hostname); err != nil {
			diag = diagsErrFromErr("creating hostname succeeded, loading free ssl certificate failed", err)
		}
	}

	if forceSSL := d.Get(keyHostnameForceSSL).(bool); forceSSL {
		err = clt.PullZone.SetForceSSL(ctx, pullZoneID, &bunny.SetForceSSLOptions{
			Hostname: hostnameOpt.Hostname,
			ForceSSL: &forceSSL,
		})
		if err != nil {
			diag = append(diag, diagsErrFromErr("creating hostname succeeded, enabling force_ssl failed", err)...)
		}
	}

	hostname, err := resourceHostnameGetByName(ctx, clt, pullZoneID, *hostnameOpt.Hostname)
	if err != nil {
		return append(diag, diagsErrFromErr("creating hostname succeeded, retrieving it from api failed", err)...)
	}

	if err := hostnameToResource(hostname, d); err != nil {
		return append(diag, diagsErrFromErr("converting hostname api type to terraform resource failed", err)...)
	}

	return diag
}

func loadFreeCertRetry(ctx context.Context, clt *bunny.Client, timeout time.Duration, hostname string) error {
	const (
		stateWaitingForDNSRecord = "waiting_for_dns_record"
		stateDone                = "certificate_loaded"
	)

	stateConf := resource.StateChangeConf{
		Pending:    []string{stateWaitingForDNSRecord},
		Target:     []string{stateDone},
		Timeout:    timeout,
		MinTimeout: loadFreeCertMinDelay,
		Refresh: func() (interface{}, string, error) {
			err := clt.PullZone.LoadFreeCertificate(ctx, hostname)
			if err != nil {
				if apiErr, ok := err.(*bunny.APIError); ok {
					if strings.Contains(strings.ToLower(apiErr.Message), "is not pointing to our servers") {
						logger.Infof("cname dns record missing for hostname %q", hostname)

						return "", stateWaitingForDNSRecord, nil
					}

					return nil, "", err
				}
			}

			// StateChangeConf seems to require that a non-nil
			// result is returned to consider the state change as successful.
			// Return an "" instead of nil as result.
			return "", stateDone, nil
		},
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceHostnameGetByName(ctx context.Context, clt *bunny.Client, pullZoneID int64, hostname string) (*bunny.Hostname, error) {
	pz, err := clt.PullZone.Get(ctx, pullZoneID)
	if err != nil {
		return nil, fmt.Errorf("retrieving pull zone failed: %w", err)
	}

	for _, pzHostname := range pz.Hostnames {
		if pzHostname.Value == nil {
			logger.Warnf("bunny.net api returned pull zone (%d) with an hostname element with nil value", pullZoneID)
			continue
		}

		if *pzHostname.Value == hostname {
			if pzHostname.ID == nil {
				return nil, fmt.Errorf("found hostname with name %q of pull zone (%d) but id is nil", hostname, pullZoneID)
			}

			return pzHostname, nil
		}
	}

	return nil, errors.New("hostname not found")
}

func resourceDataToAddCustomHostnameOption(d *schema.ResourceData) *bunny.AddCustomHostnameOptions {
	hostname := d.Get(keyHostnameHostname).(string)
	return &bunny.AddCustomHostnameOptions{
		Hostname: &hostname,
	}
}

func resourceHostnameDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	pullZoneID := int64(d.Get(keyHostnamePullZoneID).(int))
	hostnameOpt := hostnameFromResource(d)

	return diag.FromErr(clt.PullZone.RemoveCustomHostname(ctx, pullZoneID, hostnameOpt))
}

func hostnameFromResource(d *schema.ResourceData) *bunny.RemoveCustomHostnameOptions {
	hostname := d.Get(keyHostnameHostname).(string)
	return &bunny.RemoveCustomHostnameOptions{
		Hostname: &hostname,
	}
}

func resourceHostnameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clt := meta.(*bunny.Client)

	hostnameID, err := getIDAsInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	pullZoneID := int64(d.Get(keyHostnamePullZoneID).(int))

	hostname, err := resourceHostnameGetByID(ctx, clt, pullZoneID, hostnameID)
	if err != nil {
		diagsErrFromErr("could not fetch hostname from provider", err)
	}

	if err := hostnameToResource(hostname, d); err != nil {
		return diagsErrFromErr("converting api hostname to resource data failed", err)
	}

	return nil
}

func resourceHostnameGetByID(ctx context.Context, clt *bunny.Client, pullZoneID, hostnameID int64) (*bunny.Hostname, error) {
	pz, err := clt.PullZone.Get(ctx, pullZoneID)
	if err != nil {
		return nil, fmt.Errorf("retrieving pull zone failed: %w", err)
	}

	for _, hostname := range pz.Hostnames {
		if hostname.ID == nil {
			logger.Warnf("bunny.net api returned hostname with nil ID for pull zone: %d", pullZoneID)
			continue
		}

		if *hostname.ID == hostnameID {
			return hostname, nil
		}
	}

	return nil, fmt.Errorf("pull zone with id %d, has no hostname with id: %d", pullZoneID, hostnameID)
}

func hostnameToResource(hostname *bunny.Hostname, d *schema.ResourceData) error {
	if hostname.ID == nil {
		return errors.New("id is empty")
	}

	d.SetId(strconv.FormatInt(*hostname.ID, 10))

	if err := d.Set(keyHostnameHostname, hostname.Value); err != nil {
		return err
	}
	logger.Debugf("hostnameToResource %d, forcessl: %v", *hostname.ID, *hostname.ForceSSL)
	if err := d.Set(keyHostnameForceSSL, hostname.ForceSSL); err != nil {
		return err
	}
	if err := d.Set(keyHostnameIsSystemHostname, hostname.IsSystemHostname); err != nil {
		return err
	}
	if err := d.Set(keyHostnameHasCertificate, hostname.HasCertificate); err != nil {
		return err
	}

	return nil
}

func resourceHostnameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange(keyHostnameForceSSL) {
		// nothing to do, all other attributes have ForceNew enabled
		return nil
	}

	clt := meta.(*bunny.Client)

	pullZoneID := int64(d.Get(keyHostnamePullZoneID).(int))
	hostname := d.Get(keyHostnameHostname).(string)
	forceSSL := d.Get(keyHostnameForceSSL).(bool)

	err := clt.PullZone.SetForceSSL(ctx, pullZoneID, &bunny.SetForceSSLOptions{
		Hostname: &hostname,
		ForceSSL: &forceSSL,
	})

	if err != nil {
		return diagsErrFromErr("setting force ssl failed", err)
	}

	return nil
}
