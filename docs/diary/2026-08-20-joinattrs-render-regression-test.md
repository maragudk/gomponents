# Diary: Regression test that `JoinAttrs` doesn't render element children

`components.JoinAttrs` inspects its children to find attributes with a given name. It must
inspect attribute-type nodes without rendering element-type ones, since rendering an element
during inspection would be both wasteful and observable through side effects. Nothing in the
test suite pins that property down. This task adds a regression test for it.

The test is salvaged from PR #293, which was closed as superseded by PR #335.

## Step 1: Compare PRs 293 and 335, pick a winner, salvage the leftover test

**Author:** main

### Prompt Context

**Verbatim prompt:** "Have a look at PRs 293 (ours) and 335 (contributor). Which is the better fix? DO NOT write anything as a github comment unless I say otherwise." Followed by "I've closed 293, add a superseded comment to it" and "Yeah, please do. Pull main first, 335 merged".

**Interpretation:** Evaluate two competing fixes for the same bug, recommend one, then record the
outcome on the losing PR and open a follow-up for the one piece of it worth keeping.

**Inferred intent:** Markus wanted a real technical comparison rather than a preference, and wanted
nothing said publicly until he had decided. Once he merged #335 and closed #293, the remaining goal
was to not silently lose the good idea buried in the closed PR.

### What I did

Read both PRs and the current `/components/components.go`. Both fixed the same bug: `JoinAttrs`
unwrapped `g.Group` only one level deep, so an attribute nested deeper fell through to the result
verbatim and `renderChild` emitted it as a duplicate attribute. Browsers keep the first occurrence,
so the value was silently dropped.

Both PRs made the same core change — turn `processNode` into a recursive closure that unwraps groups
at any depth, matching `renderChild` in `/gomponents.go`, which already recurses into groups without
limit.

Verified the two against merged `main` by fetching each PR ref into a scratch worktree. PR #335
merged cleanly, passed `go test -shuffle on ./components/...` at 100% coverage, and passed
`golangci-lint run ./components/...` with zero issues. PR #293 reported `CONFLICTING` and produced
conflicts in both files:

```
CONFLICT (content): Merge conflict in components/components.go
CONFLICT (content): Merge conflict in components/components_test.go
```

Recommended #335. Markus merged it as `ce1fbc2` and closed #293, and I posted a comment on #293
explaining the supersession.

### Why

PR #293 dated from February and predated #303, which added boolean, empty, and whitespace-only
attribute handling and extracted `processNode` as a closure along the way. Most of what remained in
#293 was that same refactor arriving late. Rebasing it would have produced essentially #335's diff.
#335 also documented the semantics more precisely — a `g.Group` is transparent and doesn't count as
a level, while a child element still does — and tested the composition pattern that actually
triggers the bug, where one component forwards `g.Group(children)` into `JoinAttrs` and another does
the same.

### What worked

Fetching both PR refs into throwaway worktrees and actually merging them against `main` settled the
question faster than reading diffs did. The conflict in #293 and the clean run in #335 were the
decisive facts, and neither was visible from the diffs alone.

### What didn't work

Nothing failed outright. The one wasted motion was reading #293's diff carefully before checking its
merge state — the diff was against a months-old base, so it described code that no longer exists and
was misleading until I noticed the base commit.

### What I learned

`extractAttrValue` in `/components/components.go` skips a node without rendering it in two distinct
cases: when the node's `Type()` reports `g.ElementType`, and when the node doesn't implement the
`nodeTypeDescriber` interface at all. Both branches avoid a `Render` call, and a regression test
worth having should cover both rather than just the first.

### What was tricky

Judging a PR whose diff no longer applies takes care: the diff shows the author's intent against a
base that has moved, so "this PR re-does an existing refactor" is only visible by comparing against
current `main` rather than by reading the patch.

### What warrants review

Nothing to review yet in this step — no code changed. The comparison conclusion is recorded on
[PR #293](https://github.com/maragudk/gomponents/pull/293) for anyone who wants to check the
reasoning.

### Future work

Add the regression test itself, which is the remaining piece of #293 worth keeping. It should assert
that `JoinAttrs` never calls `Render` on a node it skips, covering both the explicit `ElementType`
case and the no-`Type()`-method case.

## Step 2: Add the regression test

**Author:** main (requirements) and the builder sub-agent (implementation)

### Prompt Context

**Verbatim prompt:** "Yeah, please do. Pull main first, 335 merged"

**Interpretation:** Open the follow-up PR for the salvaged test, based on `main` with #335 already
merged rather than on the pre-merge state.

**Inferred intent:** Land the leftover value from #293 as a small, self-contained PR now that the
bug it sat next to is fixed, so the property is guarded before anyone touches `JoinAttrs` again.

### What I did

Created a worktree branched off `ce1fbc2` and wrote requirements, then delegated the implementation.
The result is in `/components/components_test.go`: two recording node types and four subtests, 84
lines added, no production code touched.

The two node types cover the two distinct skip branches in `extractAttrValue`. `wiretap` implements
`Type()` returning `g.ElementType`; `mole` has no `Type()` method at all. Each sets a `rendered`
flag and writes a tell string when rendered. The four subtests place a skip node at top level, one
group deep, and two groups deep, each sandwiched between `g.Text("before")` and `g.Text("after")`
siblings.

Verified independently of the builder's own report: `make test` passes at 100% coverage across every
package, `golangci-lint run` reports 0 issues, and `git diff ce1fbc2 -- components/components.go` is
empty. Opened as [PR #342](https://github.com/maragudk/gomponents/pull/342).

### Why

`JoinAttrs` decides whether a child matches by rendering it and string-matching the output, so the
guard that skips non-attribute nodes before rendering is load-bearing. Without a test, a future
refactor of `extractAttrValue` could start rendering every child and nothing would go red — the
joined output would still be correct, and only the wasted work and the side effects would differ.

### What worked

Requiring a recording node rather than #293's panicking one. The recorded flag makes the assertion
visible to a reader, and it allowed a second assertion the panic version couldn't express: after
checking `rendered` is false, the test renders the result and confirms `rendered` flips to true.
That guards the guard — it proves the recorder actually works and the test isn't passing because the
node was never involved at all.

Requiring the deliberate-failure check paid for itself as a specification: it forced the test to be
falsifiable rather than merely green.

### What didn't work

Nothing failed. The deliberate-failure check behaved exactly as intended — deleting the three-line
guard at `/components/components.go:141-144` failed precisely the four new subtests and no others:

```
--- FAIL: TestJoinAttrs/does_not_render_a_node_it_skips_because_its_Type_is_ElementType
    components_test.go:168: wiretap was rendered while JoinAttrs was still deciding what to do with it
```

Reverting from a `sed -i.bak` backup left `/components/components.go` byte-identical to `ce1fbc2`.

### What I learned

A "does not happen" test needs a paired positive assertion or it can pass for the wrong reason. Here
the pairing is cheap: assert the flag is false right after `JoinAttrs` returns, then render and
assert it is true. Without the second half, deleting the node from the call entirely would still
leave the test green.

### What was tricky

Deciding what `mole` should be. The no-`Type()`-method branch is easy to overlook because it looks
like defensive coding rather than a real case, but `renderChild` treats such nodes as element
content, so they genuinely reach `extractAttrValue` and genuinely must not be rendered there.

### What warrants review

The four subtests are near-identical and could be table-driven, which `/CLAUDE.md` prefers where
appropriate. They were left expanded because the two node types differ in interface shape — one has
`Type()`, one doesn't — so a table would need an interface plus a reset closure and would likely
read worse. Worth a second opinion.

Reviewers should also confirm the PR touches only `/components/components_test.go`.

### Future work

None. The property is guarded and the salvage from #293 is complete.
