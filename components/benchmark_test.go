//go:build go1.24

package components_test

import (
	"bufio"
	"io"
	"strconv"
	"sync"
	"testing"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func BenchmarkJoinAttrs(b *testing.B) {
	flat := func() []g.Node {
		return []g.Node{
			Class("inline-flex items-center"),
			ID("badge"),
			Class("px-2.5 py-0.5 rounded-full"),
			Type("button"),
			Class("bg-green-100 text-green-800"),
		}
	}

	nested := func() []g.Node {
		return []g.Node{
			Class("inline-flex items-center"),
			g.Group{
				ID("badge"),
				g.Group{
					Class("px-2.5 py-0.5 rounded-full"),
					Class("bg-green-100 text-green-800"),
				},
			},
			Type("button"),
		}
	}

	noMatch := func() []g.Node {
		return []g.Node{ID("badge"), Type("button"), Lang("en"), Title("hat")}
	}

	cases := []struct {
		Name     string
		Children func() []g.Node
	}{
		{Name: "flat", Children: flat},
		{Name: "nested groups", Children: nested},
		{Name: "no match", Children: noMatch},
	}

	for _, c := range cases {
		b.Run(c.Name+"/construct", func(b *testing.B) {
			children := c.Children()

			for b.Loop() {
				_ = JoinAttrs("class", children...)
			}
		})

		b.Run(c.Name+"/construct and render", func(b *testing.B) {
			children := c.Children()

			for b.Loop() {
				_ = JoinAttrs("class", children...).Render(io.Discard)
			}
		})
	}
}

// staticTree is a medium-sized tree with nothing dynamic in it: a document head and a navigation
// with a couple of dozen links, the kind of thing [Cached] is for. It is built anew on each call,
// like a component called per request.
func staticTree() g.Node {
	links := make([]g.Node, 24)
	for i := range links {
		links[i] = Li(A(Href("/docs/section-"+strconv.Itoa(i)), Class("nav-link"), g.Text("Section "+strconv.Itoa(i))))
	}
	return g.Group{
		Head(
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			TitleEl(g.Text("Party hats")),
			Link(Rel("stylesheet"), Href("/static/app.css")),
			Link(Rel("icon"), Href("/static/favicon.ico")),
			Script(Src("/static/app.js"), Defer()),
		),
		Nav(Class("navbar"), Ul(Class("nav"), g.Group(links))),
	}
}

// rwMutexCache is a [Cache] on a map guarded by a [sync.RWMutex], the simplest safe implementation.
type rwMutexCache struct {
	mu   sync.RWMutex
	html map[string]string
}

func (c *rwMutexCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	html, ok := c.html[key]
	return html, ok
}

func (c *rwMutexCache) Set(key, html string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.html[key] = html
}

// syncMapCache is a [Cache] on a [sync.Map], whose loads take no lock.
type syncMapCache struct{ m sync.Map }

func (c *syncMapCache) Get(key string) (string, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (c *syncMapCache) Set(key, html string) {
	c.m.Store(key, html)
}

func BenchmarkCached(b *testing.B) {
	// The buffered writer makes the many short writes of a direct render cost something, unlike
	// [io.Discard]. 2 KiB is a typical size for the buffer in front of an HTTP response.
	writers := []struct {
		Name string
		New  func() io.Writer
	}{
		{Name: "discarded", New: func() io.Writer { return io.Discard }},
		{Name: "buffered", New: func() io.Writer { return bufio.NewWriterSize(io.Discard, 2048) }},
	}

	// Two caches, since the cached path costs whatever the cache costs: the RWMutex one contends
	// on its lock across CPUs, and the sync.Map one doesn't.
	caches := []struct {
		Name string
		New  func() Cache
	}{
		{Name: "rwmutex", New: func() Cache { return &rwMutexCache{html: map[string]string{}} }},
		{Name: "syncmap", New: func() Cache { return &syncMapCache{} }},
	}

	for _, writer := range writers {
		// Built once and rendered repeatedly, to separate the cost of rendering from the cost
		// of building the tree.
		b.Run("direct/render pre-built/"+writer.Name, func(b *testing.B) {
			tree := staticTree()
			w := writer.New()

			for b.Loop() {
				_ = tree.Render(w)
			}
		})

		// Called on every iteration, like a component called per request, so the tree is built
		// and rendered every time.
		b.Run("direct/construct and render/"+writer.Name, func(b *testing.B) {
			w := writer.New()

			for b.Loop() {
				_ = staticTree().Render(w)
			}
		})

		// The same, from every CPU at once, like a server rendering requests concurrently.
		b.Run("direct/construct and render parallel/"+writer.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				w := writer.New()

				for pb.Next() {
					_ = staticTree().Render(w)
				}
			})
		})

		for _, impl := range caches {
			// Cached builds and renders the tree on the first call only, so every cached row
			// measures the hit path. This one renders a single node repeatedly.
			b.Run("cached "+impl.Name+"/render pre-built/"+writer.Name, func(b *testing.B) {
				tree := staticTree()
				node := Cached(impl.New(), "tree", func() g.Node { return tree })
				w := writer.New()

				for b.Loop() {
					_ = node.Render(w)
				}
			})

			// This one calls Cached on every iteration, like a component called per request.
			b.Run("cached "+impl.Name+"/construct and render/"+writer.Name, func(b *testing.B) {
				store := impl.New()
				w := writer.New()

				for b.Loop() {
					_ = Cached(store, "tree", staticTree).Render(w)
				}
			})

			// This is where a shared lock or cache line in the cached path shows up as contention.
			b.Run("cached "+impl.Name+"/construct and render parallel/"+writer.Name, func(b *testing.B) {
				store := impl.New()

				b.RunParallel(func(pb *testing.PB) {
					w := writer.New()

					for pb.Next() {
						_ = Cached(store, "tree", staticTree).Render(w)
					}
				})
			})
		}
	}
}
