package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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
		} else {
			log.Println("Next slide")
		}

		// redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	handlePrevious := func(w http.ResponseWriter, r *http.Request) {
		if err := prevSlide(); err != nil {
			log.Println(err)
		} else {
			log.Println("Previous slide")
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
		findGoogleSlidesWindow()
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func findGoogleSlidesWindow() (string, error) {
	pidsOutput, err := exec.Command("xdotool", "search", "--name", "Google Slides").Output()
	if err != nil {
		return "", err
	}
	pids := strings.Split(string(pidsOutput), "\n")
	for _, pid := range pids {
		name, err := exec.Command("xdotool", "getwindowname", pid).Output()
		if err != nil {
			return "", err
		}
		// avoid presenter view
		if strings.Contains(string(name), "Google Slides") && !strings.Contains(string(name), "Presenter view") {
			return pid, nil
		}
	}
	return "", errors.New("Could not find Google Slides window")
}

func sendKeyToGoogleSlides(key string) error {
	pid, err := findGoogleSlidesWindow()
	if err != nil {
		return err
	}

	cmd := exec.Command("xdotool", "windowactivate", "--sync", pid, "key", "--clearmodifiers", key)
	err = cmd.Run()
	return err
}

func nextSlide() error {
	err := sendKeyToGoogleSlides("Down")
	return err
}

func prevSlide() error {
	err := sendKeyToGoogleSlides("Up")
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
