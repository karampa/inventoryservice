package product

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/karampa/inventoryservice/database"
)

// var productMap = struct {
// 	sync.RWMutex
// 	m map[int]Product
// }{m: make(map[int]Product)}

// func init() {
// 	fmt.Println("Loading products...")
// 	prodMap, err := loadProductMap()
// 	productMap.m = prodMap
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%d products were loaded", len(productMap.m))
// }

// func loadProductMap() (map[int]Product, error) {
// 	filename := "products.json"
// 	_, err := os.Stat(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("file [%s] does not exist", filename)
// 	}
// 	file, _ := ioutil.ReadFile(filename)
// 	productList := make([]Product, 0)
// 	err = json.Unmarshal([]byte(file), &productList)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	prodMap := make(map[int]Product)
// 	for i := 0; i < len(productList); i++ {
// 		prodMap[productList[i].ProductID] = productList[i]
// 	}
// 	return prodMap, nil
// }

func getProduct(id int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	row := database.DBConn.QueryRowContext(ctx, `SELECT productID,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products 
	WHERE productID = ?`, id)

	var product = &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return product, nil

}

func GetTopTenProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := database.DBConn.QueryContext(ctx, `SELECT productID,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products ORDER BY quantityOnHand DESC LIMIT 10`)
	if err != nil {
		log.Println(err.Error())
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
}

func removeProduct(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DBConn.ExecContext(ctx, `DELETE FROM products WHERE productID = ?`, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	results, err := database.DBConn.QueryContext(ctx, `SELECT productID,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products`)
	if err != nil {
		log.Println(err.Error())
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

// func getProductids() []int {
// 	productMap.RLock()
// 	productids := []int{}
// 	for key := range productMap.m {
// 		productids = append(productids, key)
// 	}
// 	productMap.RUnlock()
// 	sort.Ints(productids)
// 	return productids
// }

// func getNextProductID() int {
// 	productIDs := getProductids()
// 	return productIDs[len(productIDs)-1] + 1
// }

func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	log.Println("Updating db.... ")
	_, err := database.DBConn.ExecContext(ctx, `UPDATE products SET 
			manufacturer=?,
			sku=?,
			upc=?,
			pricePerUnit=CAST(? AS DECIMAL(13,2)),
			quantityOnHand=?,
			productName=?
			WHERE productID=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DBConn.ExecContext(ctx, `INSERT INTO products ( 
			manufacturer,
			sku,
			upc,
			pricePerUnit,
			quantityOnHand,
			productName) VALUES( ?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

// func addOrUpdateProduct(product Product) (int, error) {
// 	// if the product id is set, update, otherwise add
// 	addOrUpdateID := -1
// 	if product.ProductID > 0 {
// 		oldProduct, err := getProduct(product.ProductID)
// 		if err != nil {
// 			return addOrUpdateID, err
// 		}
// 		// if it exists, replace it, otherwise return error
// 		if oldProduct == nil {
// 			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
// 		}
// 		addOrUpdateID = product.ProductID
// 	} else {
// 		addOrUpdateID = getNextProductID()
// 		product.ProductID = addOrUpdateID
// 	}
// 	productMap.Lock()
// 	productMap.m[addOrUpdateID] = product
// 	productMap.Unlock()
// 	return addOrUpdateID, nil
// }
