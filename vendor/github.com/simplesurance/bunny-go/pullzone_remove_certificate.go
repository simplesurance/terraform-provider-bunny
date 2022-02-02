package bunny

import (
	"context"
	"fmt"
)

// RemoveCertificateOptions represents the request parameters for the Remove
// Certificate API Endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_removecertificate
type RemoveCertificateOptions struct {
	Hostname *string `json:"Hostname,omitempty"`
}

// RemoveCertificate represents the Remove Certificate API Endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_removecertificate
func (s *PullZoneService) RemoveCertificate(ctx context.Context, pullZoneID int64, options *RemoveCertificateOptions) error {
	req, err := s.client.newDeleteRequest(fmt.Sprintf("/pullzone/%d/removeCertificate", pullZoneID), options)
	if err != nil {
		return err
	}

	return s.client.sendRequest(ctx, req, nil)
}
