package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/vimolicious/receipt-processor/api/middleware"
	"github.com/vimolicious/receipt-processor/data/models"
	"github.com/vimolicious/receipt-processor/data/repositories"
	"github.com/vimolicious/receipt-processor/data/transform"
)

const MAX_RECEIPT_BYTES int64 = 1 << 20 // 1 MiB

type ReceiptController struct {
	receiptRepository repositories.ReceiptRepository
}

func NewReceiptController(rr repositories.ReceiptRepository) *ReceiptController {
	newReceiptController := &ReceiptController{
		receiptRepository: rr,
	}
	return newReceiptController
}

func (rc *ReceiptController) AddRouteHandlers(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST /receipts/process",
		middleware.LogRoute(rc.processReceiptHandler),
	)
	mux.HandleFunc(
		"GET /receipts/{id}/points",
		middleware.LogRoute(rc.getPointsHandler),
	)
}

type getPointsResponse struct {
	Points int `json:"points"`
}

func (rc *ReceiptController) getPointsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	receipt, err := rc.receiptRepository.ReceiptById(id)
	if err != nil {
		http.Error(w, "No receipt found for that ID", http.StatusNotFound)
		return
	}

	res, err := json.Marshal(getPointsResponse{
		Points: receipt.Points,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

type processReceiptResponse struct {
	Id string `json:"id"`
}

func (rc *ReceiptController) processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_RECEIPT_BYTES)

	decoder := json.NewDecoder(r.Body)

	var receiptModel models.Receipt

	err := decoder.Decode(&receiptModel)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError
		_, receiptError := err.(models.ReceiptError)

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf(
				"Request body JSON has bad syntax at position %d",
				syntaxError.Offset,
			)
			http.Error(w, msg, http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			msg := fmt.Sprintf("Request body is empty")
			http.Error(w, msg, http.StatusBadRequest)

		case errors.As(err, &unmarshalError):
			msg := fmt.Sprintf(
				"Request body has invalid value for '%s' field at position %d",
				unmarshalError.Field,
				unmarshalError.Offset,
			)
			http.Error(w, msg, http.StatusBadRequest)

		case err.Error() == "http: request body too large":
			msg := "Request body is too big"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		case receiptError:
			msg := fmt.Sprintf("Receipt error: %s", err.Error())
			http.Error(w, msg, http.StatusBadRequest)

		default:
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		return
	}

	receipt, err := transform.ReceiptModelToEntity(&receiptModel)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var total float64
	for _, item := range receipt.Items {
		total += item.Price
	}

	if math.Round(receipt.Total*100) != math.Round(total*100) {
		log.Printf(
			"%f != %f",
			math.Round(receipt.Total*100),
			math.Round(total*100),
		)
		msg := "Receipt error: wrong value in 'total'"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = rc.receiptRepository.AddReceipt(receipt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(processReceiptResponse{Id: receipt.Id.String()})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/text")
	w.Write(res)
}
