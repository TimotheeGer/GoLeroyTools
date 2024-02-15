package main

import (
	"GoLeroyTools/leroyTools/handlers"
	"fmt"
)

func main() {
	
	fmt.Printf("Go Leroy Tools\n")
	
	urls, err := handlers.RetrieveURLsFromJSON()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	handlers.ProcessURLs(urls)
}
