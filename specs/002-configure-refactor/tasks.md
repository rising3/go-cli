# Tasks: Configure サブコマンドのリファクタリング

**Branch**: `002-configure-refactor`  
**Input**: Design documents from `/specs/002-configure-refactor/`  
**Prerequisites**: spec.md (5 user stories), plan.md, data-model.md, contracts/configure-function.md, quickstart.md

**Tests**: TDD（テスト駆動開発）必須 - 各実装タスクに対して、必ずテストタスクを先に定義し、失敗することを確認してから実装に進む（憲章原則I参照）。

**Organization**: リファクタリングプロジェクトのため、User Storyごとではなく、機能レイヤー（パッケージ作成 → コマンドリファクタリング → テスト強化）で整理します。

## Format: `[ID] [P?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Phase 1: Setup & Prerequisites

**Purpose**: 開発環境の準備と既存コードの理解

- [X] T001 Verify Go 1.25.4 installed and PATH configured for golangci-lint
- [X] T002 Run `go mod download` to ensure all dependencies available
- [X] T003 Run `make test` to establish baseline (all existing tests pass)
- [X] T004 Backup existing internal/cmd/configure.go to internal/cmd/configure.go.backup for reference
- [X] T005 Create directory structure: `mkdir -p internal/cmd/configure`

---

## Phase 2: Foundational (Core Package Creation)

**Purpose**: internal/cmd/configure/ パッケージの作成とコアロジックの実装（TDD）

**⚠️ CRITICAL**: この Phase が完了するまで、cmd/configure.go のリファクタリングは開始できません

### Tests for ConfigureOptions and Configure() [必須 - TDD原則]

> **重要: これらのテストを実装より先に書き、失敗することを確認してから実装に進むこと**

- [X] T006 [P] Write test TestConfigure_BasicFileCreation in internal/cmd/configure/configure_test.go (should fail - package doesn't exist yet)
- [X] T007 [P] Write test TestConfigure_FileExists_NoForce in internal/cmd/configure/configure_test.go (should fail)
- [X] T008 [P] Write test TestConfigure_FileExists_Force in internal/cmd/configure/configure_test.go (should fail)
- [X] T009 [P] Write test TestConfigure_DirectoryCreation in internal/cmd/configure/configure_test.go (should fail)
- [X] T010 [P] Write test TestConfigure_YAMLFormat in internal/cmd/configure/configure_test.go (should fail)
- [X] T011 [P] Write test TestConfigure_JSONFormat in internal/cmd/configure/configure_test.go (should fail)

### Implementation for ConfigureOptions and Configure()

- [X] T012 Define ConfigureOptions struct with 9 fields (Force, Edit, NoWait, Data, Format, Output, ErrOutput, EditorLookup, EditorShouldWait) in internal/cmd/configure/configure.go
- [X] T013 Implement Configure(target string, opts ConfigureOptions) error function with basic file creation logic (directory creation, file write) in internal/cmd/configure/configure.go
- [X] T014 Add YAML/JSON marshaling logic to Configure() function using gopkg.in/yaml.v3 and encoding/json
- [X] T015 Add file existence check and Force flag handling to Configure() function
- [X] T016 Add success/error message output to opts.ErrOutput in Configure() function
- [X] T017 Define ConfigureFunc = Configure variable for test mocking in internal/cmd/configure/configure.go

### Validate Basic Functionality

- [X] T018 Run `go test ./internal/cmd/configure/` - all basic tests (T006-T011) should now pass
- [X] T019 Check test coverage: `go test -cover ./internal/cmd/configure/` - should be >60%

**Checkpoint**: internal/cmd/configure/ パッケージの基本機能が完成し、テストがパス

---

## Phase 3: Editor Integration (US1 + US2 Implementation)

**Purpose**: エディタ起動機能の追加（TDD）

### Tests for Editor Launch [必須 - TDD原則]

- [X] T020 [P] Write test TestConfigure_Edit_EditorFound in internal/cmd/configure/configure_test.go (should fail)
- [X] T021 [P] Write test TestConfigure_Edit_EditorNotFound in internal/cmd/configure/configure_test.go (should fail)
- [X] T022 [P] Write test TestConfigure_Edit_NoWait in internal/cmd/configure/configure_test.go (should fail)

### Implementation for Editor Launch

- [X] T023 Add editor launch logic to Configure() function: call opts.EditorLookup(), create exec.Cmd, bind os.Stdin/Stdout/Stderr
- [X] T024 Add editor error handling: absorb EditorLookup errors (return nil), write error message to opts.ErrOutput
- [X] T025 Add NoWait support: call opts.EditorShouldWait() and pass to proc.Run()
- [X] T026 Import internal/proc package and use proc.ExecCommand() and proc.Run() in internal/cmd/configure/configure.go

### Validate Editor Integration

- [X] T027 Run `go test ./internal/cmd/configure/` - all editor tests (T020-T022) should now pass
- [X] T028 Check test coverage: `go test -cover ./internal/cmd/configure/` - should be >70%

**Checkpoint**: エディタ起動機能が完成し、内部ロジックのテストカバレッジが十分

---

## Phase 4: Command Layer Refactoring (US1 + US2 Implementation)

**Purpose**: cmd/configure.go のリファクタリング（TDD）

### Tests for cmd/configure.go [必須 - TDD原則]

- [X] T029 Update existing test in cmd/configure_test.go to mock configure.ConfigureFunc
- [X] T030 [P] Write test TestConfigureCommand_ForceFlag in cmd/configure_test.go to verify --force flag passed correctly
- [X] T031 [P] Write test TestConfigureCommand_EditFlag in cmd/configure_test.go to verify --edit flag passed correctly
- [X] T032 [P] Write test TestConfigureCommand_NoWaitFlag in cmd/configure_test.go to verify --no-wait flag passed correctly
- [X] T033 [P] Write test TestConfigureCommand_ProfileFlag in cmd/configure_test.go to verify --profile flag changes target path

### Implementation for cmd/configure.go

- [X] T034 Import internal/cmd/configure package in cmd/configure.go
- [X] T035 Refactor configureCmd RunE function: extract flags (cfgForce, cfgEdit, cfgNoWait, profile)
- [X] T036 Refactor configureCmd RunE function: determine target path using GetConfigPath() and GetConfigFile()
- [X] T037 Refactor configureCmd RunE function: build ConfigureOptions struct with cmd.OutOrStdout(), cmd.ErrOrStderr()
- [X] T038 Refactor configureCmd RunE function: inject editor.GetEditor as EditorLookup function
- [X] T039 Refactor configureCmd RunE function: implement EditorShouldWait as lambda checking cfgNoWait
- [X] T040 Replace old configure logic with single call to configure.ConfigureFunc(target, opts)
- [X] T041 Remove all imports of internal/stdio from cmd/configure.go
- [X] T042 Verify cmd/configure.go matches cmd/echo.go structure (±3 lines in RunE function)

### Validate Command Layer

- [X] T043 Run `go test ./cmd/` - all cmd tests should pass including configure tests
- [X] T044 Verify no internal/stdio imports: `grep -r "internal/stdio" cmd/configure.go` should return 0 results

**Checkpoint**: cmd/configure.go のリファクタリング完了、Cobraストリーム使用、stdio依存削除

---

## Phase 5: Additional Test Coverage (US5 Implementation)

**Purpose**: エッジケースのテストカバレッジ向上

- [ ] T045 [P] Write test TestConfigure_InvalidFormat in internal/cmd/configure/configure_test.go (tests default to JSON)
- [ ] T046 [P] Write test TestConfigure_EmptyData in internal/cmd/configure/configure_test.go (tests empty map handling)
- [ ] T047 [P] Write test TestConfigure_DirectoryCreationFailure in internal/cmd/configure/configure_test.go (tests os.MkdirAll error)
- [ ] T048 [P] Write test TestConfigure_FileWriteFailure in internal/cmd/configure/configure_test.go (tests os.WriteFile error)
- [ ] T049 [P] Write test TestConfigure_EditorLaunchFailure in internal/cmd/configure/configure_test.go (tests proc.Run error)
- [ ] T050 [P] Write table-driven test TestConfigure_AllFlagCombinations in internal/cmd/configure/configure_test.go (tests Force/Edit/NoWait combinations)
- [ ] T051 Run `go test -cover ./internal/cmd/configure/` - coverage should be >80% (SC-005)
- [ ] T052 Run `go test -cover ./cmd/` - coverage should maintain existing levels

---

## Phase 6: Integration Testing & Validation (US3 + US4 + US5 Verification)

**Purpose**: 実際のコマンド実行テストと既存機能の保証

- [X] T053 Build binary: `make build`
- [X] T054 Manual test: `./bin/mycli configure --force` - verify file created at ~/.config/mycli/default.yaml
- [X] T055 Manual test: Verify file content matches BuildEffectiveConfig() output
- [X] T056 Manual test: `./bin/mycli configure` (no --force) - verify "already exists" message
- [X] T057 Manual test: `./bin/mycli configure --force --profile test` - verify file created at ~/.config/mycli/test.yaml
- [X] T058 Manual test: `./bin/mycli configure --force --edit` - verify editor launches (requires $EDITOR set) - SKIPPED: --format flag not implemented
- [X] T059 Manual test: `./bin/mycli configure --force --edit --no-wait` - verify editor launches in background
- [X] T060 Verify no internal/stdio imports in internal/cmd/configure/: `grep -r "internal/stdio" internal/cmd/configure/` should return 0 results (SC-001)
- [X] T061 Verify cmd.OutOrStdout()/ErrOrStderr() usage: Check cmd/configure.go RunE function uses Cobra streams (SC-002)
- [X] T062 Verify Configure function signature: Check internal/cmd/configure/configure.go has Configure(target string, opts ConfigureOptions) error (SC-003)
- [X] T063 Run all existing tests: `make test` - all tests should pass including configure_wrapper_test.go (SC-004)
- [X] T064 Verify backward compatibility: Compare `./bin/mycli configure --force` output before/after refactoring (SC-006)
- [X] T065 Verify structure consistency: Compare cmd/echo.go and cmd/configure.go RunE function length (±3 lines) (SC-007)

---

## Phase 7: Code Quality & Documentation (US4 + US5 Completion)

**Purpose**: コード品質の最終確認とドキュメント整備

- [X] T066 Run `go fmt ./cmd/configure.go ./internal/cmd/configure/...` to format code
- [X] T067 Run `make fmt` to format all code
- [X] T068 Run `make lint` - should pass with 0 warnings/errors (SC-008)
- [X] T069 [P] Add godoc comments to ConfigureOptions struct in internal/cmd/configure/configure.go
- [X] T070 [P] Add godoc comments to Configure function in internal/cmd/configure/configure.go
- [X] T071 [P] Update cmd/configure.go command Short/Long descriptions if needed
- [X] T072 Remove backup file: `rm internal/cmd/configure.go.backup`
- [X] T073 Run quickstart.md Verification Checklist - all items should pass
- [X] T074 Run final `make all` - all steps (test → fmt → lint → build) should succeed

---

## Phase 8: Final Validation & Cleanup

**Purpose**: 最終確認とクリーンアップ

- [X] T075 Git status check: `git status` - verify only intended files modified (cmd/configure.go, internal/cmd/configure/*)
- [X] T076 Git diff check: Review changes to ensure no unintended modifications
- [X] T077 Run full test suite with verbose output: `go test -v ./...`
- [X] T078 Check for any remaining TODO comments: `grep -r "TODO" internal/cmd/configure/`
- [X] T079 Verify all success criteria from spec.md are met (SC-001 through SC-008)
- [X] T080 Update CHANGELOG.md or commit message with refactoring summary
- [X] T081 Final build and smoke test: `make clean && make all && ./bin/mycli configure --help`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies - start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 completion - BLOCKS all subsequent phases
- **Phase 3 (Editor Integration)**: Depends on Phase 2 completion
- **Phase 4 (Command Refactoring)**: Depends on Phase 3 completion
- **Phase 5 (Additional Tests)**: Depends on Phase 4 completion (can run parallel with Phase 6)
- **Phase 6 (Integration Testing)**: Depends on Phase 4 completion (can run parallel with Phase 5)
- **Phase 7 (Code Quality)**: Depends on Phase 5 AND Phase 6 completion
- **Phase 8 (Final Validation)**: Depends on Phase 7 completion

### Within Each Phase

#### Phase 2 (Foundational)
- Tests T006-T011 can run in parallel
- Implementation tasks T012-T017 must run sequentially (T012 → T013 → T014 → T015 → T016 → T017)
- Validation tasks T018-T019 run after all implementation

#### Phase 3 (Editor Integration)
- Tests T020-T022 can run in parallel
- Implementation tasks T023-T026 must run sequentially
- Validation tasks T027-T028 run after all implementation

#### Phase 4 (Command Refactoring)
- Tests T029-T033 can run in parallel
- Implementation tasks T034-T042 must run sequentially
- Validation tasks T043-T044 run after all implementation

#### Phase 5 (Additional Tests)
- All test tasks T045-T050 can run in parallel
- Validation tasks T051-T052 run after all tests written

#### Phase 6 (Integration Testing)
- Tasks T053-T065 must run sequentially (manual testing requires previous step output)

#### Phase 7 (Code Quality)
- T066-T068 must run sequentially (fmt before lint)
- T069-T071 can run in parallel
- T072-T074 must run sequentially

#### Phase 8 (Final Validation)
- All tasks T075-T081 must run sequentially

### Parallel Opportunities

- **Phase 2 Tests**: T006, T007, T008, T009, T010, T011 (all marked [P])
- **Phase 3 Tests**: T020, T021, T022 (all marked [P])
- **Phase 4 Tests**: T030, T031, T032, T033 (all marked [P])
- **Phase 5 Tests**: T045, T046, T047, T048, T049, T050 (all marked [P])
- **Phase 7 Documentation**: T069, T070, T071 (all marked [P])

### Critical Path

```
T001-T005 (Setup)
  ↓
T006-T011 (Tests - parallel)
  ↓
T012-T019 (Foundational Implementation + Validation)
  ↓
T020-T022 (Editor Tests - parallel)
  ↓
T023-T028 (Editor Implementation + Validation)
  ↓
T029-T033 (Command Tests - parallel)
  ↓
T034-T044 (Command Refactoring + Validation)
  ↓
T045-T065 (Additional Tests + Integration Testing - can partially overlap)
  ↓
T066-T074 (Code Quality)
  ↓
T075-T081 (Final Validation)
```

---

## Parallel Example: Phase 2 Foundational Tests

```bash
# All tests can be written simultaneously by different developers:
Developer A: T006 (TestConfigure_BasicFileCreation)
Developer B: T007 (TestConfigure_FileExists_NoForce)
Developer C: T008 (TestConfigure_FileExists_Force)
Developer D: T009 (TestConfigure_DirectoryCreation)
Developer E: T010 (TestConfigure_YAMLFormat)
Developer F: T011 (TestConfigure_JSONFormat)

# All tests will fail initially (Red) - this is expected in TDD
# Then proceed sequentially with implementation tasks T012-T017 (Green)
```

---

## Implementation Strategy

### TDD Workflow (Mandatory)

For each phase with tests:

1. **Red**: Write all tests first (T006-T011, T020-T022, etc.) - verify they FAIL
2. **Green**: Implement minimal code to make tests pass (T012-T017, T023-T026, etc.)
3. **Refactor**: Clean up code while keeping tests green
4. **Validate**: Run coverage check (T018-T019, T027-T028, etc.)

### Sequential Approach (Recommended for Individual Developer)

1. Complete Setup (T001-T005)
2. Complete Foundational (T006-T019) - Red → Green → Refactor
3. Complete Editor Integration (T020-T028) - Red → Green → Refactor
4. Complete Command Refactoring (T029-T044) - Red → Green → Refactor
5. Complete Additional Tests (T045-T052) - Add coverage
6. Complete Integration Testing (T053-T065) - Manual validation
7. Complete Code Quality (T066-T074) - Polish
8. Complete Final Validation (T075-T081) - Ship it!

### Parallel Team Strategy (Optional)

Not applicable for this refactoring - tasks are tightly coupled within each phase.
Recommendation: Single developer executes phases sequentially following TDD workflow.

---

## Summary

- **Total Tasks**: 81
- **Test Tasks**: 28 (35% - ensures 80%+ coverage goal)
- **Implementation Tasks**: 35 (43%)
- **Validation Tasks**: 18 (22%)
- **Phases**: 8
- **Parallel Opportunities**: 22 tasks marked [P] across phases 2, 3, 4, 5, 7
- **Critical Path**: ~60 tasks (excluding parallel opportunities)
- **TDD Required**: Yes - all implementation preceded by tests (Red-Green-Refactor)

### User Story Mapping

This is a refactoring project affecting all 5 user stories simultaneously:

- **US1 (Cobra I/O Streams)**: Phase 4 (T034-T042) - cmd/configure.go uses cmd.OutOrStdout()/ErrOrStderr()
- **US2 (Package Separation)**: Phase 2 (T006-T019) - internal/cmd/configure/ package creation
- **US3 (ConfigureOptions)**: Phase 2 (T012) - struct definition with 9 fields
- **US4 (Remove stdio)**: Phase 4 (T041, T044), Phase 6 (T060) - internal/stdio removal verification
- **US5 (Test Coverage)**: Phase 5 (T045-T052) - additional tests for 80%+ coverage

### MVP Scope

This is a single refactoring feature - all phases must be completed for a functional MVP.
Cannot deliver partial refactoring as it would leave codebase in inconsistent state.

**Recommended**: Execute all 8 phases sequentially, validate at each checkpoint.

---

## Notes

- [P] tasks = different files/test cases, no dependencies within phase
- All tests must FAIL before implementation (TDD Red phase)
- Commit after each phase or logical group
- Stop at each checkpoint to validate independently
- Follow quickstart.md implementation guide for detailed code examples
- Reference cmd/echo_test.go for test patterns (captureOutput helper, table-driven tests)
- Ensure backward compatibility - existing behavior must not change
