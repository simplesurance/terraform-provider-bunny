package bunny

import (
	"context"
	"fmt"
)

// PullZoneAddCustomCertificateOptions are the request parameters for the Add Custom Certificate API Endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_addcertificate
type PullZoneAddCustomCertificateOptions struct {
	Hostname       string `json:"Hostname"`
	Certificate    []byte `json:"Certificate"`
	CertificateKey []byte `json:"CertificateKey"`
}

// AddCustomCertificate represents the Add Custom Certificate API Endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_addcertificate
func (s *PullZoneService) AddCustomCertificate(ctx context.Context, pullZoneID int64, options *PullZoneAddCustomCertificateOptions) error {
	req, err := s.client.newPostRequest(fmt.Sprintf("/pullzone/%d/addCertificate", pullZoneID), options)
	if err != nil {
		return err
	}

	return s.client.sendRequest(ctx, req, nil)
}
