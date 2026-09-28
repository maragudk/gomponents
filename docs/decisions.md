# Project Decisions

This document records significant architectural and design decisions made throughout the project's development.

## 2026-09-28: gomponents caches rendered HTML, but owns no cache

Issue #323 asked for a way to avoid re-rendering large element trees that never change. PR #325 proposed `Static(node) Node`, which rendered once into a `sync.Once` closure. Over several iterations the design moved to a caller-owned slot, then to a caller-owned cache, and the narrow static-only version was dropped entirely.

Context: gomponents builds its element tree at runtime on every request, so a large tree that is identical every time pays its full render cost every time. Code-generating libraries like templ and quicktemplate turn static HTML into string literals and have nothing to cache; `html/template` caches the parse but re-walks the tree on every `Execute`. The gap is real, but there is no library to copy an API from.

Alternatives considered:

- A `Static(node)` or `Static(f)` helper held in a package-level variable. Rejected because it cannot be written inline inside a component tree, which is how every other gomponents helper is used.
- A caller-owned `*string` slot synchronized by a package-global mutex inside the library. Rejected after measurement: the global lock put an atomic read-modify-write on one process-wide cache line into every render of every such node, costing 85 ns per render on ten cores against 0.73 ns for a lock-free design, on a feature whose entire justification is render cost.
- A caller-owned struct wrapping `atomic.Value`. Fast and correct, but it only ever caches content guaranteed never to change, and it requires the library to export and own a synchronized type.
- A library-owned global map keyed by a caller-supplied string, and keying on the call site via `runtime.Callers`. Both rejected for reintroducing global state, and the latter for being magic in a library whose pitch is that it is just Go.

Decision: ship `Cached(cache Cache, key string, f func() Node) Node`, where `Cache` is a two-method interface the caller implements. The library holds no cache, no lock, and no global state, and ships no cache implementation. Keying and invalidation belong to the caller, who is the only party that knows when content changes and the only party with the request in scope from which a key can be derived; a node renders to an `io.Writer` and has no access to either.

Tradeoffs: the simplest case gets harder, since someone who only wants a static `<head>` rendered once must now supply a cache rather than declare a variable, and the fast path becomes whatever the caller's cache costs rather than a single inlined atomic load. Both were accepted deliberately: anyone reaching for a cache can spare the few lines, and one component that covers static and mostly-static content is worth more than a narrow one that only covers content that can never change. The `Cache` interface deliberately has no error returns, because a failed read is indistinguishable from a miss and a failed write simply means rendering again, so both degrade correctly, and an interface implemented by users cannot grow a method later without breaking every implementer.
