package entities

import (
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Receipt struct {
	Items            []Item
	Retailer         string
	PurchaseDateTime time.Time
	Total            float64
	Points           int
	Id               uuid.UUID
}

func CountPoints(r *Receipt) int {
	var points int

	for _, c := range r.Retailer {
		alphanumeric := unicode.IsLetter(c) || unicode.IsDigit(c)
		if alphanumeric {
			// One point for every alphanumeric character in the retailer name.
			points += 1
		}
	}

	if math.Mod(r.Total, 1.0) < 0.01 {
		// 50 points if the total is a round dollar amount with no cents.
		points += 50
	}

	if math.Mod(r.Total, 0.25) < 0.01 {
		// 25 points if the total is a multiple of 0.25.
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += len(r.Items) / 2 * 5

	for _, item := range r.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)

		if len(trimmedDesc)%3 == 0 {
			// If the trimmed length of the item description is a multiple of 3,
			// multiply the price by 0.2 and round up to the nearest integer.
			// The result is the number of points earned.
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	if r.PurchaseDateTime.Day()%2 == 1 {
		// 6 points if the day in the purchase date is odd.
		points += 6
	}

	twoPM := time.Date(
		r.PurchaseDateTime.Year(),
		r.PurchaseDateTime.Month(),
		r.PurchaseDateTime.Day(),
		14,
		0,
		0,
		0,
		r.PurchaseDateTime.Location(),
	)

	fourPM := twoPM.Add(2 * time.Hour)

	if r.PurchaseDateTime.After(twoPM) && r.PurchaseDateTime.Before(fourPM) {
		// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
		points += 10
	}

	return points
}
