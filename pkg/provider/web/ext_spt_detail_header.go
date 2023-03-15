package web

import (
	"context"
	"net/http"
)

type sptDetailHeaderClient struct {
	c *Client
}

type SptDetailHeaderFindRequest struct {
	MasaPajak   int         `json:"masaPajak"`
	TahunPajak  int         `json:"tahunPajak"`
	RevisionNo  int         `json:"revNo"`
	KodeFormSpt KodeFormSpt `json:"kdFormSpt"`
}

type SptDetailHeader struct {
	JumlahRecord int         `json:"jmlRecord"`
	JumlahDpp    int         `json:"jmlDpp"`
	JumlahPpn    int         `json:"jmlPpn"`
	JumlahPpnbm  int         `json:"jmlPPnbm"`
	KodeForm     KodeFormSpt `json:"kdForm"`
}

type SptDetailHeaderFindResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    *SptDetailHeader `json:"data"`
}

func (c *sptDetailHeaderClient) Find(ctx context.Context, request SptDetailHeaderFindRequest) (response SptDetailHeaderFindResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptDetailHeaderFindRequest, SptDetailHeaderFindResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptDetailHeaderFind,
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
