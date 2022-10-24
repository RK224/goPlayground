package urlshort

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

type BaseHandlerParams struct {
	Base     string
	BasePath string
}

func GetBaseHandler(pathsToUrls map[string]string, fallback http.Handler, base string) http.HandlerFunc {
	if base == "db" {
		initializeDb(pathsToUrls)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		pathKey := r.URL.Path
		var redirectUrl string
		var ok bool
		switch base {
		case "db":
			redirectUrl, ok = getUrlForPathFromDb(pathKey)
		case "map":
			redirectUrl, ok = pathsToUrls[pathKey]
		}

		if ok {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func initializeDb(pathsToUrls map[string]string) {
	db, err := bolt.Open("bolt.db", 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("PathToUrl"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		for path, url := range pathsToUrls {
			err = b.Put([]byte(path), []byte(url))
			if err != nil {
				return fmt.Errorf("Unable to insert value %s for key %s", url, path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getUrlForPathFromDb(path string) (url string, ok bool) {
	db, err := bolt.Open("bolt.db", 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PathToUrl"))
		if b == nil {
			return errors.New("The request bucket does not exist")
		}
		v := b.Get([]byte(path))
		if v != nil {
			url = string(v)
			ok = true
		}
		return nil
	})
	return
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
	return GetBaseHandler(pathToUrls, fallback, "map"), nil
}
