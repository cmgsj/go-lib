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
	//go:embed all:docs
	docs embed.FS
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
	initParams := make(map[string]string)

	for _, schema := range schemas {
		schemaURL := fmt.Sprintf("%s/schemas/%s", prefix, schema.Name)
		overrides[schemaURL] = schema.Content
		initParams[schemaURL] = schema.Name
	}

	docsFS, err := fs.Sub(docs, "docs")
	if err != nil {
		return nil, err
	}

	initializerTmpl, err := template.New("swagger-initializer").Parse(initializer)
	if err != nil {
		return nil, err
	}

	var initializerBuf bytes.Buffer

	err = initializerTmpl.Execute(&initializerBuf, initParams)
	if err != nil {
		return nil, err
	}

	initURL := fmt.Sprintf("%s/swagger-initializer.js", prefix)
	overrides[initURL] = initializerBuf.Bytes()

	return &docsHandler{
		docs:      http.StripPrefix(prefix, http.FileServer(http.FS(docsFS))),
		overrides: overrides,
	}, nil
}

type docsHandler struct {
	docs      http.Handler
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

	h.docs.ServeHTTP(w, r)
}
