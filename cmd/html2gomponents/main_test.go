package main

import (
	"os"
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	entries, err := os.ReadDir("testdata")
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
			out := readTestData(t, "_out/"+name+".go")

			r := strings.NewReader(in)
			var w strings.Builder
			err := start(r, &w)

			if err != nil {
				t.Fatal(err)
			}
			if out != w.String() {
				t.Fatalf("expected %q, got %q", out, w.String())
			}
		})
	}
}

func readTestData(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile("testdata/" + path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
