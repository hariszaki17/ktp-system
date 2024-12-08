package handlers

import (
	"net/http"
	"time"

	"golang-app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetUserRequestTracker(c *gin.Context, db *gorm.DB) {
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
		return
	}

	var ktpRequest models.KTPRequest
	if err := db.Where("email = ?", email).First(&ktpRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var reqTrackers []models.RequestTracking
	if err := db.Where("request_id = ?", ktpRequest.ID).Find(&reqTrackers).Order("sequence ASC").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reqTrackers)
}

func Verify(c *gin.Context, db *gorm.DB) {
	email, exists := c.Get("email")
	if !exists || email == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please login"})
		return
	}

	var verifyReq models.VerifyReq

	// Bind JSON to the KTPRequest struct
	if err := c.ShouldBindJSON(&verifyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if verifyReq.ReqID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reqID required"})
	}

	var ktpRequest models.KTPRequest
	if err := db.Where("id = ?", verifyReq.ReqID).First(&ktpRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var lastTracker models.RequestTracking
	if err := db.Where("request_id = ?", ktpRequest.ID).First(&lastTracker).Order("sequence desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pic models.User
	if err := db.Where("email = ?", email).First(&pic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var status models.KTPRequestStatus
	newReqTracker := &models.RequestTracking{
		RequestID: ktpRequest.ID,
		UpdatedAt: time.Now().Local().String(),
		Remarks:   "",
		Sequence:  lastTracker.Sequence + 1,
		PICName:   pic.Name,
		Role:      pic.Role,
	}
	if !verifyReq.IsVerified {
		if role == "rt" {
			status = models.StatusRejectedRT
		} else if role == "rw" {
			status = models.StatusRejectedRW
		} else if role == "kelurahan" {
			status = models.StatusRejectedKel
		} else if role == "kecamatan" {
			status = models.StatusRejectedKec
		}
	} else {
		if role == "rt" {
			status = models.StatusVerifiedRT
		} else if role == "rw" {
			status = models.StatusVerifiedRW
		} else if role == "kelurahan" {
			status = models.StatusVerifiedKel
		} else if role == "kecamatan" {
			status = models.StatusKTPProcessed
		}
	}

	ktpRequest.Status = status
	newReqTracker.Status = status

	if err := db.Save(&ktpRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&newReqTracker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newReqTracker)
}
