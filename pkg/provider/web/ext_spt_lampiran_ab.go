package web

import (
	"context"
	"net/http"
)

type sptLampiranAbClient struct {
	c *Client
}

/*
Data types NOT yet known, use `map[string]any` for now.
*/
type SptLampiranAb map[string]any

type SptLampiranAbFindRequest struct {
	MasaPajak  int `json:"masaPajak"`
	TahunPajak int `json:"tahunPajak"`
	RevisionNo int `json:"revNo"`
}
type SptLampiranAbFindResponse struct {
	Data []SptLampiranAb
}

func (c *sptLampiranAbClient) Find(ctx context.Context, request SptLampiranAbFindRequest) (response SptLampiranAbFindResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptLampiranAbFindRequest, SptLampiranAbFindResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptLampiranAbFind,
		Request: request,
		ValidateResponse: validateResponse(
			vr_checkIfRedirect,
			"",
		),
		BodyDecoder: bodyDecoderJSON,
	})
	if err != nil {
		return
	}
	return
}
