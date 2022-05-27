package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"inceptiondb/collection"
)

func insertItem(collections map[string]*collection.Collection) interface{} {

	type Item map[string]interface{}

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		collectionName := getParam(ctx, "collection_name")
		collection := collections[collectionName]

		jsonReader := json.NewDecoder(r.Body)
		jsonWriter := json.NewEncoder(w)

		for {
			item := map[string]interface{}{}
			err := jsonReader.Decode(&item)
			if err != nil {
				// TODO: handle error properly
				fmt.Println("ERROR:", err.Error())
				return
			}
			err = collection.Insert(item)
			if err != nil {
				// TODO: handle error properly
				fmt.Println("ERROR:", err.Error())
				return
			}

			jsonWriter.Encode(item)
		}

	}
}
