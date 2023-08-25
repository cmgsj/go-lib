package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

var (
	//go:embed sqldb.tmpl
	sqldb string

	tmpl = template.Must(template.New("sqldb").Parse(sqldb))
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	var packageName string
	var outputPath string

	flag.StringVar(&packageName, "p", packageName, "sqlc package name")
	flag.StringVar(&outputPath, "o", outputPath, "output path")
	flag.Parse()

	if packageName == "" {
		return fmt.Errorf("package name is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	dir := filepath.Join(outputPath, "sqldb")
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	file := filepath.Join(dir, "sqldb.go")
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { err = f.Close() }()

	err = tmpl.Execute(f, packageName)
	if err != nil {
		return err
	}

	fmt.Println("generated", f.Name())

	return nil

}
