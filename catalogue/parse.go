package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	Name  string  `firestore:"name"`
	Unit  string  `firestore:"unit,omitempty"`
	Stock float64 `firestore:"stock,omitempty"`
	// (source : price) pairs
	Prices map[string]float64 `firestore:"prices,omitempty"`
}

/*
Reads a CSV file of products in categories into a map. The keys are the category names and the values are the member product structs
The CSV file should have the following format:
- The first row is a header row with the names of the columns
- The first row of a category is a single column with the category name
Cols:
- The first column is the product name
- The remaining columns are prices from different sources
*/
func ReadCSV(file_path string) (map[string][]Product, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Price sources headers
	price_sources := lines[0][3:]

	log.Println("second line:", lines[1])

	categories := make(map[string][]Product)
	var category string
	for _, line := range lines[1:] {
		// category "headers", i.e. category name rows
		if len(line) == 1 || (strings.ToUpper(line[0]) == line[0] && line[0] != "") {
			category = line[0]
			log.Println("Category:", category)

		} else
		// otherwise ignore incomplete lines <=> no name and/or unit
		if len(line) > 1 && line[0] != "" && line[1] != "" && line[2] != "" {
			prices := make(map[string]float64, len(line[1:]))

			// iterate over columns after the product name and unit
			var price float64
			for i, raw_price := range line[3:] {
				// for now, ignore empty prices
				if raw_price == "" {
					price = 0.0
				} else {
					price, err = strconv.ParseFloat(raw_price, 64)
					if err != nil {
						return nil, err
					}
				}
				prices[price_sources[i]] = price
			}

			stock_num, err := strconv.ParseFloat(line[2], 64)
			if err != nil {
				log.Fatalln("Error converting stock from CSV (a string) to number:", err)
				return nil, err
			}
			product := Product{
				Name:   line[0],
				Unit:   line[1],
				Stock:  stock_num,
				Prices: prices,
			}

			categories[category] = append(categories[category], product)
		}
	}
	return categories, nil
}

/*
func main() {
	file_path := filepath.Join("..", "assets", "products.csv")

	rows, err := ReadCSV(file_path)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	fmt.Println("CSV Data:")
	for _, row := range rows {
		fmt.Printf("%+v\n", row)
	}
}
*/
