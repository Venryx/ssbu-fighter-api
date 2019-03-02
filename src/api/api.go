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
 */

package main

import (
    "fmt"
    "io"
    "encoding/csv"
    "strings"
    "os"
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)


const sourcePath string = "data/frame-data"
var fighters[]Fighter


type Fighter struct {
    Name        string `json:"name,omitempty"`
    Action      string `json:"action,omitempty"`
    Frames      *Frames `json:"frames,omitempty"`
}

type Frames struct {
    Startup     string `json:"startup,omitempty"`
    TotalFrames string `json:"totalframes,omitempty"`
    LandingLag  string `json:"landinglag,omitempty"`
}


// Display all from the fighters var
func GetFrameData(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GET params: ", r.URL.Query())
    json.NewEncoder(w).Encode(fighters)
}


// Display a fighter
func GetFighter(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["name"]
    fmt.Fprintf(w, "name: " + name)
    // json.NewEncoder(w).Encode(fighters)
}


// Main Function
func main() {

    /* Determine port on system. */
    port := ":" + os.Getenv("PORT")
    if port == "" {
        log.Fatal("$PORT must be set")
    }

    dir, _ := os.Open(sourcePath)
    files, _ := dir.Readdir(-1)

    for _, file := range files {
            fileName := file.Name()
        filePath := sourcePath + "/" + fileName
        f, _ := os.Open(filePath)
        defer f.Close()

        go read_frame_data(fileName, f)
    }

    router := mux.NewRouter()
    router.HandleFunc("/api", GetFrameData).Methods("GET")
    router.HandleFunc("/api/{name}", GetFighter).Methods("GET")

    log.Fatal(http.ListenAndServe(port, router))

    fmt.Println("Finished main()")
}

func read_frame_data(name string, file io.Reader) {
    records, _ := csv.NewReader(file).ReadAll()
    for _, row := range records {
        fighter_csv := strings.Split(name, " - ")[1]
        name := strings.Replace(fighter_csv, ".csv", "", -1)
        fighters = append(fighters, Fighter{Name: name, Action: row[0], Frames: &Frames{Startup: row[1], TotalFrames: row[2], LandingLag: row[3]}})
        // fmt.Println(fighter, row)
        // fmt.Println(name, row)
    }
}