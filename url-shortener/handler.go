package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-yaml/yaml"
)

type entry struct {
	Key   string `json:"path" yaml:"path"`
	Value string `json:"url" yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path
		if address, present := pathsToUrls[key]; present {
			http.Redirect(w, r, address, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func buildMap(entries []entry) map[string]string {
	m := make(map[string]string)
	for _, ent := range entries {
		m[ent.Key] = ent.Value
	}
	return m
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) (entries []entry, err error) {
	err = yaml.UnmarshalStrict(yml, &entries)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	fmt.Println(pathMap)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsn []byte) (entries []entry, err error) {
	err = json.Unmarshal(jsn, &entries)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return
}
