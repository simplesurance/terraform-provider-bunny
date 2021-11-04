package bunny

import "context"

const (
	// DefaultPaginationPage is the default value that is used for
	// PullZonePaginationOptions.Page if it is unset.
	DefaultPaginationPage = 1
	// DefaultPaginationPerPage is the default value that is used for
	// PullZonePaginationOptions.PerPage if it is unset.
	DefaultPaginationPerPage = 1000
)

// PullZones represents the response of the List Pull Zone API endpoint.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_index
type PullZones struct {
	Items []*PullZone `json:"Items,omitempty"`
	PullZonePaginationReply
}

// PullZonePaginationReply represents the pagination information contained in a
// Pull Zone List API endpoint response.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_index
type PullZonePaginationReply struct {
	CurrentPage  *int32 `json:"CurrentPage"`
	TotalItems   *int32 `json:"TotalItems"`
	HasMoreItems *bool  `json:"HasMoreItems"`
}

// PullZonePaginationOptions specifies optional parameters for List APIs.
type PullZonePaginationOptions struct {
	// Page the page to return
	Page int32 `url:"page,omitempty"`
	// PerPage how many entries to return per page
	PerPage int32 `url:"per_page,omitempty"`
}

// List retrieves the Pull Zones.
// If opts is nil, DefaultPaginationPerPage and DefaultPaginationPage will be used.
// if opts.Page or or opts.PerPage is < 1, the related DefaultPagination values are used.
//
// Bunny.net API docs: https://docs.bunny.net/reference/pullzonepublic_index
func (s *PullZoneService) List(ctx context.Context, opts *PullZonePaginationOptions) (*PullZones, error) {
	var res PullZones

	// Ensure that opts.Page is >=1, if it isn't bunny.net will send a
	// different response JSON object, that contains only a single
	// PullZone, without items and paginations fields.
	// Enforcing opts.page =>1 ensures that we always unmarshal into the
	// same struct.
	if opts == nil {
		opts = &PullZonePaginationOptions{
			Page:    DefaultPaginationPage,
			PerPage: DefaultPaginationPerPage,
		}
	} else {
		opts.ensureConstraints()
	}

	req, err := s.client.newGetRequest("/pullzone", opts)
	if err != nil {
		return nil, err
	}

	if err := s.client.sendRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *PullZonePaginationOptions) ensureConstraints() {
	if p.Page < 1 {
		p.Page = DefaultPaginationPage
	}

	if p.PerPage < 1 {
		p.PerPage = DefaultPaginationPerPage
	}
}
