package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Receipt struct {
	Items        *[]Item `json:"items"`
	Retailer     *string `json:"retailer"`
	PurchaseDate *string `json:"purchaseDate"`
	PurchaseTime *string `json:"purchaseTime"`
	Total        *string `json:"total"`
}

type ReceiptError error

var receiptStringPattern = regexp.MustCompile(`^[\w\s\-&]+$`)
var receiptTimePattern = regexp.MustCompile(`^[0-2]\d:[0-5]\d$`)
var receiptDatePattern = regexp.MustCompile(`^\d{4}\-[01]\d\-[0-3]\d$`)
var receiptPricePattern = regexp.MustCompile(`^\d+\.\d{2}$`)

func (r *Receipt) UnmarshalJSON(b []byte) error {
	type RawReceipt Receipt
	var parsedReceipt RawReceipt

	if err := json.Unmarshal(b, &parsedReceipt); err != nil {
		return err
	}

	// Check for missing fields
	missingFields := make([]string, 0, 5)

	if parsedReceipt.Items == nil {
		missingFields = append(missingFields, "items")
	}

	if parsedReceipt.Retailer == nil {
		missingFields = append(missingFields, "retailer")
	}

	if parsedReceipt.PurchaseDate == nil {
		missingFields = append(missingFields, "purchaseDate")
	}

	if parsedReceipt.PurchaseTime == nil {
		missingFields = append(missingFields, "purchaseTime")
	}

	if parsedReceipt.Total == nil {
		missingFields = append(missingFields, "total")
	}

	if len(missingFields) > 0 {
		missingFieldsList := strings.Join(missingFields, ", ")

		return ReceiptError(fmt.Errorf("missing fields: %s", missingFieldsList))
	}

	// Check regular expressions
	invalidFields := make([]string, 0, 5)

	if !receiptStringPattern.MatchString(*parsedReceipt.Retailer) {
		invalidFields = append(invalidFields, "retailer")
	}

	if !receiptDatePattern.MatchString(*parsedReceipt.PurchaseDate) {
		invalidFields = append(invalidFields, "purchaseDate")
	}

	if !receiptTimePattern.MatchString(*parsedReceipt.PurchaseTime) {
		invalidFields = append(invalidFields, "purchaseTime")
	}

	if !receiptPricePattern.MatchString(*parsedReceipt.Total) {
		invalidFields = append(invalidFields, "total")
	}

	if len(invalidFields) > 0 {
		invalidFieldsList := strings.Join(missingFields, ", ")

		return ReceiptError(fmt.Errorf("invalid fields: %s", invalidFieldsList))
	}

	*r = Receipt(parsedReceipt)

	return nil
}
