package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Item struct {
	ShortDescription *string `json:"shortDescription"`
	Price            *string `json:"price"`
}

func (i *Item) UnmarshalJSON(b []byte) error {
	type RawItem Item
	var parsedItem RawItem

	if err := json.Unmarshal(b, &parsedItem); err != nil {
		return err
	}

	// Check for missing fields
	missingFields := make([]string, 0, 2)

	if parsedItem.ShortDescription == nil {
		missingFields = append(missingFields, "shortDescription")
	}

	if parsedItem.Price == nil {
		missingFields = append(missingFields, "price")
	}

	if len(missingFields) > 0 {
		missingFieldsList := strings.Join(missingFields, ", ")

		return ReceiptError(fmt.Errorf(
			"missing fields in at least one item: %s",
			missingFieldsList,
		))
	}

	// Check regular expressions
	invalidFields := make([]string, 0, 2)

	if !receiptStringPattern.MatchString(*parsedItem.ShortDescription) {
		invalidFields = append(invalidFields, "shortDescription")
	}

	if !receiptPricePattern.MatchString(*parsedItem.Price) {
		invalidFields = append(invalidFields, "price")
	}

	if len(invalidFields) > 0 {
		invalidFieldsList := strings.Join(invalidFields, ", ")

		return ReceiptError(fmt.Errorf(
			"invalid fields in at least one item: %s",
			invalidFieldsList,
		))
	}

	*i = Item(parsedItem)

	return nil
}
