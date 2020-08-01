package main

import (
	"encoding/json"
	"log"
	"net/http"
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
		{"ProductID":1,"Manufacturer":"Alfa Romeo","Sku":"4567qHJK","Upc":"234567893456","PricePerUnit":"99.99","QuantityOnHand":800,"ProductName":"QUADRIFOGLIO badge"},
		{"ProductID":2,"Manufacturer":"Alfa Romeo","Sku":"44444AF","Upc":"3132123123321","PricePerUnit":"200.50","QuantityOnHand":400,"ProductName":"QUADRIFOGLIO Keys"},
		{"ProductID":3,"Manufacturer":"Alfa Romeo","Sku":"33333AF","Upc":"8798798546546","PricePerUnit":"40.44","QuantityOnHand":400,"ProductName":"QUADRIFOGLIO Cups"}
]`

	err := json.Unmarshal([]byte(productJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	}
}

func main() {
	http.HandleFunc("/products", productsHandler)
	// port, ServeMux nil is the default
	http.ListenAndServe(":5000", nil)
}
