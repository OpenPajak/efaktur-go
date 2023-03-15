package web_mock

var rawResponse_SptHeaderList = `{
	"status" : 0,
	"message" : "OK",
	"data" : [ {
	  "npwp" : "XXXXXXXXX6X2XXX",
	  "fgStatusRekam" : 0,
	  "fgStatusTransfer" : 1,
	  "masa1" : 2,
	  "masa2" : 2,
	  "tahun" : 2023,
	  "revNo" : 0,
	  "idNtte" : null,
	  "keterangan" : "SUKSES POSTING",
	  "tglTerima" : null,
	  "fgLbkbn" : null,
	  "nilaiLbkbn" : null
	}, {
	  "npwp" : "XXXXXXXXX6X2XXX",
	  "fgStatusRekam" : 1,
	  "fgStatusTransfer" : 1,
	  "masa1" : 1,
	  "masa2" : 1,
	  "tahun" : 2023,
	  "revNo" : 0,
	  "idNtte" : "8XXXXXXXXXXXXXXXXXX5",
	  "keterangan" : "SUKSES LAPOR",
	  "tglTerima" : "2023-02-27",
	  "fgLbkbn" : "2",
	  "nilaiLbkbn" : 99999999999
	} ]
}`

var rawResponse_SptHeaderCreate = `{
	"status" : 1,
	"message" : "OK",
	"data" : null
}`

var rawResponse_SptHeaderDelete = `{
	"status" : 1,
	"message" : "OK",
	"data" : null
}`

var rawResponse_SptHeaderCetak = `<REAL PDF CONTENT HERE>`

var rawResponse_SptDetailHeaderFind = `{
	"status" : 1,
	"message" : "OK",
	"data" : {
	  "jmlRecord" : 1,
	  "jmlDpp" : 99999999999,
	  "jmlPpn" : 10999999999,
	  "jmlPpnbm" : 0,
	  "kdForm" : "A2"
	}
}`

var rawResponse_SptDetailFind = `{
	"status" : 1,
	"message" : "OK",
	"data" : [ {
	  "namaLt" : "NAMA ORANG PRIBADI, SI",
	  "npwpPasporLt" : "XXXXXXXXX6X9XXX",
	  "nomor" : "01000XXXXXXXXXX5",
	  "tanggal" : "06/12/2022",
	  "jmlDpp" : 99999999999,
	  "jmlPpn" : 10999999999,
	  "jmlPpnbm" : 0,
	  "noRef" : null,
	  "ket" : null
	} ]
}`

// GZIP of CSV semi-colon delimiter.
var rawResponse_PrepopulatedDownload = []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 0, 77, 142, 57, 14, 195, 48, 12, 4, 251, 124, 137, 213, 58, 58, 172, 131, 52, 33, 81, 112, 201, 255, 255, 34, 74, 17, 35, 205, 98, 138, 193, 96, 19, 83, 11, 94, 163, 148, 233, 54, 32, 19, 109, 22, 74, 217, 53, 74, 206, 16, 43, 36, 23, 95, 195, 19, 154, 173, 65, 140, 9, 87, 84, 52, 50, 156, 75, 30, 254, 218, 253, 103, 137, 222, 74, 2, 6, 161, 239, 53, 239, 187, 214, 160, 84, 23, 119, 156, 30, 244, 65, 85, 249, 195, 131, 105, 31, 121, 143, 24, 138, 225, 232, 241, 245, 1, 130, 125, 90, 154, 160, 0, 0, 0}

var rawResponse_ProfilePkpGet = `{
	"status" : 0,
	"message" : null,
	"data" : [ {
	  "alamat" : "JL ASELI'D BUMI I NO 1 RT 001 RW 001, KECAMATAN, KOTA GOTHAM",
	  "idCabangWp" : null,
	  "jabatanSpt" : "DIREKTUR",
	  "klu" : "72102",
	  "kodePos" : null,
	  "kppAdm" : "6X9",
	  "masaBuku1" : "1",
	  "masaBuku2" : "12",
	  "nama" : "PT AAAAAAAAA BBBBBB CCCCCCCCC",
	  "noFax" : "031567890",
	  "noHp" : "628134567890",
	  "noTelepon" : "031567890",
	  "npwp" : "XXXXXXXXX6X9XXX",
	  "penandatanganFaktur" : "",
	  "penandatanganSpt" : "NXXXAXX",
	  "tempatPenandatangan" : "KOTA GOTHAM"
	} ]
}`

var rawResponse_ProfileSaveOrUpdate = `{
	"status" : 0,
	"message" : null,
	"data" : [ {
	  "alamat" : null,
	  "idCabangWp" : 11111111100000001111,
	  "jabatanSpt" : "DIREKTUR",
	  "klu" : null,
	  "kodePos" : null,
	  "kppAdm" : null,
	  "masaBuku1" : null,
	  "masaBuku2" : null,
	  "nama" : null,
	  "noFax" : null,
	  "noHp" : null,
	  "noTelepon" : null,
	  "npwp" : "XXXXXXXXX6X9XXX",
	  "penandatanganFaktur" : "NXXXAXX",
	  "penandatanganSpt" : "NXXXAXX",
	  "tempatPenandatangan" : null
	} ]
}`

var rawResponse_SptInfoCekSkpKp_FOUND = `{
	"status" : 1,
	"message" : "SKPKP Ditemukan",
	"data" : "SPT anda merupakan SPT dengan status pembetulan Apabila atas SPT yang dibetulkan pernah diterbitkan SKPPKP, pastikan informasi terkait SKPPKP telah diinput ke dalam SPT sesuai dengan PER-04/PJ/2021"
}`

var rawResponse_SptInfoCekSkpKp_NOTFOUND = `{
	"status" : 0,
	"message" : "SKPKP Tidak Ditemukan",
	"data" : false
}`

var rawResponse_SptIndukFind = `{
	"status" : 0,
	"message" : "OK",
	"data" : [ {
	  "attribute1" : "AAAAAAAAA BBBBBB CCCCCCCCC",
	  "attribute2" : "JL ASELI'D BUMI I NO 1 RT 001 RW 001, KECAMATAN, KOTA GOTHAM",
	  "attribute3" : "031567890",
	  "attribute4" : "628134567890",
	  "attribute5" : "72102",
	  "attribute6" : "XXXXXXXXX6X9XXX",
	  "attribute7" : 12,
	  "attribute8" : 12,
	  "attribute9" : 2022,
	  "attribute10" : "1",
	  "attribute11" : "12",
	  "attribute12" : 1,
	  "attribute13" : "satu",
	  "attribute14" : 0,
	  "attribute15" : 0,
	  "attribute16" : 450000000,
	  "attribute17" : 49500000,
	  "attribute18" : 0,
	  "attribute19" : 0,
	  "attribute20" : 0,
	  "attribute21" : 0,
	  "attribute22" : 0,
	  "attribute23" : 0,
	  "attribute24" : 450000000,
	  "attribute25" : 49500000,
	  "attribute26" : 0,
	  "attribute27" : 450000000,
	  "attribute28" : 49500000,
	  "attribute29" : 0,
	  "attribute30" : 29071013,
	  "attribute31" : 20428987,
	  "attribute32" : "20428987",
	  "attribute33" : "0",
	  "attribute34" : null,
	  "attribute35" : null,
	  "attribute36" : 0,
	  "attribute37" : 0,
	  "attribute38" : 0,
	  "attribute39" : 0,
	  "attribute40" : 0,
	  "attribute41" : 0,
	  "attribute42" : 0,
	  "attribute43" : null,
	  "attribute44" : null,
	  "attribute45" : 0,
	  "attribute46" : 0,
	  "attribute47" : 0,
	  "attribute48" : 0,
	  "attribute49" : 0,
	  "attribute50" : 0,
	  "attribute51" : 0,
	  "attribute52" : 0,
	  "attribute53" : 0,
	  "attribute54" : 0,
	  "attribute55" : null,
	  "attribute56" : null,
	  "attribute57" : 0,
	  "attribute58" : null,
	  "attribute59" : null,
	  "attribute60" : 0,
	  "attribute61" : 0,
	  "attribute62" : 0,
	  "attribute63" : "0",
	  "attribute64" : "0",
	  "attribute65" : null,
	  "attribute66" : null,
	  "attribute67" : 1,
	  "attribute68" : 1,
	  "attribute69" : 1,
	  "attribute70" : 1,
	  "attribute71" : 1,
	  "attribute72" : 1,
	  "attribute73" : 0,
	  "attribute74" : null,
	  "attribute75" : 0,
	  "attribute76" : null,
	  "attribute77" : 0,
	  "attribute78" : 0,
	  "attribute79" : null,
	  "attribute80" : null,
	  "attribute81" : "KOTA GOTHAM",
	  "attribute82" : "2023-03-09",
	  "attribute83" : 1,
	  "attribute84" : 0,
	  "attribute85" : null,
	  "attribute86" : "NXXXAXX",
	  "attribute87" : "DIREKTUR",
	  "signature" : null,
	  "certSn" : null
	} ]
}`

var rawResponse_SptLampiranAB = `{
	"status" : 0,
	"message" : "OK",
	"data" : [ {
	  "attribute1" : "AAAAAAAAA BBBBBB CCCCCCCCC",
	  "attribute2" : "XXXXXXXXX6X9XXX",
	  "attribute3" : 12,
	  "attribute4" : 12,
	  "attribute5" : 2022,
	  "attribute6" : 1,
	  "attribute7" : 0,
	  "attribute8" : 450000000,
	  "attribute9" : 49500000,
	  "attribute10" : 0,
	  "attribute11" : 0,
	  "attribute12" : 0,
	  "attribute13" : 0,
	  "attribute14" : 450000000,
	  "attribute14P" : 450000000,
	  "attribute15" : 49500000,
	  "attribute15P" : 49500000,
	  "attribute16" : 0,
	  "attribute16P" : 0,
	  "attribute17" : 0,
	  "attribute18" : 0,
	  "attribute19" : 0,
	  "attribute20" : 0,
	  "attribute21" : 0,
	  "attribute22" : 0,
	  "attribute23" : 0,
	  "attribute24" : 0,
	  "attribute25" : 0,
	  "attribute26" : 0,
	  "attribute27" : 0,
	  "attribute28" : 0,
	  "attribute29" : 264282089,
	  "attribute30" : 29071013,
	  "attribute31" : 0,
	  "attribute32" : 0,
	  "attribute33" : 0,
	  "attribute34" : 0,
	  "attribute35" : 264282089,
	  "attribute36" : 29071013,
	  "attribute37" : 0,
	  "attribute38" : 29071013,
	  "attribute39" : 0,
	  "attribute40" : null,
	  "attribute41" : null,
	  "attribute42" : 0,
	  "attribute43" : 0,
	  "attribute44" : 0,
	  "attribute45" : 29071013,
	  "signature" : null,
	  "certSn" : null
	} ]
}`

var rawResponse_SigningAgentWpPilotingStatusGet = `{
	"status" : 0,
	"message" : "Belum implementasi signing agent",
	"data" : null
}`
