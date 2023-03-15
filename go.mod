module github.com/OpenPajak/efaktur-go

go 1.20

require (
	github.com/dmarkham/enumer v1.5.8
	github.com/h2non/gock v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	golang.org/x/net v0.8.0
	software.sslmate.com/src/go-pkcs12 v0.2.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/h2non/parth v0.0.0-20190131123155-b4df798d6542 // indirect
	github.com/pascaldekloe/name v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace software.sslmate.com/src/go-pkcs12 => github.com/ii64/go-pkcs12 v0.2.1-0.20220610195639-426fe5a3b19a

replace github.com/dmarkham/enumer => github.com/ii64/enumer v1.5.4
