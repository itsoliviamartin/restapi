package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//album represents data about a record album.
// struct tags specify the format and name of the contents that are seruakuzed into JSON.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data
var albums = []album{
	{ID: "1", Title: "Clouds", Artist: "NF", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Lost", Artist: "NF", Price: 39.99},
}

// assign handler function to an endpoint path "/albums"
func main() {
	// initialize gin router
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	// associate /albums/:id path with getAlbumsByID func
	// colon preceding an item in the path means the item
	// is a path param
	router.GET("/albums/:id", getAlbumByID)
	// router.DELETE("/albums/:artist", deleteAlbumByArtist)
	// attach the router to an http.Server and start server
	router.Run("localhost:9000")
}

// getAlbums responds with the list of all albums as JSON.
// gin.Context carries JSON request details, validates and serializes JSON.
func getAlbums(c *gin.Context) {
	// seralizes stcture into JSON and added to JSON response
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
}

// getAlbumByID locates the album whose ID values matches the ID
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	// retrieves the id path param from URL
	id := c.Param("id")

	// loop over list of albums, looking for an album that matches "id" with ID parameter
	for _, a := range albums {
		if a.ID == id {
			// if album ID matches param id, album is serialized to JSOn and returned
			// as a response of 200 OK
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	// return HTTP 404 error if album isn't found
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
