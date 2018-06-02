package main

import (
        "fmt"
        "net/http"
        // "encoding/json"
)

type Player struct {
        name     string
        school   string
        position string
        year      Year
}

type Year int

var players []Player

func main() {
        // log.Fatal("asdfasf")
        // json.Marshal()
        fmt.Print("we have ignition")
        http.ListenAndServe(":8080", nil)

}
