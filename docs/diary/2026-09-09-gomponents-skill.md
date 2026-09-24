# Diary: A gomponents skill for coding agents

Write a skill that teaches coding agents how to use gomponents, validate it with evals against runs without the skill, take it through review, and ship it as a Claude Code plugin. The work ran from 2026-09-09 to 2026-09-17 and ended in PR #358, the plugin manifests in `/.claude-plugin`, and three follow-up issues on the library.

## Step 1: Read the project and draft the skill

**Author:** main

### Prompt Context

**Verbatim prompt:** "I'd like to create a "gomponents" skill. Read through everything in this project, but avoid the now-deleted LLMs.md and AGENTS.md changes that's still in the git workspace, those might color your content. When you've read everything, use your skill-creator skill to create the skill to save directly in this repo (probably under a skills/ dir?). One thing to avoid is enumerating all available HTML elements and attributes, that's pretty redundant. Any questions before you start?"
**Interpretation:** Read every source, test, and doc in the repo except the staged deletions, then write a `SKILL.md` under `/skills` that documents how to use the library without listing elements and attributes, using the skill-creator process.
**Inferred intent:** Give agents a compact, accurate guide that produces idiomatic gomponents code, grounded in the library itself rather than in prior AI-written docs.

### What I did

Read `/gomponents.go`, `/html/elements.go`, `/html/attributes.go`, `/components/components.go`, `/http/handler.go`, `/x/slices/slices.go`, every test file, the example app under `/internal/examples/app`, the README, and the two existing diaries. Asked two questions before writing: audience (app developers, not contributors) and import style (dot imports, as in the README). Drafted `/skills/gomponents/SKILL.md` at about 300 lines: the `Node` model, imports, the core helpers, the `html` package's naming rules and the element-versus-attribute clashes, `components.HTML5`, `Classes`, `JoinAttrs`, `http.Adapt`, patterns, gotchas, and testing. Markus asked me to take the frontmatter description from his existing personal gomponents skill, since it had been tuned for triggering, and not to load that skill otherwise so it wouldn't colour the content.

### Why

The skill had to be derived from the library, so reading everything came first. Not enumerating elements pushes the agent toward `go doc` for names, which stays correct as the library grows.

### What worked

Asking the two scoping questions up front. Both answers shaped every example in the file.

### What didn't work

Nothing failed in this step.

### What I learned

`JoinAttrs` flattening nested groups landed on `main` after the v1.3.0 release, so a skill written against `main` describes behaviour a `go get` user doesn't have yet. Markus decided the skill documents `main` and that he'd cut a release when it lands.

### What was tricky

Keeping the description's trigger text while dropping its reference to a personal Go skill that other users don't have.

### What warrants review

The frontmatter description. It's deliberately pushy so agents load the skill on any HTML work.

### Future work

Evals, to find out what the skill changes.

## Step 2: Build an eval harness and run three rounds

**Author:** main

### Prompt Context

**Verbatim prompt:** "Run all three, and make sure they use the new skill, not the old one from fabrik" followed by "The promtps are very specific, that's not how I would prompt. E.g. "Convert this HTML into a gomponents component in Go", then the HTML code. And you're mentioning the "views" package in the prompt, intentional?"
**Interpretation:** Run each test prompt twice, once with the skill and once without, in subagents that must not load any other skill, then rewrite the prompts to sound like a real request.
**Inferred intent:** Measure what the skill actually changes in generated code, with prompts that don't leak the answer.

### What I did

Set up the skill-creator workspace outside the repo, in the session scratchpad. Each eval spawned two subagents on a fresh Go module: one told to read only the skill file, one told to use nothing but its own knowledge and `go doc`. Graded with a small script for the mechanical checks (vet and test pass, import path, deprecated aliases, dot imports, package name) and by reading the source for the rest. Round one used prompts I'd over-specified, naming a `views` package and demanding go.mod, tests, and httptest; Markus called that out, and round three used three one-line prompts instead. The eval viewer needed three patches to be usable at all: it bound to localhost while the reviewer was on another machine, it listed only top-level output files while the Go sources sat in subdirectories, and one baseline test file contained a literal `</script>`, which terminated the viewer's inline script and left the page blank.

### Why

Without a baseline there's no way to know whether the skill helps or just costs tokens. Natural prompts matter because a prompt that says "put it in package views" overrides the skill's own advice.

### What didn't work

The viewer refused to start under the system Python: `TypeError: unsupported operand type(s) for |: 'type' and 'NoneType'` at `def build_run(root: Path, run_dir: Path) -> dict | None:`, because macOS ships Python 3.9. Running it through `uv run --python 3.12` fixed that. `python -m scripts.aggregate_benchmark` failed with `command not found: python`; the machine has `python3` only. The blank page took two wrong guesses (caching, nested files) before the `</script>` cause showed up; the fix was escaping `</` in the embedded JSON.

### What worked

Grading with probes that run the output, not just read it. The `</script>` bug was found by grepping the output directories for the string once the page data proved correct.

### What I learned

Baselines fetch and read library source from the module cache when unsure, and get things right by doing so. That makes "did it read the source" a better signal of skill value than "was the output correct".

### What was tricky

Subagents inherit every installed skill. The harness had to forbid loading any skill so the personal gomponents skill couldn't leak into either configuration.

### What warrants review

The eval prompts and assertions live in the scratchpad, not the repo. They were deliberately kept out of the tree.

### Future work

Iteration 4, with assertions agreed one at a time.

## Step 3: Real-app conventions and a concision pass

**Author:** main

### Prompt Context

**Verbatim prompt:** "And make sure the skill file itself is concise and uses fairly simple language. Use a writer subagent for a pass if you need to." followed by "And honestly, keep the skill concise and on-subject. We're not going to cover all kinds of use cases. We want to document gomponents, how to use it effectively, common patterns."
**Interpretation:** Fold in the conventions from one of Markus's production apps without naming it, then cut the file down to the library, effective use, and common patterns.
**Inferred intent:** A skill that reads like documentation, not a style guide for one codebase.

### What I did

Read the views and handlers of one of Markus's production apps to learn the conventions the library docs don't state: components in a package named `html`, exported `XxxPage(props)` functions, unexported camelCase building blocks named after the element they wrap, and tests that check decisions with `strings.Contains`. Folded those in, then ran a writer subagent over the prose and cut the long profile-page example and the app-architecture bullets myself, landing at 300 lines.

### Why

The eval feedback was almost entirely about taste the library can't express: package name, import style, test shape. The app showed what the taste actually is.

### What worked

The writer pass on prose; it shortened sentences without touching code. Compiling every Go block in the skill against `main` in a scratch module, which later caught a real error.

### What didn't work

The writer couldn't hit the 280-line target because the file was mostly code blocks and tables; it reported that honestly instead of cutting content. The structural cut had to be mine.

### What I learned

A skill's length is set by its examples. Trimming prose saves lines only at the margin.

### What was tricky

Encoding conventions from a private app in a public file without naming the app or copying anything identifying.

### What warrants review

Nothing in the skill names or quotes the app.

### Future work

None.

## Step 4: Decide what the skill must fix, one item at a time

**Author:** main

### Prompt Context

**Verbatim prompt:** "Ideally you'd want a considerably worse performance in the no-skill version, right? Maybe look harder for those which the skill fixes." then "Let's make a list of gotchas that the skill could help with. Let's iterate on the list together." then "Let's take them one at a time"
**Interpretation:** Build the assertion list from observed without-skill behaviour, agreeing each item individually, then run the evals against that list.
**Inferred intent:** An eval harness that measures the skill's value, not generic code quality.

### What I did

Listed fifteen candidate behaviours, eight observed in baseline runs and seven predicted from the skill's gotchas, and went through them one per message. Kept the observed ones (deprecated `TitleAttr`, aliased imports, wrong package, exported `Button` colliding with the element, whole-page golden tests, hand-rolled `WriteHeader`, hand-built class strings, and reading library source). Dropped the whitespace, entity, script-body, and standalone-`Group` items as too detailed or unrealistic. Added a nil-user navbar probe and a nested-component `JoinAttrs` probe. Ran iteration four with ten subagents; with the skill 33 of 37 assertions passed, without it 25 of 37.

### Why

Markus wanted a bigger gap between the configurations, and the way to get one honestly is to assert on things the baseline gets wrong, not on style it happens to share.

### What worked

Making "read the library source" a pass/fail assertion. Every baseline did it; three with-skill runs did too, and each of those pointed at a gap in the skill.

### What didn't work

Two of the three with-skill lookups turned out to be library questions, not skill gaps. The `Adapt` lookup led to the finding that a wrapped error loses its status code, since `Adapt` uses a type switch rather than `errors.As`; filed as #353. The `Content-Type` question became #351 after a probe showed fragments starting with `<tr>`, `<li>`, `<span>`, `<form>`, or `<button>` are served as `text/plain`, and that `httptest.ResponseRecorder` reports `text/plain` for every element fragment because `El` writes `<` on its own first.

### What I learned

The session model changed between rounds, so cross-round comparisons are unreliable; only within-round comparisons mean anything.

### What was tricky

Deciding that `go doc` counts as allowed while reading source in `pkg/mod` counts as a failure. Markus's rule: `go doc` is fine but shouldn't be needed if the skill is right.

### What warrants review

Issues #351 and #353. The skill assumes both behave correctly, per Markus's call, so fixing them needs no skill change.

### Future work

Fixes for #351 and #353 in the library.

## Step 5: The PR and three review rounds

**Author:** main

### Prompt Context

**Verbatim prompt:** "yes, open the PR" and, three times, "/fabrik:address-code-review"
**Interpretation:** Open PR #358 from the fork, then work through each round of inline review comments one at a time, triaging every comment with Markus before changing anything.
**Inferred intent:** A skill that reads the way Markus would write it.

### What I did

Opened #358 with one file. Round one had 13 comments: say "component" not "view", drop the contributor policy, replace "lazy" (component calls run immediately; only rendering is deferred), move the core table up next to a fuller example, fold conditionals and gotchas into one Common patterns section, move linting to `/skills/gomponents/references/linting.md`. Round two had 13 more: open with gomponents as a DSL for HTML in valid Go, make `JoinAttrs` the one way callers add classes, drop the class-string helper and the props, plain-values, and element-function patterns, keep htmx out. Round three had 26: "like a DSL", note that gomponents doesn't check for valid HTML by design, "unsanitized" everywhere user content is mentioned, all six `HTML5Props` fields in the layout example, `Adapt` bullets for pages versus fragments, `go doc` instead of pkg.go.dev, Codex install commands behind `<details>` in the README. A fourth, one-comment round dropped the email from the plugin manifests. Replied to and resolved every thread; filed #359 for public test helpers when Markus suggested it.

### Why

Every comment was a taste call only the maintainer can make, and the address-code-review process keeps me from acting on one before it's decided.

### What worked

Compiling the skill's examples after each round. It caught `type Link struct` colliding with the html package's `Link` element: `vet: html/navbar.go:9:6: Link already declared through dot-import of package html`. The example type is now `NavLink`.

### What didn't work

The first push failed with `Please make sure you have the correct access rights and the repository exists.` because the `fork` remote had vanished from the checkout; re-adding it fixed the push, and the first `gh pr create` then failed with `Head sha can't be blank, Base sha can't be blank, No commits between maragudk:main and maragubot:add-gomponents-skill` because nothing had been pushed yet. On one comment I resolved the thread before Markus had said apply; he pointed it out, and the recording stood only because he agreed.

### What I learned

The re-export pattern for testing unexported blocks, `var Card = card` in an internal test file, collides with dot-imported element names exactly as an exported function would, so it works for `card` and not for `button`.

### What was tricky

Three rounds moved the same paragraphs several times. Keeping each edit to what was agreed, and nothing more, took discipline.

### What warrants review

The final structure: mental model, common patterns, imports and layout, then the three packages, then testing.

### Future work

#359, public test helpers, which would shrink the testing section.

## Step 6: Package the skill as a plugin

**Author:** main

### Prompt Context

**Verbatim prompt:** "Maybe we should make it a Claude plugin like in ../fabrik?" then "Actually, first, figure out whether the Claude plugin system relies on Git tags or Github releases, because we can't be making those in this repo. They're tied to the Go module."
**Interpretation:** Add the plugin and marketplace manifests so the skill installs with two commands, but only after confirming the plugin system won't need tags that would collide with the Go module's.
**Inferred intent:** Distribution without a second release process.

### What I did

Checked the local plugin cache and the docs. Installed plugins record a commit SHA, not a tag; a plugin with no `version` uses the short commit hash and updates on every commit; the plugin CLI's own `tag` command would produce `gomponents--v1.0.0`, which can't collide with `v1.x`. Added `/.claude-plugin/plugin.json` and `/.claude-plugin/marketplace.json` with no version, validated them with `claude plugin validate .`, and wrote the README install section.

### Why

The repo already had the plugin layout: a `skills/` directory at the root. Two JSON files made it installable.

### What worked

`claude plugin validate` for every manifest change. It's what showed that `owner` is required and its email isn't: `✘ Found 1 error: ❯ owner: Invalid input` without the field, a pass with only a name.

### What didn't work

Nothing failed beyond that validation probe, which was the point of running it.

### What I learned

Third-party marketplaces have auto-update off by default; only Anthropic's are on. Users of the plugin get updates through `claude plugin update` unless they turn it on.

### What was tricky

Deciding the marketplace name. fabrik's is `maragu`, and a second marketplace with that name would clash for anyone with both, so this one is `gomponents`.

### What warrants review

`/.claude-plugin/plugin.json` has no version on purpose.

### Future work

A README note on enabling auto-update, if users are expected to stay current.

## Step 7: Merge and verify the install

**Author:** main

### Prompt Context

**Verbatim prompt:** "merge 358" then "Can you see the new plugin?"
**Interpretation:** Merge the PR with a merge commit, clean up, then confirm the installed plugin exposes the skill.
**Inferred intent:** Ship it and prove it works from a user's seat.

### What I did

The first merge attempt was refused: `Pull request maragudk/gomponents#358 is not mergeable: the head branch is not up to date with the base branch.` Merged `main` into the branch, pushed, enabled auto-merge with a merge commit, and it landed as `aa563d7` once CI passed. Pulled `main`, deleted the branch locally and on the fork. Markus installed the plugin; the skill loaded from the plugin cache in full. In the same session it showed without a description, which a restart fixed: a subagent's fresh system prompt listed `gomponents:gomponents` with the complete trigger text, and Markus then removed his personal copy so only one gomponents skill remains.

### Why

Merging was an explicit, standalone instruction naming the PR, which is the one case the standing "don't merge" rule allows.

### What worked

`gh pr merge --merge --auto` after updating the branch.

### What didn't work

The blank description after a mid-session `/reload-plugins`. It's a listing artifact of reloading, not a manifest problem; the cached frontmatter parsed correctly.

### What I learned

The plugin cache holds the whole repository, since the plugin source is the repo root, so an install is a few megabytes rather than the skill's few kilobytes.

### What was tricky

Nothing.

### What warrants review

That the personal skill in fabrik is gone, so the two never fire together.

### Future work

The release Markus plans to cut now that the skill documents `main`.
