package web

import (
	"context"
	"net/http"
	"strconv"
)

type signingAgentClient struct {
	c *Client
}

type SigningAgentGetWpPilotingStatusRequest struct {
	Masa  int
	Tahun int
}
type SigningAgentGetWpPilotingStatusResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	// Data type is still unknown as the server constantly sending:
	// "Belum implementasi signing agent" during implementation.
	Data any `json:"data"`
}

func (c *signingAgentClient) GetWpPilotingStatus(ctx context.Context, request SigningAgentGetWpPilotingStatusRequest) (response SigningAgentGetWpPilotingStatusResponse, err error) {
	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet, EndpointSigningAgentWpPilotingStatusGet,
		nil,
	); err != nil {
		return
	}

	qs := req.URL.Query()
	qs.Set("masa", strconv.Itoa(request.Masa))
	qs.Set("tahun", strconv.Itoa(request.Tahun))
	req.URL.RawQuery = qs.Encode()

	var resp *http.Response
	if resp, err = c.c.do(req, prepareReqOptMainView); err != nil {
		return
	}
	defer resp.Body.Close()

	if err = validateResponse(vr_checkIfRedirect, "")(resp); err != nil {
		return
	}

	if err = bodyDecoderJSON(resp.Body, &response); err != nil {
		return
	}
	return
}
