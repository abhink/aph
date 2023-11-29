package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

// registerer is an interface that must be satisfied by a type providing register functionality.
// This is a generic type bounded by unsiged or floating numeric types.
type registerer[K constraints.Unsigned | constraints.Float] interface {
	Put(ctx context.Context, ds []K) error
	Withdraw(ctx context.Context, amount int) ([]K, error)
}

// RequestData is used to decode incoming request data.
type RequestData struct {
	TransactionID     *string        `json:"transaction_id"`
	Denominations     []denomination `json:"denominations"`
	FinishTransaction bool           `json:"finish_transaction"`
	TotalPrice        int            `json:"total_price"`
	// items selected, etc.
}

// Transaction holds state of a single transaction.
// More information about transactions in README.
type Transaction struct {
	ID            string         `json:"id"`
	TotalPrice    int            `json:"total_price"`
	TotalInserted []denomination `json:"total_inserted"`
	TotalReturned []denomination `json:"total_returned"`
	Completed     bool           `json:"completed"`
	// other fields that track selected items
}

func main() {
	m := map[denomination]int{
		cent:        10,
		twentyCents: 10,
		fiftyCents:  10,
		one:         10,
		two:         10,
		five:        10,
		ten:         10,
		twenty:      10,
		fifty:       10,
	}

	reg := NewRegister(m, orderedDenominations)
	p := paymentProcessor{
		trStore:  make(map[string]*Transaction),
		register: reg,
	}// normally, I would be using an existing library here
	http.HandleFunc("/transaction", p.transactionHandler())

	// for authentication, middleware based routing can be used as a simple yet effective solution,
	// e.g. access to admin panel can be guarded by another middleware handler
	//
	// secureEndpoint(http.HandleFunc("/admin", p.admin()), ...)
	//
	// with `secureEndpoint` being a function that chains the handler func provided to it as an argument.

	// Start the server and listen on port 8080
	// todo: create custom server with tweaked timeouts. This typically looks like:
	// server := &http.Server{
	// 	Addr:              cfg.Address,
	// 	Handler:           router,
	// 	ReadTimeout:       cfg.ReadTimeout,
	// 	ReadHeaderTimeout: cfg.ReadTimeout,
	// 	WriteTimeout:      cfg.WriteTimeout,
	// 	IdleTimeout:       cfg.IdleTimeout,
	// }
	port := 8080
	fmt.Printf("Server is listening on port %d...\n", port)

	// todo: implement listener for OS signals, graceful shutdown
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	// Check for errors when starting the server
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

// paymentProcessor is the type that processes transaction requests
// this type is for demo only, NOT concurrent access safe
type paymentProcessor struct {
	trStore  map[string]*Transaction
	register registerer[denomination]
	// locks and other bits
}

// transactionHandler is a method that returns the handler for managing transaction.
// todo: refactor for better readability
func (p *paymentProcessor) transactionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    // ok to work on both POST and PATCH
		if r.Method != http.MethodPost && r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var requestData RequestData
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}
		var transaction *Transaction

		if requestData.TransactionID != nil {
			tr, ok := p.trStore[*requestData.TransactionID]
			if !ok {
				http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
				return
			}
			tr.TotalInserted = append(tr.TotalInserted, requestData.Denominations...)
			w.WriteHeader(http.StatusOK)
			transaction = tr
		} else {
			newID := uuid.New().String()

      
			// todo: 
			//   1. implement a more robust store
			//   2. put this store behind an interface
			p.trStore[newID] = &Transaction{
				ID:            newID,
				TotalInserted: requestData.Denominations,
				TotalPrice:    0,
				Completed:     false,
			}
      w.WriteHeader(http.StatusCreated)
			transaction = p.trStore[newID]
		}

		if requestData.FinishTransaction {
			transaction.Completed = true
			transaction.TotalReturned, err = processPaymentInCents(context.Background(), p.register, requestData.TotalPrice, transaction.TotalInserted)
			if err != nil {
				if errors.Is(err, ErrInsufficientFunds) || errors.Is(err, ErrRegisterEmpty) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		_ = json.NewEncoder(w).Encode(transaction)
	}
}
