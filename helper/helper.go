package helper

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

func GenerateNIK(placeCode string, dob string) (string, error) {
	// Validate placeCode
	if len(placeCode) != 6 {
		return "", fmt.Errorf("Invalid place code")
	}

	// Parse date of birth
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return "", fmt.Errorf("Invalid date of birth")
	}

	// Format date of birth as DDMMYY
	day := fmt.Sprintf("%02d", parsedDOB.Day())
	month := fmt.Sprintf("%02d", parsedDOB.Month())
	year := fmt.Sprintf("%02d", parsedDOB.Year()%100)

	// Generate random serial number (6 digits)
	serialNumber := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Combine all parts to form the NIK
	nik := placeCode + day + month + year + serialNumber

	// Generate the checksum (for simplicity, let's use the last 4 digits of serialNumber)
	// In real applications, this might be a more complex checksum algorithm
	checksum := serialNumber[len(serialNumber)-4:]

	return nik + checksum, nil
}
