package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server running on localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal")
	}
}
