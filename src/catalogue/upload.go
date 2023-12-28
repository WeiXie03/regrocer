package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

const COLLECTION_NAME = "catalogue"

/*
Uploads the products given in the map to a collection named the given date in the given Firestore client.
The map keys are the category names and the values are the member product structs, as defined in parse.go
Specifically, the name of the collection will be the date in the format "YYYY-MM-DD".
*/
func add_prods_as_entity(ctx context.Context, client *firestore.Client, collection_name string, date time.Time, prods_categories_map map[string][]Product) error {
	date_str := date.Format("2006-01-02")
	// Firestore is a NoSQL database => everything inside is a key : value pair
	// Have to throw the flat array as value into a map, just set key to document name
	_, err := client.Collection(collection_name).Doc(date_str).Set(ctx, prods_categories_map)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

func main() {
	fmt.Println("In main in upload.go")

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatalln("Environment variable PROJECT_ID is not set")
	}

	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	file_path := filepath.Join("..", "assets", "products.csv")
	rows, err := ReadCSV(file_path)
	if err != nil {
		log.Fatalln("Error reading CSV:", err)
		return
	}
	log.Println("Read CSV successfully")
	// fmt.Print(rows)
	formattedJSON, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		log.Fatalln("Error formatting JSON:", err)
		return
	}
	fmt.Printf("%s\n", formattedJSON)

	fmt.Println("Adding products to Firestore")
	err = add_prods_as_entity(ctx, client, COLLECTION_NAME, time.Now(), rows)
	if err != nil {
		log.Fatalln("Error adding products to Firestore:", err)
		return
	}
	fmt.Println("Added products to Firestore")
}
