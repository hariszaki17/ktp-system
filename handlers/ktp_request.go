package handlers

import (
	"fmt"
	"net/http"
	"time"

	"golang-app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllKTPRequests handles the GET request to fetch all KTP requests
func GetAllKTPRequests(c *gin.Context, db *gorm.DB) {
	status := models.StatusSubmitted
	role, exists := c.Get("role")
	if exists {
		if role == "rw" {
			status = models.StatusVerifiedRT
		} else if role == "kelurahan" {
			status = models.StatusVerifiedRW
		} else if role == "kecamatan" {
			status = models.StatusVerifiedKel
		} else {
			status = models.StatusSubmitted
		}
	}
	var ktpRequests []*models.KTPRequest
	if status != "" {
		if err := db.Where("status = ?", status).Find(&ktpRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := db.Find(&ktpRequests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	for _, ktpRequest := range ktpRequests {
		var kk models.KartuKeluarga
		if err := db.Where("no_kk = ?", ktpRequest.NoKK).First(&kk).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ktpRequest.KK = kk
	}

	c.JSON(http.StatusOK, ktpRequests)
}

func GetKTPRequests(c *gin.Context, db *gorm.DB) {
	email, _ := c.Get("email")
	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
	}

	var ktpRequest models.KTPRequest

	err := db.Debug().
		Where("email = ?", email).
		First(&ktpRequest).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var kk models.KartuKeluarga

	err = db.Debug().
		Where("no_kk = ?", ktpRequest.NoKK).
		First(&kk).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ktpRequest.KK = kk

	c.JSON(http.StatusOK, ktpRequest)
}

func CreateKTPRequest(c *gin.Context, db *gorm.DB) {
	var ktpRequest models.KTPRequest

	// Bind JSON to the KTPRequest struct
	if err := c.ShouldBindJSON(&ktpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get email from context
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
		return
	}

	ktpRequest.Email = email.(string)
	ktpRequest.Status = models.StatusSubmitted

	var exisistingKTPReq []models.KTPRequest
	if err := db.Where("email = ?", email).Find(&exisistingKTPReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(exisistingKTPReq) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot request ktp multiple times"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save to database
	if err := db.Create(&ktpRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ktpRequestTrack := &models.RequestTracking{
		RequestID: ktpRequest.ID,
		Status:    models.StatusSubmitted,
		Sequence:  1,
		Role:      user.Role,
		PICName:   user.Name,
		UpdatedAt: time.Now().Local().String(),
		Remarks:   "",
	}
	if err := db.Create(&ktpRequestTrack).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var kk models.KartuKeluarga
	if err := db.Where("no_kk = ?", ktpRequest.NoKK).First(&kk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ktpRequest.KK = kk

	c.JSON(http.StatusOK, ktpRequest)
}
