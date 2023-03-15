package web

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"io"
	"net/http"
)

type prepopulatedClient struct {
	c *Client
}

type PrepopulatedDownloadRequest struct {
	// MasaPajak int enum encoded as string
	MasaPajak PrepopulatedMasaPajak `json:"masaPajak,string"`
	// TahunPajak int
	TahunPajak int `json:"tahunPajak"`
	// JenisDokumen int enum encoded as string
	JenisDokumen PrepopulatedJenisDokumen `json:"jenisDokumen,string"`
}

// Download prepopulated CSV
func (c *prepopulatedClient) Download(ctx context.Context, request PrepopulatedDownloadRequest) (reader *csv.Reader, err error) {
	if reader, err = restInvoke(c.c, restInvokeArgs[PrepopulatedDownloadRequest, *csv.Reader]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointPrepopulatedDownload,
		Request: request,
		// don't close body (on restInvokes' defer)
		// as we are doing streaming read from the body.
		DontCloseBody: true,
		// Validate the response
		ValidateResponse: validateResponse(
			vr_checkIfRedirect|
				vr_checkContentType|
				vr_checkContentDisposition_Attachment_Filename,
			"application/zip", // zip? Zip app capable of opening GZIP compressed data.
		),
		BodyDecoder: func(reader io.Reader, response any) (err error) {
			// response body is a GZIP'd CSV with semi-colon delimiter.
			var zr *gzip.Reader
			zr, err = gzip.NewReader(reader)
			if err != nil {
				return
			}
			csvReader := csv.NewReader(zr)
			csvReader.Comma = ';'
			*(response.(**csv.Reader)) = csvReader
			return
		},
	}); err != nil {
		return
	}

	return
}
