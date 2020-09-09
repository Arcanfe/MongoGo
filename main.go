package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Cart struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Items []*Category    `json:"items" bson:"items"`
}

type Item struct {
	Name string `json:"_id,omitempty" bson:"_id,omitempty"`
	Cost int `json:"cost" bson:"cost"`
	Desc string `json:"desc" bson:"desc"`
}

var client *mongodb.Client

func CreateCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var cart Cart
	json.NewDecoder(request.Body).Decode(&cart)
	collection := client.Database("apiCarts").Collection("carts")
	ctx, error := context.WithTimeout(context.Background(), 10*time.Second)
	result, error := collection.InsertOne(ctx, cart)
	json.NewEncoder(response).Encode(result)

}

func GetCarts(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var carts []Cart
	collection := client.Database("apiCarts").Collection("carts")
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer curson.Close(ctx)
	for cursor.Next(ctx) {
		var cart Cart
		cursor.Decode(&cart)
		carts = append(people, cart)
	}
	if err := cursor.Error(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(cart)
}

func GetCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var cart Cart
	collection := client.Database("apiCarts").Collection("carts")
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Cart{ID: id}).Decode(&cart)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(cart)
}

func DeleteCart(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var cart Cart
	collection := client.Database("apiCarts").Collection("carts")
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.DeleteOne(ctx, Cart{ID: id}).Decode(&cart)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(cart)
}

func main() {
	fmt.Println("Starting the application")
	ctx, error := context.WithTimeout(context.Background(), 10*time.Second)
	client, error := mongo.Connect(ctx, "mongodb://localhost:27017")
	router := mux.NewRouter()
	router.HandleFunc("/cart", CreateCart).Methods("POST")
	router.HandleFunc("/carts", GetCarts).Methods("GET")
	router.HandleFunc("/cart/{id}", GetCart).Methods("GET")
	router.HandleFunc("/cart/{id}", DeleteCart).Methods("DELETE")
	http.ListenAndServe(":10000", router)

}
