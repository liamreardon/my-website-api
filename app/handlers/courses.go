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

// COURSE HANDLERS

// Get all courses
func GetCourses(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	var courses []models.Course
	collection := client.Database("liamreardonio").Collection("Courses")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// Get cursor for given collection
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	// Iterate through each document and add to projects slice
	for cursor.Next(ctx) {
		var course models.Course
		err := cursor.Decode(&course)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
		}
		courses = append(courses, course)
	}

	if err := cursor.Err(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	err = cursor.Close(ctx)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":"List of courses",
		"courses":courses,
	})
}

// Add new Course
func CreateCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	course, err := validateCourseRequest(r)
	if len(err) > 0 {
		respondError(w, http.StatusBadRequest, map[string]interface{}{
			"message":"Invalid request body",
			"error":err,
		})
		return
	}
	collection := client.Database("liamreardonio").Collection("Courses")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, _ := collection.InsertOne(ctx, course)
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":"Created course!",
		"result":res,
	})
}

// Update a course
func UpdateCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	collection := client.Database("liamreardonio").Collection("Courses")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	course := getCourseOr404(collection, ctx, title, w, r)
	if course == nil {
		return
	}

	res, err := collection.UpdateOne(ctx, bson.M{"title":title}, bson.D{
		{"$set", bson.D{
			{"title", course.Title},
		}},
	})

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":"Course updated!",
		"result":res,
	})
}

// Delete a course
func DeleteCourse(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	collection := client.Database("liamreardonio").Collection("Courses")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := getDocument(collection, ctx, title)
	if len(err) > 0 {
		respondError(w, http.StatusNotFound, map[string]interface{}{
			"message":"Course with that title doesn't exist",
			"error":err,
		})
		return
	}

	res, _ := collection.DeleteOne(ctx, result)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":"Course deleted!",
		"result":res,
	})
}

// Get course from db if exists, check for multiple error conditions
func getCourseOr404(collection *mongo.Collection, ctx context.Context, title string, w http.ResponseWriter, r *http.Request) *models.Course {

	res, err := validateCourseRequest(r)
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
			"message":"Course with that title doesn't exist",
			"error":err,
		})
		return nil
	}
	return res
}

// Validate request body and return as course model
func validateCourseRequest(r *http.Request) (*models.Course, map[string]interface{}) {
	var c models.Course

	rules := govalidator.MapData{
		"title":     []string{"required"},
	}

	opts := govalidator.Options{
		Request: r,
		Data:    &c,
		Rules:   rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{"validationError": e}
		return &c, err
	}

	return &c, map[string]interface{}{}
}