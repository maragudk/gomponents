# Diary: Zero-allocation singletons for childless elements and boolean attributes in `html`

Make `html.Br()`, `html.Div()`, and every other element helper return a shared, package-level node when called with no children, and make every boolean attribute helper (`html.Async()`, `html.Checked()`, ...) return a shared node always. These nodes have no dynamic input, so they can be built once at init and reused, taking those calls from 1 allocation to 0. Nodes with children or values keep allocating: the allocation is the closure boxed into the `Node` interface, and there is no way around that for a tree of interfaces.

## Step 1: Brainstorm and requirements

**Author:** main

### Prompt Context

**Verbatim prompt:** "See the readme for some benchmark simplified results. I wonder if we could avoid allocations for elements and attributes by having a big set/map of all elements in the core package, which is then used by all the helper functions? Let's brainstorm."
**Interpretation:** Markus wants to reduce the per-node allocations shown in the README benchmark table, and his first idea is a lookup table of element nodes in the core package.
**Inferred intent:** Fewer allocations in real pages, without changing the API or making the core package more complex.

### What I did

Read `/gomponents.go`, the README performance FAQ, and the March 2026 autoresearch diary at `/docs/diary/2026-03-25-benchmark-optimization.md`, then walked Markus through why a name table alone cannot remove the allocation: `El` allocates exactly one closure that captures both the name and the children slice header, and the `Node` interface boxes it. The diary from March already tried struct nodes and pre-rendered tag bytes and found closures to be as good as it gets.

What a table *can* do is remove the allocation for nodes with no dynamic input. We settled on package-level singletons in the `html` package: one per element helper, returned when `len(children) == 0`, and one per boolean attribute helper, returned always. Markus asked whether sharing a node is safe if someone mutates it; it is, because a `NodeFunc` is a func value with only immutable captures and `Node` exposes nothing but `Render`.

Requirements agreed:

- Scope is the `html` package only. Core, `components/`, and `Doctype` are untouched. Output is byte-for-byte identical. No new exported API.
- All 114 element helpers (every function in `/html/elements.go` except `Doctype`) and all 17 zero-arg boolean attribute helpers in `/html/attributes.go`.
- Each singleton var sits directly below its helper. Naming: `<tag>El` for elements and `<name>Attr` for attributes, which also sidesteps the Go keywords `select`, `map`, and `defer`. Where two exported helpers render the same tag (`Cite`/`CiteEl`, `Data`/`DataEl`, `Form`/`FormEl`, ...), they share one var named after the tag.
- Tests: the existing tables in `/html/elements_test.go` and `/html/attributes_test.go` gain childless cases with expected output, a `testing.AllocsPerRun` subtest asserting zero allocations for every helper in both tables, and a small test rendering a singleton twice to pin that sharing is stateless. Coverage stays at 100%.
- Benchmarks: `make benchmark` before and after with `-count=10` through `benchstat`. `BenchmarkRealisticPage` is the headline; core `BenchmarkEl`/`BenchmarkAttr` are controls and must not move. Two new benchmarks in `/html/benchmark_test.go` for the childless element and boolean attribute cases, and matching rows in the README table.

### Why

The user's idea was aimed at the right symptom but the wrong mechanism. Rather than build a table that would not deliver, we narrowed to the subset of nodes where a precomputed value is actually possible and made the win measurable and test-enforced.

### What worked

Reading the March diary first. It saved re-running experiments that had already been tried and discarded, and gave a concrete explanation for why the closure is the floor.

### What didn't work

Nothing failed in this step; it was design only.

### What I learned

Go avoids boxing an interface value only for pointer-shaped values, zero-sized values, and small integers. Any node carrying a string or a slice allocates when it becomes a `Node`, which is why singletons are the only route to zero for this library's design.

### What was tricky

Deciding whether to include all elements or only void ones. Markus chose all of them so the rule is uniform and easy to explain, accepting one extra length check per call.

### What warrants review

The requirements above are the contract for the builder; the implementation steps that follow will say what to look at in the code.

### Future work

None implied yet. If the realistic page barely moves, that is a fine result to report; the per-call numbers stand on their own.

## Step 2: Implement the singletons, measure, and self-review

**Author:** builder-singletons

### Prompt Context

**Verbatim prompt:** "Make every element helper in `html/elements.go` return a shared package-level singleton node when called with no children, and make every zero-arg boolean attribute helper in `html/attributes.go` always return a shared singleton. These calls go from 1 allocation to 0. Everything else keeps working exactly as today; rendered output must be byte-for-byte identical."
**Interpretation:** Turn the design from Step 1 into code, with the benchmark evidence and the test enforcement that the requirements ask for.
**Inferred intent:** A measurable, test-enforced allocation win in the `html` package that nobody has to think about again.

### What I did

Transformed all 111 element helpers in `/html/elements.go` that call `g.El` directly, each gaining an `if len(children) == 0 { return <tag>El }` guard and a `var <tag>El = g.El("<tag>")` two lines below its closing brace. The three deprecated shims `CiteEl`, `FormEl`, and `LabelEl` were left untouched: they already delegate with `return Cite(children...)`, so a childless call forwards a nil slice and lands on the same singleton, which is how 111 vars cover 114 helpers. All 17 zero-argument boolean attribute helpers in `/html/attributes.go` now return a `<name>Attr` var unconditionally.

Both element tables in `/html/elements_test.go` gained an `Expected` column carrying the literal childless rendering, and every row now runs three subtests: the pre-existing with-children assertion, the exact childless string, and a `testing.AllocsPerRun` check. `/html/attributes_test.go` got the same alloc subtest for the 17 boolean attributes. The shared helpers `sink` and `assertNoAllocs` live at the bottom of `/html/elements_test.go`. `TestSharedNodes` renders one shared element and one shared attribute into eight separate writers concurrently and asserts every writer got the same bytes.

`/html/benchmark_test.go` gained `BenchmarkElement` (childless `Br()`) and `BenchmarkBooleanAttribute` (`Async()`), both in the construct-and-render-to-`io.Discard` shape of the core `BenchmarkEl`. `/README.md` gained two rows in the performance FAQ table.

### Why

The length check is the cheapest possible discriminator, and putting the var directly below its helper keeps the file readable as a flat list of elements rather than splitting it into a wall of declarations at the top. Leaving the deprecated shims alone avoids duplicating the guard in three places for no behavioural difference.

### The measurements

Baseline came from a pristine export of `HEAD` into the scratchpad, with only the two new benchmark functions copied in so the new benchmarks have a before as well. Both runs were `go test -bench . -benchmem -count=10 ./...`, compared with `benchstat`. The `html` package:

```
                                                 │  before.txt  │              after.txt              │
                                                 │    sec/op    │   sec/op     vs base                │
Element/element_without_children-10                12.965n ± 1%   5.884n ± 1%  -54.62% (p=0.000 n=10)
BooleanAttribute/boolean_attribute-10               9.669n ± 0%   4.768n ± 0%  -50.69% (p=0.000 n=10)
RealisticPage/construct_and_render/discarded-10     67.41µ ± 0%   67.14µ ± 1%   -0.39% (p=0.023 n=10)
RealisticPage/construct_and_render/buffered-10      90.72µ ± 1%   89.81µ ± 0%   -1.00% (p=0.000 n=10)
RealisticPage/render_pre-built_tree/discarded-10    43.28µ ± 1%   43.14µ ± 0%   -0.33% (p=0.007 n=10)
RealisticPage/render_pre-built_tree/buffered-10     66.72µ ± 0%   66.33µ ± 0%   -0.59% (p=0.004 n=10)

                                                 │ before.txt  │                after.txt                 │
                                                 │  allocs/op  │  allocs/op   vs base                     │
Element/element_without_children-10                 1.000 ± 0%    0.000 ± 0%  -100.00% (p=0.000 n=10)
BooleanAttribute/boolean_attribute-10               1.000 ± 0%    0.000 ± 0%  -100.00% (p=0.000 n=10)
RealisticPage/construct_and_render/discarded-10    3.044k ± 0%   3.042k ± 0%    -0.07% (p=0.000 n=10)
RealisticPage/construct_and_render/buffered-10     3.044k ± 0%   3.042k ± 0%    -0.07% (p=0.000 n=10)
RealisticPage/render_pre-built_tree/discarded-10    18.00 ± 0%    18.00 ± 0%         ~ (p=1.000 n=10) ¹
RealisticPage/render_pre-built_tree/buffered-10     18.00 ± 0%    18.00 ± 0%         ~ (p=1.000 n=10) ¹
```

Bytes per op went 48 to 0 and 24 to 0 for the two new benchmarks. The realistic page contains exactly two childless nodes, so it drops from 3044 to 3042 allocations: real, reproducible, and not worth anything. That is the honest headline.

### What didn't work

Two things, and the second cost real time.

The first baseline run was started against the live worktree before I had edited anything, on the assumption that `go test ./...` builds every test binary up front. Fifteen minutes in, the output showed it was still inside `BenchmarkAttr` in the root package, meaning the `html` package had not run yet and would be measured against test files I had already changed. Killed it with `pkill -f "go test -bench"` and started over from a pristine export of `HEAD` in the scratchpad, which guarantees the baseline sources cannot drift under the run.

The second was the controls. Comparing the two full runs, the core benchmarks moved even though `gomponents.go` and `gomponents_benchmark_test.go` are byte-identical in both trees (`diff` confirms it):

```
El/normal_elements-10                                17.27n ± 1%   16.93n ± 1%   -1.94% (p=0.005 n=10)
Attr/construct_and_render/discarded/boolean-10       11.06n ± 0%   10.66n ± 1%   -3.57% (p=0.000 n=10)
Text/construct_and_render/discarded/short_text...    14.72n ± 1%   13.92n ± 2%   -5.37% (p=0.000 n=10)
```

The requirements say to stop and find out why when the controls move, so I did. Re-running just the core benchmarks in both trees back to back, with nothing else on the machine:

```
                      │ ctrl-baseline.txt │         ctrl-worktree.txt          │
El/normal_elements-10         15.84n ± 1%   15.99n ± 2%  +0.95% (p=0.014 n=10)
Raw/raw_element-10            2.784n ± 0%   2.767n ± 1%  -0.63% (p=0.002 n=10)
geomean                       6.642n        6.653n       +0.16%
```

The two trees agree within a percent, and `El` reads 15.8-16.0 ns here against 16.9-17.3 ns in the long runs, so the machine drifts several percent between sessions and `p=0.000` inside one benchstat pair is not evidence of a real change. Allocations and bytes per op for the controls are identical in every comparison, which is the number that actually matters. Nothing in this change can reach the core package.

### What I learned

`benchstat`'s p-values only speak to variance within the samples it was given. Two runs separated by twenty minutes on a laptop share a systematic offset that no amount of `-count` will surface, so a control that "moves significantly" needs a back-to-back re-measurement before it means anything.

Also: `testing.AllocsPerRun` asserting zero is a test that fails open. If the compiler eliminates the call being measured, the count is zero and the test passes for the wrong reason. Storing the result in a package-level `sink` and calling through a func value in a table makes the call indirect and unremovable. Both reviewers verified this by hand, reverting `Br` and `Async` to their old bodies and confirming the tests fail with "expected 0 allocations but got 1".

### What was tricky

Deciding what to do about the `Realistic full page` row in the README. It is now literally stale by two allocations out of 3044, but the instruction was to touch it only if the numbers moved meaningfully, and they did not. I left it, and flagged it rather than quietly deciding for Markus. The row's nanoseconds already came from an older run than either of mine, so refreshing only the allocation count would make it less internally consistent, not more.

The `g.Group(children)` helpers (`Abbr`, `B`, `H1`, and the rest) were allocating twice when childless, since `Group(nil)` gets boxed into a `Node` on its way into `El`. The red run caught that: "expected 0 allocations but got 2". They now cost nothing like the others.

### What I found during self-review

Two competing reviewers went over the diff. Both independently confirmed the tag-to-var mapping mechanically across all 114 helpers, both mutation-tested the alloc assertions, and both confirmed sharing is safe: a node is a closure over an immutable string and a nil slice, `Node` exposes only `Render`, and node identity was never observable because comparing two `Node` values holding func types panics. Neither found a correctness defect.

They agreed on one weakness in my tests. `TestSharedNodes` as first written rendered a node into two writers and compared them, which any node passes, singleton or not, so it pinned nothing. I rewrote it to render the shared node from eight goroutines into eight writers, which does say something a fresh node would not trivially satisfy, and which the race detector can judge. `go test -race ./html/` passes.

They also both caught a factual error in Step 1 above: it lists `Data`/`DataEl` as a pair of exported helpers rendering the same tag. There is no `Data` element helper; `Data(name, v string)` is the data-attribute function in `/html/attributes.go`, and `DataEl` has its own `dataEl` var. The real pairs are exactly `Cite`/`CiteEl`, `Form`/`FormEl`, and `Label`/`LabelEl`, which the code handles correctly. Step 1 is not mine to edit, so the correction lives here.

The remaining single-reviewer notes I did not act on: the `Expected` column duplicating what `Name` could derive (the requirements ask for literal expected strings on purpose), the README row labels using call syntax (the format Markus specified), and refreshing the realistic page row (covered above).

### What warrants review

Three judgement calls rather than three risks.

First, the deprecated shims `CiteEl`, `FormEl`, and `LabelEl` reach their singleton by delegation rather than carrying their own length check. Both reviewers agreed this satisfies the requirement and is simpler; it is worth a nod that this reading of "share one var" is the intended one.

Second, the `Realistic full page` README row still says 3044 allocations where the branch produces 3042.

Third, `sink` in `/html/elements_test.go` is a package-level variable written from two test files. It exists solely to keep `testing.AllocsPerRun` honest, and nothing reads it, which looks odd until you know why.

To validate: `go test -shuffle on ./...`, `go test -coverprofile=cover.out -shuffle on ./...` (100.0% for `gomponents`, `components`, `html`, `http`, `x/slices`), `go test -race ./html/`, `golangci-lint run`, and `go test -bench 'Element|BooleanAttribute' -benchmem ./html/` to see the zeroes.

### Future work

None implied. If anyone ever wants the realistic page number to move, the remaining 3042 allocations are nodes with children or values, and Step 1 already explains why those cannot go to zero without changing what a `Node` is.
