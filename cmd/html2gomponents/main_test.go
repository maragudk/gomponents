package main

import (
	"embed"
	"strings"
	"testing"

	"maragu.dev/is"
)

//go:embed testdata
var testdata embed.FS

func TestStart(t *testing.T) {
	entries, err := testdata.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		name = strings.TrimSuffix(name, ".html")

		t.Run(name, func(t *testing.T) {
			in := readTestData(t, name+".html")
			out := readTestData(t, name+".go")

			r := strings.NewReader(in)
			var w strings.Builder
			err := start(r, &w)

			is.NotError(t, err)
			is.Equal(t, out, w.String())
		})
	}
}

func readTestData(t *testing.T, path string) string {
	t.Helper()

	b, err := testdata.ReadFile("testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
