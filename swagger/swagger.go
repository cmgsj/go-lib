package swagger

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"text/template"
)

var (
	//go:embed all:dist
	dist embed.FS
	//go:embed swagger-initializer.js
	initializer string
)

type Schema struct {
	Name    string
	Content []byte
}

func Docs(prefix string, schemas ...Schema) http.Handler {
	hanler, err := NewDocsHandler(prefix, schemas...)
	if err != nil {
		panic(err)
	}

	return hanler
}

func NewDocsHandler(prefix string, schemas ...Schema) (http.Handler, error) {
	prefix = strings.TrimSuffix(prefix, "/")
	prefix = strings.TrimSuffix(prefix, "/*")

	overrides := make(map[string][]byte)
	initializerParams := make(map[string]string)

	for _, schema := range schemas {
		schemaURL := fmt.Sprintf("%s/schemas/%s", prefix, schema.Name)

		overrides[schemaURL] = schema.Content

		initializerParams[schemaURL] = schema.Name
	}

	distFS, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, err
	}

	initializerTmpl, err := template.New("swagger-initializer.js").Parse(initializer)
	if err != nil {
		return nil, err
	}

	var initializerBuf bytes.Buffer

	err = initializerTmpl.Execute(&initializerBuf, initializerParams)
	if err != nil {
		return nil, err
	}

	initializerURL := fmt.Sprintf("%s/swagger-initializer.js", prefix)

	overrides[initializerURL] = initializerBuf.Bytes()

	return &docsHandler{
		dist:      http.StripPrefix(prefix, http.FileServer(http.FS(distFS))),
		overrides: overrides,
	}, nil
}

type docsHandler struct {
	dist      http.Handler
	overrides map[string][]byte
}

func (h *docsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	content, ok := h.overrides[r.URL.Path]
	if ok {
		_, err := w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.dist.ServeHTTP(w, r)
}
