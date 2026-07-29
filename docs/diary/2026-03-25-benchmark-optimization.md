# Diary: Benchmark optimization via autoresearch

Autonomous experiment loop to optimize gomponents benchmark performance. Started with a synthetic benchmark (10K elements), pivoted to a realistic full-page dashboard benchmark mid-session. Ended with two PRs submitted upstream.

## Step 1: Initial profiling and failed experiments (iterations 1-5)

### Prompt Context

**Verbatim prompt:** "Benchmark numbers" (in response to autoresearch setup asking what to optimize)
**Interpretation:** Optimize the gomponents benchmark numbers, specifically ns/op.
**Inferred intent:** Make the library faster for real-world usage patterns.

### What I did

Ran `go test -bench . -benchmem ./...` to identify benchmarks. Focused on `BenchmarkLargeHTMLDocument` (synthetic: 10K div/span elements with class attributes). Profiled with `-cpuprofile` and `-memprofile` to understand the hot path.

Tried five experiments:
1. **Combined tag writes** (`<` + name into one string) -- added 20K allocs from string concat, 10% slower.
2. **Separate attrs/elements at creation time** -- failed tests because Groups containing attributes weren't unwrapped.
3. **Struct-based element + pre-rendered attributes** -- 12% less memory but 5% slower ns/op. Struct method dispatch slower than closure calls.
4. **Pre-render attribute strings only** -- 9% less memory but ns/op unchanged. Write count reduction didn't matter because `strings.Builder` writes are already very cheap.
5. **CPU profiling deep dive** -- discovered 73% of CPU was in GC/madvise/runtime, only 10% in actual rendering code. This reframed the entire optimization strategy.

### Why

Needed to understand the performance profile before making targeted changes. The profiling in step 5 was the key insight: reducing allocations is the only path to meaningful improvement.

### What worked

The profiling approach. Understanding that 73% of CPU was GC-related immediately eliminated entire categories of optimization (reducing write calls, inlining functions, etc.) and pointed to allocation reduction as the only viable strategy.

### What didn't work

- String concatenation at El creation time (experiment 1): pre-computing `"<div"` and `"</div>"` strings added one allocation per element, which is worse than the two separate writes it replaced.
- Struct-based nodes (experiment 3): `*element` with method receivers was marginally slower than `NodeFunc` closures despite identical allocation counts. Go closures are well-optimized.
- Pre-rendering attribute strings (experiment 4): reducing 5 write calls to 1 per attribute had no measurable impact because `strings.Builder.WriteString` is essentially just a memcpy into a growing buffer.

### What I learned

- Go closures are surprisingly efficient -- a closure capturing two strings is roughly equivalent to a struct with two string fields in terms of allocation cost and dispatch speed.
- `strings.Builder` write calls are so cheap that reducing write count doesn't move the needle. The bottleneck is allocation count, not I/O operations.
- `template.HTMLEscapeString("foo")` returns the input string unchanged (no allocation) when there are no special characters to escape.

### What was tricky

Experiment 2 (separating attrs/elements at creation time) seemed promising but failed because `Group` nodes can contain both attribute and element children, and the pre-separation didn't recursively unwrap Groups. The existing two-pass design in `renderChild` handles this correctly but can't be easily short-circuited.

### What warrants review

Nothing from this step made it into the final PRs. All experiments were discarded.

### Future work

The allocation profile breakdown (per benchmark iteration): ~10K `attrFunc` closures, ~10K `NodeFunc` closures, ~10K variadic `[]Node` slices, ~10K variadic `[]string` slices. The variadic slices were identified as the most promising target.

## Step 2: The variadic escape breakthrough (iteration 6)

### Prompt Context

**Verbatim prompt:** (continuation of autoresearch loop)
**Interpretation:** Find a way to reduce the 40K allocations per benchmark iteration.
**Inferred intent:** Eliminate unnecessary heap allocations in the hot path.

### What I did

Split `Attr(name string, value ...string)` into three functions:
- `Attr` (public, variadic) dispatches to one of:
- `booleanAttr(name string)` -- captures just the name string
- `valueAttr(name, value string)` -- captures name and extracted value

The key: by extracting `value[0]` from the variadic slice *before* creating the closure, the `[]string` variadic parameter no longer escapes to the heap. The closure captures plain `string` values instead of a `[]string`.

Result on the old synthetic benchmark: **-25% allocations (40,028 → 30,028), -7.8% ns/op**.

### Why

The `Attr` closure previously captured the `value ...string` parameter, which is a `[]string` slice. Even though only `value[0]` was ever used, the slice escaped to the heap because the closure captured it. By extracting the string before closure creation, the compiler can see the slice doesn't escape.

### What worked

Go's escape analysis is precise enough that extracting `value[0]` before the closure creation eliminates the heap allocation for the variadic slice. Verified with `go build -gcflags='-m=1'` -- the `leaking param` annotation disappeared for the value parameter.

### What didn't work

Nothing -- this was a clean win. The only test change needed was adjusting the error-write-count test (`TestEl/returns_render_error_on_cannot_write`) because `booleanAttr` does fewer writes than the original `Attr` with zero values.

### What I learned

Go's escape analysis tracks whether captured variables escape through closures. If a closure captures a slice but you extract a scalar from the slice before creating the closure, the slice may not need to escape. This is a general optimization pattern applicable anywhere variadic parameters are captured by closures.

### What was tricky

The same technique doesn't work for `El(name string, children ...Node)` because the closure genuinely needs the entire `children` slice to iterate over during rendering. You can't extract "just one child" because elements can have any number of children.

### What warrants review

PR #305: https://github.com/maragudk/gomponents/pull/305

The public API is unchanged. `Attr` still accepts variadic `...string`. The internal dispatch to `booleanAttr`/`valueAttr` is invisible to callers.

### Future work

The remaining 30K allocations are split between closures (unavoidable with current design) and variadic `[]Node` slices in element constructors (would require API changes to eliminate).

## Step 3: Pivot to realistic benchmark (iteration 8)

### Prompt Context

**Verbatim prompt:** "if you rebase on origin/main you've got a new, more realistic benchmark"
**Interpretation:** Upstream added a better benchmark; rebase and use it.
**Inferred intent:** Optimize against a benchmark that reflects actual usage patterns.

### What I did

Rebased `autoresearch` branch on `upstream/main` which included `BenchmarkRealisticPage` -- a full dashboard page with navigation, sidebar, content cards, a 50-row data table, and footer. Uses `components.HTML5`, `components.Classes`, `Map`, `If`, `Text`, `Textf`, `Raw`, and many element/attribute combinations.

Two sub-benchmarks:
- `construct_and_render` -- builds tree from scratch and renders (the realistic case)
- `render_pre-built_tree` -- renders a pre-built tree (edge case, not typical)

Established new baseline on upstream/main: **279,700 ns/op, 4,002 allocs, 359,968 B** for construct+render.

Re-measured experiment 6 against the new benchmark: **263,400 ns/op, 3,074 allocs, 345,066 B** -- a solid **-5.8% ns/op, -23.2% allocs** on the realistic workload.

### Why

The synthetic benchmark (10K identical elements) was dominated by GC pressure from sheer allocation volume. The realistic benchmark better represents actual library usage: diverse elements, mixed attributes, text nodes, conditional rendering, mapped collections.

### What worked

The experiment 6 optimization carried over well to the realistic benchmark, confirming it's a genuine improvement rather than a synthetic-benchmark-specific win.

### What didn't work

Nothing -- the rebase was clean and the measurements were straightforward.

### What I learned

Always validate optimizations against realistic workloads. The synthetic benchmark had 40K allocs per iteration where GC dominated at 73% of CPU. The realistic benchmark has ~3-4K allocs where the application code is a larger fraction of total time.

### What was tricky

Getting the correct baseline comparison required measuring upstream/main in a separate worktree to avoid contamination from our changes.

### What warrants review

The decision to focus exclusively on `construct_and_render` rather than `render_pre-built_tree` was made because the typical gomponents usage pattern builds a fresh tree per HTTP request with dynamic data.

### Future work

With the realistic benchmark established, smaller optimizations (1-3% each) become worth pursuing since they compound.

## Step 4: Switch and renderChild optimizations (iterations 9-10)

### Prompt Context

**Verbatim prompt:** (continuation of autoresearch loop)
**Interpretation:** Find additional micro-optimizations now that the big allocation win is landed.
**Inferred intent:** Squeeze out remaining performance from the hot path.

### What I did

**Iteration 9:** Replaced `isVoidElement` map lookup with a switch statement. The 16 void HTML elements are a fixed set defined by the HTML5 spec. A switch generates more efficient code than a map hash+lookup for small fixed sets.

Result: **-1.3% ns/op** (260,500 vs 263,400). Also eliminates a package-level `map[string]struct{}` allocation.

**Iteration 10:** Simplified `renderChild` to use a single `nodeTypeDescriber` type assertion instead of two separate checks for `ElementType` and `AttributeType`. The new logic: assert once, check the type, return early if it doesn't match.

Result: **-1.0% ns/op** (258,200 vs 260,500).

### Why

After the big allocation win, the CPU profile still showed runtime/GC at ~80%. These micro-optimizations target the remaining ~20% of application code in the hot path.

### What worked

Both changes are simple, clean, and produce consistent (if small) improvements across multiple benchmark runs. The switch statement also reduced the source code by 16 lines.

### What didn't work

Nothing failed, but the gains are small. The CPU profile confirms we're approaching diminishing returns -- application code is only ~5-8% of total CPU time.

### What I learned

For small fixed sets of strings, Go switch statements are faster than map lookups. The compiler can generate efficient comparison code (potentially a jump table or binary search) without the overhead of hashing.

### What was tricky

Nothing particularly tricky. Both changes were straightforward refactors with clear semantics.

### What warrants review

PR #306 (void element switch): https://github.com/maragudk/gomponents/pull/306

Iteration 10 (renderChild simplification) was not submitted as a PR -- the gain is marginal and the code change, while logically equivalent, alters the control flow enough that it warrants more scrutiny relative to its benefit.

### Future work

The pre-escape attribute values optimization (iteration 11) showed a massive -18.8% improvement on the render-only path but was negligible for construct+render. It was kept on the autoresearch branch but not submitted as a PR since the typical usage pattern is construct+render.

## Step 5: PR submission and communication

### Prompt Context

**Verbatim prompt:** "Create separate PRs for exp 6 and 9 only, towards upstream/main. Make sure to document pre and post benchmark results in the PR description"
**Interpretation:** Submit the two most impactful, clean optimizations as upstream PRs with benchmark evidence.
**Inferred intent:** Get the optimizations merged into gomponents.

### What I did

Created two PRs against `maragudk/gomponents`:
- **PR #305** (`reduce-attr-allocs`): Cherry-picked experiment 6 commit, clean commit message, full benchmark before/after in description.
- **PR #306** (`void-element-switch`): Cherry-picked experiment 9 commit. Initially submitted without raw benchmark output -- fixed after review by measuring the change in isolation using a git worktree and adding 5-run before/after numbers.

Posted about both PRs on Bluesky.

### Why

Only experiments 6 and 9 were selected because they have the clearest cost/benefit ratio and are the most defensible changes. Experiments 10-11 have smaller gains relative to their code change complexity.

### What worked

Using git worktrees to measure each change in isolation gave clean, comparable numbers without benchmark contamination from other changes.

### What didn't work

First attempt at the Bluesky post exceeded the 300-character limit (540 chars). Took three attempts to fit the message.

### What I learned

When submitting performance PRs, always include raw benchmark output (not just percentages). Reviewers want to see the actual numbers and run counts to judge statistical significance.

### What was tricky

PR #306 needed to be measured in isolation (only the switch change, on top of upstream/main without experiment 6). The cumulative measurements from the autoresearch loop weren't suitable for a standalone PR description.

### What warrants review

Both PRs are open and awaiting review:
- https://github.com/maragudk/gomponents/pull/305
- https://github.com/maragudk/gomponents/pull/306

### Future work

If both PRs are merged, experiments 10 (renderChild simplification) and 11 (pre-escape attribute values) could be submitted as follow-up PRs. The pre-escape change is particularly interesting for use cases where templates are built once and rendered many times (e.g., static site generators).
