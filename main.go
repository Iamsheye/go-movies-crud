package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movies struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movies

func findMovieById(id string) (*Movies, error) {
	for _, movie := range movies {
		if movie.ID == id {
			return &movie, nil
		}
	}
	return nil, errors.New("movie not found")
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	movie, err := findMovieById(vars["id"])

	if err != nil {
		movieError := map[string]string{
			"status": strconv.Itoa(http.StatusNotFound),
			"error":  err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(movieError)
		return
	}

	json.NewEncoder(w).Encode(movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movies
	json.NewDecoder(r.Body).Decode(&movie)

	movies = append(movies, movie)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idFound := false
	var movie Movies
	json.NewDecoder(r.Body).Decode(&movie)

	for i, v := range movies {
		if v.ID == vars["id"] {
			idFound = true
			movies[i] = movie
			break
		}
	}

	if idFound {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusOK),
			"message": vars["id"] + " succesfully updated",
		})
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  strconv.Itoa(http.StatusNotFound),
			"message": vars["id"] + " not found",
		})
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idFound := false
	var newMovies []Movies

	for _, movie := range movies {
		if movie.ID != vars["id"] {
			newMovies = append(newMovies, movie)
		} else {
			idFound = true
		}
	}

	movies = newMovies

	if !idFound {
		movieError := map[string]string{
			"status": strconv.Itoa(http.StatusNotFound),
			"error":  "movie not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(movieError)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func main() {
	const PORT int = 5000
	movies = append(movies, Movies{ID: "1", Isbn: "438227", Title: "Black Panther", Director: &Director{Firstname: "Sheye", Lastname: "Majek"}})

	movies = append(movies, Movies{ID: "2", Isbn: "45444", Title: "Black Adam", Director: &Director{Firstname: "Warner", Lastname: "Broski"}})

	r := mux.NewRouter()

	r.HandleFunc("/movies", getAllMovies).Methods("GET")  //Done
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET") //Done
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE") //Done

	fmt.Printf("Server listening on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), r))

	fmt.Println(movies)
}
