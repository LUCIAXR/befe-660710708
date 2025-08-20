package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Movie struct
type Movie struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Director string  `json:"director"`
	Genre    string  `json:"genre"`
	Year     int     `json:"year"`
	Rating   float64 `json:"rating"`
}

// In-memory database (ในโปรเจคจริงใช้ database)
var movies = []Movie{
	{ID: "1", Title: "Inception", Director: "Christopher Nolan", Genre: "Sci-Fi", Year: 2010, Rating: 8.8},
	{ID: "2", Title: "Parasite", Director: "Bong Joon-ho", Genre: "Thriller", Year: 2019, Rating: 8.6},
	{ID: "3", Title: "Spirited Away", Director: "Hayao Miyazaki", Genre: "Animation", Year: 2001, Rating: 8.5},
}

func getMovies(c *gin.Context) {
	yearQuery := c.Query("year")

	if yearQuery != "" {
		filter := []Movie{}
		for _, movie := range movies {
			if fmt.Sprint(movie.Year) == yearQuery {
				filter = append(filter, movie)
			}
		}
		c.JSON(http.StatusOK, filter)
		return
	}
	c.JSON(http.StatusOK, movies)
}

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/movies", getMovies)
	}

	r.Run(":8080")
}
