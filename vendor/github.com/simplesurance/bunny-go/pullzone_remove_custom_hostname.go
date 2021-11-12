package bunny

import (
	"context"
	"fmt"
)

// RemoveCustomHostnameOptions  represents the message that is sent to the
// Remove Custom Hostname API Endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_removehostname
type RemoveCustomHostnameOptions struct {
	// Hostname is the hostname that is removed. (Required)
	Hostname *string `json:"Hostname,omitempty"`
}

// RemoveCustomHostname removes a custom hostname from the Pull Zone.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_removehostname
func (s *PullZoneService) RemoveCustomHostname(ctx context.Context, pullZoneID int64, opts *RemoveCustomHostnameOptions) error {
	req, err := s.client.newDeleteRequest(fmt.Sprintf("pullzone/%d/removeHostname", pullZoneID), opts)
	if err != nil {
		return err
	}

	return s.client.sendRequest(ctx, req, nil)
}
