package web

import (
	"context"
	"encoding/json"
	"net/http"
)

type profileClient struct {
	c *Client
}

type Profile struct {
	Npwp      string `json:"npwp"`
	Nama      string `json:"nama"`
	NoFax     string `json:"noFax"`
	NoHP      string `json:"noHp"`
	NoTelepon string `json:"noTelepon"`
	Alamat    string `json:"alamat"`
	// KodePos  any    `json:"kodePos"`

	KLU    string `json:"klu"`
	KppAdm string `json:"kppAdm"`
	// IDCabangWp int64  `json:"idCabangWp"`

	MasaBuku1 string `json:"masaBuku1"`
	MasaBuku2 string `json:"masaBuku2"`

	JabatanSpt          string `json:"jabatanSpt"`
	PenandatanganFaktur string `json:"penandatanganFaktur"`
	PenandatanganSpt    string `json:"penandatanganSpt"`
	TempatPenandatangan string `json:"tempatPenandatangan"`
}

type ProfileGetResponse struct {
	// Status  int        `json:"status"`
	// Message string     `json:"message"`
	Data []*Profile `json:"data"`
}

func (p *ProfileGetResponse) GetOne() *Profile {
	if len(p.Data) <= 0 {
		return nil
	}
	return p.Data[0]
}

// Get profile PKP
func (c *profileClient) Get(ctx context.Context) (response ProfileGetResponse, err error) {
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, EndpointProfilePkpGet, nil); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = c.c.do(req, prepareReqOptMainView); err != nil {
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return
	}

	return
}

type ProfileUpdated struct {
	JabatanSpt          string `json:"jabatanSpt"`
	Npwp                string `json:"npwp"`
	PenandatanganFaktur string `json:"penandatanganFaktur"`
	PenandatanganSpt    string `json:"penandatanganSpt"`
	// IDCabangWp          int    `json:"idCabangWp"`
}

type ProfileSaveOrUpdateRequest struct {
	PenandatanganFaktur string `json:"penandatanganFaktur"`
	PenandatanganSpt    string `json:"penandatanganSpt"`
	JabatanSpt          string `json:"jabatanSpt"`
}
type ProfileSaveOrUpdateResponse struct {
	Data []*ProfileUpdated
}

// SaveOrUpdate save or update profile PKP
func (c *profileClient) SaveOrUpdate(ctx context.Context, request ProfileSaveOrUpdateRequest) (response ProfileSaveOrUpdateResponse, err error) {
	response, err = restInvoke(c.c, restInvokeArgs[ProfileSaveOrUpdateRequest, ProfileSaveOrUpdateResponse]{
		Ctx:     ctx,
		Method:  http.MethodPost,
		Url:     EndpointProfileSaveOrUpdate,
		Request: request,
		ValidateResponse: validateResponse(
			vr_checkIfRedirect,
			"",
		),
		BodyDecoder: bodyDecoderJSON,
	})
	return
}
