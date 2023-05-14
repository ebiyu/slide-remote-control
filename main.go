package main

import (
	"log"
	"net/http"
	"os/exec"
)

func main() {
	// define handlers
	handleNext := func(w http.ResponseWriter, _ *http.Request) {
		if err := nextSlide(); err != nil {
			log.Fatal(err)
		}
	}
	handlePrevious := func(w http.ResponseWriter, _ *http.Request) {
		if err := prevSlide(); err != nil {
			log.Fatal(err)
		}
	}

	// define routes
	http.HandleFunc("/next/", handleNext)
	http.HandleFunc("/prev/", handlePrevious)

	// start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func nextSlide() error {
	cmd := exec.Command("xdotool", "search", "Google Slides", "windowactivate", "--sync", "key", "--clearmodifiers", "Down")
	err := cmd.Run()
	return err
}

func prevSlide() error {
	cmd := exec.Command("xdotool", "search", "Google Slides", "windowactivate", "--sync", "key", "--clearmodifiers", "Up")
	err := cmd.Run()
	return err
}
