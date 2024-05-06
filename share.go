package proton

import (
	"context"

	"github.com/go-resty/resty/v2"
)

func (c *Client) ListShares(ctx context.Context, all bool) ([]ShareMetadata, error) {
	var res struct {
		Shares []ShareMetadata
	}

	if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		if all {
			r.SetQueryParam("ShowAll", "1")
		}

		return r.SetResult(&res).Get("/drive/shares")
	}); err != nil {
		return nil, err
	}

	return res.Shares, nil
}

func (c *Client) GetShare(ctx context.Context, shareID string) (Share, error) {
	var res struct {
		Share
	}

	if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.SetResult(&res).Get("/drive/shares/" + shareID)
	}); err != nil {
		return Share{}, err
	}

	return res.Share, nil
}

func (c *Client) GetMainShare(ctx context.Context) (Share, error) {
	volume, err := c.GetActiveVolume(ctx)
	if err != nil {
		return Share{}, err
	}

	return c.GetShare(ctx, volume.Share.ShareID)
}

// CreateShare create a new share and return the new shareID.
func (c *Client) CreateShare(ctx context.Context, volumeID string, req CreateShareReq) (CreateShareRes, error) {
	var res struct {
		Share CreateShareRes
	}

	if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.SetResult(&res).SetBody(req).Post("/drive/volumes/" + volumeID + "/shares")
	}); err != nil {
		return CreateShareRes{}, err
	}

	return res.Share, nil
}

// DeleteShare deletes a share. If force is true, any attached member or shareURLs is deleted. If force is not set to true, and the
// share has any member of shareURL attached, an error 422 with body code 2005 will be returned
func (c *Client) DeleteShare(ctx context.Context, shareID string, force bool) error {
	return c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		if force {
			r.SetQueryParam("Force", "1")
		}

		return r.Delete("/drive/shares/" + shareID)
	})
}
