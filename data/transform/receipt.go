package transform

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/vimolicious/receipt-processor/data/entities"
	"github.com/vimolicious/receipt-processor/data/models"
)

func ReceiptEntityToModel(r *entities.Receipt) (*models.Receipt, error) {
	purchaseDate := r.PurchaseDateTime.Format("2006-01-02")
	purchaseTime := r.PurchaseDateTime.Format("15:04")
	total := fmt.Sprintf("%.2f", r.Total)
	items := make([]models.Item, len(r.Items))

	for i, ri := range r.Items {
		item, err := ItemEntityToModel(&ri)
		if err != nil {
			return nil, err
		}

		items[i] = *item
	}

	receipt := models.Receipt{
		Items:        &items,
		Retailer:     &r.Retailer,
		PurchaseDate: &purchaseDate,
		PurchaseTime: &purchaseTime,
		Total:        &total,
	}

	return &receipt, nil
}

func ReceiptModelToEntity(r *models.Receipt) (*entities.Receipt, error) {
	purchaseDateTime, err := time.Parse(
		"2006-01-02 15:04",
		fmt.Sprintf("%s %s", *r.PurchaseDate, *r.PurchaseTime),
	)
	if err != nil {
		return nil, err
	}

	total, err := strconv.ParseFloat(*r.Total, 64)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	items := make([]entities.Item, len(*r.Items))
	for i, ri := range *r.Items {
		item, err := ItemModelToEntity(&ri)
		if err != nil {
			return nil, err
		}

		items[i] = *item
	}

	receipt := entities.Receipt{
		Items:            items,
		Retailer:         *r.Retailer,
		PurchaseDateTime: purchaseDateTime,
		Total:            total,
		Id:               id,
	}

	receipt.CountPoints()

	return &receipt, nil
}
