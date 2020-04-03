package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/liamreardon/my-website-api/app/models"
	"github.com/thedevsaddam/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// PROJECTS HANDLERS //

// Get all projects
func GetProjects(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	var projects []models.Project
	collection := client.Database("liamreardonio").Collection("Projects")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Get cursor for given collection
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	// Iterate through each document and add to projects slice
	for cursor.Next(ctx) {
		var project models.Project
		err := cursor.Decode(&project)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
		}
		projects = append(projects, project)
	}

	if err := cursor.Err(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	err = cursor.Close(ctx)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":"List of projects",
		"projects":projects,
	})
}

// Add new project
func CreateProject(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	project, err := validateProjectRequest(r)
	if len(err) > 0 {
		respondError(w, http.StatusBadRequest, map[string]interface{}{
			"message":"Invalid request body",
			"error":err,
		})
		return
	}
	collection := client.Database("liamreardonio").Collection("Projects")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, _ := collection.InsertOne(ctx, project)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":"Created project!",
		"result":res,
	})
}

// Update a Project
func UpdateProject(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	collection := client.Database("liamreardonio").Collection("Projects")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	project := getProjectOr404(collection, ctx, title, w, r)
	if project == nil {
		return
	}

	res, err := collection.UpdateOne(ctx, bson.M{"title":title}, bson.D{
		{"$set", bson.D{
			{"title", project.Title},
			{"description", project.Description},
			{"img", project.Img},
			{"link", project.Link},
			{"tools", project.Tools},
		}},
	})

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":"Project updated!",
		"result":res,
	})
}

// Delete a Project
func DeleteProject(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	collection := client.Database("liamreardonio").Collection("Projects")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := getDocument(collection, ctx, title)
	if len(err) > 0 {
		respondError(w, http.StatusNotFound, map[string]interface{}{
			"message":"Project with that title doesn't exist",
			"error":err,
		})
		return
	}

	res, _ := collection.DeleteOne(ctx, result)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":"Project deleted!",
		"result":res,
	})
}

// Get project from db if exists, check for multiple error conditions
func getProjectOr404(collection *mongo.Collection, ctx context.Context, title string, w http.ResponseWriter, r *http.Request) *models.Project {

	res, err := validateProjectRequest(r)
	if len(err) > 0 {
		respondError(w, http.StatusInternalServerError, map[string]interface{}{
			"message":"Invalid request body",
			"error":err,
		})
		return nil
	}

	var result bson.M
	error := collection.FindOne(ctx, bson.M{"title":title}).Decode(&result)
	if error != nil {
		respondError(w, http.StatusNotFound, map[string]interface{}{
			"message":"Project with that title doesn't exist",
			"error":err,
		})
		return nil
	}
	return res
}

// Validate request body and return as project model
func validateProjectRequest(r *http.Request) (*models.Project, map[string]interface{}) {
	var p models.Project

	rules := govalidator.MapData{
		"title": 		  []string{"required"},
		"description":    []string{"required"},
		"img":      	  []string{"required"},
		"link":    	      []string{"required"},
		"tools":    	  []string{"required"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &p,
		Rules:   rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{"validationError": e}
		return &p, err
	}
	
	return &p, map[string]interface{}{}
}
