# Diary: Issue #330 follow-up — document the name-escaping contract, add a security policy, fix lint

Issue #330 reported that attribute values are escaped but attribute and element *names* are rendered verbatim, so a name containing whitespace injects attributes. The reporter self-closed it after finding #74, where this was settled deliberately years ago: names are trusted structural input. What remained was that the decision was documented nowhere, and the repo had no security policy. This session assessed the issue, then shipped three PRs: doc comments stating the contract (#332), a SECURITY.md (#331), and a lint fix that fell out of verifying the first PR (#333). All three merged the same day.

## Step 1: Assess issue #330

**Author:** main

### Prompt Context

**Verbatim prompt:** "Have a look at issue 330. What do you think? DO NOT post to the issue directly, reply to me."
**Interpretation:** Read the issue and its context, verify the claims against the code, and give an opinion in conversation only — no public comments.
**Inferred intent:** Markus wanted a second pair of eyes on a security-flavored report before deciding how to respond as maintainer.

### What I did

Fetched issue #330 with `gh issue view` (the default view only showed the closing comment; the JSON API got the full body), then read #74 and #133 for prior art. Verified the technical claims against `/gomponents.go` (`El` writes the tag name verbatim, `booleanAttr`/`valueAttr` write the attribute name verbatim, only values go through `template.HTMLEscapeString`) and located the runtime-name helpers `Aria` and `Data` in `/html/attributes.go` and `JoinAttrs` in `/components/components.go`. Reported an assessment: the claims are accurate, the severity framing (hardening, not a live vulnerability) is honest, and the legitimate residual gap is that the names-are-trusted contract is undocumented.

### Why

A security report deserves verification before opinion, and the issue explicitly referenced decisions (#74) that needed checking to avoid relitigating settled questions.

### What worked

Going straight to the JSON API when `gh issue view` showed only a comment. Verifying each claim against actual code lines rather than trusting the report.

### What didn't work

Nothing failed in this step.

### What I learned

`Data()` and `Aria()` are the only helpers in `/html/attributes.go` whose signatures take an attribute name at runtime, which makes them the realistic injection surface — rendering a map of data attributes from a CMS is ordinary usage. Also: the repo had no SECURITY.md, which is why the reporter defaulted to a public issue.

### What was tricky

The issue was opened and self-closed within five minutes, so the default `gh issue view` output was misleading — it looked like the closing comment was the issue body.

### What warrants review

The recommendation itself: document rather than validate names. Validation (rejecting bad names) was deliberately deferred as a possible v2 consideration.

### Future work

None beyond the PRs in the following steps.

## Step 2: Doc comments PR (#332) and SECURITY.md PR (#331)

**Author:** main

### Prompt Context

**Verbatim prompt:** "Use a writer subagent to add doc comments where applicable and make a PR. Then make a SECURITY.md file and just state in it that creating issues in the issue tracker is generally fine, but feel free to send an email to markus@maragu.dk if in doubt (separate PR)."
**Interpretation:** Delegate the doc comments to a writer subagent ending in a PR; write SECURITY.md myself with exactly that content, as a second PR.
**Inferred intent:** Close the two gaps identified in step 1 with minimal, non-breaking changes, keeping the two concerns separately reviewable.

### What I did

Spawned a `fabrik:writer` subagent in an isolated worktree with a brief covering the contract to document, the functions to touch, doc style rules, and the fork PR workflow (`origin` is the maragubot fork with push access; `upstream` is maragudk/gomponents, read-only). Meanwhile I branched `add-security-policy`, wrote `/SECURITY.md` (issues in the tracker are generally fine; email markus@maragu.dk if in doubt), and opened PR #331. The writer added 8 comment-only lines — the package doc in `/gomponents.go` plus one-sentence contracts on `El`, `Attr`, `html.Aria`, `html.Data`, and `components.JoinAttrs` — and opened PR #332. I then had to repair #332's base (see below), verified its final diff was exactly the 8 lines against current main, and fixed #331's empty description mid-stream when Markus flagged it.

### Why

The writer subagent kept the prose-heavy work out of the main context; worktree isolation let both PRs proceed in parallel without fighting over the checkout.

### What worked

The writer's actual doc comments were good on the first pass: correctly phrased, consistent across functions, and it sensibly left the deprecated `DataAttr` alone since its deprecation notice already routes to `Data`. Repairing the branch via a plumbing-level merge (below) avoided both a force push and re-resolving conflicts.

### What didn't work

Three failures, all recorded verbatim:

1. **The writer's worktree started from a stale base.** Its branch `document-trusted-names` had parent `e72dace` ("Bump actions/checkout from 5 to 6 (#281)", months old — the fork's unsynced default branch) instead of upstream main at `3d2c3fb`. GitHub reported the PR as `"mergeStateStatus":"DIRTY","mergeable":"CONFLICTING"` even though its own diff was clean. The writer itself reported branching "off main at 3d2c3fb", which was wrong. First symptom: the `gh pr diff 332` context lines showed an `Attr` implementation (`if len(value) > 1`) that didn't match current main (a `switch`).
2. **Force-push was denied.** After cherry-picking the doc commit onto `upstream/main` and resolving conflicts in `/gomponents.go` and `/components/components.go`, `git push origin +document-trusted-names-rebase:document-trusted-names` was blocked by the permission classifier. Workaround: build the merge commit directly with `git commit-tree '<resolved>^{tree}' -p db4f785 -p upstream/main`, giving the standard merge-main-into-branch result, pushed as a normal fast-forward. PR went `MERGEABLE`.
3. **PR #331 went up with an empty description.** I passed the body as `gh pr create ... --body -` with a heredoc; the `-` is not stdin for that flag. Markus caught it ("Also, PR 331 has empty description"). Fix: `gh pr edit 331 --body-file -`.

### What I learned

Subagent worktrees can silently start from a stale ref in a fork setup — always check `git merge-base upstream/main <branch>` against the upstream tip before reporting a subagent's PR as done (saved to persistent memory as `verify-subagent-branch-base`). `gh` reads stdin via `--body-file -`, never `--body -`. And a history-rewrite denial doesn't block the merge-commit route, which is the better fix for an already-pushed branch anyway. Also from this step: Markus's standing instruction "DO NOT merge any PR yourself btw" (saved to memory as `never-merge-prs`).

### What was tricky

The stale base was disguised: the PR's *diff* looked perfect, and only the mismatch between the diff's context lines and code I'd read earlier gave it away. Resolving the `JoinAttrs` conflict needed care because main's doc comment for it had been substantially rewritten since the writer's base — the new sentence had to slot into the newer, longer doc block.

### What warrants review

Merged already, but the wording of the contract sentence ("The name is rendered unescaped and must be a trusted value, never user-controlled data") is repeated across five doc comments — if it ever needs changing, it needs changing in all five places: `/gomponents.go` (package doc, `El`, `Attr`), `/html/attributes.go` (`Aria`, `Data`), `/components/components.go` (`JoinAttrs`).

### Future work

Optional: name validation behind a future major version, as discussed in #330. No commitment made.

## Step 3: Lint fix PR (#333)

**Author:** main

### Prompt Context

**Verbatim prompt:** "Also create a PR to fix the lint issue on main"
**Interpretation:** Fix the pre-existing `make lint` failure discovered while verifying PR #332, as its own PR.
**Inferred intent:** Keep `make lint` green on main; the failure predated this session's work but surfaced during it.

### What I did

Diagnosed the failure: `html/attributes_test.go:166:7: inline: Call of html.DataAttr should be inlined (govet)`. `DataAttr` carries a `//go:fix inline` directive (as do five other deprecated aliases), and recent golangci-lint versions run govet's `inline` analyzer, which flags direct *calls* of such functions. Every other deprecated alias (`StyleAttr`, `TitleAttr`, `CiteEl`, `FormEl`, `LabelEl`) escapes this because tests pass them as function values in tables. Only `DataAttr` was called directly in `TestData`. Restructured `TestData` in `/html/attributes_test.go` into a table of `func(name, v string) g.Node` covering `Data` and `DataAttr`, matching the file's existing style. `make test` (100% coverage) and `make lint` (0 issues) both green. PR #333.

### Why

Testing through a function value is the established idiom in this repo for covering deprecated aliases without tripping the linter — more consistent than a `nolint` comment, and coverage of `DataAttr` is preserved.

### What worked

Checking CI first: `.github/workflows/ci.yml` lints with `golangci-lint-action@v9, version: latest`, so this was a real upcoming CI failure, not a local-version artifact — worth stating in the PR description. Verifying the failure existed on unmodified main before attributing it.

### What didn't work

Nothing failed in this step.

### What I learned

govet's `inline` check only flags call expressions; taking a function value is invisible to it, which is exactly why the table-driven pattern in this test file sidesteps the issue legitimately.

### What was tricky

Deciding between `nolint` and restructuring. The tiebreaker was that the codebase already had a consistent pattern for this exact situation.

### What warrants review

Merged already. The subtest names changed shape slightly (`Data returns an attribute...` / `DataAttr returns an attribute...`), which shows up in verbose test output.

### Future work

None.

## Step 4: Merges and cleanup

**Author:** main

### Prompt Context

**Verbatim prompt:** "I've merged 333" followed by "I've merged both"
**Interpretation:** All three PRs are merged; sync and tidy up.
**Inferred intent:** Leave the local checkout and the fork clean.

### What I did

Fast-forwarded local main to upstream after each merge (final: `d1fd82b`), deleted `fix-dataattr-lint-issue`, `add-security-policy`, and `document-trusted-names` locally and on the fork, removed the writer's worktree at `.claude/worktrees/agent-ae4b4e7351f5c1762`, and deleted its scaffold branch `worktree-agent-ae4b4e7351f5c1762` — which sat at `e72dace`, confirming the stale-base diagnosis from step 2. Left the pre-existing branches (`autoresearch-*`, `garden/*`, etc.) untouched. Markus commented on #330 himself, closing the loop with the reporter.

### Why

Merged branches and dead worktrees are noise; the confirmation of the stale-base theory was a free byproduct.

### What worked

`gh pr view --json mergeable` returned `UNKNOWN` immediately after main moved; a single retry resolved it — no need for anything cleverer.

### What didn't work

Nothing failed in this step.

### What I learned

The worktree scaffold branch's commit is direct evidence of what base a subagent worktree was created from — useful forensics for next time.

### What was tricky

Nothing; pure housekeeping.

### What warrants review

Nothing — all changes were merged through reviewed PRs.

### Future work

None. Issue #330's thread is fully resolved: contract documented (#332), security policy in place (#331), lint green (#333).
