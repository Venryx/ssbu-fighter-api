/* api.go:
 * Authors: Anthony Luc (aluc)
 *          Donald Luc (dluc)
 *          Michael Wang (mwang6)
 * Workflow:
 *  (1) Pull data from Google sheet [main|backup].
 *  (2) Store into (MySQL?) database.
 *  (3) Handle RESTful API requests.
 *
 * Source:
 *   https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
 *   https://mholt.github.io/json-to-go/
 */

package main

import (
    "fmt"
    "io"
    "encoding/csv"
    "strings"
    "os"
    "time"
    "encoding/json"
    "log"
    "net/http"

	"github.com/gorilla/handlers" 
	"github.com/gorilla/mux"
)


const sourcePath string = "data/frame-data"
var fighters[]Fighter


type Fighter struct {
    Name        string `json:"name,omitempty"`
    ID          string `json:"id,omitempty"`
    Frames      *Frames `json:"frames,omitempty"`
}

type Frames struct {
    Action      string `json:"action,omitempty"`
    Startup     string `json:"startup,omitempty"`
    TotalFrames string `json:"totalframes,omitempty"`
    LandingLag  string `json:"landinglag,omitempty"`
}


// Display all from the fighters var
func GetFrameData(w http.ResponseWriter, r *http.Request) {
    // With Query
    query := r.URL.Query()
    action := query.Get("action")
    if len(query) != 0 {
        found := false
        for _, fighter := range fighters {
            if fighter.Frames != nil {
                if fighter.Frames.Action == action {
                    json.NewEncoder(w).Encode(fighter)
                    found = true
                }
            }
        }
        if found == false {
            fmt.Fprintf(w, "Action not found!")
        }
    } else { // Without Query
        json.NewEncoder(w).Encode(fighters)
    }
}


// Display a fighter
func GetFighter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]

    query := r.URL.Query()
    action := query.Get("action")
    if len(query) == 0 { // Without Query
        found := false
        for _, fighter := range fighters {
            if fighter.Name == name || fighter.ID == name {
                json.NewEncoder(w).Encode(fighter)
                found = true
            }
        }
        if found == false {
            fmt.Fprintf(w, "Name/ID not found!")
        }
    } else { // With Query
        found := false
        for _, fighter := range fighters {
            // fmt.Fprintf(w, fighter.Name)
            if fighter.Name == name || fighter.ID == name {
                if fighter.Frames != nil {
                    if fighter.Frames.Action == action {
                        json.NewEncoder(w).Encode(fighter)
                        found = true
                    }
                }
            }
        }
        if found == false {
            fmt.Fprintf(w, "Action not found!")
        }
    }
}


// Main Function
func main() {

    /* Determine port on system. */
    port := ":" + os.Getenv("PORT")
    if port == "" {
        log.Fatal("$PORT must be set")
    }

	/* Load data from CSV. */
    dir, _ := os.Open(sourcePath)
    files, _ := dir.Readdir(-1)
    for _, file := range files {
            fileName := file.Name()
        filePath := sourcePath + "/" + fileName
        f, _ := os.Open(filePath)
        defer f.Close()

        go read_frame_data(fileName, f)
        time.Sleep(10 * time.Millisecond)
   
    }

	/* Add routes. */
    router := mux.NewRouter()
    router.HandleFunc("/api", GetFrameData).Methods("GET")
    router.HandleFunc("/api/{name}", GetFighter).Methods("GET")

	/* Handle CORS.
	 * Source: https://stackoverflow.com/questions/40985920/making-golang-gorilla-cors-handler-work
	 */
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))	

    fmt.Println("Finished main()")
}

func read_frame_data(name string, file io.Reader) {
    records, _ := csv.NewReader(file).ReadAll()
    for _, row := range records {
        id := strings.Split(name, " - ")[0]
        fighter_csv := strings.Split(name, " - ")[1]
        name := strings.Replace(fighter_csv, ".csv", "", -1)
        fighters = append(fighters, Fighter{Name: name, ID: id, Frames: &Frames{Action: row[0], Startup: row[1], TotalFrames: row[2], LandingLag: row[3]}})
    }
}
