package provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	bunny "github.com/simplesurance/bunny-go"
)

func init() {
	resource.AddTestSweepers("pullzones", &resource.Sweeper{
		Name: "pullzones",
		F:    sweepPullZones,
	})
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sweepPullZones deletes all Pull Zones at the provider that have a naming starting with resourcePrefix.
func sweepPullZones(_ string) error {
	clt := newAPIClient()

	for page := int32(0); ; page++ {
		pullzones, err := clt.PullZone.List(context.Background(), &bunny.PullZonePaginationOptions{
			Page:    page,
			PerPage: 1000,
		})
		if err != nil {
			return fmt.Errorf("listing pull zones failed: %w", err)
		}

		for _, pz := range pullzones.Items {
			deletePullZone(clt, pz, resourcePrefix)
		}

		if !*pullzones.HasMoreItems {
			return nil
		}
	}
}

// deletePullZone deletes the Pull Zone if it's name starts with namePrefix.
func deletePullZone(clt *bunny.Client, pz *bunny.PullZone, namePrefix string) {
	if pz.ID == nil {
		log.Printf("ignoring pull zone with nil ID: %+v", pz)
		return
	}

	if pz.Name == nil {
		log.Printf("ignoring pull zone with nil name: %+v", pz)
		return
	}

	if !strings.HasPrefix(*pz.Name, resourcePrefix) {
		log.Printf("ignoring pull zone %d (%s) without name prefix %s", *pz.ID, *pz.Name, namePrefix)
		return
	}

	err := clt.PullZone.Delete(context.Background(), *pz.ID)
	if err != nil {
		log.Printf("deleting pull zone %d (%s) failed: %s", *pz.ID, *pz.Name, err)
		return
	}

	log.Printf("deleted pull zone %d (%s)", *pz.ID, *pz.Name)
}
