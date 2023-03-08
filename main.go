package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Item struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    *Category `json:"category"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var items []Item

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

}

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range items {
		if item.ID == params["item_id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID = strconv.Itoa(rand.Intn(100000))
	items = append(items, item)
	json.NewEncoder(w).Encode(item)
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parms := mux.Vars(r)
	var updatedItem Item
	_ = json.NewDecoder(r.Body).Decode(&updatedItem)
	updatedItem.ID = parms["item_id"]

	for index, item := range items {
		if item.ID == parms["item_id"] {
			items = append(items[:index], items[index+1:]...)
			items = append(items, updatedItem)
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range items {
		if item.ID == params["item_id"] {
			items = append(items[:index], items[index+1:]...)
			resp := struct {
				Message string `json:"message"`
			}{Message: item.Title + " is deleted"}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	items = append(items, Item{ID: "1", Title: "item 1", Description: "This is an organic item. YOu can use it", Category: &Category{ID: "1", Name: "CAT_1"}})
	items = append(items, Item{ID: "2", Title: "item 2", Description: "This is an organic item 2. YOu can use it too", Category: &Category{ID: "2", Name: "CAT_2"}})

	r.HandleFunc("/items", itemsHandler).Methods("GET")
	r.HandleFunc("/item/{item_id}", getItemHandler).Methods("GET")
	r.HandleFunc("/item", createItemHandler).Methods("POST")
	r.HandleFunc("/item/{item_id}", createItemHandler).Methods("PUT")
	r.HandleFunc("/item/{item_id}", deleteItemHandler).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
