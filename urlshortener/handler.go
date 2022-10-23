package urlshort

import (
	"encoding/json"
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

type pathToUrlYmlConf struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYaml(yamlByteString []byte) ([]pathToUrlYmlConf, error) {
	var pathToUrlYmlConfs []pathToUrlYmlConf
	err := yaml.Unmarshal(yamlByteString, &pathToUrlYmlConfs)
	return pathToUrlYmlConfs, err
}

func buildMapFromYml(pathToUrlYmlConfs []pathToUrlYmlConf) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pathToUrl := range pathToUrlYmlConfs {
		pathToUrls[pathToUrl.Path] = pathToUrl.Url
	}
	return pathToUrls
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
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMapFromYml(parsedYaml)
	return MapHandler(pathToUrls, fallback), nil
}

type pathToUrl struct {
	Path string
	Url  string
}

func parseJson(jsn []byte) ([]pathToUrl, error) {
	var pathToUrls []pathToUrl
	err := json.Unmarshal(jsn, &pathToUrls)
	return pathToUrls, err
}

func buildMapFromJsn(pathToUrls []pathToUrl) map[string]string {
	pathToUrlsMap := make(map[string]string)
	for _, pathToUrl := range pathToUrls {
		pathToUrlsMap[pathToUrl.Path] = pathToUrl.Url
	}
	return pathToUrlsMap
}

func JsonHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(json)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMapFromJsn(parsedJson)
	return MapHandler(pathToUrls, fallback), nil
}
