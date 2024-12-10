package models

import (
	"github.com/jinzhu/gorm"
)

type KTPRequestStatus string

const (
	StatusSubmitted    KTPRequestStatus = "submitted"
	StatusVerifiedRT   KTPRequestStatus = "verified_rt"
	StatusRejectedRT   KTPRequestStatus = "rejected_rt"
	StatusVerifiedRW   KTPRequestStatus = "verified_rw"
	StatusRejectedRW   KTPRequestStatus = "rejected_rw"
	StatusVerifiedKel  KTPRequestStatus = "verified_kel"
	StatusRejectedKel  KTPRequestStatus = "rejected_kel"
	StatusRejectedKec  KTPRequestStatus = "rejected_kec"
	StatusKTPProcessed KTPRequestStatus = "ktp_processed"
	StatusKTPCompleted KTPRequestStatus = "completed"
)

type KTPRequest struct {
	ID           uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	Email        string           `json:"email"`
	NamaWarga    string           `json:"namaWarga"`
	AlamatWarga  string           `json:"alamatWarga"`
	NoKK         string           `json:"noKk" gorm:"column:no_kk"` // Ensure column name matches
	TglPengajuan string           `json:"tglPengajuan"`
	Status       KTPRequestStatus `json:"status"`
	KK           KartuKeluarga    `json:"kk" gorm:"foreignKey:NoKK;references:no_kk"` // Ensure foreign key matches column name
}

type KartuKeluarga struct {
	NoKK      string `json:"noKk" gorm:"primaryKey;column:no_kk"` // Ensure column name matches
	NamaKK    string `json:"namaKk"`
	Alamat    string `json:"alamat"`
	AnggotaKK string `json:"anggotaKk"`
}

type RequestTracking struct {
	ID        uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestID uint             `json:"requestId"`
	Status    KTPRequestStatus `json:"status"`
	Role      string           `json:"role"`
	PICName   string           `json:"picName"`
	UpdatedAt string           `json:"updatedAt"`
	Remarks   string           `json:"remarks"`
	Sequence  int              `json:"sequence"`
}

type TransaksiPembuatanKTP struct {
	NoTransaksi  string `json:"noTransaksi" gorm:"primaryKey"`
	NoKTP        string `json:"noKtp"`
	NamaWarga    string `json:"namaWarga"`
	AlamatWarga  string `json:"alamatWarga"`
	TglTransaksi string `json:"tglTransaksi"`
	Status       string `json:"status"`
}

type LaporanPembuatanKTP struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	NoKTP       string `json:"noKtp"`
	NamaWarga   string `json:"namaWarga"`
	AlamatWarga string `json:"alamatWarga"`
	TglLaporan  string `json:"tglLaporan"`
}

type KTP struct {
	NIK              string `json:"nik" gorm:"primaryKey"`
	NamaLengkap      string `json:"namaLengkap"`
	TanggalLahir     string `json:"tanggalLahir"`
	Alamat           string `json:"alamat"`
	JenisKelamin     string `json:"jenisKelamin"`
	Agama            string `json:"agama"`
	StatusPerkawinan string `json:"statusPerkawinan"`
	Email            string `json:"email"`
}

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&KTPRequest{})
	db.AutoMigrate(&RequestTracking{})
	db.AutoMigrate(&KartuKeluarga{})
	db.AutoMigrate(&TransaksiPembuatanKTP{})
	db.AutoMigrate(&LaporanPembuatanKTP{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&KTP{})
}
