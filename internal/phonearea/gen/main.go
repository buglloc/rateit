package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"
)

//go:embed out.gotmpl
var outTmpl string

func fatalf(msg string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "areagen: "+msg, a...)
	os.Exit(1)
}

func isDigitsOnly(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

type AreaInfo struct {
	Country  string
	Template string
}

func main() {
	var outPath, pkgName string
	flag.StringVar(&outPath, "out", "phone_areas.go", "output filename")
	flag.StringVar(&pkgName, "pkg-name", "phonearea", "output pkg name")
	flag.Parse()

	var phoneMap map[string]string
	rsp, err := resty.New().R().
		SetResult(&phoneMap).
		Get("http://country.io/phone.json")
	if err != nil {
		fatalf("unable to request country.io: %w", err)
	}

	if !rsp.IsSuccess() {
		fatalf("non-200 status: %s", rsp.Status())
	}

	var areas []AreaInfo
	for c, a := range phoneMap {
		if !isDigitsOnly(a) || a == "" {
			//some piece of shit
			continue
		}

		var tmpl strings.Builder
		tmpl.WriteByte('+')
		tmpl.WriteString(a)
		for i := tmpl.Len(); i <= 11; i++ {
			tmpl.WriteByte('#')
		}

		areas = append(areas, AreaInfo{
			Country:  c,
			Template: tmpl.String(),
		})
	}

	tmpl := template.Must(
		template.New("errs").Parse(outTmpl),
	)

	out, err := os.Create(outPath)
	if err != nil {
		fatalf("unable to create output file: %v", err)
	}
	defer func() { _ = out.Close() }()

	err = tmpl.Execute(out, map[string]interface{}{
		"PkgName": pkgName,
		"Areas":   areas,
	})
	if err != nil {
		fatalf("render failed: %v", err)
	}
}
