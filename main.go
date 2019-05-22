package main

import (
	"fmt"
	"go-card/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// Create card account
	router.HandleFunc("/api/card/new", controllers.CreateAccount).Methods("POST")

	// add funds to card
	router.HandleFunc("/api/card/addfunds", controllers.AddFunds).Methods("POST")

	// create a transaction
	router.HandleFunc("/api/transaction/new", controllers.CreateTransaction).Methods("POST")

	// capture funds
	router.HandleFunc("/api/transaction/capture", controllers.CaptureTransaction).Methods("POST")

	// Change authorised funds
	router.HandleFunc("/api/transaction/change", controllers.ChangeTransaction).Methods("POST")

	// refund transaction.
	router.HandleFunc("/api/transaction/refund", controllers.RefundTransaction).Methods("POST")

	// Show balance
	router.HandleFunc("/api/card/balance", controllers.GetBalance).Methods("POST")

	// get statement

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
