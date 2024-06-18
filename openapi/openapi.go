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
	docs   embed.FS
	docsFS = http.FileServer(http.FS(docs))

	//go:embed swagger-initializer.tmpl
	swaggerInitializer     string
	swaggerInitializerTmpl = template.Must(template.New("swagger-initializer").Parse(swaggerInitializer))
)

type Schema struct {
	Name        string
	ContentJSON []byte
	ContentYAML []byte
}

func Docs(route string, schemas ...Schema) http.Handler {
	route = strings.TrimSuffix(route, "/")

	schemaNames := make(map[string]string)
	schemaContents := make(map[string][]byte)

	for _, schema := range schemas {
		if len(schema.ContentJSON) > 0 {
			url := fmt.Sprintf("%s/schemas/%s.json", route, schema.Name)
			schemaNames[url] = schema.Name
			schemaContents[url] = schema.ContentJSON
		}

		if len(schema.ContentYAML) > 0 {
			url := fmt.Sprintf("%s/schemas/%s.yaml", route, schema.Name)
			schemaNames[url] = schema.Name
			schemaContents[url] = schema.ContentYAML
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc(fmt.Sprintf("%s/swagger-initializer.js", route), func(w http.ResponseWriter, r *http.Request) {
		err := swaggerInitializerTmpl.Execute(w, schemaNames)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc(fmt.Sprintf("%s/", route), func(w http.ResponseWriter, r *http.Request) {
		content, ok := schemaContents[r.URL.Path]
		if ok {
			w.Write(content)
			return
		}
		docsFS.ServeHTTP(w, r)
	})

	return mux
}
