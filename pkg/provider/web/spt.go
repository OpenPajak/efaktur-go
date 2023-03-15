package web

// DO NOT ADD Un/marshaller.
//go:generate go run github.com/dmarkham/enumer -type=PrepopulatedJenisDokumen -trimprefix=PrepopulatedJenisDokumen_ -sql
//go:generate go run github.com/dmarkham/enumer -type=PrepopulatedMasaPajak -trimprefix=PrepopulatedMasaPajak_ -sql
// With json Un/mashaller implementation
//go:generate go run github.com/dmarkham/enumer -type=KodeFormSpt -trimprefix=KodeFormSpt_ -json -text -yaml -sql

// Source: https://web-efaktur.pajak.go.id/app/views/prepopulated/download_csv.html?20200528
type PrepopulatedJenisDokumen int

const (
	// UNSELECTED
	PrepopulatedJenisDokumen_UNSELECTED PrepopulatedJenisDokumen = 0 // exist

	// FPM - Faktur Pajak Masukan
	PrepopulatedJenisDokumen_FPM PrepopulatedJenisDokumen = 1

	// PIB - Pemberitahuan Impor Barang
	PrepopulatedJenisDokumen_PIB PrepopulatedJenisDokumen = 2

	// PEB - Pemberitahuan Ekspor Barang
	PrepopulatedJenisDokumen_PEB PrepopulatedJenisDokumen = 3

	// CUKAI - Cukai
	PrepopulatedJenisDokumen_CUKAI PrepopulatedJenisDokumen = 4

	// BC40 - Pemberitahuan pemasukan barang asal Tempat Lain Dalam Daerah Pabean ke TPB
	// yang selanjutnya disebut BC 4.0 adalah pemberitahuan pabean untuk
	// pemasukan barang asal Tempat Lain Dalam Daerah Pabean ke TPB.
	PrepopulatedJenisDokumen_BC40 PrepopulatedJenisDokumen = 5

	// BC25 - Pemberitahuan Impor Barang dari TPB yang selanjutnya disebut dengan BC 2.5
	// adalah pemberitahuan pabean untuk pengeluaran barang impor dari TPB untuk impor untuk dipakai.
	PrepopulatedJenisDokumen_BC25 PrepopulatedJenisDokumen = 6

	// BC27 - Pemberitahuan pengeluaran barang dari TPB ke TPB lainnya yang selanjutnya disebut BC 2.7
	// adalah pemberitahuan pengeluaran barang untuk diangkut dari TPB ke TPB lainnya.
	PrepopulatedJenisDokumen_BC27 PrepopulatedJenisDokumen = 7

	// BC41 - Pemberitahuan pengeluaran barang asal Tempat Lain Dalam Daerah Pabean dari TPB
	// yang selanjutnya disebut BC 4.1 adalah pemberitahuan pabean
	// untuk pengeluaran barang asal Tempat Lain Dalam Daerah Pabean dari TPB.
	PrepopulatedJenisDokumen_BC41 PrepopulatedJenisDokumen = 8
)

// Source: https://web-efaktur.pajak.go.id/app/views/prepopulated/download_csv.html?20200528
type PrepopulatedMasaPajak int

const (
	PrepopulatedMasaPajak_UNSELECTED PrepopulatedMasaPajak = 0 // exist
	PrepopulatedMasaPajak_January    PrepopulatedMasaPajak = 1
	PrepopulatedMasaPajak_February   PrepopulatedMasaPajak = 2
	PrepopulatedMasaPajak_March      PrepopulatedMasaPajak = 3
	PrepopulatedMasaPajak_April      PrepopulatedMasaPajak = 4
	PrepopulatedMasaPajak_May        PrepopulatedMasaPajak = 5
	PrepopulatedMasaPajak_June       PrepopulatedMasaPajak = 6
	PrepopulatedMasaPajak_July       PrepopulatedMasaPajak = 7
	PrepopulatedMasaPajak_August     PrepopulatedMasaPajak = 8
	PrepopulatedMasaPajak_September  PrepopulatedMasaPajak = 9
	PrepopulatedMasaPajak_October    PrepopulatedMasaPajak = 10
	PrepopulatedMasaPajak_November   PrepopulatedMasaPajak = 11
	PrepopulatedMasaPajak_December   PrepopulatedMasaPajak = 12
)

// Source: https://web-efaktur.pajak.go.id/app/views/spt/lampiran_detail.html?20201026
type KodeFormSpt int

const (
	// Unselected - unselected
	// KodeFormSpt_UNSELECTED KodeFormSpt = 0 // "0"

	// A1 - Daftar Ekspor
	KodeFormSpt_A1 KodeFormSpt = 1 // "A1"

	// A2 - PK atas Penyerahan Dalam Negeri
	KodeFormSpt_A2 KodeFormSpt = 2 // "A2"

	// B1 - PM Dapat Dikreditkan atas Impor
	KodeFormSpt_B1 KodeFormSpt = 3 // "B1"

	// B2 - PM Dapat Dikreditkan atas Perolehan Dalam Negeri
	KodeFormSpt_B2 KodeFormSpt = 4 // "B2"

	// B3 - PM Tidak Dapat Dikreditkan atau Mendapat Fasilitas
	KodeFormSpt_B3 KodeFormSpt = 5 // "B3"
)
