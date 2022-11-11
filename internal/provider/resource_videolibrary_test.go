package provider

import (
	"context"
	"regexp"

	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	ptr "github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	bunny "github.com/simplesurance/bunny-go"
)

type videoLibraryWanted struct {
	TerraformResourceName string
	bunny.VideoLibrary
	Name               string
	ReplicationRegions []string
}

func checkBasicVideoLibraryAPIState(wanted *videoLibraryWanted) resource.TestCheckFunc {
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

		vl, err := clt.VideoLibrary.Get(context.Background(), int64(id), &bunny.VideoLibraryGetOpts{IncludeAccessKey: true})
		if err != nil {
			return fmt.Errorf("fetching video library with id %d from api client failed: %w", id, err)
		}

		if err := stringsAreEqual(wanted.Name, vl.Name); err != nil {
			return fmt.Errorf("name of created videolibrary differs: %w", err)
		}

		return nil
	}
}

func checkVideoLibraryNotExists(videoLibraryName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clt := newAPIClient()

		var page int32

		for {
			videolibraries, err := clt.VideoLibrary.List(context.Background(), &bunny.PaginationOptions{
				Page:    page,
				PerPage: 1000,
			})
			if err != nil {
				return fmt.Errorf("listing videolibraries failed: %w", err)
			}

			for _, vl := range videolibraries.Items {
				if vl.Name == nil {
					return fmt.Errorf("got videolibrary from api with empty Name: %+v", vl)
				}

				if videoLibraryName == *vl.Name {
					return &resource.UnexpectedStateError{
						State:         "exists",
						ExpectedState: []string{"not exists"},
					}

				}

				if !*videolibraries.HasMoreItems {
					return nil
				}

				page++
			}
		}
	}
}

func TestAccVideoLibrary_basic(t *testing.T) {
	attrs := videoLibraryWanted{
		TerraformResourceName: "bunny_videolibrary.mytest1",
		Name:                  randResourceName(),
		ReplicationRegions:    []string{"NY", "BR"},
	}

	tf := fmt.Sprintf(`
resource "bunny_videolibrary" "mytest1" {
	name = "%s"
	replication_regions = %s
}
`,
		attrs.Name,
		tfStrList(attrs.ReplicationRegions),
	)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check:  checkBasicVideoLibraryAPIState(&attrs),
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
		CheckDestroy: checkVideoLibraryNotExists(attrs.Name),
	})
}

func TestAccVideoLibrary_full(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_videolibrary." + resourceName

	// set fields to different values then their defaults, to be able to test if the settings are applied
	attrs := bunny.VideoLibrary{
		Name:               ptr.ToString(randResourceName()),
		ReplicationRegions: []string{"NY", "BR"},
	}

	tf := fmt.Sprintf(`
resource "bunny_videolibrary" "%s" {
	name = "%s"
	replication_regions = %s
}
`,
		resourceName,

		ptr.GetString(attrs.Name),
		tfStrList(attrs.ReplicationRegions),
	)

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: tf,
				Check:  checkVlState(t, fullResourceName, &attrs),
			},
			{
				ResourceName:      fullResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  tf,
				Destroy: true,
			},
		},
		CheckDestroy: checkVideoLibraryNotExists(fullResourceName),
	})
}

func TestVideoLibraryChangingImmutableFieldsFails(t *testing.T) {
	const resourceName = "mytest1"
	const fullResourceName = "bunny_videolibrary." + resourceName
	videoLibraryName := randResourceName()

	attrs := bunny.VideoLibrary{
		Name:               ptr.ToString(videoLibraryName),
		ReplicationRegions: []string{"NY", "LA"},
	}

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			// create videolibrary
			{
				Config: fmt.Sprintf(`
resource "bunny_videolibrary" "mytest1" {
	name = "%s"
	replication_regions = %s
}
`,
					videoLibraryName,
					tfStrList(attrs.ReplicationRegions),
				),
				Check: checkVlState(t, fullResourceName, &attrs),
			},
			// replace a replication_region
			{
				Config: fmt.Sprintf(`
resource "bunny_videolibrary" "mytest1" {
	name = "%s"
	replication_regions = ["SYD"]
}
`,
					videoLibraryName,
				),
				Check:       checkVlState(t, fullResourceName, &attrs),
				ExpectError: regexp.MustCompile(".*'replication_regions' cant not be mutated.*"),
			},
			// remove replication_region
			{
				Config: fmt.Sprintf(`
resource "bunny_videolibrary" "mytest1" {
	name = "%s"
}
`,
					videoLibraryName,
				),
				Check:       checkVlState(t, fullResourceName, &attrs),
				ExpectError: regexp.MustCompile(".*'replication_regions' can not be mutated.*"),
			},
		},
		CheckDestroy: checkVideoLibraryNotExists(fullResourceName),
	})
}

func checkVlState(t *testing.T, resourceName string, wanted *bunny.VideoLibrary) resource.TestCheckFunc {
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

		vl, err := clt.VideoLibrary.Get(context.Background(), int64(id), &bunny.VideoLibraryGetOpts{IncludeAccessKey: true})
		if err != nil {
			return fmt.Errorf("fetching video library with id %d from api client failed: %w", id, err)
		}

		diff := vlDiff(t, wanted, vl)
		if len(diff) != 0 {
			return fmt.Errorf("wanted and actual state differs:\n%s", strings.Join(diff, "\n"))
		}

		return nil

	}
}

// videoLibraryDiffIgnoredFields contains a list of fieldsnames in a bunny.VideoLibrary struct that are ignored by vlDiff.
var videoLibraryDiffIgnoredFields = map[string]struct{}{
	"ID":               {}, // computed field
	"VideoCount":       {}, // computed field
	"TrafficUsage":     {}, // computed field
	"StorageUsage":     {}, // computed field
	"DateCreated":      {}, // computed field
	"APIKey":           {}, // computed field
	"ReadOnlyAPIKey":   {}, // computed field
	"HasWatermark":     {}, // computed field
	"PullZoneID":       {}, // computed field
	"StorageZoneID":    {}, // computed field
	"APIAccessKey":     {}, // computed field
	"PullZoneType":     {}, // computed field
	"AllowedReferrers": {}, // computed field
	"BlockedReferrers": {}, // computed field
}

func vlDiff(t *testing.T, a, b interface{}) []string {
	t.Helper()
	return diffStructs(t, a, b, videoLibraryDiffIgnoredFields)
}
