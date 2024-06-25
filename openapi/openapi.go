package openapi

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
	swaggerDocs embed.FS
	//go:embed swagger-initializer.js
	swaggerInit string
)

type Schema struct {
	Name    string
	Content []byte
}

func SwaggerDocs(prefix string, schemas ...Schema) http.Handler {
	hanler, err := NewSwaggerDocsHandler(prefix, schemas...)
	if err != nil {
		panic(err)
	}

	return hanler
}

func NewSwaggerDocsHandler(prefix string, schemas ...Schema) (http.Handler, error) {
	prefix = strings.TrimSuffix(prefix, "/")
	prefix = strings.TrimSuffix(prefix, "/*")

	overrides := make(map[string][]byte)
	initParams := make(map[string]string)

	for _, schema := range schemas {
		schemaURL := fmt.Sprintf("%s/schemas/%s", prefix, schema.Name)
		overrides[schemaURL] = schema.Content
		initParams[schemaURL] = schema.Name
	}

	docsFS, err := fs.Sub(swaggerDocs, "docs")
	if err != nil {
		return nil, err
	}

	initTmpl, err := template.New("swagger-initializer").Parse(swaggerInit)
	if err != nil {
		return nil, err
	}

	var initBuf bytes.Buffer

	err = initTmpl.Execute(&initBuf, initParams)
	if err != nil {
		return nil, err
	}

	initURL := fmt.Sprintf("%s/swagger-initializer.js", prefix)
	overrides[initURL] = initBuf.Bytes()

	return &swaggerDocsHandler{
		docs:      http.StripPrefix(prefix, http.FileServer(http.FS(docsFS))),
		overrides: overrides,
	}, nil
}

type swaggerDocsHandler struct {
	docs      http.Handler
	overrides map[string][]byte
}

func (h *swaggerDocsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
