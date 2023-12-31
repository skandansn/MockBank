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

func RemoveDuplicatesUint(elements []uint) []uint {
	encountered := map[uint]bool{}
	result := []uint{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

func ParseStringAsFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("Error converting string to float64:", err)
		return 0
	}
	return f
}
