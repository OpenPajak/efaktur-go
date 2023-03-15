package web

import (
	"context"
	"io"
	"net/http"
)

type sptHeaderClient struct {
	c *Client
}

type SptHeader struct {
	Npwp string `json:"npwp"`

	FgStatusRekam    int `json:"fgStatusRekam"`
	FgStatusTransfer int `json:"fgStatusTransfer"`

	Masa1      int `json:"masa1"`
	Masa2      int `json:"masa2"`
	Tahun      int `json:"tahun"`
	RevisionNo int `json:"revNo"`

	// Known value:
	// - "SUKSES POSTING"
	// - "SUKSES LAPOR"
	Keterangan string `json:"keterangan"`

	IDNpwpTandaTanganElektronik *string `json:"idNtte"`
	TglTerima                   *string `json:"tglTerima"`
	FgLbkbn                     *string `json:"fgLbkbn"`
	NilaiLbkbn                  *int    `json:"nilaiLbkbn"`
}

type SptHeaderList []*SptHeader

type SptHeaderListRequest struct {
	TahunPajak int `json:"tahunPajak"`
}

type SptHeaderListResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    SptHeaderList `json:"data"`
}

// List all SPT header
func (c *sptHeaderClient) List(ctx context.Context, request SptHeaderListRequest) (response SptHeaderListResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptHeaderListRequest, SptHeaderListResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptHeaderList,
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
	// No way to check if this return success: Status == 0 but Message == "OK"
	return
}

type SptHeaderCreateRequest struct {
	Masa1      int `json:"masa1"`
	RevisionNo int `json:"revNo"`
	Tahun      int `json:"tahun"`
}
type SptHeaderCreateResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	// Data "data" always null. Presumably because this is an async RPC call, or EDA
	// where the action is queued for creation and no more than enqueue status
	// can be retrieved to the client.
	// Comment this since we don't really need it.
	// Data any `json:"data"`
}

// Create new SPT header (known as "Posting SPT").
func (c *sptHeaderClient) Create(ctx context.Context, request SptHeaderCreateRequest) (response SptHeaderCreateResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptHeaderCreateRequest, SptHeaderCreateResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptHeaderCreate,
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

	// { "status" : 1, "message" : "OK", "data" : null }
	if response.Status != 1 {
		err = ErrUnsuccessfulAction
		return
	}

	return
}

type SptHeaderDeleteRequest struct {
	SptHeader
}

type SptHeaderDeleteResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`

	// Data "data" always null. Presumably because this is an async RPC call, or EDA
	// where the action is queued for deleteion and no more than enqueue status
	// can be retrieved to the client.
	// Comment this since we don't really need it.
	// Data any `json:"data"`
}

// Delete existing SPT header that's not yet been submitted (aka. "Lapor").
func (c *sptHeaderClient) Delete(ctx context.Context, request SptHeaderDeleteRequest) (response SptHeaderDeleteResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[SptHeaderDeleteRequest, SptHeaderDeleteResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointSptHeaderDelete,
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

	// { "status" : 1, "message" : "OK", "data" : null }
	if response.Status != 1 {
		err = ErrUnsuccessfulAction
		return
	}

	return
}

type SptHeaderCetakRequest struct {
	MasaPajak  int `json:"masaPajak"`
	TahunPajak int `json:"tahunPajak"`
	RevisionNo int `json:"revNo"`
}

// Cetak cetak SPT Induk remotely (responded with PDF buffer)
func (c *sptHeaderClient) Cetak(ctx context.Context, request SptHeaderCetakRequest) (reader io.Reader, err error) {
	if reader, err = restInvoke(c.c, restInvokeArgs[SptHeaderCetakRequest, io.Reader]{
		Ctx:           ctx,
		Method:        http.MethodPost,
		Url:           EndpointSptHeaderCetak,
		Request:       request,
		DontCloseBody: true,
		ValidateResponse: validateResponse(
			vr_checkIfRedirect|
				vr_checkContentType|vr_checkContentDisposition_Attachment_Filename,
			"application/x-pdf",
		),
		BodyDecoder: func(reader io.Reader, response any) error {
			*(response.(*io.Reader)) = reader
			return nil
		},
	}); err != nil {
		return
	}

	return
}
