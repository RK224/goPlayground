package urlshort

import (
	"encoding/json"
	"errors"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathKey := r.URL.Path
		redirectUrl, ok := pathsToUrls[pathKey]
		if ok {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type pathToUrl struct {
	Path string
	Url  string
}

func parseConf(conf []byte, extension string) ([]pathToUrl, error) {
	var (
		pathToUrls []pathToUrl
		err        error
	)

	switch extension {
	case "yaml":
		err = yaml.Unmarshal(conf, &pathToUrls)
	case "json":
		err = json.Unmarshal(conf, &pathToUrls)
	default:
		err = errors.New("Unsupported file type")
	}
	return pathToUrls, err
}

func buildMap(pathToUrls []pathToUrl) map[string]string {
	pathToUrlMap := make(map[string]string)
	for _, pathToUrl := range pathToUrls {
		pathToUrlMap[pathToUrl.Path] = pathToUrl.Url
	}
	return pathToUrlMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseConf(yml, "yaml")
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(parsedYaml)
	return MapHandler(pathToUrls, fallback), nil
}

func JsonHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseConf(jsn, "json")
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(parsedJson)
	return MapHandler(pathToUrls, fallback), nil
}
