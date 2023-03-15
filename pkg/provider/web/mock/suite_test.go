package web_mock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/OpenPajak/efaktur-go/pkg/provider/web"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrOf[T any](value T) *T {
	return &value
}

func testNewClient(t *testing.T) *web.Client {
	// tcert, err := web.PKCS12ToTLSCertificateFromFile("../fix.p12", "")
	// require.NoError(t, err)
	client, err := web.NewClient(web.ClientOptions{
		Transport: http.DefaultTransport,
	})
	require.NoError(t, err)
	return client
}

func matcherBodyJson[T any](handler func(data T) (valid bool, err error)) func(req *http.Request, rg *gock.Request) (valid bool, err error) {
	return func(req *http.Request, rg *gock.Request) (valid bool, err error) {
		var data T
		var b bytes.Buffer
		body := io.TeeReader(req.Body, &b)
		if err = json.NewDecoder(body).Decode(&data); err != nil {
			return
		}
		print(b.String())
		return handler(data)
	}
}

func matcherBodyFormData(handler func(values url.Values) (valid bool, err error)) func(req *http.Request, rg *gock.Request) (valid bool, err error) {
	return func(req *http.Request, rg *gock.Request) (valid bool, err error) {
		var query []byte
		var b bytes.Buffer
		body := io.TeeReader(req.Body, &b)
		if query, err = io.ReadAll(body); err != nil {
			return
		}
		print(b.String())

		qs, err := url.ParseQuery(string(query))
		if err != nil {
			return
		}

		return handler(qs)
	}
}

func TestMain(m *testing.M) {
	defer gock.Off()
	defer gock.RestoreClient(http.DefaultClient)

	// using http scheme.
	web.DefaultHost = "efaktur-web-mock-suite"
	web.BaseURL = "http://" + web.DefaultHost
	// build endpoint with mock base URL
	web.BuildEndpoints()

	println("===== START =====")
	m.Run()
	println("===== END =====")
}

func TestClientLogin(t *testing.T) {
	check := func(t *testing.T) (err error) {
		client := testNewClient(t)
		err = client.Login(context.Background(), "ASDFC")
		return
	}

	t.Run("success", func(t *testing.T) {
		defer gock.Off()
		gock.New(web.EndpointLogin).
			MatchType("application/x-www-form-urlencoded").
			Post("/").
			AddMatcher(matcherBodyFormData(func(values url.Values) (valid bool, err error) {
				valid = values.Get("j_password") == "ASDFC"
				return
			})).
			Reply(http.StatusFound).
			SetHeader("Location", web.BaseURL+"/")
		assert.NoError(t, check(t))

	})
	t.Run("fail", func(t *testing.T) {
		defer gock.Off()
		gock.New(web.EndpointLogin).
			MatchType("application/x-www-form-urlencoded").
			Post("/").
			Reply(http.StatusFound).
			SetHeader("Location", web.BaseURL+"/login")
		assert.Error(t, check(t))
	})
}

func TestClientLogout(t *testing.T) {
	check := func(t *testing.T) (err error) {
		client := testNewClient(t)
		err = client.Logout(context.Background())
		return
	}
	t.Run("success", func(t *testing.T) {
		defer gock.Off()
		gock.New(web.EndpointLogout).
			Get("/").
			Reply(http.StatusFound).
			SetHeader("Location", web.BaseURL+"/login")
		assert.NoError(t, check(t))
	})
	t.Run("fail", func(t *testing.T) {
		defer gock.Off()
		gock.New(web.EndpointLogout).
			Get("/").
			Reply(http.StatusFound).
			SetHeader("Location", web.BaseURL+"/")
		assert.Error(t, check(t))
	})
}

func TestClientSptHeaderList(t *testing.T) {
	defer gock.Off()
	req := web.SptHeaderListRequest{
		TahunPajak: 2023,
	}
	gock.New(web.EndpointSptHeaderList).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptHeaderListRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptHeaderList)

	client := testNewClient(t)
	resp, err := client.SptHeader.List(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, web.SptHeaderListResponse{
		Status:  0,
		Message: "OK",
		Data: []*web.SptHeader{
			{
				Npwp:                        "XXXXXXXXX6X2XXX",
				FgStatusRekam:               0,
				FgStatusTransfer:            1,
				Masa1:                       2,
				Masa2:                       2,
				Tahun:                       2023,
				RevisionNo:                  0,
				Keterangan:                  "SUKSES POSTING",
				IDNpwpTandaTanganElektronik: nil,
				TglTerima:                   nil,
				FgLbkbn:                     nil,
				NilaiLbkbn:                  nil,
			},
			{
				Npwp:                        "XXXXXXXXX6X2XXX",
				FgStatusRekam:               1,
				FgStatusTransfer:            1,
				Masa1:                       1,
				Masa2:                       1,
				Tahun:                       2023,
				RevisionNo:                  0,
				Keterangan:                  "SUKSES LAPOR",
				IDNpwpTandaTanganElektronik: ptrOf("8XXXXXXXXXXXXXXXXXX5"),
				TglTerima:                   ptrOf("2023-02-27"),
				FgLbkbn:                     ptrOf("2"),
				NilaiLbkbn:                  ptrOf(99999999999),
			},
		},
	}, resp)

}

func TestClientSptHeaderCreate(t *testing.T) {
	defer gock.Off()
	req := web.SptHeaderCreateRequest{
		Masa1:      2,
		RevisionNo: 0,
		Tahun:      2023,
	}

	gock.New(web.EndpointSptHeaderCreate).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptHeaderCreateRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptHeaderCreate)

	client := testNewClient(t)
	resp, err := client.SptHeader.Create(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, web.SptHeaderCreateResponse{
		Status:  1,
		Message: "OK",
	}, resp)
}

func TestClientSptHeaderDelete(t *testing.T) {
	defer gock.Off()
	sptHeader := web.SptHeader{
		Npwp:                        "XXXXXXXXX6X2XXX",
		FgStatusRekam:               0,
		FgStatusTransfer:            1,
		Masa1:                       2,
		Masa2:                       2,
		Tahun:                       2023,
		RevisionNo:                  0,
		Keterangan:                  "SUKSES POSTING",
		IDNpwpTandaTanganElektronik: nil,
		TglTerima:                   nil,
		FgLbkbn:                     nil,
		NilaiLbkbn:                  nil,
	}

	gock.New(web.EndpointSptHeaderDelete).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptHeaderDeleteRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, web.SptHeaderDeleteRequest{
				SptHeader: sptHeader,
			})
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptHeaderDelete)

	client := testNewClient(t)
	resp, err := client.SptHeader.Delete(context.Background(), web.SptHeaderDeleteRequest{
		SptHeader: sptHeader,
	})
	assert.NoError(t, err)
	assert.Equal(t, web.SptHeaderDeleteResponse{
		Status:  1,
		Message: "OK",
	}, resp)

}

func TestClientSptHeaderCetak(t *testing.T) {
	defer gock.Off()
	req := web.SptHeaderCetakRequest{
		MasaPajak:  1,
		TahunPajak: 2023,
		RevisionNo: 0,
	}

	gock.New(web.EndpointSptHeaderCetak).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptHeaderCetakRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/x-pdf").
		SetHeader("Content-Disposition", "attachment; filename=INDUK.pdf").
		BodyString(rawResponse_SptHeaderCetak)

	client := testNewClient(t)
	reader, err := client.SptHeader.Cetak(context.Background(), req)
	assert.NoError(t, err)

	pdfContent, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, []byte(rawResponse_SptHeaderCetak), pdfContent)
}

func TestClientSptDetailHeaderFind(t *testing.T) {
	defer gock.Off()
	req := web.SptDetailHeaderFindRequest{
		MasaPajak:   12,
		TahunPajak:  2022,
		RevisionNo:  1,
		KodeFormSpt: web.KodeFormSpt_A2,
	}

	gock.New(web.EndpointSptDetailHeaderFind).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptDetailHeaderFindRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptDetailHeaderFind)

	client := testNewClient(t)
	resp, err := client.SptDetailHeader.Find(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, web.SptDetailHeaderFindResponse{
		Status:  1,
		Message: "OK",
		Data: &web.SptDetailHeader{
			JumlahRecord: 1,
			JumlahDpp:    99999999999,
			JumlahPpn:    10999999999,
			JumlahPpnbm:  0,
			KodeForm:     web.KodeFormSpt_A2,
		},
	}, resp)
}

func TestClientSptDetailFind(t *testing.T) {
	defer gock.Off()
	req := web.SptDetailFindRequest{
		MasaPajak:   12,
		TahunPajak:  2022,
		RevisionNo:  0,
		KodeFormSpt: web.KodeFormSpt_B2,
		PageNum:     1,
		PageSize:    10,
	}

	gock.New(web.EndpointSptDetailFind).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptDetailFindRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptDetailFind)

	client := testNewClient(t)
	resp, err := client.SptDetail.Find(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, web.SptDetailFindResponse{
		Status:  1,
		Message: "OK",
		Data: []*web.SptDetail{
			{
				NamaLt:       "NAMA ORANG PRIBADI, SI",
				NpwpPasporLt: "XXXXXXXXX6X9XXX",
				Nomor:        "01000XXXXXXXXXX5",
				Tanggal:      "06/12/2022",
				JumlahDpp:    99999999999,
				JumlahPpn:    10999999999,
				JumlahPpnbm:  0,
				NoRef:        nil,
				Keterangan:   nil,
			},
		},
	}, resp)
}

func TestClientPrepopulatedDownload(t *testing.T) {
	defer gock.Off()
	// {"masaPajak":"12","tahunPajak":2022,"jenisDokumen":"1"}
	req := web.PrepopulatedDownloadRequest{
		MasaPajak:    web.PrepopulatedMasaPajak_December,
		TahunPajak:   2022,
		JenisDokumen: web.PrepopulatedJenisDokumen_FPM,
	}

	gock.New(web.EndpointPrepopulatedDownload).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.PrepopulatedDownloadRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		SetHeader("Content-Type", "application/zip").
		SetHeader("Content-Disposition", "[attachment; filename=FPM-XXXXXXXXXXXXXXX202212.zip]").
		Body(bytes.NewReader(rawResponse_PrepopulatedDownload))

	client := testNewClient(t)
	csvReader, err := client.Prepopulated.Download(context.Background(), req)
	assert.NoError(t, err)

	contents, err := csvReader.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t,
		[][]string{
			{"FM", "KD_JENIS_TRANSAKSI", "FG_PENGGANTI", "NOMOR_FAKTUR", "MASA_PAJAK", "TAHUN_PAJAK", "TANGGAL_FAKTUR", "NPWP", "NAMA", "ALAMAT_LENGKAP", "JUMLAH_DPP", "JUMLAH_PPN", "JUMLAH_PPNBM", "IS_CREDITABLE"},
		},
		contents)
}

func TestClientProfilePkpGet(t *testing.T) {
	defer gock.Off()
	gock.New(web.EndpointProfilePkpGet).
		Get("/").
		Reply(http.StatusOK).
		BodyString(rawResponse_ProfilePkpGet)

	client := testNewClient(t)
	resp, err := client.Profile.Get(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, web.ProfileGetResponse{
		Data: []*web.Profile{
			{
				Npwp:                "XXXXXXXXX6X9XXX",
				Nama:                "PT AAAAAAAAA BBBBBB CCCCCCCCC",
				NoFax:               "031567890",
				NoHP:                "628134567890",
				NoTelepon:           "031567890",
				Alamat:              "JL ASELI'D BUMI I NO 1 RT 001 RW 001, KECAMATAN, KOTA GOTHAM",
				KLU:                 "72102",
				KppAdm:              "6X9",
				MasaBuku1:           "1",
				MasaBuku2:           "12",
				JabatanSpt:          "DIREKTUR",
				PenandatanganFaktur: "",
				PenandatanganSpt:    "NXXXAXX",
				TempatPenandatangan: "KOTA GOTHAM",
			},
		},
	}, resp)
}

func TestClientProfileSaveOrUpdate(t *testing.T) {
	defer gock.Off()
	req := web.ProfileSaveOrUpdateRequest{
		PenandatanganFaktur: "NXXXAXX",
		PenandatanganSpt:    "NXXXAXX",
		JabatanSpt:          "DIREKTUR",
	}

	gock.New(web.EndpointProfileSaveOrUpdate).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.ProfileSaveOrUpdateRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_ProfileSaveOrUpdate)

	client := testNewClient(t)
	resp, err := client.Profile.SaveOrUpdate(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, web.ProfileSaveOrUpdateResponse{
		Data: []*web.ProfileUpdated{
			{
				JabatanSpt:          "DIREKTUR",
				Npwp:                "XXXXXXXXX6X9XXX",
				PenandatanganFaktur: "NXXXAXX",
				PenandatanganSpt:    "NXXXAXX",
			},
		},
	}, resp)
}

func TestClientSptInfoCekSkpKp(t *testing.T) {
	req := web.SptInfoSkpKpCekRequest{
		Masa:       1,
		Tahun:      2023,
		RevisionNo: 1,
	}

	check := func(t *testing.T) error {
		client := testNewClient(t)
		resp, err := client.SptInfo.CekSkpKp(context.Background(), req)
		if err != nil {
			return err
		}
		if resp.Status == 0 {
			return web.ErrUnsuccessfulAction
		}
		return nil
	}
	jMatcher := matcherBodyJson(func(data web.SptInfoSkpKpCekRequest) (valid bool, err error) {
		valid = reflect.DeepEqual(data, req)
		return
	})

	t.Run("found", func(t *testing.T) {
		defer gock.Off()
		gock.New(web.EndpointSptInfoCekSkpKp).
			Post("/").
			AddMatcher(jMatcher).
			Reply(http.StatusOK).
			BodyString(rawResponse_SptInfoCekSkpKp_FOUND)
		assert.NoError(t, check(t))
	})
	t.Run("notfound", func(t *testing.T) {
		defer gock.Off()
		gock.New(web.EndpointSptInfoCekSkpKp).
			Post("/").
			AddMatcher(jMatcher).
			Reply(http.StatusOK).
			BodyString(rawResponse_SptInfoCekSkpKp_NOTFOUND)
		assert.Error(t, check(t))
	})
}

func TestClientSptIndukFind(t *testing.T) {
	defer gock.Off()
	req := web.SptIndukFindRequest{
		MasaPajak:  12,
		TahunPajak: 2022,
		RevisionNo: 1,
	}

	gock.New(web.EndpointSptIndukFind).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptIndukFindRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptIndukFind)

	client := testNewClient(t)
	resp, err := client.SptInduk.Find(context.Background(), req)
	assert.NoError(t, err)

	// TODO: check this
	_ = resp
}

func TestClientSptLampiranAB(t *testing.T) {
	defer gock.Off()
	req := web.SptLampiranAbFindRequest{
		MasaPajak:  12,
		TahunPajak: 2022,
		RevisionNo: 1,
	}

	gock.New(web.EndpointSptLampiranAbFind).
		Post("/").
		AddMatcher(matcherBodyJson(func(data web.SptLampiranAbFindRequest) (valid bool, err error) {
			valid = reflect.DeepEqual(data, req)
			return
		})).
		Reply(http.StatusOK).
		BodyString(rawResponse_SptLampiranAB)

	client := testNewClient(t)
	resp, err := client.SptLampiranAB.Find(context.Background(), req)
	assert.NoError(t, err)

	// TODO: check this
	_ = resp
}

func TestClientSigningAgentWpPilotingStatusGet(t *testing.T) {
	defer gock.Off()
	gock.New(web.EndpointSigningAgentWpPilotingStatusGet).
		Get("/").
		MatchParams(map[string]string{
			"masa": "1", "tahun": "2023",
		}).
		Reply(http.StatusOK).
		BodyString(rawResponse_SigningAgentWpPilotingStatusGet)

	client := testNewClient(t)
	resp, err := client.SigningAgent.GetWpPilotingStatus(context.Background(), web.SigningAgentGetWpPilotingStatusRequest{
		Masa:  1,
		Tahun: 2023,
	})
	assert.NoError(t, err)

	// TODO: check this
	_ = resp
}
