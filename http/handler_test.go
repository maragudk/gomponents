package http_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	g "github.com/maragudk/gomponents"
	ghttp "github.com/maragudk/gomponents/http"
)

func TestAdapt(t *testing.T) {
	t.Run("renders a node to the response writer", func(t *testing.T) {
		h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return g.El("div"), nil
		})
		code, body := get(t, h)
		if code != http.StatusOK {
			t.Fatal("status code is", code)
		}
		if body != "<div></div>" {
			t.Fatal("body is", body)
		}
	})

	t.Run("renders nothing when returning nil node", func(t *testing.T) {
		h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return nil, nil
		})
		code, body := get(t, h)
		if code != http.StatusOK {
			t.Fatal("status code is", code)
		}
		if body != "" {
			t.Fatal(`body is`, body)
		}
	})

	t.Run("errors with 500 if node cannot render", func(t *testing.T) {
		h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return erroringNode{}, nil
		})
		code, body := get(t, h)
		if code != http.StatusInternalServerError {
			t.Fatal("status code is", code)
		}
		if body != "error rendering node: don't want to\n" {
			t.Fatal(`body is`, body)
		}
	})

	t.Run("errors with status code if error implements StatusCode method and renders node", func(t *testing.T) {
		h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return g.El("div"), statusCodeError{http.StatusTeapot}
		})
		code, body := get(t, h)
		if code != http.StatusTeapot {
			t.Fatal("status code is", code)
		}
		if body != "<div></div>" {
			t.Fatal(`body is`, body)
		}
	})

	t.Run("errors with 500 if other error and renders node", func(t *testing.T) {
		h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
			return g.El("div"), errors.New("")
		})
		code, body := get(t, h)
		if code != http.StatusInternalServerError {
			t.Fatal("status code is", code)
		}
		if body != "<div></div>" {
			t.Fatal(`body is`, body)
		}
	})
}

type erroringNode struct{}

func (n erroringNode) Render(io.Writer) error {
	return errors.New("don't want to")
}

type statusCodeError struct {
	code int
}

func (e statusCodeError) Error() string {
	return http.StatusText(e.code)
}

func (e statusCodeError) StatusCode() int {
	return e.code
}

func get(t *testing.T, h http.Handler) (int, string) {
	t.Helper()

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(recorder, request)
	result := recorder.Result()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	return result.StatusCode, string(body)
}

func ExampleAdapt() {
	h := ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		return g.El("div"), nil
	})
	mux := http.NewServeMux()
	mux.Handle("/", h)
}
