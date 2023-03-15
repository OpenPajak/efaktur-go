package web

import (
	"context"
	"net/http"
)

type sptIndukClient struct {
	c *Client
}

type SptIndukFindRequest struct {
	MasaPajak  int `json:"masaPajak"`
	TahunPajak int `json:"tahunPajak"`
	RevisionNo int `json:"revNo"`
}

/*
Data types NOT yet known, use `map[string]any` for now.
*/
type SptInduk map[string]any

type SptIndukFindResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []SptInduk `json:"data"`
}

func (c *sptIndukClient) Find(ctx context.Context, request SptIndukFindRequest) (response SptIndukFindResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptIndukFindRequest, SptIndukFindResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptIndukFind,
		Request: request,
		ValidateResponse: validateResponse(
			vr_checkIfRedirect,
			"",
		),
		BodyDecoder: bodyDecoderJSON,
	})
	return
}
