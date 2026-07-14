---
name: audit-deps
description: Scan this Go module's dependencies for known vulnerabilities with govulncheck, upgrade any that have a non-breaking fix available, leave breaking ones untouched, and produce a report of what was vulnerable/fixed/still-needs-manual-upgrade. Use when asked to check for vulnerable dependencies, run a dependency security audit, update vulnerable packages, or generate a vulnerability report for this repo.
---

# Vulnerability scan & update

Uses the repo's existing tooling — `govulncheck` is already declared as a `tool` dependency in
`go.mod` and wired up as `make govulncheck`. Do not write ad-hoc parsing scripts; use
`govulncheck -json` piped through `jq` for structured data, and plain `go`/`git` commands to
apply and verify each upgrade.

## 0. Preconditions

- Run `git status --porcelain`. If the tree isn't clean, stop and ask the user how to proceed —
  do not stash or discard their work silently.
- Confirm Go toolchain matches `.tool-versions` (`go version`). If missing, tell the user to
  `asdf install golang <version>` rather than trying to work around it.

## 1. Scan

Run the scan with JSON output so results can be parsed reliably instead of scraping text:

```bash
go tool govulncheck -json ./... > /tmp/govulncheck.json
```

Extract the actionable facts with `jq`:

```bash
# Vulnerabilities that actually affect this module (not just present in go.sum)
jq -s '[.[] | select(.finding.trace) ]' /tmp/govulncheck.json
# Module-level fix info
jq -s '[.[] | select(.osv) | {id: .osv.id, modules: [.osv.affected[].package.name], fixed: [.osv.affected[].ranges[]?.events[]? | select(.fixed) | .fixed]}]' /tmp/govulncheck.json
```

Build a list of `{module, current_version, vulnerable_id, fixed_version}` tuples. A module can
appear multiple times (multiple advisories) — take the highest fixed version per module.

If the list is empty, report "no known vulnerabilities found" and stop — no need to touch
`go.mod`.

## 2. Attempt an upgrade per vulnerable module

Process modules **one at a time** so a bad upgrade can be isolated and reverted without affecting
the others. For each module `M` with fixed version `V`:

```bash
go get M@V
go mod tidy
go build ./... && go vet ./cmd/... ./pkg/...
go test ./...
```

- If all of the above succeed: keep the change, record `M` as **Fixed** (old → new version).
- If `go get`/`go mod tidy` themselves fail (e.g. no such version, incompatible Go version
  requirement), or `go build`/`go vet`/`go test` fail: revert just this module's change with
  `git checkout -- go.mod go.sum` (restoring the pre-attempt state, since each module is handled
  in isolation) and record `M` as **Breaking / needs manual upgrade**, including the tail of the
  failing command's output as the reason.

Re-run `go mod tidy` once at the end after all successful upgrades are applied together, to make
sure `go.sum` is fully consistent.

## 3. Verify

Re-run the scan to confirm which advisories are resolved and which remain:

```bash
go tool govulncheck -json ./... > /tmp/govulncheck-after.json
```

Also run `make lint` if available, since golangci-lint can catch issues govulncheck/build don't.

## 4. Report

Produce a markdown report (print it in the response; only write it to a file if the user asks
for one) with three sections:

1. **Vulnerable dependencies found** — table of module, installed version, advisory ID(s)
   (e.g. GO-2024-XXXX), and whether it's a direct or indirect dependency.
2. **Fixed automatically** — module, old version → new version, advisory IDs resolved.
3. **Breaking changes — needs manual decision** — module, current version, fixed version,
   advisory IDs still open, and a one-line reason the automatic upgrade was reverted (compile
   error, test failure, etc).

If section 3 is non-empty, explicitly ask the user (e.g. via AskUserQuestion, one question per
module or grouped if related) whether they want to proceed with each breaking upgrade. Do not
force the upgrade yourself — breaking changes need the user's call. Do not commit anything;
leave the working tree for the user to review and commit.
