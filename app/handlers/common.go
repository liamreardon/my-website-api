package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// Response in JSON format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}

// Error response in JSON format
func respondError(w http.ResponseWriter, code int, msg interface{}) {
	respondJSON(w, code, msg)
}

// Get document in db if exists
func getDocument(collection *mongo.Collection, ctx context.Context, title string) (bson.M, string) {
	var result bson.M
	err := collection.FindOne(ctx, bson.M{"title":title}).Decode(&result)
	if err != nil {
		return nil, "Document doesn't exist"
	}
	return result, ""
}
