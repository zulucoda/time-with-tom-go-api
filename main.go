package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ProductID      int    `json: "productId"`
	Manufacturer   string `json: "manufacturer"`
	Sku            string `json: "sku"`
	Upc            string `json: "upc"`
	PricePerUnit   string `json: "pricePerUnit"`
	QuantityOnHand int    `json: "quantityOnHand"`
	ProductName    string `json: "productName"`
}

var productList []Product

func init() {
	productJSON := `[
  {
    "ProductID": 1,
    "Manufacturer": "Alfa Romeo",
    "Sku": "4567qHJK",
    "Upc": "234567893456",
    "PricePerUnit": "99.99",
    "QuantityOnHand": 800,
    "ProductName": "QUADRIFOGLIO badge"
  },
  {
    "ProductID": 2,
    "Manufacturer": "Alfa Romeo",
    "Sku": "44444AF",
    "Upc": "3132123123321",
    "PricePerUnit": "200.50",
    "QuantityOnHand": 400,
    "ProductName": "QUADRIFOGLIO Keys"
  },
  {
    "ProductID": 3,
    "Manufacturer": "Alfa Romeo",
    "Sku": "33333AF",
    "Upc": "8798798546546",
    "PricePerUnit": "40.44",
    "QuantityOnHand": 400,
    "ProductName": "QUADRIFOGLIO Cups"
  }
]`

	err := json.Unmarshal([]byte(productJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextId() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

func findProductByID(productID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductID == productID {
			return &product, i
		}
	}
	return nil, 0
}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler; middleware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware finished; %s", time.Since(start))
	})
}

func main() {

	productListHandler := http.HandlerFunc(productsHandler)
	productItemHandler := http.HandlerFunc(productHandler)

	http.Handle("/products", middlewareHandler(productListHandler))
	http.Handle("/products/", middlewareHandler(productItemHandler))
	// port, ServeMux nil is the default
	http.ListenAndServe(":5000", nil)
}
