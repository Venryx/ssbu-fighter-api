/**
 * api.go
 * Source:
 *   https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
 */

package api

import (
//    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

// our main function
func api() {
    router := mux.NewRouter()
    log.Fatal(http.ListenAndServe(":8000", router))
}
