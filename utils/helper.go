package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandomChar(length int) string {
	b := make([]byte, length)
	for i := range b {
		if i == 0 {
			b[i] = charset[seededRand.Intn(25)]
		} else {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
	}
	return string(b)
}

func SlugGenerator(name string) string {
	text := []byte(strings.ToLower(name))

	regE := regexp.MustCompile("[[:space:]]")
	text = regE.ReplaceAll(text, []byte("-"))

	final := string(text) + "-" + RandomChar(3)
	final = strings.ReplaceAll(final, ".", "")

	return strings.ToLower(final)
}

func PackageInvoiceGenerator() string {
	date := time.Now().Unix()
	char := RandomChar(3)

	return "P-" + strconv.FormatInt(date, 10) + char
}

func HospitalityInvoiceGenerator() string {
	date := time.Now().Unix()
	char := RandomChar(3)

	return "H-" + strconv.FormatInt(date, 10) + char
}

func DateToDatetime(date string) time.Time {
	dates := strings.Split(date, "-")
	year, _ := strconv.Atoi(dates[0])
	month, _ := strconv.Atoi(dates[1])
	day, _ := strconv.Atoi(dates[2])

	result := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())

	return result
}
