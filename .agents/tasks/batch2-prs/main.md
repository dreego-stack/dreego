# Batch 2 — 3 Parallel PRs (v0.1 Blockers, merge after Batch 1)

## Workflow per PR (TDD Loop)
Test → Code → Docs → Review → Push. No merges. PR creation only. Changelog is NEVER edited manually (CI handles it via pr.md).

## PR6: docs-accessibility (feat/docs-accessibility)
- Worktree: .worktrees/docs-accessibility
- Todos: docs-drift.1 + accessibility.1
- Files: _docs/*, README.md, cli/dreego/* (diagnostics), core/codegen_*.go (blueprint templates)
- Status: coded
- Coder: coder (ses_pr-queue)
- Done: Implemented core/a11y_check.go (Diagnostic, a11yCheck, a11yDiagnostics for img-alt and input-label). Wired a11y warnings into generate.go. Added SourceText field to TemplateNode + setSourceText helper; parser sets Pos on NodeText. CLI errors.go reformats parse/lex errors to file:line:col + cause + Fix:. Blueprints: landing layout has main/nav/skip-link/slot; Hero uses section. Docs: rewrote testing.md (no stale counts/⬜), added _docs/accessibility.md, fixed README (147 tests removed), sitemap+index updated. pr.md created (patch). Todo files git rm'd. NOT committed/pushed.

## PR7: error-propagation (feat/error-propagation)
- Worktree: .worktrees/error-propagation
- Todos: error-propagation.1
- Files: core/generate.go, core/forms.go, core/session*.go, core/codegen*.go, core/context.go
- Status: changes-requested
- Coder: coder (ses_ff08be596ffeRnTxAYWk1Uq2t0)
- Done: TDD loop complete. Tests: core/error_propagation_test.go (generate abort on walk/readdir/source-read/layout failures, session/form error reachability via c.SessionError()/c.FormError(), generic 500 no-disclosure in NewHandler/Recovery, CSRF store failure fails loudly, Listen error propagation), core/codegen_error_propagation_test.go (self-close Render error propagation, form re-render error check, generated handlers generic 500), _tests/go/error_propagation_test.go (E2E generic 500 no-disclosure for panic/component/form-action/session/CSRF, form bind error field via c.Errors("_form")). Code: generate.go/generate_components.go propagate walk/readdir/read errors with path+cause; context.go adds FormError()/SessionError() and records store/parse failures; middleware_csrf.go fails loudly with generic 500 + logs on store errors; response.go/NewHandler/codegen_page/codegen_error/codegen_form emit generic 500 (no err.Error() to clients); codegen_form no longer discards render errors; codegen_template self-close Render error panics (reached by Recovery); middleware_compress.go no longer commits 200 on panic with gzip. Docs: _docs/forms.md, _docs/runtime.md. Golden files updated. gofmt fixed 3 pre-existing violations. Verify: build/vet/core/CLI/_tests-go all green, check-core-deps pass. pr.md created (version: patch). Todo file deleted. NOT committed/pushed.
- Review: CHANGES-REQUESTED by reviewer (ses_ff020a60cffenVkF4fdANXwA0F) — Blocker: 3 permission-based tests skip as root; CI (golang:1.22-alpine Dockerfile, root) treats any '--- SKIP:' in core output as FAIL (test.sh:78-82) → CI red. Warning: Compress panic fix lacks the mandatory regression test. See reviewer session for details.
- Fixes applied (coder, ses_review-fixes): BLOCKER fixed — removed all os.Geteuid()==0 / chmod-0000 skip guards from TestRunAbortsOnReadDirFailure, TestRunAbortsOnWalkError, TestRunAbortsOnLayoutReadFailure in core/error_propagation_test.go:78-139. Replaced the chmod approach (root-bypassed) with broken symlinks: each "blocked" path now uses os.Symlink to a missing target, so os.ReadFile fails with ENOENT deterministically for root AND non-root. Note: the task's "regular file instead of dir" suggestion does NOT trigger the error paths in the current code (WalkDir returns nil for non-dirs before os.ReadDir is called), so broken-symlink → ReadFile-failure was used instead — it exercises the same error-propagation code (generate.go:115-118 ReadFile path, generate.go:480-491 findLayout ReadFile path) and is root-safe. WARNING fixed — added core/middleware_compress_panic_test.go: TestCompressPanicDoesNotCorruptResponse stacks Recovery(nil) over Compress() over a panicking handler with Accept-Encoding: gzip; asserts status 500, no corrupt gzip body, Content-Encoding header removed. Verifies the middleware_compress.go:18-26 recover→Del→re-panic fix. Verification (single run): go build ./... ok; go test ./core/... ok; go test ./_tests/go/ -parallel 1 -p 1 ok (60s); gofmt -l core/ _tests/ empty; go test -v ./core/ grep -c SKIP = 0; all 3 previously-skipping tests PASS; new Compress panic test PASS. NOT committed/pushed. next: shell commit + reviewer re-review.

## PR8: reference-apps (feat/reference-apps)
- Worktree: .worktrees/reference-apps
- Todos: reference-apps.1
- Files: new fixture apps under _tests/fixtures/, dreegotest/fixture.go, _docs/reference-apps.md
- Status: approved
- Coder: done — 4 fixture apps (_tests/fixtures/{hello,forms,components,plugin}), dreegotest.ServeFixture helper, _tests/go/reference_apps_test.go (4 tests), _docs/reference-apps.md + index/sitemap entries, pr.md (version: patch), todo file git-rm'd. All tests green.
- Reviewer: APPROVED (ses reviewer, read-only pass 2) — previous blocker (go.mod replace `../..` → `../../..`) fixed in all 4 fixtures. 8 review dimensions pass: (1) fixtures use only public APIs (dreego.New, gen.Register, app.Listen, app.SetSessionStore, dreego.NewCookieStore, app.Register); (2) 4 areas covered (minimal/forms+sessions/components/plugin); (3) ServeFixture does generate→build→start→HTTP-verify per app; (4) plugin fixture uses realistic Register(app, Options) pattern; (5) documented in _docs/reference-apps.md + index + sitemap; (6) all 4 go.mod line 7 = `replace github.com/dreego-stack/dreego => ../../..`; (7) coding rules ok (max 300 lines, stdlib+core only, no comments in fixtures, godoc comments in dreegotest consistent with package); (8) pr.md valid (version: patch, 3 changelog lines). Notes (non-blocking): forms counter uses len(string) as count (semantically odd but works); generated gen/ artifacts in worktree are gitignored (dreego/gen/) — fresh clone relies on `dreego generate` to create the dir, which it does (generate.go:285). Ready to commit + push + PR.
next: shell
## Merge strategy
- Batch 2 branches are based on main @ 9396310 (before Batch 1 merges)
- After Batch 1 PRs (#43-#47) merge, rebase Batch 2 branches onto new main
- Conflicts expected to be minimal (different files/functions)

## Hand-off protocol
- Coder writes test + code + docs, marks status "review-pending" in this file
- Reviewer reviews, marks "approved" or "changes-requested"
- Shell commits + pushes + creates PR, marks "pr-created"
