package main

import (
	"log"
	"net/http"

	"github.com/vimolicious/receipt-processor/api/controllers"
	"github.com/vimolicious/receipt-processor/data/repositories/inmemory"
)

func main() {
	receiptRepo := inmemory.NewInMemoryReceiptRepository()
	receiptController := controllers.NewReceiptController(receiptRepo)

	mux := http.NewServeMux()

	receiptController.AddRouteHandlers(mux)

	log.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", mux)
}
