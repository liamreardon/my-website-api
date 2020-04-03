package models

// Course struct
type Course struct {
	Title string `json:"title"`
}

// Project struct
type Project struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Img string `json:"img"`
	Link string `json:"link"`
	Tools string `json:"tools"`
}