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

func (c *cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	html, ok := c.html[key]
	return html, ok
}

func (c *cache) Set(key, html string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.html[key] = html
}

var c = &cache{html: map[string]string{}}

// page renders the head once per locale, since only the title depends on it.
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
