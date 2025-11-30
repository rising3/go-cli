# Tasks: Cat サブコマンド実装

**Input**: Design documents from `/specs/004-cat-subcommand/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: このプロジェクトではTDD（テスト駆動開発）が必須です。すべての実装タスクに対して、必ずテストタスクを先に定義し、Red-Green-Refactorサイクルに従います（憲章原則I参照）。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `- [ ] [ID] [P?] [Story?] Description with file path`

- **Checkbox**: `- [ ]` ALWAYS at start
- **[ID]**: Sequential task number (T001, T002, T003...)
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story label (e.g., US1, US2) - ONLY for user story phases
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create directory structure: `internal/cmd/cat/` for cat subcommand implementation
- [X] T002 Verify Go 1.25.4 environment and download dependencies with `go mod download`
- [X] T003 [P] Verify golangci-lint 2.6.2 is installed and PATH is configured correctly

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Create `cmd/cat_wrapper_test.go` - wrapper test for Cobra command layer ✅ COMPLETED
- [X] T005 [P] Define Options struct in `internal/cmd/cat/options.go` per contracts/options.md
- [X] T006 [P] Define Formatter interface in `internal/cmd/cat/formatter.go` per contracts/formatter.md
- [X] T007 [P] Define Processor interface in `internal/cmd/cat/processor.go` per contracts/processor.md
- [X] T008 Create control character map constant in `internal/cmd/cat/formatter.go` (ASCII 0-31 excluding 9,10 + ASCII 127)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - ファイル内容の基本表示 (Priority: P1) 🎯 MVP

**Goal**: ファイルの内容を標準出力に表示。複数ファイルを順番に連結して出力。

**Independent Test**: `mycli cat testfile.txt` を実行してファイル内容が標準出力に表示される。

### Tests for User Story 1 [必須 - TDD原則]

> **重要: これらのテストを実装より先に書き、失敗することを確認してから実装に進むこと（Red-Green-Refactorサイクル）**

- [X] T009 [P] [US1] Create `internal/cmd/cat/formatter_test.go` - TestFormatLine_NoOptions (plain text, no formatting)
- [X] T010 [P] [US1] Create `internal/cmd/cat/processor_test.go` - TestProcessFile_Success (basic file read)
- [X] T011 [P] [US1] Add test TestProcessFile_MultipleFiles in `internal/cmd/cat/processor_test.go`
- [X] T012 [P] [US1] Create `cmd/cat_test.go` - TestCatCommand_BasicFile (Cobra integration test)

### Implementation for User Story 1

- [X] T013 [P] [US1] Implement DefaultFormatter.FormatLine() basic logic in `internal/cmd/cat/formatter.go` (no options, just return line)
- [X] T014 [P] [US1] Implement NewDefaultFormatter() factory in `internal/cmd/cat/formatter.go`
- [X] T015 [US1] Implement DefaultProcessor.ProcessFile() in `internal/cmd/cat/processor.go` with bufio.Scanner (32KB buffer)
- [X] T016 [US1] Implement DefaultProcessor.processReader() helper in `internal/cmd/cat/processor.go`
- [X] T017 [US1] Implement NewDefaultProcessor() factory in `internal/cmd/cat/processor.go`
- [X] T018 [US1] Create `cmd/cat.go` - Cobra command definition with basic RunE function
- [X] T019 [US1] Add catCmd to rootCmd in `cmd/cat.go` init() function
- [X] T020 [US1] Implement multiple file handling loop in `cmd/cat.go` RunE function
- [X] T021 [US1] Add error handling for file errors with stderr output in `cmd/cat.go`
- [X] T022 [US1] Implement exit code 1 when any file has error in `cmd/cat.go`

**Checkpoint**: At this point, User Story 1 should be fully functional - `mycli cat file.txt` works

---

## Phase 4: User Story 2 - 標準入力からの読み込み (Priority: P2)

**Goal**: 標準入力からデータを読み込んで出力。`-` が指定された場合も標準入力から読み込む。

**Independent Test**: `echo "test" | mycli cat` を実行して標準入力からの読み込みが機能する。

### Tests for User Story 2 [必須 - TDD原則]

- [X] T023 [P] [US2] Add TestProcessStdin_Success in `internal/cmd/cat/processor_test.go`
- [X] T024 [P] [US2] Add TestProcessFile_Stdin_DashArgument in `internal/cmd/cat/processor_test.go`
- [X] T025 [P] [US2] Add TestCatCommand_StdinOnly in `cmd/cat_test.go`
- [X] T026 [P] [US2] Add TestCatCommand_MixedStdinAndFiles in `cmd/cat_test.go`

### Implementation for User Story 2

- [X] T027 [US2] Implement DefaultProcessor.ProcessStdin() in `internal/cmd/cat/processor.go`
- [X] T028 [US2] Add dash ("-") detection in ProcessFile() to call ProcessStdin() in `internal/cmd/cat/processor.go`
- [X] T029 [US2] Update `cmd/cat.go` RunE to handle empty args (stdin mode)
- [X] T030 [US2] Add support for `-` argument in file list in `cmd/cat.go`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - 行番号表示オプション (-n) (Priority: P3)

**Goal**: すべての行（空行を含む）の先頭に行番号を付加して表示。

**Independent Test**: `mycli cat -n test.txt` を実行して各行に行番号が付加される。

### Tests for User Story 3 [必須 - TDD原則]

- [X] T031 [P] [US3] Add TestFormatLine_NumberAll in `internal/cmd/cat/formatter_test.go`
- [X] T032 [P] [US3] Add TestFormatLine_NumberAll_EmptyLine in `internal/cmd/cat/formatter_test.go`
- [X] T033 [P] [US3] Add TestFormatLine_NumberAll_Overflow (999,999+ lines) in `internal/cmd/cat/formatter_test.go`
- [X] T034 [P] [US3] Add TestNewOptions_NumberFlag in `internal/cmd/cat/options_test.go`
- [X] T035 [P] [US3] Add TestCatCommand_NumberFlag in `cmd/cat_test.go`

### Implementation for User Story 3

- [X] T036 [US3] Implement line numbering logic in FormatLine() in `internal/cmd/cat/formatter.go` (`%6d  ` format)
- [X] T037 [US3] Implement overflow handling (lineNum % 1000000) in `internal/cmd/cat/formatter.go`
- [X] T038 [US3] Update processReader() to track lineNum and pass to FormatLine() in `internal/cmd/cat/processor.go`
- [X] T039 [US3] Create `internal/cmd/cat/options_test.go` for options tests
- [X] T040 [US3] Implement NewOptions() factory in `internal/cmd/cat/options.go`
- [X] T041 [US3] Add `-n/--number` flag definition in `cmd/cat.go` init()
- [X] T042 [US3] Update RunE to call NewOptions() and pass to processor in `cmd/cat.go`

**Checkpoint**: Line numbering with `-n` flag works independently

---

## Phase 6: User Story 4 - 空行スキップ行番号表示 (-b) (Priority: P4)

**Goal**: 空行以外の行にのみ行番号を付加して表示。

**Independent Test**: `mycli cat -b test.txt` を実行して空行以外に行番号が付加される。

### Tests for User Story 4 [必須 - TDD原則]

- [X] T043 [P] [US4] Add TestFormatLine_NumberNonBlank_EmptyLine in `internal/cmd/cat/formatter_test.go`
- [X] T044 [P] [US4] Add TestFormatLine_NumberNonBlank_NonEmptyLine in `internal/cmd/cat/formatter_test.go`
- [X] T045 [P] [US4] Add TestNewOptions_NumberNonBlankFlag in `internal/cmd/cat/options_test.go`
- [X] T046 [P] [US4] Add TestNewOptions_NumberConflict (both -n and -b) in `internal/cmd/cat/options_test.go`

### Implementation for User Story 4

- [X] T047 [US4] Implement NumberNonBlank logic in FormatLine() - skip numbering if isEmpty in `internal/cmd/cat/formatter.go`
- [X] T048 [US4] Update processReader() to detect isEmpty (len(line) == 0) in `internal/cmd/cat/processor.go`
- [X] T049 [US4] Implement `-n` and `-b` conflict resolution in NewOptions() in `internal/cmd/cat/options.go`
- [X] T050 [US4] Add `-b/--number-nonblank` flag definition in `cmd/cat.go` init()

**Checkpoint**: Line numbering with `-b` flag and `-n`/`-b` conflict resolution work

---

## Phase 7: User Story 5 - 行末表示オプション (-E) (Priority: P4)

**Goal**: 各行の末尾に `$` 記号を表示。

**Independent Test**: `mycli cat -E test.txt` を実行して各行末に `$` が付加される。

### Tests for User Story 5 [必須 - TDD原則]

- [X] T051 [P] [US5] Add TestFormatLine_ShowEnds in `internal/cmd/cat/formatter_test.go`
- [X] T052 [P] [US5] Add TestFormatLine_ShowEnds_EmptyLine in `internal/cmd/cat/formatter_test.go`
- [X] T053 [P] [US5] Add TestNewOptions_ShowEndsFlag in `internal/cmd/cat/options_test.go`

### Implementation for User Story 5

- [X] T054 [US5] Implement ShowEnds logic in FormatLine() - append "$" in `internal/cmd/cat/formatter.go`
- [X] T055 [US5] Add `-E/--show-ends` flag definition in `cmd/cat.go` init()
- [X] T056 [US5] Update NewOptions() to handle ShowEnds flag in `internal/cmd/cat/options.go`

**Checkpoint**: Line end marker with `-E` flag works

---

## Phase 8: User Story 6 - タブ文字表示オプション (-T) (Priority: P5)

**Goal**: すべてのタブ文字を `^I` として表示。

**Independent Test**: `mycli cat -T test.txt` を実行してタブが `^I` として表示される。

### Tests for User Story 6 [必須 - TDD原則]

- [X] T057 [P] [US6] Add TestFormatLine_ShowTabs in `internal/cmd/cat/formatter_test.go`
- [X] T058 [P] [US6] Add TestFormatLine_ShowTabs_MultipleTabs in `internal/cmd/cat/formatter_test.go`
- [X] T059 [P] [US6] Add TestNewOptions_ShowTabsFlag in `internal/cmd/cat/options_test.go`

### Implementation for User Story 6

- [X] T060 [US6] Implement ShowTabs logic in FormatLine() - replace "\t" with "^I" in `internal/cmd/cat/formatter.go`
- [X] T061 [US6] Add `-T/--show-tabs` flag definition in `cmd/cat.go` init()
- [X] T062 [US6] Update NewOptions() to handle ShowTabs flag in `internal/cmd/cat/options.go`

**Checkpoint**: Tab visualization with `-T` flag works

---

## Phase 9: User Story 7 - 非表示文字表示オプション (-v) (Priority: P5)

**Goal**: タブと改行以外の非表示文字を `^` 記法で表示。

**Independent Test**: 制御文字を含むファイルで `mycli cat -v test.txt` を実行して制御文字が可視化される。

### Tests for User Story 7 [必須 - TDD原則]

- [X] T063 [P] [US7] Add TestFormatLine_ShowNonPrinting_ControlChars in `internal/cmd/cat/formatter_test.go`
- [X] T064 [P] [US7] Add TestFormatLine_ShowNonPrinting_DEL in `internal/cmd/cat/formatter_test.go`
- [X] T065 [P] [US7] Add TestFormatLine_ShowNonPrinting_NoControlChars in `internal/cmd/cat/formatter_test.go`
- [X] T066 [P] [US7] Add TestNewOptions_ShowNonPrintingFlag in `internal/cmd/cat/options_test.go`

### Implementation for User Story 7

- [X] T067 [US7] Implement ShowNonPrinting logic in FormatLine() using controlCharMap in `internal/cmd/cat/formatter.go`
- [X] T068 [US7] Update buildControlCharMap() to include all mappings (ASCII 0-31 except 9,10 + ASCII 127) in `internal/cmd/cat/formatter.go`
- [X] T069 [US7] Add `-v/--show-nonprinting` flag definition in `cmd/cat.go` init()
- [X] T070 [US7] Update NewOptions() to handle ShowNonPrinting flag in `internal/cmd/cat/options.go`

**Checkpoint**: Control character visualization with `-v` flag works

---

## Phase 10: User Story 8 - 複数オプション組み合わせ (-A) (Priority: P5)

**Goal**: `-A` オプションは `-v -E -T` の組み合わせと同等。

**Independent Test**: `mycli cat -A test.txt` を実行してすべての非表示文字が可視化される。

### Tests for User Story 8 [必須 - TDD原則]

- [X] T071 [P] [US8] Add TestFormatLine_AllOptions in `internal/cmd/cat/formatter_test.go`
- [X] T072 [P] [US8] Add TestNewOptions_ShowAllFlag in `internal/cmd/cat/options_test.go`
- [X] T073 [P] [US8] Add TestNewOptions_ShowAll_EquivalentToVET in `internal/cmd/cat/options_test.go`

### Implementation for User Story 8

- [X] T074 [US8] Implement `-A` expansion in NewOptions() - set ShowNonPrinting, ShowEnds, ShowTabs to true in `internal/cmd/cat/options.go`
- [X] T075 [US8] Add `-A/--show-all` flag definition in `cmd/cat.go` init()
- [X] T076 [US8] Add integration test for `-A` vs `-vET` equivalence in `cmd/cat_test.go`

**Checkpoint**: All user stories should now be independently functional

---

## Phase 11: Error Handling & Edge Cases

**Purpose**: Comprehensive error handling and edge case coverage

- [X] T077 [P] Add TestProcessFile_NotExist in `internal/cmd/cat/processor_test.go`
- [X] T078 [P] Add TestProcessFile_IsDirectory in `internal/cmd/cat/processor_test.go`
- [ ] T079 [P] Add TestProcessFile_PermissionDenied in `internal/cmd/cat/processor_test.go`
- [X] T080 [P] Add TestCatCommand_PartialError (some files fail) in `cmd/cat_test.go`
- [X] T081 [P] Add TestCatCommand_EmptyFile in `cmd/cat_test.go`
- [X] T082 [P] Add TestProcessFile_BinaryFile in `internal/cmd/cat/processor_test.go`
- [X] T083 Ensure proper error messages format "cat: filename: error" in `cmd/cat.go`
- [X] T084 Verify exit code 1 behavior for any error in `cmd/cat.go`

---

## Phase 12: BATS Integration Tests

**Purpose**: End-to-end validation with actual file I/O

- [X] T085 Create `integration_test/cat.bats` with setup/teardown functions
- [X] T086 [P] Add BATS test "cat displays file content" in `integration_test/cat.bats`
- [X] T087 [P] Add BATS test "cat with -n flag numbers all lines" in `integration_test/cat.bats`
- [X] T088 [P] Add BATS test "cat with -b flag numbers nonempty lines" in `integration_test/cat.bats`
- [X] T089 [P] Add BATS test "cat with -E flag shows line ends" in `integration_test/cat.bats`
- [X] T090 [P] Add BATS test "cat with -T flag shows tabs" in `integration_test/cat.bats`
- [X] T091 [P] Add BATS test "cat with -v flag shows control chars" in `integration_test/cat.bats`
- [X] T092 [P] Add BATS test "cat with -A flag shows all" in `integration_test/cat.bats`
- [X] T093 [P] Add BATS test "cat from stdin" in `integration_test/cat.bats`
- [X] T094 [P] Add BATS test "cat multiple files" in `integration_test/cat.bats`
- [X] T095 [P] Add BATS test "cat nonexistent file error" in `integration_test/cat.bats`
- [X] T096 [P] Add BATS test "cat directory error" in `integration_test/cat.bats`
- [X] T097 [P] Add BATS test "cat partial error continues" in `integration_test/cat.bats`
- [X] T098 Run all BATS tests with `cd integration_test && bats cat.bats`

---

## Phase 13: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T099 Add comprehensive help message and examples to `cmd/cat.go` Long field
- [ ] T100 Add package documentation comments to all exported functions in `internal/cmd/cat/*.go`
- [X] T101 [P] Run `go test -cover ./internal/cmd/cat/` - verify 100% coverage goal
- [X] T102 [P] Run `go test -cover ./cmd/` - verify cat tests pass
- [X] T103 Run `make fmt` - ensure gofmt compliance
- [X] T104 Run `make lint` - ensure golangci-lint (govet) passes
- [X] T105 Run `make build` - ensure binary builds successfully
- [X] T106 Run `make test` - ensure all tests pass
- [X] T107 Run `make all` - verify complete pipeline (test → fmt → lint → build)
- [X] T108 Manual validation: Follow quickstart.md workflow to verify developer experience
- [ ] T109 Performance test: Verify 1MB file processes in <100ms
- [ ] T110 Performance test: Verify 1GB file uses <100MB memory

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-10)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → P4 → P5)
- **Error Handling (Phase 11)**: Can start after US1 (basic file processing)
- **BATS Tests (Phase 12)**: Should wait for most user stories to be complete
- **Polish (Phase 13)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after US1 - Uses same Processor interface
- **User Story 3 (P3)**: Can start after US1 - Adds formatting to basic display
- **User Story 4 (P4)**: Can start after US3 - Extends line numbering
- **User Story 5 (P4)**: Can start after US1 - Independent formatting feature
- **User Story 6 (P5)**: Can start after US1 - Independent formatting feature
- **User Story 7 (P5)**: Can start after US1 - Independent formatting feature
- **User Story 8 (P5)**: Depends on US5, US6, US7 completion

### Within Each User Story (TDD Cycle)

1. **Red**: Write tests first - they MUST fail
2. **Green**: Implement minimal code to pass tests
3. **Refactor**: Improve code while keeping tests green

### Parallel Opportunities

- **Setup**: All tasks marked [P] can run in parallel
- **Foundational**: T005, T006, T007, T008 can run in parallel
- **User Story Tests**: All test tasks within a story marked [P] can run in parallel
- **User Story Models/Components**: Tasks marked [P] can run in parallel
- **Different User Stories**: Can be worked on in parallel by different team members after Foundation
- **BATS Tests**: All BATS test tasks (T086-T097) can be created in parallel
- **Final Checks**: T101, T102 can run in parallel

---

## Parallel Example: User Story 3 (Line Numbering)

```bash
# Red Phase - Launch all tests together:
Task T031: "TestFormatLine_NumberAll"
Task T032: "TestFormatLine_NumberAll_EmptyLine"
Task T033: "TestFormatLine_NumberAll_Overflow"
Task T034: "TestNewOptions_NumberFlag"
Task T035: "TestCatCommand_NumberFlag"

# Green Phase - Implement in dependency order:
Task T036: "Implement line numbering in FormatLine()"
Task T037: "Implement overflow handling"
Task T038: "Update processReader() to track lineNum"
Task T039: "Create options_test.go"
Task T040: "Implement NewOptions()"
Task T041: "Add -n flag"
Task T042: "Update RunE to use NewOptions()"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T008) - CRITICAL
3. Complete Phase 3: User Story 1 (T009-T022)
4. **STOP and VALIDATE**: Test `mycli cat file.txt` works
5. Deploy/demo if ready - this is a working cat command!

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy (MVP - basic cat!)
3. Add User Story 2 → Test independently → Deploy (stdin support!)
4. Add User Story 3 → Test independently → Deploy (line numbers!)
5. Continue for US4-US8 as needed
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers after Foundational phase completes:

- Developer A: User Story 1 (P1) - Basic file display
- Developer B: User Story 2 (P2) - Stdin support
- Developer C: User Story 3 (P3) - Line numbering
- All stories integrate independently

---

## Task Completion Summary

**Total Tasks**: 110
- Phase 1 (Setup): 3 tasks
- Phase 2 (Foundational): 5 tasks (1 completed ✅)
- Phase 3 (US1): 14 tasks
- Phase 4 (US2): 8 tasks
- Phase 5 (US3): 12 tasks
- Phase 6 (US4): 8 tasks
- Phase 7 (US5): 6 tasks
- Phase 8 (US6): 6 tasks
- Phase 9 (US7): 8 tasks
- Phase 10 (US8): 6 tasks
- Phase 11 (Error Handling): 8 tasks
- Phase 12 (BATS): 14 tasks
- Phase 13 (Polish): 12 tasks

**Parallel Opportunities**: 52 tasks marked [P] can be parallelized
**Independent Test Criteria**: Each user story has clear acceptance tests
**Suggested MVP Scope**: Phase 1-3 (US1 only) = 22 tasks for working basic cat command

---

## Notes

- ✅ Task T004 already completed: `cmd/cat_wrapper_test.go` created
- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **TDD is mandatory**: Write tests first (Red), implement (Green), refactor
- Verify tests fail before implementing
- Commit after each logical task group
- Stop at any checkpoint to validate story independently
- Run `make all` frequently to catch issues early
