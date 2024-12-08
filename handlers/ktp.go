package handlers

import (
	"net/http"
	"time"

	"golang-app/helper"
	"golang-app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateKTP(c *gin.Context, db *gorm.DB) {
	var ktpReq models.CreateKTPReq

	if err := c.ShouldBindJSON(&ktpReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get email from context
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
		return
	}

	nik, err := helper.GenerateNIK(ktpReq.PlaceCode, ktpReq.TanggalLahir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ktp := &models.KTP{
		NIK:              nik,
		NamaLengkap:      ktpReq.NamaLengkap,
		TanggalLahir:     ktpReq.TanggalLahir,
		Alamat:           ktpReq.Alamat,
		JenisKelamin:     ktpReq.JenisKelamin,
		Agama:            ktpReq.Agama,
		StatusPerkawinan: ktpReq.StatusPerkawinan,
		Email:            ktpReq.Email,
	}
	if err := db.Create(&ktp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	laporanKTP := &models.LaporanPembuatanKTP{
		NoKTP:       nik,
		NamaWarga:   ktp.NamaLengkap,
		AlamatWarga: ktp.Alamat,
		TglLaporan:  time.Now().Local().String(),
	}
	if err := db.Create(&laporanKTP).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, laporanKTP)
}


func GetKTP(c *gin.Context, db *gorm.DB) {
	// Get email from context
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
		return
	}

	var ktp models.KTP
	if err := db.Where("email = ?", email).First(&ktp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ktp)
}
