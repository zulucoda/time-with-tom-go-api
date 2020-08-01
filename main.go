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
		return

	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductID = getNextId()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {

	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, listItemIndex := findProductByID(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJSON)
		return

	case http.MethodPut:
		var updateProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updateProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updateProduct
		productList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
