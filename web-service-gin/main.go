package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gin-gonic/gin"
)

//album represents data about a record album.
// struct tags specify the format and name of the contents that are seruakuzed into JSON.
type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

// albums slice to seed record album data
var albums = []album{
	{ID: "1", Title: "Clouds", Artist: "NF", Price: "56"},
	{ID: "2", Title: "Jeru", Artist: "Gerry", Price: "17"},
	{ID: "3", Title: "Real", Artist: "Olivia", Price: "39"},
}

// assign handler function to an endpoint path "/albums"
func main() {
	// returns an empty context, used to derive other contexts
	ctx := context.Background()
	// returns address of loader that deserializes openapi3 context document
	loader := &openapi3.Loader{Context: ctx}
	// loads OAS document and deserializes openapi3 it with the loader. Returns
	// as a doc or blank identifier
	doc, _ := loader.LoadFromFile("openapi3/petstore.json")
	// validates the doc and assigns it as a blank identifier
	_ = doc.Validate(ctx)
	// defines a new router for the doc
	router, _ := legacyrouter.NewRouter(doc)
	// httpReq, _ := http.NewRequest(http.MethodGet, "/inventory", nil)
	httpReq, _ := http.NewRequest(http.MethodGet, "/v2/pet/findByStatus?status=sold", nil)

	// Find route
	route, pathParams, _ := router.FindRoute(httpReq)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}

	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = bytes.NewBufferString(`{}`)
	)

	log.Println("Response:", respStatus)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}

	if respBody != nil {
		data, _ := json.Marshal(respBody)
		responseValidationInput.SetBodyBytes(data)
	}

	// write handler and creating server

	// http.ListenAndServe(":8080", router)

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}

	// switch reqtype {
	// case "GET":
	// 	if reqtype.URL == "/inventory" {
	// 		httpReq, _ := http.NewRequest(http.MethodGet, "/inventory", nil)
	// 	}
	// 	if reqtype.URL == "/albums/:id" {
	// 		httpReq, _ := http.NewRequest(http.MethodGet, "/albums/:id", nil)
	// 	}
	// 	if reqtype.URL == "/albums" {
	// 		httpReq, _ := http.NewRequest(http.MethodGet, "/albums/", nil)
	// 	}
	// case "POST":
	// 	httpReq, _ := http.NewRequest(http.MethodPost, "/albums", nil)
	// }

	// initialize gin router
	// router := gin.Default()
	// router.GET("/inventory", getAlbums)
	// router.POST("/albums", postAlbums)
	// router.GET("/albums/:id", getAlbumByPath)
	// router.GET("/albums", getAlbumByQuery)
	// router.Run("localhost:9000")
}

// getAlbums responds with the list of all albums as JSON.
// gin.Context carries JSON request details, validates and serializes JSON.
func getAlbums(c *gin.Context) {
	// seralizes stcture into JSON and added to JSON response
	c.Request.Header.Set("Content-Type", "application/json")
	log.Println(c.GetHeader("Content-Type"))
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON recieved in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// call BindJSON to bind the recieved JSON to newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// add new album to the albums slice
	albums = append(albums, newAlbum)
	// 201 status code to reponse with newly created album to JSON
	c.IndentedJSON(http.StatusCreated, newAlbum)
	c.String(200, "\nSuccess")
	log.Println(newAlbum)
	c.Request.Header.Set("Host", "localhost:9000")
	log.Println(c.Request.Header.Get("Host"))
	c.Request.Header.Add("id", newAlbum.ID)
	c.Request.Header.Add("price", newAlbum.Price)
	log.Println(c.Request.Header.Get("id"))
	log.Println(c.Request.Header.Get("price"))
	log.Println(c.Request.Header.Get("App-id"))
}

// getAlbumByID locates the album whose ID values matches the ID
// parameter sent by the client, then returns that album as a response.
func getAlbumByPath(c *gin.Context) {
	// retrieves the id path param from URL
	id := c.Param("id")

	// loop over list of albums, looking for an album that matches "id" with ID parameter
	for _, a := range albums {
		if a.ID == id {
			// if album ID matches param id, album is serialized to JSOn and returned
			// as a response of 200 OK
			c.IndentedJSON(http.StatusOK, a)
			c.String(200, "\nSuccess")
			c.Request.Header.Set("Host", "localhost:9000")
			log.Println(c.Request.Header.Get("Host"))
			c.Request.Header.Add("id", a.ID)
			log.Println(c.Request.Header.Get("id"))
			return
		}
	}

	// return HTTP 404 error if album isn't found
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// getAlbumByQuery locates the album whose ID, title, artist or price matches the
// query parameter sent by the client, then returns that album as a response.
func getAlbumByQuery(c *gin.Context) {
	// retrieves the id path param from URL
	id := c.Query("id")
	title := c.Query("title")
	artist := c.Query("artist")
	price := c.Query("price")

	// loop over list of albums, looking for an album that matches "id" with ID parameter
	for _, a := range albums {
		if a.ID == id || a.Title == title || a.Artist == artist || a.Price == price {
			// if album ID matches param id, album is serialized to JSOn and returned
			// as a response of 200 OK
			c.IndentedJSON(http.StatusOK, a)
			c.String(200, "\nSuccess")
			return
		}
	}
	// return HTTP 404 error if album isn't found
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
