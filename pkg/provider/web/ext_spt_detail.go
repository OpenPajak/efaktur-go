package web

import (
	"context"
	"net/http"
)

type sptDetailClient struct {
	c *Client
}

type SptDetailFindRequest struct {
	MasaPajak   int         `json:"masaPajak"`
	TahunPajak  int         `json:"tahunPajak"`
	RevisionNo  int         `json:"revNo"`
	KodeFormSpt KodeFormSpt `json:"kdFormSpt"`
	PageNum     int         `json:"pageNum"`
	PageSize    int         `json:"pageSize"`
}

type SptDetail struct {
	NamaLt       string `json:"namaLt"`
	NpwpPasporLt string `json:"npwpPasporLt"`
	Nomor        string `json:"nomor"`
	Tanggal      string `json:"tanggal"`
	JumlahDpp    int    `json:"jmlDpp"`
	JumlahPpn    int    `json:"jmlPpn"`
	JumlahPpnbm  int    `json:"jmlPpnbm"`
	NoRef        any    `json:"noRef"`
	Keterangan   any    `json:"ket"`
}

type SptDetailFindResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    []*SptDetail `json:"data"`
}

func (c *sptDetailClient) Find(ctx context.Context, request SptDetailFindRequest) (response SptDetailFindResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptDetailFindRequest, SptDetailFindResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptDetailFind,
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
	if response.Status != 1 {
		err = ErrUnsuccessfulAction
		return
	}
	return
}
