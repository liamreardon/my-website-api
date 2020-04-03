# Simple Go RESTful API

> Simple RESTful API used to create and retrieve portfolio data for my personal website, built using Go and MongoDB.

## Quickstart

``` bash
# Install mux router
go get -u github.com/gorilla/mux
```

``` bash
# Install MongoDB Driver
go get go.mongodb.org/mongo-driver
```

``` bash
# Install govalidator
go get github.com/thedevsaddam/govalidator
```

``` bash
# Replace values in config/config.go with your own DB Uri and Port.
func GetConfig() *Config {
	uri, exists := os.LookupEnv("DB_URI")
	if exists {
		return &Config{
			DbURI: uri,
			Port:  ":8080",
		}
	}
	return &Config{}
}

```

``` bash
# Run
go run main.go
```

## Endpoints

### - Projects -
### Get All Projects
``` bash
GET /api/projects
```
### Add a Project
``` bash
POST /api/projects
```
### Update a Project
``` bash
PUT /api/projects/:title

# Request body
{
	"Title":"Example",
	"Description":"Example",
	"Img":"Example",
	"Link":"Example",
	"Tools":"Example"
}
```
### Delete a Project
``` bash
DELETE /api/projects/:title
```

### - Courses - 
### Get All Courses
``` bash
GET /api/courses
```
### Add a Course
``` bash
POST /api/courses
```
### Update a Course
``` bash
PUT /api/courses/:title

# Request body
{
	"Title":"Example"
}
```
### Delete a Course
``` bash
DELETE /api/courses/:title
```

## TODO
* [ ] Write unit tests for all packages
