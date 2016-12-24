package main

import (
	"database/sql"
	"flag"
)

var (
	portNumber = flag.Int("p", 8081, "Port yang digunakan oleh aplikasi")
)

func main() {
	// Parse flags
	flag.Parse()

	// Create backend
	backEnd := NewBackEnd(*portNumber)
	defer backEnd.Close()

	// Serve app
	backEnd.ServeApp()
}

func checkError(err error) {
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
}
