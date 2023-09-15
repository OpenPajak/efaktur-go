package web

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"software.sslmate.com/src/go-pkcs12"
)

const (
	ctApplicationJson              = "application/json"
	ctApplicationWwwFormUrlencoded = "application/x-www-form-urlencoded"

	refererForLogin    = "https://web-efaktur.pajak.go.id/login"
	refererForMainView = "https://web-efaktur.pajak.go.id/"
)

var (
	prepareReqOptLogin = prepareRequestOptions{
		contentType: ctApplicationWwwFormUrlencoded,
		referer:     refererForLogin,
	}
	prepareReqOptMainView = prepareRequestOptions{
		contentType: "",
		referer:     refererForMainView,
	}
	prepareReqOptRestAPI = prepareRequestOptions{
		contentType: ctApplicationJson,
		referer:     refererForMainView,
	}
)

func PKCS12ToTLSCertificateFromMemory(pfxData []byte, password string) (tlsCert *tls.Certificate, clientCAs []*x509.Certificate, err error) {
	var (
		privateKey any
		cert       *x509.Certificate
		caCerts    []*x509.Certificate
	)
	if privateKey, cert, caCerts, err = pkcs12.DecodeChain(pfxData, password); err != nil {
		err = errors.Wrap(err, "decode chain")
		return
	}
	clientCAs = caCerts

	tlsCert = &tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}
	return
}

func PKCS12ToTLSCertificateFromFile(path string, password string) (cert *tls.Certificate, clientCAs []*x509.Certificate, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()

	var content []byte
	if content, err = io.ReadAll(f); err != nil {
		return
	}

	cert, clientCAs, err = PKCS12ToTLSCertificateFromMemory(content, password)
	return
}

type Client struct {
	client    *http.Client
	userAgent string

	SptHeader       *sptHeaderClient
	SptDetail       *sptDetailClient
	SptDetailHeader *sptDetailHeaderClient
	Prepopulated    *prepopulatedClient
	Profile         *profileClient
	SptInduk        *sptIndukClient
	SptLampiranAB   *sptLampiranAbClient
	SigningAgent    *signingAgentClient
	SptInfo         *sptInfoClient
}

type ClientOptions struct {
	UserAgent             string
	TLSCertificate        *tls.Certificate
	TLSClientCAs          []*x509.Certificate
	TLSInsecureSkipVerify bool

	// Transport overrides http Transport TLS configuraton
	// specified for given [`TLSCertificate`] in the option.
	// If you set this, then you'll need to add your TLS certificate
	// and private key yourself.
	Transport http.RoundTripper
}

func (opt *ClientOptions) validate() {
	if opt.UserAgent == "" {
		opt.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"
	}
}

func NewClient(opts ClientOptions) (*Client, error) {
	opts.validate()
	c := &Client{
		userAgent: opts.UserAgent,
	}

	var transport = opts.Transport
	// Transport option is not specified, configure transport TLS config.
	if transport == nil {
		certPool := x509.NewCertPool()
		for _, clientCA := range opts.TLSClientCAs {
			certPool.AddCert(clientCA)
		}
		trans := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{},
				ClientAuth:   tls.VerifyClientCertIfGiven,
				ClientCAs:    certPool,

				InsecureSkipVerify: opts.TLSInsecureSkipVerify,
			},
			// Doc: https://pkg.go.dev/net/http#pkg-overview
			// > Programs that must disable HTTP/2 can do so by setting
			// > Transport.TLSNextProto (for clients)
			// > or Server.TLSNextProto (for servers) to a non-nil,
			// > empty map.
			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		}
		transport = trans

		// Add certificate if it's specified. Enforce?
		if opts.TLSCertificate != nil {
			trans.TLSClientConfig.Certificates = append(trans.TLSClientConfig.Certificates,
				*opts.TLSCertificate,
			)
		}
	}

	// Initialize cookiejar
	cjar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	c.client = &http.Client{
		Transport: transport,
		Jar:       cjar,
	}

	// DO NOT FOLLOW REDIRECT.
	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// setup clients
	c.SptHeader = &sptHeaderClient{c}
	c.SptDetail = &sptDetailClient{c}
	c.SptDetailHeader = &sptDetailHeaderClient{c}
	c.Prepopulated = &prepopulatedClient{c}
	c.Profile = &profileClient{c}
	c.SptInduk = &sptIndukClient{c}
	c.SptLampiranAB = &sptLampiranAbClient{c}
	c.SigningAgent = &signingAgentClient{c}
	c.SptInfo = &sptInfoClient{c}

	return c, nil
}

// amd64: 8*2 * 2 strings = 32 bytes
type prepareRequestOptions struct {
	contentType string
	referer     string
}

func (c *Client) prepareRequest(req *http.Request, opts prepareRequestOptions) {
	if opts.contentType != "" {
		req.Header.Set("Content-Type", opts.contentType)
	}
	if opts.referer != "" {
		req.Header.Set("Referer", opts.referer)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Origin", "https://web-efaktur.pajak.go.id")
}

func (c *Client) do(req *http.Request, opts prepareRequestOptions) (resp *http.Response, err error) {
	c.prepareRequest(req, opts)
	resp, err = c.client.Do(req)
	return
}

func (c *Client) Login(ctx context.Context, password string) (err error) {
	var data = url.Values{}
	data.Set("j_password", password)

	var req *http.Request
	if req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost, EndpointLogin,
		strings.NewReader(data.Encode()),
	); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = c.do(req, prepareReqOptLogin); err != nil {
		return
	}
	// close the body before HTTP response passed into the error for encapsulation.
	defer resp.Body.Close()

	// check if redirect path is `/login`
	// check if it has `error` in query string, no matter what the value is
	// a successful action indicated with no `error` query string.
	if err = validateResponse(
		vr_checkIfRedirect|vr_checkIfRedirect_MustRedirect,
		"",
	)(resp); !errors.Is(err, ErrLoginRequired) {
		return
	} else {
		// we're trying to perform login, so if it's redirect back to the
		// same page then the action is failed.
		err = ErrUnsuccessfulAction
	}

	return
}

func (c *Client) Logout(ctx context.Context) (err error) {
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, EndpointLogout, nil); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = c.do(req, prepareReqOptMainView); err != nil {
		return
	}
	// close the body before HTTP response passed into the error for encapsulation.
	defer resp.Body.Close()

	// expected to redirect to /login so this will need to be ErrLoginRequired
	if err = validateResponse(
		vr_checkIfRedirect|vr_checkIfRedirect_MustRedirect,
		"",
	)(resp); !errors.Is(err, ErrLoginRequired) {
		err = ErrUnsuccessfulAction
		return
	} else {
		err = nil
	}
	return
}

type restInvokeValidateResponse func(resp *http.Response) (err error)
type restInvokeBodyDecoder func(reader io.Reader, response any) (err error)

type restInvokeArgs[Req any, Resp any] struct {
	Ctx     context.Context
	Method  string
	Url     string
	Request Req

	ValidateResponse restInvokeValidateResponse
	BodyDecoder      restInvokeBodyDecoder
	DontCloseBody    bool
}

func restInvoke[Req any, Resp any](c *Client, args restInvokeArgs[Req, Resp]) (response Resp, err error) {
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(&args.Request); err != nil {
		return
	}

	var req *http.Request
	if req, err = http.NewRequestWithContext(
		args.Ctx,
		args.Method,
		args.Url,
		&buf,
	); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = c.do(req, prepareReqOptRestAPI); err != nil {
		return
	}
	if !args.DontCloseBody {
		defer resp.Body.Close()
	}

	// response validator
	if f := args.ValidateResponse; f != nil {
		if err = f(resp); err != nil {
			return
		}
	}

	// body decoder
	if f := args.BodyDecoder; f != nil {
		if err = f(resp.Body, &response); err != nil {
			return
		}
	}

	return
}

const (
	vr_checkContentType                            = 1 << 0
	vr_checkContentDisposition_Attachment_Filename = 1 << 1
	vr_checkIfRedirect_ToLogin                     = 1 << 2
	vr_checkIfRedirect_HasErrorQueryString         = 1 << 3
	vr_checkIfRedirect_MustRedirect                = 1 << 4
	vr_checkIfRedirect                             = vr_checkIfRedirect_ToLogin | vr_checkIfRedirect_HasErrorQueryString
)

func validateResponse(
	opts int,
	expectedContentType string,
) func(resp *http.Response) (err error) {
	return func(resp *http.Response) (err error) {
		var (
			u         *url.URL
			mediaType string
			params    map[string]string
		)

		// 1
		if opts&(vr_checkIfRedirect|vr_checkIfRedirect_MustRedirect) != 0 {
			if opts&vr_checkIfRedirect_MustRedirect != 0 && resp.StatusCode != http.StatusFound {
				err = ErrUnsuccessfulAction
				return
			} else if resp.StatusCode == http.StatusFound {
				if u, err = validateRedirectionHTTPResponse(resp); err != nil {
					return
				}
				switch {
				case opts&vr_checkIfRedirect_ToLogin != 0 && strings.HasPrefix(u.Path, "/login"):
					// Possible path: "/login;jsessionid=XXXXXXX-XXXXXXXXXXXXXXXX.node1"
					err = &ErrInvalidResponse{
						resp,
						"login required",
						ErrLoginRequired,
					}
					return
				case opts&vr_checkIfRedirect_HasErrorQueryString != 0 && u.Query().Has("error"):
					err = &ErrInvalidResponse{
						resp,
						"unsuccessful action",
						ErrUnsuccessfulAction,
					}
					return
				}
			}
		}

		// 2
		if opts&vr_checkContentType != 0 {
			contentType := resp.Header.Get("Content-Type")
			if !strings.EqualFold(contentType, expectedContentType) {
				err = &ErrInvalidResponse{resp, fmt.Sprintf("mismatch ct: %s", contentType), nil}
				return
			}
		}

		// 3
		if opts&vr_checkContentDisposition_Attachment_Filename != 0 {
			contentDisposition := resp.Header.Get("Content-Disposition")
			// Trim from "[attachment; filename=FPM-XXXXXXXXXXXXXXX202303.zip]"
			contentDisposition = strings.Trim(contentDisposition, "[]")
			if mediaType, params, err = mime.ParseMediaType(contentDisposition); err != nil {
				err = &ErrInvalidResponse{
					resp,
					fmt.Sprintf("invalid content disposition: %s", contentDisposition),
					err,
				}
				return
			}
			if !strings.EqualFold(mediaType, "attachment") {
				err = &ErrInvalidResponse{
					resp,
					fmt.Sprintf("mismatch mediaType: %s", mediaType),
					nil,
				}
				return
			}

			fileName, exist := params["filename"]
			if !exist || fileName == "" {
				err = &ErrInvalidResponse{
					resp,
					"invalid mime params",
					nil,
				}
				return
			}
		}

		return
	}

}

func validateRedirectionHTTPResponse(resp *http.Response) (nextURL *url.URL, err error) {
	if resp.StatusCode != http.StatusFound {
		err = &ErrInvalidResponse{resp, "status code mismatch", nil}
		return
	}

	locationValue := resp.Header.Get("Location")
	if locationValue == "" {
		err = &ErrInvalidResponse{resp, "invalid location header", nil}
		return
	}

	if nextURL, err = url.Parse(locationValue); err != nil {
		err = &ErrInvalidResponse{resp, "parse location url", nil}
		return
	}

	// SECURITY: check redirect Host, ensure it's still the same domain name.
	if nextURL.Host != DefaultHost {
		err = &ErrInvalidResponse{resp, "redirect host mismatch", nil}
		return
	}

	return
}

// bodyDecoderJSON decode JSON from io.Reader to dst (a pointer).
func bodyDecoderJSON(reader io.Reader, dst any) (err error) {
	if err = json.NewDecoder(reader).Decode(dst); err != nil {
		return
	}
	return
}
