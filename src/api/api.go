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

//    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

const sourcePath string = "data/frame-data"
var frameDataDict = map[string]map[string]string{}


// our main function
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
        fmt.Println(fileName)
        filePath := sourcePath + "/" + fileName
        f, _ := os.Open(filePath)
        defer f.Close()

        go read_frame_data(fileName, f)
    }

    router := mux.NewRouter()
    log.Fatal(http.ListenAndServe(port, router))

    fmt.Println("Finished main()")
}

func read_frame_data(name string, file io.Reader) {
    records, _ := csv.NewReader(file).ReadAll()
    for _, row := range records {
        fighter_csv := strings.Split(name, " - ")[1]
        fighter := strings.Replace(fighter_csv, ".csv", "", -1)
        fmt.Println(fighter, row)
        // fmt.Println(name, row)
    }
}