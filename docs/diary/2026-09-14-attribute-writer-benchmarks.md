# Diary: Benchmark rendering attributes to different writers

PR #354. Grew out of reviewing #352, a one-line change to how attribute values are escaped and written, whose only evidence was one `io.Discard` sub-benchmark. The review found the change looked different on other writers, but the numbers came from throwaway benchmarks that nobody else could run. Markus asked for the benchmark to become part of the suite so #352 could be judged on something reproducible. Over one day the benchmark went from a separate function with four writers and four hand-written values to `BenchmarkAttr` as a single table over construction mode, writer, and generated values.

All measurements below are on an Apple M4 with Go 1.27.1, interleaved base-versus-PR runs of prebuilt test binaries, compared with benchstat.

## Step 1: Open the benchmark PR

**Author:** main

### Prompt Context

**Verbatim prompt:** "Yes, open the benchmark PR to main"
**Interpretation:** Turn the writer-variant benchmark I had used to review #352 into a proper addition to `/gomponents_benchmark_test.go` and open a PR against main.
**Inferred intent:** Give #352's author and Markus something runnable that shows the effects the `io.Discard`-only benchmark hides, so the discussion on #352 can move from "show me" to numbers.

### What I did

Created a worktree on a new `attr-writer-benchmarks` branch from `origin/main`, since the local `main` was six commits stale. Added `BenchmarkAttrRender` to `/gomponents_benchmark_test.go`: a pre-built attribute rendered to `io.Discard`, a 2048-byte `bufio.Writer`, and both behind a `writeOnly` wrapper type that only has a `Write` method. Four values: nothing to escape, the short escaping value from the existing benchmark, a JSON-ish value, and prose. Pushed, ran the benchmark on main and on #352's head, and put the resulting table in the PR body.

The first push went to the `maragubot/gomponents` fork out of habit. Markus then gave the account write access to `maragudk/gomponents` and asked for the fork remote to be removed, so the branch moved to `origin`, the fork branch was deleted, and the remote was removed from the checkout.

### Why

`io.Discard` implements `io.StringWriter` and its writes are free, so it cannot show the two things the #352 change affects: the cost of many small writes into a buffer, and the path taken by writers that only implement `Write`, which is what middleware wrapping `http.ResponseWriter` looks like.

### What worked

Compiling both test binaries once and alternating them across rounds gave stable numbers with tight confidence intervals. Keeping raw benchstat input in the scratchpad meant tables could be regenerated when names changed.

### What didn't work

The first attempt at a helper test used `v[:min(len(v), 40)]`. It failed to compile because `/go.mod` declares `go 1.18` and the test file inherits that language version:

```
./gomponents.go:240:17: cannot range over len(s) (value of type int): requires go1.22 or later (-lang was set to go1.18; check go.mod)
```

That one was from an earlier experiment with `for i := range len(s)`; the `min` builtin failed the same way. Anything written in this repo's test files has to be Go 1.18 syntax.

### What I learned

The cross-binary comparison showed the write-only rows moving by up to fifteen percent for a code path that is logically identical on both sides. A single-binary A/B of the two paths showed them equal to the nanosecond. Adding a function shifts everything after it in the binary, which changes instruction-cache line boundaries, and on Apple Silicon that swings tight loops by five to fifteen percent, consistently for a given binary, so benchstat reports it with a confident p-value. The tell is identical allocations plus the gap vanishing when both versions share one binary.

### What was tricky

Deciding what to compare against. #352's merge base was older than `origin/main`, so the PR body compares main to #352's head after its merge of main, and the v1 comparison cherry-picks the single original commit onto current main.

### What warrants review

Whether `writeOnly` wrapping `io.Discard` is a fair stand-in for a middleware-wrapped response writer. It reproduces the one property that matters, no `WriteString` method, and nothing else.

### Future work

The same value axis for `BenchmarkText`, which only has short, short-escaping, and prose values. The writer axis matters less there because `Text` escapes at construction.

## Step 2: First review round

**Author:** main

### Prompt Context

**Verbatim prompt:** "/fabrik:address-code-review" (three invocations over the round), then "Leave tailwind out of it, just use a long string that doesn't need esaping"
**Interpretation:** Walk through Markus's inline comments one at a time, propose a fix for each, apply after agreement.
**Inferred intent:** Make the benchmark cases say what they measure rather than what they look like.

### What I did

Renamed the JSON and prose values to "many characters needing escaping" and "long value with few characters needing escaping", with a synthetic value for the former. Added a long value needing no escaping, as plain text rather than the Tailwind class list I first proposed. Folded the separate function into `BenchmarkAttr` as `render pre-built/<writer>/<value>` sub-benchmarks next to the three existing ones. Each change was one commit, and the PR body table was re-measured and rewritten after each.

### Why

Markus's comments: "Why is this one special? Because it's longer than the line just above? Maybe don't use JSON as an example, name it after which characteristic we're benchmarking." The characteristic is the count of escaped bytes and the length, not the format. The long clean value fills the gap where the scan itself is the dominant cost, which is exactly the path #352's first version regressed and its `io.Discard` escaping benchmark could not see.

### What worked

Presenting one comment with a concrete proposal and waiting. Every round was a one-word answer.

### What didn't work

Nothing failed in this step.

### What I learned

The long clean value exposed the scan regression in #352 v1 directly: +11% to +16% on every writer, where the short value showed only +3%.

### What was tricky

Keeping the old three sub-benchmark names alive alongside the table was tempting for continuity with numbers quoted on #352, and wrong. Half-table, half-hand-written reads worse than either.

### What warrants review

That the generated dense value, `strings.Repeat("\"hat\" & ", 6)` at the time, was a fair "many characters" case. It was later replaced anyway.

### Future work

None beyond what later steps did.

## Step 3: One table over construction, writer, and value

**Author:** main

### Prompt Context

**Verbatim prompt:** "Before we do that, should we add prebuilt/not and bool attr/name-value attr to the table and just have BenchmarkAttr be one large table-driven test?" then "Let's do it." then "/fabrik:address-code-review" for the resulting comment.
**Interpretation:** Make construction mode and the boolean attribute axes of the same table, replacing the hand-written sub-benchmarks entirely.
**Inferred intent:** One place to read every attribute number, so an optimization that trades constructor cost for render cost shows both sides at once.

### What I did

Rewrote `BenchmarkAttr` as nested loops over writers and values, with a `construct and render` and a `render pre-built` sub-benchmark per pair, and a boolean row that calls `Attr("hat")` with no value. 48 sub-benchmarks. Markus's inline comment "No package-level state please, keep tests independent" removed the package-level sink variable I had added; a local declared in the sub-benchmark outside the loop does the same job.

### Why

The construct axis keeps the closure allocation in view. CI runs `go test -bench . -benchmem ./...` once, so the table adds about 25 seconds to the benchmark job, which Markus accepted.

### What worked

Measuring the sink question before proposing. In one binary, `attr(v).Render(w)` inline reports 0 allocations and 9.7 ns, a local declared inside the loop or outside it reports 48 B and 1 allocation, and a package-level variable the same. The local form is what the code now uses.

### What didn't work

The first table version constructed the node inline and reported 0 B/op for `construct and render` rows that had reported 48 B/op in the old hand-written form. Escape analysis kept the closure on the stack because nothing retained it. That would have hidden the constructor allocation the axis exists to show. Caught by comparing against the old numbers before pushing.

### What I learned

`b.Loop` keeps the loop body's results alive but does not promise heap allocation. Whether the constructed `Node` escapes depends on whether the interface `Render` call is devirtualized after inlining, and a plain local assignment is enough to stop that.

### What was tricky

The three original sub-benchmark names had to go, which means numbers quoted on #352 against `BenchmarkAttr/name-value_attributes_needing_escaping` no longer have a counterpart. That is a one-time cost.

### What warrants review

The comment above the local sink explains it in terms of pages keeping their nodes, which is the modelling reason, not the compiler mechanism. A second opinion in Step 5 called that the wrong explanation. The claim it rests on, that construction cost disappears without the sink, I verified as true.

### Future work

None.

## Step 4: A grid of generated values

**Author:** main

### Prompt Context

**Verbatim prompt:** Inline comment "Streamline these?" with a six-item list of short and long values at three escaping levels, then "Maybe even make the strings the same, and the same length, but change characters around according to escaping level", then "Instead of stride, maybe a fraction? Makes it more readable. Or (n, escapes int)?", then "Sounds good. But keep n powers of two? Or is 192 significant?"
**Interpretation:** Replace the hand-written values with a short-by-long times none-by-little-by-much grid, generated from one base string so only the escaped byte count varies within a length.
**Inferred intent:** Make row-to-row deltas attributable to exactly one thing.

### What I did

Added a `value(n, escapes int) string` helper that repeats "party hat " to n bytes and swaps `escapes` evenly spread bytes for `"`, `&`, and `'` in turn. The grid is 16 and 256 bytes with 0, 1 or 16, and 4 or 64 escapes. 56 sub-benchmarks. 192 had been the length of the prose value and meant nothing, so the long value became 256, which also sits above the 48-byte threshold where `strings.Replacer` switches counting strategy while 16 sits below it.

### Why

With hand-written values, "many characters" was 48 bytes and "long, few characters" 184, so length and density changed together. Generated values hold length fixed within a row pair.

### What worked

Printing the generated values once through a throwaway test to confirm lengths, escape counts, and first-escape positions before benchmarking them.

### What didn't work

The throwaway test used the `min` builtin and did not compile under the `go 1.18` language version. Rewrote it without.

### What I learned

The generated grid reproduced the one real trade-off in #352 inside the PR's own benchmark: `buffered / long value, much escaping` is +22% for #352's current head, because streaming through the replacer issues 129 writes into the buffer where main issued one. At one byte in sixteen the same length wins by 29%. Markus's reading: an attribute value with a lot of escaping is not the common case, so the trade-off is right.

### What was tricky

A background benchmark run was in flight when the value set changed. I stopped it and its leftover processes before editing, since its results would have been for values that no longer existed.

### What warrants review

`value(256, 64)` puts its first escape at byte 3 and `value(256, 16)` at byte 15. Neither is a late first escape, so the case where the byte-table scan runs long before handing off to `Replace`, which rescans, is not measured. A `value(256, 1)` row would cover it.

### Future work

That row, if wanted.

## Step 5: Second opinion and per-sub-benchmark writers

**Author:** main

### Prompt Context

**Verbatim prompt:** "/fabrik:second-opinion" with the topic chosen as the PR 354 benchmark design, then "Fix finding 1 from codex, then re-benchmark and update tables in PR body", then "Make sure the tables use names to match"
**Interpretation:** Get an independent review of the table design from codex in blind mode, verify what it claims, apply the finding about shared writers, and make every table on the PR use the sub-benchmark names verbatim.
**Inferred intent:** Catch design mistakes a single reviewer would miss before the benchmark becomes the thing #352 is judged on.

### What I did

Ran codex, gpt-6-astra at high effort, read-only, without revealing any preferred conclusion. It reviewed the 48-case version before the grid landed. It would ship the full matrix, keep all four writers and both modes, and reject the old hand-written pair. Its findings, and what I checked:

Shared `bufio.Writer` values served every sub-benchmark on that row, so the fill level one left behind decided where the next one's writes crossed the 2048-byte boundary. Verified by reading the code. Fixed by making the writer table hold `New func() io.Writer` constructors, called once before each `b.Loop`.

It said the sink comment was wrong because construction still allocates without it. That is false for the inline form, as measured in Step 3; its own experiment used an in-loop local, which does allocate.

It believed the branch also swapped `escapeString` to a `strings.Replacer` and added `Text` cases. It had diffed against the stale local `main`. Against `origin/main` the branch touches one file. The caveat it built on that is void.

It measured 57.5 seconds for the 48 cases at default benchtime and called that substantial but defensible.

Then re-ran main against #352's head with per-sub-benchmark writers, rewrote the PR body table with sub-benchmark names verbatim, and re-ran the v1 comparison with the final benchmark so the v1 comment on the PR uses the same names.

### Why

Fresh writers per sub-benchmark remove a dependency on sub-benchmark order and on which sub-benchmarks a `-bench` filter selects. Matching names let a reader go from a table row to a `-bench` argument without translation.

### What worked

Verifying each of codex's factual claims before relaying. Two of five were wrong, and both would have been believed if passed on as written.

### What didn't work

Nothing failed in this step. Per-writer construction changed no number by more than noise, which is what you want from a fix that removes an artefact rather than a cost.

### What I learned

A blind second opinion is only as good as the repository state it sees. It should be told which ref is the base, not left to find `main`.

### What was tricky

Markus read `w.New()` as a method on `writeOnly`. It is a function-typed field on the anonymous table struct, and the row shadows itself as `w` inside the loop. Worth a glance when reading the code.

### What warrants review

The v1 comment on the PR now shows write-only rows at +118% and +162% with allocations 4 to 67, which is the fallback #352's current head removed. Those rows are real, unlike the write-only deltas for the current head, which are binary layout.

### Future work

The `BenchmarkText` value axis from Step 1. The `value(256, 1)` late-escape row from Step 4.
