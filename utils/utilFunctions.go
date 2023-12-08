package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func ParseStringAsInt(stringId string) uint {
	stringIdUint, err := strconv.ParseUint(stringId, 10, 0)
	if err != nil {
		fmt.Println("Error converting string to uint:", err)
		return 0
	}
	return uint(stringIdUint)
}

func GenerateRandomNumberString(length int) string {
	rand.Seed(time.Now().UnixNano())

	result := ""
	for i := 0; i < length; i++ {
		digit := rand.Intn(10)
		result += fmt.Sprintf("%d", digit)
	}

	return result
}

func GenerateCardExpiry() string {
	rand.Seed(time.Now().UnixNano())

	month := rand.Intn(12) + 1
	year := rand.Intn(10) + 2020

	return fmt.Sprintf("%d/%d", month, year)
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}
