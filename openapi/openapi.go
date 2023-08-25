package openapi

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"
)

var (
	//go:embed docs
	docs embed.FS
	//go:embed swagger-initializer.tmpl
	swaggerInitializer string

	schemas = make(map[string]Schema)

	tmpl = template.Must(template.New("swagger-initializer").Parse(swaggerInitializer))
)

type Schema struct {
	Name        string
	ContentJSON []byte
}

func RegisterSchema(schema Schema) error {
	url := fmt.Sprintf("/docs/schemas/%s.swagger.json", schema.Name)
	_, ok := schemas[url]
	if ok {
		return fmt.Errorf("openapi: schema %q already registered", schema.Name)
	}
	schemas[url] = schema
	return nil
}

func ServeDocs() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/docs/swagger-initializer.js", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, schemas)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		schema, ok := schemas[r.URL.Path]
		if ok {
			w.Write(schema.ContentJSON)
			return
		}
		http.FileServer(http.FS(docs)).ServeHTTP(w, r)
	})

	return mux
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
