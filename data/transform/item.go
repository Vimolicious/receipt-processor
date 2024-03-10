package transform

import (
	"fmt"
	"strconv"

	"github.com/vimolicious/receipt-processor/data/entities"
	"github.com/vimolicious/receipt-processor/data/models"
)

func ItemEntityToModel(i *entities.Item) (*models.Item, error) {
	price := fmt.Sprintf("%.2f", i.Price)

	item := models.Item{
		ShortDescription: &i.ShortDescription,
		Price:            &price,
	}

	return &item, nil
}

func ItemModelToEntity(i *models.Item) (*entities.Item, error) {
	price, err := strconv.ParseFloat(*i.Price, 64)
	if err != nil {
		return nil, err
	}

	item := entities.Item{
		ShortDescription: *i.ShortDescription,
		Price:            price,
	}

	return &item, nil
}
