package main

import (
	"database/sql"
	"fmt"
)

type Env struct {
    db *sql.DB
}

func main() {
	DispatchDownloads()
    fmt.Println(getActor("nm0000001"))
    fmt.Println(getConnectedActors("tt7348214"))
}
