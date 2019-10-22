package main

import (
	"flag"
	"fmt"
	urlshort "gophercises/url-shortener"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var yamlFilename = flag.String("-yml", "default.yml", "a YAML file in the format :\n- path: /some-path\n  url: url: https://www.some-url.com/demo")
	var jsonFilename = flag.String("-json", "default.json", "a YAML file in the format :\n- path: /some-path\n  url: url: https://www.some-url.com/demo")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler
	yaml, err := ioutil.ReadFile(*yamlFilename)
	if err != nil {
		log.Printf("Failed to open file %s ", *yamlFilename)
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler
	json, err := ioutil.ReadFile(*jsonFilename)
	if err != nil {
		log.Printf("Failed to open file %s ", *jsonFilename)
	}
	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
