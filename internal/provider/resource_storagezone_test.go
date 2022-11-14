package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	ptr "github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	bunny "github.com/simplesurance/bunny-go"
)

type storageZoneWanted struct {
	TerraformResourceName string
	bunny.StorageZone
	Name   string
	Region string
}

func checkBasicStorageZoneAPIState(wanted *storageZoneWanted) resource.TestCheckFunc {
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

		sz, err := clt.StorageZone.Get(context.Background(), int64(id))
		if err != nil {
			return fmt.Errorf("fetching storage-zone with id %d from api client failed: %w", id, err)
		}

		if err := stringsAreEqual(wanted.Name, sz.Name); err != nil {
			return fmt.Errorf("name of created storagezone differs: %w", err)
		}

		return nil
	}
}

func checkStorageZoneNotExists(storageZoneName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		var page int32

		for {
			storagezones, err := clt.StorageZone.List(context.Background(), &bunny.PaginationOptions{
				Page:    page,
				PerPage: 1000,
			})
			if err != nil {
				return fmt.Errorf("listing storagezones failed: %w", err)
			}

			for _, sz := range storagezones.Items {
				if sz.Name == nil {
					return fmt.Errorf("got storagezone from api with empty Name: %+v", sz)
				}

				if storageZoneName == *sz.Name {
					return &resource.UnexpectedStateError{
						State:         "exists",
						ExpectedState: []string{"not exists"},
					}

				}

				if !*storagezones.HasMoreItems {
					return nil
				}

				page++
			}
		}
	}
}

func TestAccStorageZone_basic(t *testing.T) {
	attrs := storageZoneWanted{
		TerraformResourceName: "bunny_storagezone.mytest1",
		Name:                  randResourceName(),
		Region:                "DE",
	}

	tf := fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
}
`,
		attrs.Name,
		attrs.Region,
	)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check:  checkBasicStorageZoneAPIState(&attrs),
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
		CheckDestroy: checkStorageZoneNotExists(attrs.Name),
	})
}

func TestAccStorageZone_full(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_storagezone." + resourceName

	// set fields to different values then their defaults, to be able to test if the settings are applied
	attrs := bunny.StorageZone{
		Name:               ptr.ToString(randResourceName()),
		Region:             ptr.ToString("DE"),
		ReplicationRegions: []string{"NY", "LA"},
	}

	tf := fmt.Sprintf(`
resource "bunny_storagezone" "%s" {
	name = "%s"
	region = "%s"
	replication_regions = %s
	origin_url = "%s"
	custom_404_file_path = "%s"
	rewrite_404_to_200 = %t
}
`,
		resourceName,

		ptr.GetString(attrs.Name),
		ptr.GetString(attrs.Region),
		tfStrList(attrs.ReplicationRegions),
		"http://terraform.io",
		"/custom_404.html",
		false,
	)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check:  checkSzState(t, fullResourceName, &attrs),
			},
			{
				ResourceName:            fullResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_404_file_path", "origin_url", "rewrite_404_to_200"},
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
		CheckDestroy: checkStorageZoneNotExists(fullResourceName),
	})
}

func TestChangingImmutableFieldsFails(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_storagezone." + resourceName
	storageZoneName := randResourceName()

	attrs := bunny.StorageZone{
		Name:               ptr.ToString(storageZoneName),
		Region:             ptr.ToString("NY"),
		ReplicationRegions: []string{"DE"},
	}

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			// create storagezone
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
	replication_regions = %s
}
`,
					storageZoneName,
					*attrs.Region,
					tfStrList(attrs.ReplicationRegions),
				),
				Check: checkSzState(t, fullResourceName, &attrs),
			},
			// change region
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "LA"
	replication_regions = ["DE"]
}
`,
					storageZoneName,
				),
				ExpectError: regexp.MustCompile(".*'region' is immutable.*"),
			},
			// change name
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "LA"
	replication_regions = ["DE"]
}
`,
					storageZoneName+resource.UniqueId(),
				),
				Check:       checkSzState(t, fullResourceName, &attrs),
				ExpectError: regexp.MustCompile(".*'name' is immutable.*"),
			},
			// replace a replication_region
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "NY"
	replication_regions = ["LA"]
}
`,
					storageZoneName,
				),
				Check:       checkSzState(t, fullResourceName, &attrs),
				ExpectError: regexp.MustCompile(".*'replication_regions' can be added but not removed.*"),
			},
			// remove replication_region
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "NY"
}
`,
					storageZoneName,
				),
				Check:       checkSzState(t, fullResourceName, &attrs),
				ExpectError: regexp.MustCompile(".*'replication_regions' can be added but not removed.*"),
			},
		},
		CheckDestroy: checkStorageZoneNotExists(fullResourceName),
	})
}

func TestRegionsRequiringReplicationWithoutReplicationFails(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_storagezone." + resourceName
	storageZoneName := randResourceName()

	attrs := bunny.StorageZone{
		Name:               ptr.ToString(storageZoneName),
		Region:             ptr.ToString("SYD"),
		ReplicationRegions: []string{},
	}

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
	replication_regions = %s
}
`,
					storageZoneName,
					*attrs.Region,
					tfStrList(attrs.ReplicationRegions),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`.*"SYD" region needs to have at least one replication region.*`),
			},
		},
	})
}

func TestReplicaRegionSameAsMainFails(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_storagezone." + resourceName
	storageZoneName := randResourceName()

	attrs := bunny.StorageZone{
		Name:               ptr.ToString(storageZoneName),
		Region:             ptr.ToString("SG"),
		ReplicationRegions: []string{"SG"},
	}

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
	replication_regions = %s
}
`,
					storageZoneName,
					*attrs.Region,
					tfStrList(attrs.ReplicationRegions),
				),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`.*"SG" was specified as primary and replication region.*`),
			},
		},
	})
}

func checkSzState(t *testing.T, resourceName string, wanted *bunny.StorageZone) resource.TestCheckFunc {
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

		sz, err := clt.StorageZone.Get(context.Background(), int64(id))
		if err != nil {
			return fmt.Errorf("fetching storage-zone with id %d from api client failed: %w", id, err)
		}

		diff := szDiff(t, wanted, sz)
		if len(diff) != 0 {
			return fmt.Errorf("wanted and actual state differs:\n%s", strings.Join(diff, "\n"))
		}

		return nil

	}
}

// TestAccFailedUpdateDoesNotApplychanges tests the scenario described in
// https://github.com/5-stones/terraform-provider-bunny/pull/1#discussion_r898134629
func TestAccFailedUpdateDoesNotApplyChanges(t *testing.T) {
	attrs := storageZoneWanted{
		TerraformResourceName: "bunny_storagezone.mytest1",
		Name:                  randResourceName(),
		Region:                "DE",
	}

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
	custom_404_file_path = "/error.html"
}`,
					attrs.Name,
					attrs.Region,
				),
				Check: checkBasicStorageZoneAPIState(&attrs),
			},
			// change custom_404_file_path to a value that spawns a generation error on bunny side
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
	custom_404_file_path = "abc"
}`,
					attrs.Name,
					attrs.Region,
				),
				Check:       checkBasicStorageZoneAPIState(&attrs),
				ExpectError: regexp.MustCompile(".*updating storage zone via API failed.*"),
			},
			{
				Config: fmt.Sprintf(`
resource "bunny_storagezone" "mytest1" {
	name = "%s"
	region = "%s"
	custom_404_file_path = "/error.html"
}`,
					attrs.Name,
					attrs.Region,
				),
				PlanOnly: true,
			},
		},
		CheckDestroy: checkStorageZoneNotExists(attrs.Name),
	})
}

// storageZoneDiffIgnoredFields contains a list of fieldsnames in a bunny.StorageZone struct that are ignored by szDiff.
var storageZoneDiffIgnoredFields = map[string]struct{}{
	"ID":               {}, // computed field
	"UserID":           {}, // computed field
	"Password":         {}, // computed field
	"DateModified":     {}, // computed field
	"Deleted":          {}, // computed field
	"StorageUsed":      {}, // computed field
	"FilesStored":      {}, // computed field
	"ReadOnlyPassword": {}, // computed field

	// The following fields are tested by separate testcases and ignored in
	// storage zone testcases.
	"PullZones": {},
}

func szDiff(t *testing.T, a, b interface{}) []string {
	t.Helper()
	return diffStructs(t, a, b, storageZoneDiffIgnoredFields)
}
