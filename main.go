package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

//go:embed index.html
var indexHTML string

func main() {
	useNgrok := flag.Bool("ngrok", false, "Use ngrok to expose the server to the internet")
	flag.Parse()

	// define handlers
	handleNext := func(w http.ResponseWriter, r *http.Request) {
		if err := nextSlide(); err != nil {
			log.Println(err)
		}

		// redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	handlePrevious := func(w http.ResponseWriter, r *http.Request) {
		if err := prevSlide(); err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	handleIndex := func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(indexHTML))
	}

	// define routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/next/", handleNext)
	http.HandleFunc("/prev/", handlePrevious)

	// start server
	if *useNgrok {
		token, err := getNgrokAuthToken()
		if err != nil {
			log.Fatal(err)
		}
		tun, err := ngrok.Listen(context.Background(),
			config.HTTPEndpoint(),
			ngrok.WithAuthtoken(token),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("tunnel created:", tun.URL())
		log.Fatal(http.Serve(tun, nil))
	} else {
		fmt.Println("Listening on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
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

func getNgrokAuthToken() (string, error) {
	// read file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fileName := homeDir + "/.config/ngrok/ngrok.yml"

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	type conf struct {
		AuthToken string `yaml:"authtoken"`
		Version   string `yaml:"version"`
	}

	c := conf{}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return "", err
	}

	if c.AuthToken == "" {
		return "", fmt.Errorf("authToken is empty")
	}

	return c.AuthToken, nil
}
