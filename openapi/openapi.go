package openapi

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

var (
	//go:embed docs
	docs embed.FS
	//go:embed swagger-initializer.tmpl
	swaggerInitializer string

	fs = http.FileServer(http.FS(docs))

	tmpl = template.Must(template.New("swagger-initializer").Parse(swaggerInitializer))
)

type Schema struct {
	Name        string
	ContentJSON []byte
}

func ServeDocs(route string, schemas ...Schema) http.Handler {
	route = "/" + strings.Trim(route, "/") + "/"

	registry := make(map[string]Schema)

	for _, schema := range schemas {
		registry[fmt.Sprintf("%sschemas/%s.swagger.json", route, schema.Name)] = schema
	}

	mux := http.NewServeMux()

	mux.HandleFunc(route+"swagger-initializer.js", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, registry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		schema, ok := registry[r.URL.Path]
		if ok {
			w.Write(schema.ContentJSON)
			return
		}
		fs.ServeHTTP(w, r)
	})

	return mux
}
