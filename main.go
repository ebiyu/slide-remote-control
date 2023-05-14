package main

import (
	"log"
	"net/http"
	"os/exec"
    _ "embed"
)

//go:embed index.html
var indexHTML string

func main() {
	// define handlers
	handleNext := func(w http.ResponseWriter, r *http.Request) {
		if err := nextSlide(); err != nil {
			log.Println(err)
		}

        // redirect
        http.Redirect(w, r , "/", http.StatusSeeOther)
	}
	handlePrevious := func(w http.ResponseWriter, r *http.Request) {
		if err := prevSlide(); err != nil {
			log.Println(err)
		}
        http.Redirect(w, r , "/", http.StatusSeeOther)
	}
    handleIndex := func(w http.ResponseWriter, _ *http.Request) {
        w.Write([]byte(indexHTML))
    }

	// define routes
	http.HandleFunc("/", handleIndex)
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
