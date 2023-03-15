# efakturtool

CLI utility to help maintain, dump/retrieve data from E-Faktur service with ease

## Install

```bash
go install github.com/OpenPajak/efaktur-go/cmd/efakturtool
```

## Config

Example of a configuration file, `efaktur.toml`:

```toml
[conf]
  certPath = "./path/to/a_certificate.p12"
  certPasswd = "MY_SECURE_CERTIFICATE_P12_PASSPHRASE"
  # E-Faktur/e-Nova service password
  svcPasswd = "ABCDEFGH"
```

Configuration files always getting written back, and you'll lose your comments.

## Usage

Currently implemented feature by `efakturtool`:

- [x] Lampiran (SptDetail{,Header})
    [`FPM`, `PIB`, `PEB`, `CUKAI`, `BC40`, `BC25`, `BC27`, `BC41`]
- [x] Prepopulated data
    [`A1`, `A2`, `B1`, `B2`, `B3`]
- [x] Support multiple output
    [`wide`, `json`, `yaml`, `csv`]
- [ ] Manage SptHeader (Posting/Delete/Cetak)
- [ ] Manage PKP Profile

Contributions and feature requests are welcome.
