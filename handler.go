package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type pathToURL map[string]string

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return redirectFunc(pathsToUrls, fallback)
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	m, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	return redirectFunc(m, fallback), err
}

func redirectFunc(p2u pathToURL, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := p2u[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func parseYAML(yml []byte) (pathToURL, error) {
	var pathUrls []struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}
	m := make(pathToURL)

	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return m, err
	}

	for _, pathURL := range pathUrls {
		m[pathURL.Path] = pathURL.URL
	}
	return m, err
}
