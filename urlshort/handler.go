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

func Handler(conf []byte, extension string, fallback http.Handler) (http.HandlerFunc, error) {
	parsedConf, err := parseConf(conf, extension)

	if err != nil {
		return nil, err
	}

	pathToUrls := buildMap(parsedConf)
	return MapHandler(pathToUrls, fallback), nil
}
