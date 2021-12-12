package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/karampa/inventoryservice/database"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("Loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products were loaded", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	filename := "products.json"
	_, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("file [%s] does not exist", filename)
	}
	file, _ := ioutil.ReadFile(filename)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProduct(id int) *Product {
	productMap.RLock()
	defer productMap.Unlock()
	if p, ok := productMap.m[id]; ok {
		return &p
	}
	return nil
}

func removeProduct(id int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, id)
}

func getProductList() ([]Product, error) {
	results, err := database.DBConn.Query(`SELECT productID,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)

	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		products = append(products, product)
	}
	return products, nil
	// for _, value := range productMap.m {
	// 	products = append(products, value)
	// }
	// productMap.RUnlock()
	// return products
}

func getProductids() []int {
	productMap.RLock()
	productids := []int{}
	for key := range productMap.m {
		productids = append(productids, key)
	}
	productMap.RUnlock()
	sort.Ints(productids)
	return productids
}

func getNextProductID() int {
	productIDs := getProductids()
	return productIDs[len(productIDs)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	// if the product id is set, update, otherwise add
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		// if it exists, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
