package main

import (
	"log"
	"net/http"

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
	// initialize gin router
	router := gin.Default()
	router.GET("/inventory", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByPath)
	router.GET("/albums", getAlbumByQuery)

	router.Run("localhost:9000")
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
