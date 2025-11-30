# Tasks: Configure設定構造のリファクタリング

**Input**: Design documents from `/specs/005-configure-refactor/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: このプロジェクトではTDD（テスト駆動開発）が必須です。すべての実装タスクに対して、必ずテストタスクを先に定義してください。テストは実装前に書き、失敗することを確認してから実装に進む必要があります（憲章原則I参照）。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Summary

- **Total Tasks**: 15
- **Parallel Opportunities**: 6 tasks can run in parallel
- **Estimated Time**: ~50 minutes
- **MVP Scope**: User Story 1 + User Story 2 (Core functionality)

## Implementation Strategy

1. **MVP First**: Implement User Stories 1-2 (nested structure + Config struct) for minimal viable configuration
2. **Incremental Delivery**: Add User Story 3 (defaults) and verify backward compatibility (User Story 4)
3. **Independent Testing**: Each user story has clear test criteria and can be verified independently

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Project initialization and environment verification

- [ ] T001 Verify Go 1.25.4 installation and PATH configuration with golangci-lint v2.6.2
- [ ] T002 Confirm working branch `005-configure-refactor` is checked out
- [ ] T003 Run `make test` to establish baseline (all existing tests must pass)

**Checkpoint**: Environment ready, baseline established

---

## Phase 2: User Story 1 - ネストされた設定構造のサポート (Priority: P1) 🎯 MVP Core

**Goal**: Enable generation of hierarchical YAML configuration files with nested sections (common, hoge.foo)

**Independent Test**: Execute `mycli configure --force` and verify generated YAML has nested structure with 7 fields

### Tests for User Story 1 [必須 - TDD原則]

> **重要: これらのテストを実装より先に書き、失敗することを確認してから実装に進むこと（Red-Green-Refactorサイクル）**

- [X] T004 [P] [US1] Write integration test TestConfigUnmarshal_NewStructure in cmd/root_test.go - verify Viper unmarshals complete nested YAML to Config struct
- [X] T005 [P] [US1] Write integration test TestConfigUnmarshal_BackwardCompatibility in cmd/root_test.go - verify old config files load correctly (missing fields → zero values)

### Implementation for User Story 1

- [X] T006 [US1] Define CommonConfig struct with Var1 (string) and Var2 (int) fields in cmd/root.go
- [X] T007 [US1] Define HogeConfig struct with Fuga (string) and Foo (FooConfig) fields in cmd/root.go
- [X] T008 [US1] Define FooConfig struct with Bar (string) field in cmd/root.go
- [X] T009 [US1] Add Common (CommonConfig) and Hoge (HogeConfig) fields to existing Config struct in cmd/root.go
- [X] T010 [US1] Verify all mapstructure tags use kebab-case format (client-id, common, var1, var2, hoge, fuga, foo, bar)
- [X] T011 [US1] Run integration tests (T004, T005) to verify Config struct unmarshaling works correctly

**Checkpoint**: Config struct can unmarshal nested YAML, backward compatibility maintained

---

## Phase 3: User Story 2 - Config構造体のフィールド定義とmapstructureタグの更新 (Priority: P2) 🎯 MVP Foundation

**Goal**: Ensure Config struct fields are correctly defined with proper mapstructure tags for Viper integration

**Independent Test**: Create Config instance, unmarshal from YAML, verify all nested field values

**Note**: This story is partially completed in Phase 2 (struct definitions). This phase focuses on validation and integration.

### Tests for User Story 2 [必須 - TDD原則]

> **重要: これらのテストを実装より先に書き、失敗することを確認してから実装に進むこと（Red-Green-Refactorサイクル）**

- [X] T012 [P] [US2] Write integration test TestConfigUnmarshal_PartialStructure in cmd/root_test.go - verify partial YAML (only common section) loads correctly

**Checkpoint**: Config struct fully validates with Viper, all edge cases tested

---

## Phase 4: User Story 3 - デフォルト値の設定 (Priority: P3)

**Goal**: Generate configuration files with appropriate default values for all nested fields

**Independent Test**: Read generated config file, verify each field has correct initial value

### Tests for User Story 3 [必須 - TDD原則]

> **重要: これらのテストを実装より先に書き、失敗することを確認してから実装に進むこと（Red-Green-Refactorサイクル）**

- [X] T013 [P] [US3] Write unit test TestBuildEffectiveConfig_HasAllFields in cmd/viperutils_test.go - verify map contains all keys (client-id, client-secret, common, hoge)
- [X] T014 [P] [US3] Write unit test TestBuildEffectiveConfig_CorrectDefaultValues in cmd/viperutils_test.go - verify default values ("", "", "", 123, "hello", "hello")
- [X] T015 [P] [US3] Write unit test TestBuildEffectiveConfig_YAMLMarshal in cmd/viperutils_test.go - verify map marshals to valid YAML and unmarshal preserves structure

### Implementation for User Story 3

- [X] T016 [US3] Update BuildEffectiveConfig() function in cmd/viperutils.go to return nested map with default values
- [X] T017 [US3] Run unit tests (T013-T015) to verify BuildEffectiveConfig() returns correct structure

**Checkpoint**: Generated config files have proper default values, all fields initialized

---

## Phase 5: User Story 4 - 既存機能との互換性維持 (Priority: P4)

**Goal**: Ensure all existing configure command flags and features continue to work unchanged

**Independent Test**: Run existing test suite, manually test each flag combination

### Tests for User Story 4 [必須 - TDD原則]

- [X] T018 [US4] Run existing test suite with `make test` - all tests must pass (SC-004)
- [X] T019 [US4] Write integration test TestConfigureCommand_GeneratesNewStructure in cmd/configure_test.go - execute configure --force, verify YAML structure

### Manual Verification for User Story 4

- [X] T020 [US4] Manual test: `./bin/mycli configure --force` - verify file created at ~/.config/mycli/default.yaml
- [X] T021 [US4] Manual test: `cat ~/.config/mycli/default.yaml` - verify 10-line YAML structure matches FR-005
- [X] T022 [US4] Manual test: `./bin/mycli configure --profile dev --force` - verify dev.yaml has same structure as default.yaml
- [X] T023 [US4] Manual test: Set EDITOR=cat and run `./bin/mycli configure --force --edit` - verify editor integration works

**Checkpoint**: All existing functionality preserved, no regressions detected

---

## Phase 6: Quality Assurance & Polish

**Purpose**: Final validation, code quality checks, and deliverable preparation

- [X] T024 Run `make fmt` - format all Go code with gofmt
- [X] T025 Run `export PATH="$(go env GOPATH)/bin:$PATH" && make lint` - verify zero warnings/errors (SC-005)
- [X] T026 Run `make build` - verify binary builds successfully
- [X] T027 Run `make all` - complete test → fmt → lint → build pipeline
- [X] T028 Verify Success Criteria SC-001 through SC-006 per spec.md
- [X] T029 Commit changes with message: "feat: expand Config struct to support nested configuration"

**Final Checkpoint**: All quality gates passed, feature ready for PR

---

## Dependency Graph

```
Phase 1 (Setup)
    ↓
Phase 2 (US1: Nested Structure) ← MVP Core
    ├─ T004, T005 (tests) [parallel]
    ├─ T006-T010 (implementation)
    └─ T011 (verification)
    ↓
Phase 3 (US2: Struct Validation) ← MVP Foundation
    └─ T012 (test)
    ↓
Phase 4 (US3: Default Values)
    ├─ T013, T014, T015 (tests) [parallel]
    ├─ T016, T017 (implementation)
    ↓
Phase 5 (US4: Compatibility)
    ├─ T018, T019 (tests)
    └─ T020-T023 (manual verification)
    ↓
Phase 6 (Quality Assurance)
    └─ T024-T029 (final checks)
```

## Parallel Execution Opportunities

### Within Phase 2 (User Story 1)
- **T004 + T005**: Both write different test cases in cmd/root_test.go

### Within Phase 4 (User Story 3)
- **T013 + T014 + T015**: All write different test cases in cmd/viperutils_test.go

**Total Parallelizable Tasks**: 5 out of 29 tasks (17%)

## Parallel Execution Example for User Story 1

```bash
# Developer A: Write TestConfigUnmarshal_NewStructure (T004)
# Developer B: Write TestConfigUnmarshal_BackwardCompatibility (T005)
# Both can work simultaneously on different test functions

# After both complete, proceed to implementation tasks (T006-T010) sequentially
```

## MVP Scope Definition

**Minimum Viable Product includes**:
- ✅ User Story 1: Nested structure support (Phase 2)
- ✅ User Story 2: Config struct validation (Phase 3)

**Reason**: These two stories provide the core functionality - ability to define and use nested configuration structures. Stories 3-4 are enhancements and verification but not required for basic operation.

**MVP Validation**: After completing Phase 3, you should be able to:
1. Define nested Config structs with mapstructure tags
2. Unmarshal YAML files into nested Config instances
3. Verify all fields populate correctly

## Task Count by User Story

| User Story | Priority | Task Count | Percentage |
|------------|----------|------------|------------|
| Setup | - | 3 | 10% |
| US1: Nested Structure | P1 | 8 | 28% |
| US2: Struct Validation | P2 | 1 | 3% |
| US3: Default Values | P3 | 5 | 17% |
| US4: Compatibility | P4 | 6 | 21% |
| Quality Assurance | - | 6 | 21% |
| **Total** | | **29** | **100%** |

## Success Criteria Verification Checklist

After completing all tasks, verify these criteria from spec.md:

- [ ] **SC-001**: Generated `~/.config/mycli/default.yaml` has 10 lines (2 top-level + 3 common + 4 hoge + 1 blank)
- [ ] **SC-002**: YAML parses correctly with gopkg.in/yaml.v3, all 7 fields present with expected values
- [ ] **SC-003**: After `viper.Unmarshal(&CliConfig)`, verify `CliConfig.Common.Var2 == 123` and `CliConfig.Hoge.Foo.Bar == "hello"`
- [ ] **SC-004**: `make test` exits with code 0, 100% pass rate
- [ ] **SC-005**: `make lint` exits with code 0, zero warnings/errors
- [ ] **SC-006**: Profile config `~/.config/mycli/prod.yaml` has identical structure to default.yaml

## Time Estimates

| Phase | Estimated Time | Cumulative |
|-------|----------------|------------|
| Phase 1: Setup | 3 min | 3 min |
| Phase 2: US1 | 15 min | 18 min |
| Phase 3: US2 | 5 min | 23 min |
| Phase 4: US3 | 12 min | 35 min |
| Phase 5: US4 | 10 min | 45 min |
| Phase 6: QA | 5 min | 50 min |
| **Total** | **50 min** | |

**Note**: Time estimates assume developer is familiar with Go, Viper, and the existing codebase. First-time implementation may take 1.5-2x longer.

## References

- **Spec**: [spec.md](./spec.md) - Feature requirements and user stories
- **Plan**: [plan.md](./plan.md) - Technical approach and architecture
- **Data Model**: [data-model.md](./data-model.md) - Entity definitions and relationships
- **Quickstart**: [quickstart.md](./quickstart.md) - Detailed implementation steps
- **Contracts**: [contracts/](./contracts/) - Interface specifications
- **Constitution**: `.specify/memory/constitution.md` - Quality standards and principles
