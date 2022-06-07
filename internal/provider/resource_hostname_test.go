package provider

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	ptr "github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	bunny "github.com/simplesurance/bunny-go"
)

type hostnamesWanted struct {
	TerraformPullZoneResourceName string
	PullZoneName                  string
	Hostnames                     []*bunny.Hostname
}

var hostnameDiffIgnoredFields = map[string]struct{}{
	"ID": {}, // computed field
}

func checkHostnameState(t *testing.T, wanted *hostnamesWanted) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		strID, err := idFromState(s, wanted.TerraformPullZoneResourceName)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			return fmt.Errorf("could not convert resource ID %q to int64: %w", strID, err)
		}

		pz, err := clt.PullZone.Get(context.Background(), int64(id))
		if err != nil {
			return fmt.Errorf("fetching pull-zone with id %d from api client failed: %w", id, err)
		}

		if err := stringsAreEqual(wanted.PullZoneName, pz.Name); err != nil {
			return fmt.Errorf("name of created pullzone differs: %w", err)
		}

		if len(pz.Hostnames) != len(wanted.Hostnames) {
			return fmt.Errorf("api returned pull request with %d hostnames, expected %d",
				len(pz.Hostnames), len(wanted.Hostnames),
			)
		}

		sortHostnames(wanted.Hostnames)
		sortHostnames(pz.Hostnames)

		for i := range pz.Hostnames {
			diff := diffStructs(t, wanted.Hostnames[i], pz.Hostnames[i], hostnameDiffIgnoredFields)
			if len(diff) != 0 {
				return fmt.Errorf("wanted and actual hostnames with idx %d differs:\n%s", i, strings.Join(diff, "\n"))
			}
		}

		return nil
	}
}

func sortHostnames(hostnames []*bunny.Hostname) {
	sort.Slice(hostnames, func(i, j int) bool {
		return ptr.GetString(hostnames[i].Value) < ptr.GetString(hostnames[j].Value)
	})
}

func TestAccHostname_basic(t *testing.T) {
	pzName := randPullZoneName()
	tf := fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = "google.de"
}`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check: checkHostnameState(t, &hostnamesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.pz",
					PullZoneName:                  pzName,
					Hostnames: []*bunny.Hostname{
						{
							Value:            ptr.ToString(defPullZoneHostname(pzName)),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(true),
							HasCertificate:   ptr.ToBool(true),
						},
						{
							Value:            ptr.ToString("google.de"),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
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

func TestAccHostname_addRemove(t *testing.T) {
	pzName := randPullZoneName()
	hostname1 := randHostname()
	hostname2 := randHostname()
	hostname3 := randHostname()
	hostname4 := randHostname()
	hostname5 := randHostname()
	hostname6 := randHostname()

	tfPz := fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tfPz + fmt.Sprintf(`
resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
	force_ssl = true
}

resource "bunny_hostname" "h2" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}

resource "bunny_hostname" "h3" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}
`, hostname1, hostname2, hostname3),

				Check: checkHostnameState(t, &hostnamesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.pz",
					PullZoneName:                  pzName,
					Hostnames: []*bunny.Hostname{
						{
							Value:            ptr.ToString(defPullZoneHostname(pzName)),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(true),
							HasCertificate:   ptr.ToBool(true),
						},
						{
							Value:            &hostname1,
							ForceSSL:         ptr.ToBool(true),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
						{
							Value:            &hostname2,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
						{
							Value:            &hostname3,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
					},
				}),
			},

			// Change all 3 hostname
			{
				Config: tfPz + fmt.Sprintf(`
resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}


resource "bunny_hostname" "h3" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}

resource "bunny_hostname" "h2" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}
`, hostname4, hostname5, hostname6),
				Check: checkHostnameState(t, &hostnamesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.pz",
					PullZoneName:                  pzName,
					Hostnames: []*bunny.Hostname{
						{
							Value:            ptr.ToString(defPullZoneHostname(pzName)),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(true),
							HasCertificate:   ptr.ToBool(true),
						},
						{
							Value:            &hostname4,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
						{
							Value:            &hostname5,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
						{
							Value:            &hostname6,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
					},
				}),
			},
		},
	})
}

func TestAccHostname_DefiningDuplicateHostnamesFails(t *testing.T) {
	pzName := randPullZoneName()
	tf := fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = "google.de"
}

resource "bunny_hostname" "h2" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = "google.de"
}
`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config:      tf,
				ExpectError: regexp.MustCompile(".*"), // TODO: match a more specific error string :-)
				PlanOnly:    true,
			},
		},
	})
}

func TestAccHostname_DefiningDefPullZoneHostnameFails(t *testing.T) {
	pzName := randPullZoneName()
	tf := fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = "%s"
}
`, pzName, defPullZoneHostname(pzName))

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config:      tf,
				ExpectError: regexp.MustCompile(".*"), // TODO: match a more specific error string :-)
			},
		},
	})
}

func TestAccCertificateOneof(t *testing.T) {
	pzName := randPullZoneName()
	tf := fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = "abcde.test"
	load_free_certificate = true

	certificate {
		certificate_data = "123"
		private_key_data = "456"
	}

}
`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config:      tf,
				ExpectError: regexp.MustCompile(`only one of "load_free_certificate" or "certificate" can be set`),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccCertificateCanBeSetWhenLoadFreeCertIsDisabled(t *testing.T) {
	pzName := randPullZoneName()
	tf := fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = "abcde.test"
	load_free_certificate = false

	certificate {
		certificate_data = "123"
		private_key_data = "456"
	}

}
`, pzName)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config:             tf,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCertificates(t *testing.T) {
	pzName := randPullZoneName()
	hostname := randHostname()

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q

	certificate {
		certificate_data = file("testdata/ssl.crt")
		private_key_data = file("testdata/ssl.key")
	}
}
`, pzName, hostname),
				Check: checkHostnameState(t, &hostnamesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.pz",
					PullZoneName:                  pzName,
					Hostnames: []*bunny.Hostname{
						{
							Value:            ptr.ToString(defPullZoneHostname(pzName)),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(true),
							HasCertificate:   ptr.ToBool(true),
						},
						{
							Value:            &hostname,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(true),
						},
					},
				}),
			},
			// change the certificate
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q

	certificate {
		certificate_data = file("testdata/ssl1.crt")
		private_key_data = file("testdata/ssl1.key")
	}
}
`, pzName, hostname),
				Check: checkHostnameState(t, &hostnamesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.pz",
					PullZoneName:                  pzName,
					Hostnames: []*bunny.Hostname{
						{
							Value:            ptr.ToString(defPullZoneHostname(pzName)),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(true),
							HasCertificate:   ptr.ToBool(true),
						},
						{
							Value:            &hostname,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(true),
						},
					},
				}),
			},

			// remove the certificate
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}
`, pzName, hostname),

				Check: checkHostnameState(t, &hostnamesWanted{
					TerraformPullZoneResourceName: "bunny_pullzone.pz",
					PullZoneName:                  pzName,
					Hostnames: []*bunny.Hostname{
						{
							Value:            ptr.ToString(defPullZoneHostname(pzName)),
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(true),
							HasCertificate:   ptr.ToBool(true),
						},
						{
							Value:            &hostname,
							ForceSSL:         ptr.ToBool(false),
							IsSystemHostname: ptr.ToBool(false),
							HasCertificate:   ptr.ToBool(false),
						},
					},
				}),
			},
		},
	})
}

func TestAccHostname_StateIsValidWhenCertUploadFails(t *testing.T) {
	t.Skip("disabled, because test sends 800kiB of bogus data to bunny api, which is not kind")

	pzName := randPullZoneName()
	hostname := randHostname()

	// The bunny API does not return an error if the posted data is not a
	// valid certificate.
	// To cause an API error, we post a big amount of zero bits which
	// causes a 500 error.
	// TODO: find a way to generate an error that does not require sending
	// a lot of bogus data.
	var bogusCertData [800 * 1024]byte

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				// Provoke a certificate upload error, the
				// hostname should still have been created on
				// bunny.net
				Destroy: false,
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q

	certificate {
		certificate_data = "%x"
		private_key_data = "5678"
	}
}
`, pzName, hostname, bogusCertData),
				ExpectError: regexp.MustCompile(".*uploading certificate failed.*"),
			},
			// the local terraform state should not be broken,
			// applying a change for the hostname must succeed
			{
				Config: fmt.Sprintf(`
resource "bunny_pullzone" "pz" {
	name = "%s"
	origin_url ="https://bunny.net"
}

resource "bunny_hostname" "h1" {
	pull_zone_id = bunny_pullzone.pz.id
	hostname = %q
}
`, pzName, hostname),
			},
		},
	})
}
