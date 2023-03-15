package web

import (
	"context"
	"net/http"
)

type sptInfoClient struct {
	c *Client
}

type SptInfoSkpKpCekRequest struct {
	Masa       int `json:"masa"`
	Tahun      int `json:"tahun"`
	RevisionNo int `json:"revNo"`
}
type SptInfoSkpKpCekResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// CekSkpKp return [`status`] equals 1 when there's SKPKP found, and [`data`] typed as string.
// If there's no anything found, [`status`] equals 0 and [`data`] value is `false`.
// Message:
// - "SKPKP Tidak Ditemukan"
// - "SKPKP Ditemukan"
func (c *sptInfoClient) CekSkpKp(ctx context.Context, request SptInfoSkpKpCekRequest) (response SptInfoSkpKpCekResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptInfoSkpKpCekRequest, SptInfoSkpKpCekResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptInfoCekSkpKp,
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
