package components_test

import (
	"os"
	"sync"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// cache is a Cache of rendered HTML that lives for the life of the process.
// The mutex makes it safe to use from concurrent renders.
type cache struct {
	mu   sync.RWMutex
	html map[string]string
}

func (c *cache) GetOrSet(key string, f func() (string, error)) (string, error) {
	c.mu.RLock()
	html, ok := c.html[key]
	c.mu.RUnlock()
	if ok {
		return html, nil
	}

	// Call f without holding the lock, because f may use the cache too.
	html, err := f()
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.html[key] = html
	c.mu.Unlock()
	return html, nil
}

var c = &cache{html: map[string]string{}}

// page caches the head per locale, since only the title depends on it.
func page(locale string) g.Node {
	return HTML(
		Cached(c, "head:"+locale, func() g.Node {
			return Head(TitleEl(g.Text("My site ("+locale+")")), Link(Rel("stylesheet"), Href("/app.css")))
		}),
		Body(g.Text("Hello")),
	)
}

func ExampleCached() {
	_ = page("en").Render(os.Stdout)
	// Output: <html><head><title>My site (en)</title><link rel="stylesheet" href="/app.css"></head><body>Hello</body></html>
}
