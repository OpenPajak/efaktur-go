package web

import (
	"fmt"
	"net/url"
)

var (
	DefaultHost = "web-efaktur.pajak.go.id"
	BaseURL     = "https://" + DefaultHost

	EndpointLogin  string
	EndpointLogout string

	// REST
	EndpointSptHeaderList        string
	EndpointSptHeaderCreate      string
	EndpointSptHeaderDelete      string
	EndpointSptHeaderCetak       string
	EndpointSptDetailFind        string
	EndpointSptDetailHeaderFind  string
	EndpointPrepopulatedDownload string
	EndpointProfilePkpGet        string
	EndpointProfileSaveOrUpdate  string
	EndpointSptIndukFind         string
	EndpointSptInfoCekSkpKp      string
	EndpointSptLampiranAbFind    string

	EndpointSigningAgentWpPilotingStatusGet string

	ebuilder endpointBuilder
)

// BuildEndpoints build endpoint urls based on current [`BaseURL`].
func BuildEndpoints() {
	ebuilder.Build()
}

func init() {
	ebuilder = endpointBuilder{BaseURL: &BaseURL}
	// WEB
	ebuilder.Bind(&EndpointLogin, "/j_spring_security_check") // POST
	ebuilder.Bind(&EndpointLogout, "/logout")                 // GET
	// REST
	ebuilder.Bind(&EndpointSptHeaderList, "/rest/sptHeader/list")               // POST
	ebuilder.Bind(&EndpointSptHeaderCreate, "/rest/sptHeader/posting")          // POST
	ebuilder.Bind(&EndpointSptHeaderDelete, "/rest/sptHeader/hapusSpt")         // POST
	ebuilder.Bind(&EndpointSptHeaderCetak, "/rest/sptHeader/cetak")             // POST
	ebuilder.Bind(&EndpointSptDetailFind, "/rest/sptDetail/find")               // POST
	ebuilder.Bind(&EndpointSptDetailHeaderFind, "/rest/sptDetailHeader/find")   // POST
	ebuilder.Bind(&EndpointPrepopulatedDownload, "/rest/prepopulated/download") // POST
	ebuilder.Bind(&EndpointProfilePkpGet, "/rest/profilPkp/getProfilPkp")       // GET
	ebuilder.Bind(&EndpointProfileSaveOrUpdate, "/rest/profilPkp/saveOrUpdate") // POST
	ebuilder.Bind(&EndpointSptIndukFind, "/rest/sptInduk/find")                 // POST
	ebuilder.Bind(&EndpointSptLampiranAbFind, "/rest/sptLampiranAB/find")       // POST

	// After posting SptHeader
	ebuilder.Bind(&EndpointSptInfoCekSkpKp, "/rest/sptInfo/cekSkpKp") // POST

	// unsure when this got triggered.
	ebuilder.Bind(&EndpointSigningAgentWpPilotingStatusGet, "/signing-agent/wp-piloting/status") // GET

	ebuilder.Build()
}

type endpointBuilder struct {
	BaseURL    *string
	bindMapper map[*string]string
}

func (e *endpointBuilder) mapper() map[*string]string {
	if e.bindMapper == nil {
		e.bindMapper = map[*string]string{}
	}
	return e.bindMapper
}

func (e *endpointBuilder) Bind(dst *string, path string) {
	oPath, exist := e.mapper()[dst]
	if exist {
		if oPath != path {
			panic(fmt.Sprintf("duplicate destination entry with different path [%s, %s]", oPath, path))
		}
	} else {
		e.mapper()[dst] = path
	}
}

func (e *endpointBuilder) Build() {
	u, err := url.Parse(*e.BaseURL)
	if err != nil {
		panic(fmt.Sprintf("rebuild endpoint: %s", err))
	}
	for dst, path := range e.mapper() {
		u.Path = path
		*dst = u.String() // build url
	}
}
