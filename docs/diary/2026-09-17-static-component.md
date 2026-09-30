# Diary: `components.Static` for caching static element trees

Redesign of PR #325 ("Add `Static` func", from issue #323). The PR wrapped a node in a `sync.Once` closure in the core package, which forces the caller to hold the returned node in a package-level variable and rules out inlining it like any other component. Markus did not want to merge it as-is and asked to discuss a design first.

## Step 1: Design the caller-owned slot API

**Author:** main

### Prompt Context

**Verbatim prompt:** "See PR 325. I don't want to merge it as-is. Let's discuss a design." Then, after discussion: "I want this to be the simplest possible for the library consumer, so we can take complexity inside the lib. And this is an edge case component anyway, so I think a global RMmutex is fine. Let's try it out. Make sure to have correctness tests, use synctest if needed. In particular, no races, and concurrent readers/writers should behave well."
**Interpretation:** Replace #325's render-once closure with a helper the caller can inline at the call site, keeping the cached HTML in a caller-owned `*string` slot, with a package-level `sync.RWMutex` inside the `components` package taking care of synchronization.
**Inferred intent:** Ship the feature #323 asked for with the smallest possible surface for library users, at the cost of a little internal complexity, and prove the concurrency story with tests rather than argument.

### What I did

Read PR #325, its diff, issue #323, and the two open review threads. Ran a brainstorm: the fork was what the cache is keyed on (render-once var, global key map, caller-owned slot, or a cache interface); Markus chose the caller-owned slot. I proposed a small exported struct as the slot type; Markus asked why a plain `*string` is racy, then why a mutex inside the function doesn't help, then what the downsides of a global `RWMutex` are, and settled on `*string` plus a global `RWMutex`. Created worktree `static-component` from `origin/main` at `ee8dbbd` and wrote this entry before spawning a builder.

Design handed to the builder:

- `func Static(s *string, node g.Node) g.Node` in `/components/components.go`, not in the core package (Markus's review comment on #325).
- One package-level `sync.RWMutex`. Fast path: `RLock`, read `*s`, if non-empty write it out. Slow path: `Lock`, re-check, render the child into a `strings.Builder`, store, copy, unlock. The write to the destination `io.Writer` happens outside any lock.
- A render error is returned and not stored, so the next render retries. Empty output is never cached; documented rather than fixed.
- The returned node is an element node; not for attributes.

### Why

A mutex captured inside `Static` guards nothing, because the whole point of the slot is that `Static(&head, Head(...))` is called fresh on every request and each call is a new closure; the lock has to outlive the call, so it lives at package level. `*string` over a struct because from the caller's side the ergonomics are identical (one package-level `var`, one `&`) and `*string` has no new type to learn, resets with `s = ""` in tests, and retries after an error for free. A global lock is acceptable because contention only exists during the first render of each slot after startup.

### What worked

Asking one question at a time. Each of Markus's questions removed one layer of my proposal until only the essential design was left.

### What didn't work

Nothing failed at this stage; no code has been written yet.

### What I learned

`go.mod` is at `go 1.18`, so `testing/synctest` (Go 1.25) is not available without either bumping the module's minimum Go version or guarding a test file with `//go:build go1.25`. Passed that constraint to the builder rather than deciding for them.

### What was tricky

Keeping the design honest about what the global lock costs: `RLock` on an uncontended `RWMutex` is not free and can be slower than a plain `Mutex.Lock`, so "RW" is a choice to benchmark inside the PR, not to assume.

### What warrants review

The double-checked locking in `Static`, the concurrency tests, and whether the docs are blunt enough about the tree needing to be independent of request, user, and changing state.

### Future work

Rename to `Cached` if Markus prefers the mechanism-name over the constraint-name; decided to keep `Static` for now since it is the vocabulary of #323 and #325.

## Step 2: Build `Static`, its tests, benchmark, and docs

**Author:** static-builder

### Prompt Context

**Verbatim prompt:** "You are building `components.Static` in the gomponents library, in the worktree at /Users/maragubot/Developer/gomponents/.claude/worktrees/static-component [...] API, in /components/components.go: `func Static(s *string, node g.Node) g.Node` [...] One package-level `sync.RWMutex` in `components` guards every slot. Fast path: `RLock`, read `*s`; if non-empty, unlock and write it to `w`. Slow path: `Lock`, re-check `*s` (another goroutine may have filled it), render `node` into a `strings.Builder`, store into `*s`, copy the value, `Unlock`. [...] If you find the RWMutex behaviour makes a plain `sync.Mutex` strictly better, say so in your report with a benchmark, but implement RWMutex as specified. [...] If something blocks you, stop and report the question instead of guessing on product intent."
**Interpretation:** Implement the design from Step 1 exactly as specified, with the listed tests, a benchmark, one-line doc mentions, small commits, a self-review, and then a PR.
**Inferred intent:** Ship #323 with the smallest user-facing surface and prove the concurrency story with tests, while leaving product decisions to the lead.

### What I did

Wrote the tests first in `/components/components_test.go`, watched them fail to compile, then added `Static` and a `renderStatic` helper to `/components/components.go` (commit `e21ece0`). Added `Static` to the README feature list and the `components/` line in `/AGENTS.md`, which `/CLAUDE.md` symlinks to (`d446f55`). Added `BenchmarkStatic` to `/components/benchmark_test.go` with the same medium tree rendered directly and through `Static`, pre-built and constructed per iteration, to `io.Discard` and a 2048-byte `bufio.Writer` (`b1ebd2b`). Ran the code-review skill with two competing reviewers, then fixed what they agreed on and the single-reviewer findings that were real bugs, in the commit that carries this entry.

The review fixes: `Static(&s, nil)` panicked with a nil dereference inside the write lock, while `nil` renders as nothing everywhere else in the library and `g.If(false, ...)` returns `nil`, so `renderStatic` now treats a nil node as empty output. The godoc sentence "Reset a slot to the empty string to render the tree again" recommended an unsynchronised write, so it now says that is only safe while nothing is rendering, for example in tests. The godoc also says the node is still constructed on every call, and that the first render holds a lock shared by all slots, so the node must not contain another `Static`. `html := *s` shadowed the `html` import and is now `cached`. New subtests: nil node, two call sites sharing one slot, reset, element type (`Div(Static(&slot, Class("hat")))` renders `<div> class="hat"</div>`), and a panicking node releasing the lock. The test "lets readers of a filled slot finish while another slot's first render is blocked" claimed a property the code does not have and is renamed to what it asserts. The README now says "caching the rendered HTML of static element trees", since the tree itself is rebuilt on every call.

### Why

Tests before code because the concurrency behaviour is the whole point, and the design record asked for proof rather than argument. The re-check branch under `Lock` is only reachable when a goroutine reads an empty slot, loses the race for the write lock, and then finds the slot filled; the first version of the concurrent test never reached it, because the first goroutine took the write lock before any other goroutine had read the slot, so the rest all parked in `RLock` and took the fast path afterwards. The test "renders an empty slot once when its readers wake together after another slot's first render" fixes that: 32 readers of an empty slot A queue behind a blocked first render of slot B, and when B unlocks they all wake with an empty A and race for the write lock, so all but one hit the re-check. Coverage went from 99.2% to 100% in every run after that.

### What worked

Making the child block on a channel until every goroutine had started, with an atomic counter and `runtime.Gosched` spin, gives the concurrent tests real overlap without sleeps, and they pass `go test -race -shuffle on -count=20 ./components/...` in about 1.4 seconds and under `GOMAXPROCS=1`. Two reviewers independently reproduced the same deadlock and the same data race, which made the priorities obvious.

### What didn't work

`go test -coverprofile` reported `renderStatic 91.7%` five runs out of five with the first concurrent test, for the reason above. The `render` test helper I wrote first was unused after switching to `assert.Equal`, and `golangci-lint` flagged it: `func render(t *testing.T, n g.Node) string ... unused: 1`; removed. The harness refused two compound shell commands (a Python heredoc editing a file, and a shell heredoc appending to a file) as "too complex to verify that it stays inside the worktree", so those edits went through the Edit tool instead.

### What I learned

`testing/synctest` cannot help here: its documentation lists locking a `sync.Mutex` or `sync.RWMutex` as not durably blocking, so `synctest.Wait` cannot observe goroutines parked in `RLock`, and a bubble with goroutines waiting on the lock would deadlock-panic rather than make the test deterministic. So the `//go:build go1.25` file the brief allowed was not worth adding. On the lock choice: a plain `sync.Mutex` and the `RWMutex` measure the same on the cached path, 5.0 versus 5.1 ns/op serially and 85 versus 86 ns/op with `b.RunParallel` on ten cores, because both do one atomic operation on a shared cache line per render, so neither is strictly better and the RWMutex stays as specified. `Static` saves the render, not the construction: `Head(...)` is still evaluated on every call, so "construct and render" through `Static` still allocates 252 times per op, the same as rendering directly, and takes about half the time rather than a hundredth.

### What was tricky

The deadlock the reviewers found is inherent to the specified design, not a slip: the write lock must be held across `node.Render` to guarantee the child renders exactly once, and any `Static` inside that tree then blocks on the same non-reentrant lock forever, taking every other `Static` in the process with it. The same lock-across-user-code is why a slow, failing, or empty first render of one slot delays every slot. Rendering outside the lock and only comparing-and-storing under `Lock` removes all of it, but it changes the specified slow path and turns "child invoked exactly once" into "slot filled once, every goroutine gets the same output". That is the lead's call, so this step stops before the PR with the question, and the godoc documents the restriction as the code stands.

### What warrants review

`renderStatic` and the nesting restriction in the godoc. The two choreographed concurrency tests, in particular that the `entered` counter plus `runtime.Gosched` spin cannot hang on one CPU. Whether "renders as an element, not an attribute" pinning the odd `<div> class="hat"</div>` output is the right way to document that `Static` is not for attributes.

### Future work

Decide between rendering outside the lock (nesting works, no user code under the global lock, benign duplicate renders on a cold start) and keeping exactly-once with the documented restriction. If the former, drop `renderStatic`, fold the compare-and-store into `Static`, and change the concurrent test's `renders == 1` assertion. `LLMs.md` no longer exists on main (deleted in `833b404`), so the brief's mention there went to `/AGENTS.md` instead.

## Step 3: Render outside the lock

**Author:** static-builder

### Prompt Context

**Verbatim prompt:** "Go with option (a): render outside the lock, then Lock, store only if *s is still empty, Unlock, write the local copy. Remove the nesting restriction from the godoc and replace it with one sentence saying concurrent first renders of a slot may each render the tree and the first result is kept. Add a test that nested Static works (outer and inner slots both filled, correct output). Change the 64-goroutine test to assert the slot is filled once and every goroutine gets the full output; you may additionally assert the stored value equals the child's output. Keep everything else as-is, including ExampleStatic in the test file. Re-run -race -count=20, coverage, lint, update the diary with a Step 3 for this change, then open the PR as planned with --reviewer markuswustenberg. Still no merge, no auto-merge, no comments on #325/#323. Report the PR URL and the final godoc when done."
**Interpretation:** Trade "the child renders exactly once" for "no user code runs under the shared lock", keeping the `RWMutex` and the `*string` slot, and adjust the tests and docs to match.
**Inferred intent:** Make `Static` safe to compose like any other component, including inside another `Static`, at the cost of a few duplicate renders of static content on a cold start.

### What I did

Folded `renderStatic` back into `Static` in `/components/components.go`: the fast path is unchanged, and the slow path now renders into a `strings.Builder` with no lock held, then takes `Lock`, stores only if `*s` is still empty, unlocks, and writes its own copy. A nil node skips the render, and empty output skips the store and the lock entirely. The godoc's nesting paragraph became one sentence: "Concurrent first renders of a slot may each render the tree, and the first result is kept."

Tests in `/components/components_test.go`: added "fills both slots when one Static is nested in another". The 64-goroutine test no longer blocks the child; it gates the goroutines on a `start` channel, asserts every output, that the slot holds the child's output, and that the child rendered at least once. The filled-slot test now waits for all 32 readers of slot A *before* releasing slot B's blocked first render, which is a real assertion that readers proceed during a first render; under the previous design that `Wait` would have hung. The empty-slot test became "keeps the first result when concurrent first renders of a slot all try to store it": the child blocks until all 32 goroutines are inside it, so all 32 render and then all 32 try to store, and only the first one does. That is what covers the skip-store branch, deterministically now rather than by scheduling luck. The panic test is renamed to "leaves the slot empty when the node panics", since there is no lock to release any more.

### Why

The lead chose composability over exactly-once. With the render outside the lock, the shared `RWMutex` only guards a read and a compare-and-store, so nothing a user's tree does can block another slot, nesting is fine, and a failing or empty tree no longer serialises every `Static` in the process on every request.

### What worked

The blocking-child choreography from Step 2 carried over with the counters pointed at different things: spinning on `renders` instead of `entered` gives a hard guarantee that all goroutines rendered, because nobody can store before the child returns. `go test -race -shuffle on -count=20 ./components/...` passes in 1.4 seconds, also with `GOMAXPROCS=1 -race -count=5`; coverage is 100% in six of six runs; lint is clean; the package compiles under `-lang=go1.18`.

### What didn't work

Nothing failed in this step.

### What I learned

Re-measured `BenchmarkStatic` after the change and everything is within noise of Step 2: cached path 4.9 ns discarded and 28 ns buffered, construct-and-render 2.15 versus 4.32 µs discarded and 2.18 versus 5.68 µs buffered. The cached path's cost is one `RLock`/`RUnlock` pair plus the write, whichever design holds the lock during the first render.

### What was tricky

Keeping the two remaining concurrency tests distinct once both could render on every goroutine: the free-running one asserts the user-visible contract (full output everywhere, slot filled), the blocking one asserts the internal race (everyone renders, first store wins) and exists for branch coverage as much as for behaviour.

### What warrants review

The compare-and-store in `Static`, and whether "the first result is kept" is enough of a warning for trees whose renders can legitimately differ, which the godoc already forbids.

### Future work

None beyond what Step 2 lists.

## Step 4: Take a `func() Node` so construction is skipped too

**Author:** static-builder

### Prompt Context

**Verbatim prompt:** "Markus wants the signature changed to `func Static(s *string, f func() g.Node) g.Node`, like `Iff`, so construction is skipped too. Semantics otherwise unchanged: fast path never calls f; slow path calls f() with no lock held, renders the result outside the lock, stores if still empty. A nil node returned by f renders as nothing and is not cached, same as now. Update the godoc [...], the godoc example and ExampleStatic (`Static(&head, func() Node { return Head(...) })`), README and AGENTS.md usage lines, all tests [...], and the benchmark so the "construct and render" rows go through the closure [...]. Add a diary Step 4 with the new benchmark table." Then: "One more change after the func() Node switch: move ExampleStatic and its package-level `var head string` into their own file /components/example_static_test.go as a whole-file example (only that example function plus the var and imports in the file, so pkg.go.dev renders the declaration too)."
**Interpretation:** Same cache, same lock, same slot; the second argument becomes a constructor that only runs when the slot is empty, and the example moves to a file of its own so the package-level slot shows in the rendered docs.
**Inferred intent:** Make `Static` actually free on the hot path. Step 2's benchmark showed the render was skipped but the tree was still built and allocated on every call, which is most of what #323 wanted gone.

### What I did

Changed the signature in `/components/components.go` to `Static(s *string, f func() g.Node) g.Node`; the slow path now does `if node := f(); node != nil { node.Render(&b) }` with no lock held, and nothing else moved. The godoc lost the "node is still constructed on every call" sentence and says f is only called when the slot is empty, so neither building nor rendering is repeated; the example is `Static(&head, func() Node { return Head(...) })`. Rewrote every `TestStatic` subtest in `/components/components_test.go` to pass a constructor, with the counters now counting calls to f where the test is about how often the tree is built (once and reused, twice after a render error, once despite a write error, three times for empty output, twice after a reset, exactly 32 in the store race), and the nodes that block on channels or panic are returned by f. `ExampleStatic` and its `var head string` moved to `/components/example_static_test.go`, which holds only the example, the variable, and imports, so Go renders it as a whole-file example and pkg.go.dev shows the package-level slot next to the call. In `/components/benchmark_test.go` the "construct and render" rows pass `staticTree` itself as f, and the pre-built rows pass a closure returning the pre-built tree. README and AGENTS.md only name `Static` and carry no usage snippet, so they needed no change.

### Why

`Iff` already sets the precedent for "give me a function so I can skip the work", and with a plain node argument the caller's `Head(...)` ran and allocated on every request whether or not the slot was full. The whole-file example is the fix for the review note in Step 2 that pkg.go.dev shows only the example's body, which hid the one thing the example is meant to show.

### What worked

Rewriting the test block in one go with a small script in the scratchpad, rather than a dozen edits. Passing `staticTree` directly as f in the benchmark, since it already had the right type.

### What didn't work

Nothing failed in this step. All checks pass: `go test -race -shuffle on -count=20 ./components/...` in 1.4 seconds, `GOMAXPROCS=1 -race -count=5`, coverage 100% in six of six runs, `make lint` clean, `-lang=go1.18` compiles.

### What I learned

The construction cost was the whole story for per-request use. Same machine and method as before (Apple M4, Go 1.27.1, six runs, benchstat):

| Case | Direct | `Static` |
|---|---|---|
| render pre-built, discarded | 2.15 µs, 0 allocs | 4.7 ns, 0 allocs |
| render pre-built, buffered | 3.54 µs, 0 allocs | 28.1 ns, 0 allocs |
| construct and render, discarded | 4.21 µs, 8.9 KB, 252 allocs | 4.9 ns, 0 allocs |
| construct and render, buffered | 5.85 µs, 8.9 KB, 252 allocs | 28.9 ns, 0 allocs |

The "construct and render" rows through `Static` went from about half the direct cost with the same 252 allocations to three orders of magnitude less with none, and are now the same as the pre-built rows, which is the point: after the first call, `Static` costs a read lock and one write.

### What was tricky

Deciding what the counters mean in each test once f and the node are separate. The rule used: count f calls when the test is about whether the tree is rebuilt, and keep the blocking or panicking behaviour in the node f returns, since that is where a real tree does its work.

### What warrants review

Whether the godoc's "f is only called when the slot is empty" reads clearly enough that f must be cheap to write but may be expensive to run. The whole-file example on pkg.go.dev once the PR is merged, since that is the only place it can be seen rendered.

### Future work

None.

## Step 5: Drop `Static` for a caller-owned `Cached`

**Author:** main

### Prompt Context

**Verbatim prompt:** "I wonder if it would be better to explore a design where the cache and its invalidation lives with the caller (maybe a TTL cache, or a bounded one, or one that actually uses a cache key that can be invalidated), so it's not entirely for static content. That would be a stronger component." Then, after discussion: "Definitely Cached alone, I don't want a helper for the simple case. An empty string key is just as easy to use IMO."
**Interpretation:** Replace `Static` entirely with `Cached(cache, key, f)`, where the caller supplies a cache implementing a small interface and owns keying and invalidation. No convenience wrapper for the pure-static case.
**Inferred intent:** Ship one component that covers static and mostly-static content instead of a narrow one that only covers content guaranteed never to change, and get the library out of the business of owning cache state.

### What I did

Two research passes before the decision. The first asked what stdlib type should replace the `*string` slot; the second measured what an interface parameter costs. Both are summarised below because their findings drove the pivot.

Then settled the remaining API questions: no errors on the cache interface, a nil cache renders every time rather than panicking, `Get`/`Set` rather than `Load`/`Store`, and a named exported `components.Cache` rather than an anonymous interface in the signature. Wrote this step and a `/docs/decisions.md` entry, then handed the build over.

### Why

The `*string` slot was a hack, though not for the reason first suspected. The empty-string sentinel was cosmetic; the real defect was the package-global `sync.RWMutex`, which put an atomic read-modify-write on one process-wide cache line into every render of every `Static` node in the binary. Measured end-to-end on the cached path: 5.03 ns serial and 85.2 ns on ten cores, against 4.36 ns and 0.73 ns for a slot built on `atomic.Value`. Mutex designs get worse as cores are added; atomic-load designs get better. This went unnoticed because `BenchmarkStatic` had no `b.RunParallel` case, so Step 2's comparison of `Mutex` against `RWMutex` recorded the contention number and read it as a tie between lock flavours rather than as evidence that neither belonged there.

That pointed at `StaticCache`, a struct wrapping `atomic.Value`. The pivot to `Cached` went further for a reason that only emerged from the second measurement: a keyless `Load`/`Store` interface unlocks nothing, because a `Node` renders to an `io.Writer` and has no request in scope from which to derive a key. Relocating one 16-byte process-lifetime string is not a use case. The `Vary`-style, per-tenant, and TTL cases that justify an interface all need a key, and the key can only come from the call site, where the caller does have the request. So the keyed function was always going to be a separate function; the question was whether `Static` deserved to exist beside it. Markus decided it did not, on the grounds that an empty-string key is no harder than a dedicated helper.

Errors were left off the interface because a read failure is indistinguishable from a miss and a write failure just means rendering again next time, so both degrade correctly, and an implementer who wants to log a failure does it inside their own `Get`. An interface implemented by users can never grow a method without breaking all of them, which makes the two-method shape worth protecting.

### What worked

Asking for measurement instead of accepting either my own or the researcher's reasoning about dispatch cost. The interface was found to cost about 0.5 ns serial and eight closure bytes, with no extra allocation, and the explanation mattered more than the number: the concrete version inlines to a bare atomic load, while the interface version becomes an itab call that the compiler provably cannot devirtualize, because the slot is captured by the returned closure and Go's devirtualizer does not follow proven types across a closure boundary. That is not a gap a future Go release closes. It settled the `Static` signature question, and it also showed why the cost is irrelevant for `Cached`, where the caller's own cache lookup dominates.

Checking backwards compatibility by compiling the experiment rather than reasoning about it. Adding `Cached` later would have been fully compatible, which removed all time pressure from that half of the design and left only `Static`'s own exported surface as the thing with a deadline.

### What didn't work

Nothing failed mechanically. The design did move three times — core-package closure, caller-owned `*string` slot, `atomic.Value` struct, caller-owned cache — which is the cost of having started building before the shape was settled. PR #363 now has to be substantially rewritten rather than extended.

### What I learned

A benchmark suite without a parallel case cannot see contention, and contention is the entire failure mode of a shared lock. The lock-versus-no-lock comparison was never run because both sides of the comparison that was run had a lock in them.

Go's devirtualizer does not propagate a proven concrete type into a closure, so any API that returns a closure capturing an interface pays a real itab call however obvious the concrete type is at the call site. A microbenchmark that calls the method directly rather than through the returned node will devirtualize and report the abstraction as free; that result does not transfer.

### What was tricky

Separating the two reasons the `*string` felt wrong. Markus named the sentinel; the measurement found the global lock. Fixing only the named one would have shipped the real defect.

### What warrants review

Whether `Get`/`Set` without errors is the right permanent shape, since users implement this interface and it can never grow a method. And whether dropping the pure-static convenience is the right call for the audience that issue #323 came from, which wanted static trees fast and will now have to supply a cache.

### Future work

None outstanding; the keyed design subsumes what `Static` was for.

## Step 6: Build `Cached`, its tests, benchmark, and docs

**Author:** cached-builder

### Prompt Context

**Verbatim prompt:** "You are replacing the `Static` component with a caller-owned-cache `Cached` component in the gomponents library. [...] PR #363 currently ships `Static(s *string, f func() g.Node) g.Node` plus a package-global `staticMutex`. All of that comes out. [...] In its place, in /components/components.go: `type Cache interface { Get(key string) (html string, ok bool); Set(key, html string) }` and `func Cached(cache Cache, key string, f func() g.Node) g.Node` [...] A nil node returned by `f` renders as nothing, and that empty result IS cached (unlike the old `Static`). [...] A nil `cache` renders every time and must not panic. [...] The library holds no lock and no global state. [...] Replace `BenchmarkStatic` [...] with `BenchmarkCached` [...] **Include `b.RunParallel` variants** [...] Self-review with the `fabrik:code-review` skill when implementation is done, and fix what you find. Push to `origin` branch `static-component` [...] and update PR #363 with `gh pr edit`."
**Interpretation:** Delete `Static` and its lock entirely, implement the `Cache` interface and `Cached` exactly as specified, port the test suite to a counting test cache, move the example to a whole-file `ExampleCached`, add serial and parallel benchmarks, update the two doc mentions, self-review, and rewrite the PR.
**Inferred intent:** Ship the Step 5 design with the concurrency cost made visible in the benchmark this time, and with the empty-output caching bug of the old sentinel design fixed and pinned by a test.

### What I did

Wrote `TestCached` first in `/components/components_test.go`, replacing `TestStatic` wholesale, and watched `go vet` fail with `undefined: Cached`. The tests use a `testCache` (map plus `sync.RWMutex`, with `int32` counters bumped through `atomic.AddInt32` since `atomic.Int32` is Go 1.19) and two assertion helpers, `assertCached` and `assertNotCached`. Then deleted `Static`, `staticMutex`, and the `sync` import from `/components/components.go` and added `Cache` and `Cached` (commit `1a991b4` after the rebase, `723a228` before). The example moved from the deleted `/components/example_static_test.go` to `/components/example_cached_test.go`, a whole-file example with the `cache` type, its methods, `var c`, a `page(locale)` function, and `ExampleCached`. The README bullet and the AGENTS.md `components/` line now name `Cached` (`3cd6e49`). `/docs/decisions.md` and the lead's Step 5 diary entry, both uncommitted when I started, went into the first commit since they document the pivot.

`BenchmarkCached` in `/components/benchmark_test.go` (`72d89ab`) keeps the direct-versus-cached rows from `BenchmarkStatic` and adds a `b.RunParallel` row for each, with the cached rows run against two caches: the map-and-`RWMutex` one and a `sync.Map` one. Ran the code-review skill with two competing reviewers and fixed their consensus findings (`a77d2ec`): the godoc claimed concurrent misses all `Set` "the same HTML" and that the first render wins, which is wrong when two call sites with different trees miss at once, so it now says the last `Set` is kept and trees sharing a key must render the same HTML; the godoc example called an undefined `title(locale)`; "as they are in a web server" and "the caller owns the cache" looked outward from the package; the benchmark comment attributed the 2 KiB buffer to an unexported `net/http` constant; a redundant `int32()` conversion; the benchmark borrowed the example file's `cache` type, so it now has its own `rwMutexCache`; and the loop variables `w` and `c` shadowed each other and the example's package-level `c`. Two single-reviewer test gaps were cheap enough to take: the render-error test now asserts that nothing reaches the writer, and the write-error test also covers the hit path. Rebased onto `origin/main`, which had moved by #357 (render-time escaping of `Text`), with no conflicts.

### Why

The pivot's whole justification was measured, so the benchmark had to be able to see the thing that was measured: a parallel row. Adding a `sync.Map` cache next to the `RWMutex` one turns the benchmark into the argument for the design, since it shows the cached path costs exactly what the caller's cache costs and nothing the library adds. The empty-output-is-cached test exists because that was the one behaviour the old sentinel could not express; the `ok` return from `Get` is what makes it expressible, and the test pins it so it cannot regress to "empty means miss".

### What worked

Writing the whole test block before touching the implementation, since the spec listed the cases and the previous `TestStatic` had the choreography for the concurrent ones; the two blocking-node tests from Step 3 ported with the counters pointed at the cache instead of the slot. The `ok`-based miss detection made the implementation a single straight-line function with no special cases: hit, write; miss, render, `Set`, write.

### What didn't work

The harness refused three compound shell commands as "too complex to verify that it stays inside the worktree": a Python heredoc followed by `git rm` and a `cat` heredoc, and twice a `cat >> file <<'EOF'` append followed by `go vet`. Splitting them into single-purpose commands and using the Edit tool for the appends worked, as Step 2 found. The first full benchmark run overlapped with the reviewers' `go test -race -count=5` runs and the direct rows came back at ±36%, so it was re-run on a quiet machine after the review; the cached rows were within 8% either way.

### What I learned

The parallel row confirms Step 5's measurement and puts a number on the alternative. Same machine and method as Step 4 (Apple M4, Go 1.27.1, six runs, benchstat), with the direct rows now at 251 allocations after #357:

| Case | Direct | `Cached`, `RWMutex` cache | `Cached`, `sync.Map` cache |
|---|---|---|---|
| render pre-built, discarded | 2.24 µs, 0 allocs | 9.3 ns, 0 allocs | 10.7 ns, 0 allocs |
| render pre-built, buffered | 3.64 µs, 0 allocs | 34.5 ns, 0 allocs | 37.0 ns, 0 allocs |
| construct and render, discarded | 4.23 µs, 8.8 KB, 251 allocs | 10.1 ns, 0 allocs | 10.9 ns, 0 allocs |
| construct and render, buffered | 6.33 µs, 8.8 KB, 251 allocs | 34.4 ns, 0 allocs | 37.9 ns, 0 allocs |
| construct and render, parallel, discarded | 2.60 µs, 251 allocs | 86.5 ns, 0 allocs | 2.27 ns, 0 allocs |
| construct and render, parallel, buffered | 3.43 µs, 251 allocs | 91.2 ns, 0 allocs | 7.66 ns, 0 allocs |

The `RWMutex` cache on ten cores lands at 86 ns, the same contention number Step 5 measured for the global lock, which is the point: the lock was never the library's to hold, and a caller who picks `sync.Map` gets 2.3 ns for the same render. Serially the two caches are within 2 ns of each other. Parallel rows are wall time per operation across all cores, so the direct parallel rows being faster than serial is the expected ten-way split, not a speedup.

### What was tricky

Deciding what "whichever rendered first" means once there is no lock. Sequentially it is true and the spec asked for it in the godoc. Concurrently, each miss calls `Set` and a map cache keeps the last one, so the godoc now says both, and adds the rule that makes the difference irrelevant: trees sharing a key must render the same HTML. The reviewers also both asked to drop the in-comment example because `ExampleCached` renders directly below it on pkg.go.dev and the two will drift; the brief asked for both, so it stayed, with the undefined call fixed.

### What warrants review

The godoc for `Cached`, in particular whether keeping the in-comment example next to the whole-file one is wanted. The `Cache` interface is final once released: no errors, no context, and both reviewers noted a remote cache would have to swallow failures inside `Get`, which the decisions entry accepts. The `rwMutexCache` in the benchmark duplicates the example's `cache` type by design, so the benchmark survives changes to the example. Nil `f` panics on the first miss, like `Iff` with a nil function, and is undocumented; both reviewers flagged it as low priority.

### Future work

None. The reviewers' remaining single-reviewer nits, rendering the example twice to show a hit and the nil-cache path buffering before writing, were left as-is: the second is what the brief specified.

## Step 7: Replace `Get` and `Set` with a single `GetOrSet`

**Author:** cached-builder

### Prompt Context

**Verbatim prompt:** "Change the `Cache` interface from `Get`/`Set` to a single `GetOrSet`, the usual atomic caching pattern. [...] Why: an outside adversarial review (gpt-6-astra via codex) found a contract defect, which I verified. With separate `Get` and `Set`, nothing links a miss to its later fill. Sequence: render A misses and captures old content; the app updates the content and deletes the key; render B caches the new HTML; A finishes and `Set`s the old HTML over it. So the godoc's "to render a tree again, remove its key" is unsafe. The same split also makes it impossible for a cache to share one render among concurrent misses (a lock taken in `Get` is never released if `f` errors or panics), and a cache cannot surface its own failures. Markus's response, verbatim: "Don't we just need a GetOrSet with a callback function or something?" and then "The usual atomic caching pattern". New API in /components/components.go: `type Cache interface { GetOrSet(key string, f func() (string, error)) (string, error) }` and `func Cached(cache Cache, key string, f func() g.Node) g.Node // unchanged signature` [...] The hit path was 0 allocs; check whether passing the callback through the interface now allocates per render. If it does, try building the callback once per `Cached` call rather than per render, or any other restructuring that keeps hits allocation-free, and report what you measured either way, serial and parallel, against the previous table."
**Interpretation:** Same `Cached` signature, new one-method interface, a godoc that states the implementer's contract clause by clause, tests that prove the interface makes a correct invalidating cache and a render-sharing cache possible, and a measured answer on allocations.
**Inferred intent:** Fix the contract before it is released and can never change, by giving the cache the whole miss-to-fill operation instead of two unrelated halves of it.

### What I did

Ported the test cache to `GetOrSet` first and added the new cases, then watched `go vet` fail with `*testCache does not implement components.Cache (missing method Get)`. Rewrote `Cache` and `Cached` in `/components/components.go`: `Cached` builds a `render` closure that calls `f`, renders the node into a `strings.Builder`, and returns the HTML or the error; the returned node calls `cache.GetOrSet(key, render)`, or `render()` directly for a nil cache, returns any error without writing, and otherwise writes the HTML. Ported the example cache in `/components/example_cached_test.go` and both benchmark caches in `/components/benchmark_test.go`, each calling `f` with no lock held (commit `c791ba3`, which also carries the lead's update to `/docs/decisions.md`).

`/components/components_test.go` now has four test caches: `testCache` (map and `RWMutex`, one call to `f` per miss, counts calls and fills), `failingCache` (returns its own error without calling `f`), `invalidatingCache` (a generation counter, so a result from before an invalidation is returned but not stored), and `sharingCache` (a per-key mutex held across `f` with `defer`). New tests: an error from the cache itself surfaces and `f` is not called; a node that writes half its output and then fails writes and caches nothing, with and without a cache; the stale-fill sequence from the review, deterministically; a sharing cache runs `f` once for 32 concurrent misses; and the same cache still serves the key after `f` errors and after it panics.

Ran the code-review skill again with two competing reviewers and fixed their consensus findings (`3d49597`). The single-flight test released the render as soon as one goroutine was in it, so the other 31 could arrive after the fill and merely hit; `sharingCache` now counts calls and the test waits until all 32 have entered `GetOrSet` before releasing. The contract said nothing about `f` panicking, though the panic test asserted a behaviour for it, so the `Cache` godoc now says a panic must not leave `GetOrSet` unable to serve later calls. Nothing said a tree may not contain a `Cached` node with its own cache and key, which deadlocks any cache that holds a per-key lock across `f`, as the contract allows; the `Cached` godoc now says so. `testCache` had a `Delete` that stored stale fills, violating the clause it was documented under, so it is gone and the invalidation test uses `invalidatingCache`. The godoc said "nothing is written to w" where no `w` is in scope. The recover test would have hung rather than failed, so its last renders run against a ten-second deadline.

Afterwards `/docs/decisions.md`, which held only the entry for this feature, was removed from the branch at Markus's request ("drop the decision, not needed").

### Why

With two methods the cache sees a read and, some time later, an unrelated write, and cannot tell that the write belongs to a read from before an invalidation. With one method the cache is on the stack for the whole operation, so it can record a generation before calling `f` and compare after, hold a per-key lock across `f` and release it with `defer` whatever `f` does, and return an error instead of pretending to miss. None of that is the library's to implement, but the interface has to make it possible, and the tests exist to prove that it does rather than to test the test caches.

### What worked

Writing each new test cache as the smallest thing that satisfies one clause of the contract. The stale-fill test needed no sleeps: the first render closes a channel from inside its node after `f` has already read the old content, which gives the main goroutine a happens-before edge to change the content, invalidate, and cache the new HTML before releasing the first render.

### What didn't work

`gofmt -l .` flagged `/components/components.go` after the rewrite for a trailing blank line, fixed with `gofmt -w`. My first `GetOrSet` call site named its result `html`, shadowing the `html` import in that file, the same trap Step 2 recorded; it is `rendered`. The first attempt to benchmark only the cached rows, `-bench 'BenchmarkCached/cached.*/discarded'`, matched nothing and printed only `PASS`, because slash-separated patterns match one sub-benchmark level each and the writer name is the fourth level; `-bench 'BenchmarkCached/cached//discarded'` worked.

### What I learned

The callback costs one 16-byte allocation per `Cached` call, not per render. The `render` closure captures `f` and is passed through an interface method, so it escapes and cannot live on the stack; the outer node closure still does. A node that is built once and rendered repeatedly therefore stays at zero allocations, and a component that calls `Cached` on every request pays one. Building the closure inside the node instead, per render, was measured and is strictly worse: one allocation in every row, including the reused node, at 12.2 ns against 8.5 ns. A `sync.Pool` of callbacks could remove the allocation but would be library-owned global state, which the design rules out, and would break any cache that kept `f` past its own return. Same machine and method as Step 6:

| Case | Direct | `Cached`, `RWMutex` cache | `Cached`, `sync.Map` cache |
|---|---|---|---|
| render pre-built, discarded | 2.23 µs, 0 allocs | 8.7 ns, 0 allocs (was 9.3 ns) | 10.7 ns, 0 allocs (was 10.7 ns) |
| render pre-built, buffered | 3.58 µs, 0 allocs | 31.8 ns, 0 allocs (was 34.5 ns) | 36.4 ns, 0 allocs (was 37.0 ns) |
| construct and render, discarded | 4.05 µs, 8.8 KB, 251 allocs | 13.3 ns, 16 B, 1 alloc (was 10.1 ns, 0) | 15.1 ns, 16 B, 1 alloc (was 10.9 ns, 0) |
| construct and render, buffered | 5.44 µs, 8.8 KB, 251 allocs | 37.1 ns, 16 B, 1 alloc (was 34.4 ns, 0) | 41.1 ns, 16 B, 1 alloc (was 37.9 ns, 0) |
| construct and render, parallel, discarded | 2.40 µs, 251 allocs | 95.2 ns, 16 B, 1 alloc (was 86.5 ns, 0) | 7.2 ns ±28%, 16 B, 1 alloc (was 2.27 ns, 0) |
| construct and render, parallel, buffered | 2.76 µs, 251 allocs | 92.1 ns, 16 B, 1 alloc (was 91.2 ns, 0) | 10.4 ns ±30%, 16 B, 1 alloc (was 7.66 ns, 0) |

Serially the allocation costs about 3 to 4 ns per call. The parallel `sync.Map` row shows it most, 2.3 to 7.2 ns with wide variance, because at that speed the allocator is the only shared thing left; it is still 330 times faster than rendering directly.

### What was tricky

The contract allows a cache to share one call to `f` among concurrent misses, and the obvious way to do that is a per-key lock held across `f`. That is sound only if `f` never comes back for the same key, which the library cannot prevent, so it has to be a rule on the tree rather than on the cache. One reviewer also noted that two keys nested in opposite orders on two pages would deadlock such a cache the same way; that is a property of per-key locking in the cache, not something `Cached` can rule out, and it is not in the godoc.

### What warrants review

The `Cache` godoc, clause by clause, since it is the contract every implementer will read and the interface can never change. Two sentences in it were not in the brief: that `f` may panic, and, on `Cached`, that a tree must not contain a `Cached` node with its own cache and key. The one-allocation cost on the construct-per-request path. One reviewer asked what waiters get when a shared call to `f` fails, the same error or a call of their own; the godoc leaves that to the cache, and `sharingCache` lets each waiter call `f` itself.

### Future work

None required. If the per-call allocation ever matters, the only fix that keeps the library stateless is for the component to be built once and reused, which already costs nothing.

## Step 8: Panic on a nil cache, and keep `render` outside the node

**Author:** cached-builder

### Prompt Context

**Verbatim prompt:** Markus left two inline comments on PR #363. On the `if cache != nil { ... } else { rendered, err = render() }` branch: "panic on nil cache instead in the beginning, makes no sense to have". On `render := func() (string, error) {`: "inline". The lead's instruction: "Make a nil check with a panic the first statement of `Cached` (before the node is built, so it fires when `Cached` is called, not at render time). [...] Remove the `else` branch so the node always calls `cache.GetOrSet(key, render)`. [...] Decision: skip, do NOT inline. Reason: inlining the literal into the `GetOrSet` call creates the closure on every render, which you measured at 12.2 ns / 1 alloc against 8.5 ns / 0 allocs for a node built once and rendered repeatedly. Add a one-line comment above `render` saying it is created here, once per call to `Cached`, so that rendering a reused node does not allocate."
**Interpretation:** Apply the first comment as written, and answer the second with a comment in the code instead of the change.
**Inferred intent:** A component called `Cached` that silently does not cache is a trap, so a nil cache should fail loudly and early; and the one place where the code looks needlessly indirect should say why it is that way.

### What I did

Changed the nil-cache test in `/components/components_test.go` to expect a panic from the call to `Cached` itself, with no render, and watched it fail with `expected a panic with "cache must not be nil", got <nil>`. Dropped the nil-cache half of the partial-output test. In `/components/components.go`, `Cached` now starts with `if cache == nil { panic("cache must not be nil") }`, the node always calls `cache.GetOrSet(key, render)`, the godoc paragraph about a nil cache rendering every time is replaced by "Cached panics if cache is nil.", and `render` has a comment saying it is created once per call to `Cached` so that rendering a reused node doesn't allocate. Removed the nil-cache sentence from the PR body; README and AGENTS.md never mentioned it.

### Why

The nil-cache off switch came from the Step 5 brief, for `var cache Cache; if !dev { cache = mine }`. Markus's view is that it makes no sense to have, and the panic turns a forgotten cache from a silent performance bug into an immediate one. Panicking when `Cached` is called rather than when the node renders puts the stack trace at the call site that passed nil.

### What worked

The message follows the library's one existing panic, `panic("attribute must be just name or name and value pair")` in `/gomponents.go`: lowercase, no package prefix, stating the requirement.

### What didn't work

Nothing failed in this step. `make lint`, `make test`, `go test -race -shuffle on -count=20 ./components/...`, and `go vet -gcflags=-lang=go1.18 ./components/` pass, and coverage is 100% in three of three runs.

### What I learned

The early panic does not disturb the escape analysis: the reused-node row is still 8.1 ns with no allocation and the construct-per-render row 12.9 ns with one.

### What was tricky

Nothing. The test asserts the panic without rendering, so it cannot pass if the check drifts back into the node.

### What warrants review

Whether `"cache must not be nil"` is the wanted wording. Anyone who relied on the nil off switch now needs a cache that always misses, which is three lines.

### Future work

None.
