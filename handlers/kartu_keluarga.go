package handlers

import (
	"net/http"

	"golang-app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllKTPRequests handles the GET request to fetch all KTP requests
func GetKKByKK(c *gin.Context, db *gorm.DB) {
	noKK, _ := c.Params.Get("noKK")
	if noKK == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "noKK cannot be null"})
	}

	var kk models.KartuKeluarga
	if err := db.Where("no_kk = ?", noKK).Find(&kk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kk)
}
