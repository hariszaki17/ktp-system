package models

type VerifyReq struct {
	IsVerified bool `json:"isVerified"`
	ReqID      uint `json:"reqID"`
}

type CreateKTPReq struct {
	NamaLengkap      string `json:"namaLengkap"`
	TanggalLahir     string `json:"tanggalLahir"`
	Alamat           string `json:"alamat"`
	JenisKelamin     string `json:"jenisKelamin"`
	Agama            string `json:"agama"`
	StatusPerkawinan string `json:"statusPerkawinan"`
	Email            string `json:"email"`
	ReqID            uint   `json:"reqID"`
	PlaceCode        string `json:"placeCode"`
}

type CreateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
