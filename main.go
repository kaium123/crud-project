package main

import(
	"encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Movie struct {
	ID         uint        `json:"id" gorm:"primary_key"`
    Name 	   string      `json:"name"`
    Director  *Director    `json:"director" gorm:"foreignkey:ID"`
}
type Director struct {
	Did		   uint 		`json:"did" gorm:"primary_key"`
	Firstname  string 		`json:"firstname" `
	Lastname   string		`json:"lastname" `
	ID     	   uint		    `json:"-"`
}

var db *gorm.DB

func initDB() {
    var err error
    dataSourceName := "root:@tcp(localhost:3306)/?parseTime=True"
    db, err = gorm.Open("mysql", dataSourceName)

    if err != nil {
        fmt.Println(err)
        panic("failed to connect database")
    }

    db.Exec("CREATE DATABASE movie_db")
    db.Exec("USE movie_db")

    db.AutoMigrate(&Movie{}, &Director{})
}
func createMovie(w http.ResponseWriter, r *http.Request) {
    var movie Movie
    json.NewDecoder(r.Body).Decode(&movie)

    db.Create(&movie)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movie)
}
func getMovies(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var movies []Movie
    db.Preload("Director").Find(&movies)
    json.NewEncoder(w).Encode(movies)
}
func getMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    inputMovieID := params["id"]

    var movie Movie
    db.Preload("Director").First(&movie, inputMovieID)
    json.NewEncoder(w).Encode(movie)
}
func updateMovie(w http.ResponseWriter, r *http.Request) {
    var updatedMovie Movie
    fmt.Println("df")
    json.NewDecoder(r.Body).Decode(&updatedMovie)
    db.Save(&updatedMovie)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedMovie)
}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    inputMovieID := params["id"]

    id64, _ := strconv.ParseUint(inputMovieID, 10, 64)

    idToDelete := uint(id64)

    db.Where("id = ?", idToDelete).Delete(&Director{})
    db.Where("id = ?", idToDelete).Delete(&Movie{})
    w.WriteHeader(http.StatusNoContent)
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/movies", createMovie).Methods("POST")
 
    router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
  
    router.HandleFunc("/movies", getMovies).Methods("GET")

    router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")

    router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

    initDB()

    log.Fatal(http.ListenAndServe(":8001", router))
}
