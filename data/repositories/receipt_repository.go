package repositories

import (
	"github.com/google/uuid"
	"github.com/vimolicious/receipt-processor/data/entities"
)

type ReceiptRepository interface {
	ReceiptById(uuid.UUID) (*entities.Receipt, error)
	AddReceipt(*entities.Receipt) error
}
