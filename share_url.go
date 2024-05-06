package proton

import (
	"context"

	"github.com/go-resty/resty/v2"
)

func (c *Client) CreateShareURL(ctx context.Context, shareID string, req CreateShareURLReq) (ShareURL, error) {
	var res struct {
		ShareURL ShareURL
	}

	if err := c.do(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.SetResult(&res).SetBody(req).Post("/drive/shares/" + shareID + "/urls")
	}); err != nil {
		return ShareURL{}, err
	}

	return res.ShareURL, nil
}
