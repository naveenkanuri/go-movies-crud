package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Response struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	response := Response{Message: "Movie deleted", Code: "200"}
	json.NewEncoder(w).Encode(response)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movie := Movie{}
	json.NewDecoder(r.Body).Decode(&movie)
	id := len(movies) + 1
	for index, item := range movies {
		if item.ID != strconv.Itoa(index+1) {
			id = index + 1
			break
		}
	}
	movie.ID = strconv.Itoa(id)
	movies = append(movies, movie)
	response := Response{Message: "Movie created at index " + movie.ID, Code: "200"}
	json.NewEncoder(w).Encode(response)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movie := Movie{}
	json.NewDecoder(r.Body).Decode(&movie)
	for _, item := range movies {
		if item.ID == movie.ID {
			item.Isbn = movie.Isbn
			item.Title = movie.Title
			item.Director = movie.Director
			break
		}
	}
	response := Response{Message: "Movie updated", Code: "200"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438277", Title: "Movie one", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "438144", Title: "Movie two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
