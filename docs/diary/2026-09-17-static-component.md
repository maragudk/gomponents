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
