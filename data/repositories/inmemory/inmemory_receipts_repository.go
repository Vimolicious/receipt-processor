package inmemory

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/vimolicious/receipt-processor/data/entities"
)

type InMemoryReceiptRepository struct {
	receipts map[uuid.UUID]*entities.Receipt
	mutex    sync.RWMutex
}

func NewInMemoryReceiptRepository() *InMemoryReceiptRepository {
	inMemoryRepo := InMemoryReceiptRepository{
		receipts: make(map[uuid.UUID]*entities.Receipt),
	}
	return &inMemoryRepo
}

func (r *InMemoryReceiptRepository) ReceiptById(id uuid.UUID) (*entities.Receipt, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	receipt, ok := r.receipts[id]
	if !ok {
		return nil, fmt.Errorf("No receipt with ID \"%s\"", id)
	}

	log.Printf("Receipt with ID '%s' retrieved\n", receipt.Id)

	return receipt, nil
}

func (r *InMemoryReceiptRepository) AddReceipt(receipt *entities.Receipt) error {
	r.mutex.RLock()
	_, ok := r.receipts[receipt.Id]

	if ok {
		r.mutex.RUnlock()
		return fmt.Errorf("Receipt already exists with ID \"%s\"", receipt.Id)
	}

	r.mutex.RUnlock()

	// TODO: mitigate potential DoS if too many read locks are held
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.receipts[receipt.Id] = receipt

	log.Printf("Receipt with ID '%s' saved\n", receipt.Id)

	return nil
}
