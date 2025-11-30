# Implementation Tasks: Echo サブコマンド実装

**Feature**: Echo サブコマンド実装  
**Branch**: `001-echo-subcommand`  
**Created**: 2025-11-30  
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)

---

## Implementation Strategy

### MVP First Approach
**MVP Scope**: User Story 1 (P1) - 基本的なテキスト出力のみ  
**Rationale**: P1は他の全機能の基盤となり、独立してテスト・デプロイ可能な最小価値提供

### Incremental Delivery
1. **Sprint 1**: MVP (P1) - 基本出力機能
2. **Sprint 2**: P2 (-n flag) - 改行抑制オプション
3. **Sprint 3**: P3 (-e flag) - エスケープシーケンス解釈
4. **Sprint 4**: P4 - オプション組み合わせ + Polish

各スプリント完了時に完全に動作するバージョンを提供（TDD + 品質ゲート通過）

---

## Phase 1: Setup (プロジェクト初期化)

**Goal**: 実装に必要なディレクトリ構造とテスト基盤を準備

**Blockers**: なし（既存プロジェクト構造に追加）

### Tasks

- [X] T001 internal/echo/ディレクトリを作成
- [X] T002 [P] cmd/root.goを確認し、echoコマンド登録の準備を理解

**Validation**: `internal/echo/`ディレクトリが存在し、`cmd/root.go`のコマンド登録パターンを理解済み

---

## Phase 2: Foundational (基盤整備)

**Goal**: 全User Storyで共通使用されるテスト基盤とデータ構造の定義

**Blockers**: Phase 1完了

### Tasks

- [X] T003 internal/echo/echo.goにEchoOptions構造体を定義（SuppressNewline, InterpretEscapes, Verbose, Args）
- [X] T004 [P] internal/echo/echo_test.goにEchoOptionsの基本テストを作成（構造体の初期化）
- [X] T005 [P] cmd/echo_test.goにテストヘルパー関数を作成（bytes.Bufferでstdout/stderrキャプチャ）

**Validation**: `EchoOptions`構造体が定義され、テストヘルパー関数が動作確認済み

---

## Phase 3: User Story 1 - 基本的なテキスト出力 (P1)

**Goal**: 引数をスペース区切りで標準出力に表示、デフォルトで改行を追加

**Independent Test**: `mycli echo "Hello, World!"` → `Hello, World!\n`

**Acceptance Criteria**:
- ✅ 単一引数の出力（改行付き）
- ✅ 複数引数のスペース区切り出力
- ✅ 引数なしで空行出力
- ✅ 特殊文字の正しい出力

### Tasks - Tests (Red Phase)

- [X] T006 [US1] cmd/echo_test.goに単一引数テスト追加（"Hello" → "Hello\n"）
- [X] T007 [P] [US1] cmd/echo_test.goに複数引数テスト追加（"A" "B" "C" → "A B C\n"）
- [X] T008 [P] [US1] cmd/echo_test.goに引数なしテスト追加（→ "\n"）
- [X] T009 [P] [US1] cmd/echo_test.goに特殊文字テスト追加（"!@#$%" → "!@#$%\n"）

### Tasks - Implementation (Green Phase)

- [X] T010 [US1] cmd/echo.goにCobraコマンド定義を作成（Use, Short, Long, Example）
- [X] T011 [US1] cmd/echo.goにRunE関数を実装（引数をstrings.Joinでスペース区切り結合）
- [X] T012 [US1] cmd/echo.goで出力にデフォルト改行を追加（fmt.Fprintln使用）
- [X] T013 [US1] cmd/echo.goのinit()でrootCmd.AddCommand(echoCmd)を追加

### Tasks - Refactor & Validation

- [X] T014 [US1] make allを実行してテスト・フォーマット・リント・ビルドを検証
- [X] T015 [P] [US1] bin/mycli echo "test"を手動実行し、出力を確認

**Phase Validation**: 
- ✅ 全US1テストがパス（4/4）
- ✅ `make all`成功
- ✅ SC-001達成（100ms以内の実行）

---

## Phase 4: User Story 2 - 改行抑制オプション (-n) (P2)

**Goal**: `-n`フラグで末尾の改行を抑制

**Independent Test**: `mycli echo -n "Prompt: "` → `Prompt: `（改行なし）

**Acceptance Criteria**:
- ✅ `-n`フラグで改行なし出力
- ✅ 複数引数 + `-n`フラグ
- ✅ 引数なし + `-n`フラグ（何も出力しない）

### Tasks - Tests (Red Phase)

- [X] T016 [US2] cmd/echo_test.goに-nフラグテスト追加（"-n", "Hello" → "Hello"）
- [X] T017 [P] [US2] cmd/echo_test.goに-n複数引数テスト追加（"-n", "A", "B" → "A B"）
- [X] T018 [P] [US2] cmd/echo_test.goに-n引数なしテスト追加（"-n" → ""）

### Tasks - Implementation (Green Phase)

- [X] T019 [US2] cmd/echo.goにBoolフラグ"-n"/"--no-newline"を定義（Flags().BoolP）
- [X] T020 [US2] internal/echo/echo.goにGenerateOutput関数を作成（EchoOptions受け取り、文字列返却）
- [X] T021 [US2] internal/echo/echo_test.goにGenerateOutputのユニットテスト作成
- [X] T022 [US2] cmd/echo.goのRunEでGenerateOutputを呼び出し、条件付き改行制御を実装

### Tasks - Refactor & Validation

- [X] T023 [US2] make allを実行してテスト・フォーマット・リント・ビルドを検証
- [X] T024 [P] [US2] bin/mycli echo -n "test"を手動実行し、改行なし出力を確認

**Phase Validation**:
- ✅ 全US2テストがパス（3/3）
- ✅ US1のリグレッションなし
- ✅ `make all`成功

---

## Phase 5: User Story 3 - エスケープシーケンス解釈 (-e) (P3)

**Goal**: `-e`フラグでバックスラッシュエスケープシーケンスを解釈

**Independent Test**: `mycli echo -e "Line1\nLine2\tTabbed"` → 改行とタブが正しく解釈

**Acceptance Criteria**:
- ✅ 9種類のエスケープシーケンス（`\n`, `\t`, `\\`, `\"`, `\a`, `\b`, `\c`, `\r`, `\v`）すべて動作
- ✅ `-e`未指定時はリテラル文字列として出力
- ✅ 無効なエスケープシーケンス（`\z`）はリテラル扱い

### Tasks - Tests (Red Phase)

- [X] T025 [US3] internal/echo/processor_test.goにProcessEscapes関数のテスト作成（\n → 改行）
- [X] T026 [P] [US3] internal/echo/processor_test.goに\tテスト追加（\t → タブ）
- [X] T027 [P] [US3] internal/echo/processor_test.goに\\テスト追加（\\ → \）
- [X] T028 [P] [US3] internal/echo/processor_test.goに\"テスト追加（\" → "）
- [X] T029 [P] [US3] internal/echo/processor_test.goに\aテスト追加（\a → ベル）
- [X] T030 [P] [US3] internal/echo/processor_test.goに\bテスト追加（\b → バックスペース）
- [X] T031 [P] [US3] internal/echo/processor_test.goに\cテスト追加（\c → 出力抑制、suppressNewline=true）
- [X] T032 [P] [US3] internal/echo/processor_test.goに\rテスト追加（\r → キャリッジリターン）
- [X] T033 [P] [US3] internal/echo/processor_test.goに\vテスト追加（\v → 垂直タブ）
- [X] T034 [P] [US3] internal/echo/processor_test.goに無効エスケープテスト追加（\z → \z リテラル）
- [X] T035 [US3] cmd/echo_test.goに-eフラグテスト追加（"-e", "Hello\nWorld" → "Hello\nWorld\n"）
- [X] T036 [P] [US3] cmd/echo_test.goに-e未指定テスト追加（"Hello\nWorld" → "Hello\\nWorld\n" リテラル）

### Tasks - Implementation (Green Phase)

- [X] T037 [US3] internal/echo/processor.goにProcessEscapes関数を実装（strings.Builder使用）
- [X] T038 [US3] internal/echo/processor.goで\nエスケープを実装（builder.WriteRune('\n')）
- [X] T039 [P] [US3] internal/echo/processor.goで\tエスケープを実装
- [X] T040 [P] [US3] internal/echo/processor.goで\\エスケープを実装
- [X] T041 [P] [US3] internal/echo/processor.goで\"エスケープを実装
- [X] T042 [P] [US3] internal/echo/processor.goで\aエスケープを実装
- [X] T043 [P] [US3] internal/echo/processor.goで\bエスケープを実装
- [X] T044 [US3] internal/echo/processor.goで\cエスケープを実装（即座にreturn、suppressNewline=true）
- [X] T045 [P] [US3] internal/echo/processor.goで\rエスケープを実装
- [X] T046 [P] [US3] internal/echo/processor.goで\vエスケープを実装
- [X] T047 [P] [US3] internal/echo/processor.goで無効エスケープ処理を実装（default: リテラル出力）
- [X] T048 [US3] cmd/echo.goにBoolフラグ"-e"/"--escape"を定義
- [X] T049 [US3] internal/echo/echo.goのGenerateOutputで-eフラグ時にProcessEscapesを呼び出し

### Tasks - Refactor & Validation

- [X] T050 [US3] make allを実行してテスト・フォーマット・リント・ビルドを検証
- [X] T051 [P] [US3] bin/mycli echo -e "Line1\nLine2\tTab"を手動実行し、エスケープ解釈を確認
- [X] T052 [P] [US3] bin/mycli echo "Line1\nLine2"を手動実行し、リテラル出力を確認

**Phase Validation**:
- ✅ 全US3テストがパス（28/28：12テスト + 16実装検証）
- ✅ SC-005達成（全エスケープシーケンス正しく解釈）
- ✅ US1/US2のリグレッションなし
- ✅ `make all`成功

---

## Phase 6: User Story 4 - オプション組み合わせ (P4)

**Goal**: `-n`と`-e`を同時に指定可能

**Independent Test**: `mycli echo -n -e "Hello\tWorld"` → タブ解釈、改行抑制

**Acceptance Criteria**:
- ✅ `-n -e`組み合わせで動作
- ✅ `-e -n`（順序逆）でも同じ動作
- ✅ エスケープ解釈後の改行のみ抑制（エスケープ内の\nは解釈）

### Tasks - Tests (Red Phase)

- [X] T053 [US4] cmd/echo_test.goに-n -e組み合わせテスト追加（"-n", "-e", "Tab\there" → "Tab\there"）
- [X] T054 [P] [US4] cmd/echo_test.goに-e -n組み合わせテスト追加（"-e", "-n", "Line\nNo" → "Line\nNo" 最後の改行なし）

### Tasks - Implementation (Green Phase)

- [X] T055 [US4] internal/echo/echo.goのGenerateOutputで両フラグの組み合わせロジックを実装（-e処理後、-n判定）
- [X] T056 [US4] cmd/echo_test.goで全組み合わせパターンを検証（-n, -e, -n -e, -e -n）

### Tasks - Refactor & Validation

- [X] T057 [US4] make allを実行してテスト・フォーマット・リント・ビルドを検証

**Phase Validation**:
- ✅ 全US4テストがパス（2/2）
- ✅ 全User Story（P1-P4）のリグレッションなし
- ✅ `make all`成功

---

## Phase 7: Polish & Cross-Cutting Concerns

**Goal**: エラーハンドリング、ヘルプメッセージ、パフォーマンス検証、憲章準拠の最終確認

**Blockers**: Phase 3-6完了（全User Story実装済み）

### Tasks - Error Handling & Help

- [X] T058 cmd/echo.goでSilenceUsage: falseを明示的に設定（Cobraデフォルトだが明示）
- [X] T059 cmd/echo.goのExample フィールドに2-3個の使用例を追加（FR-008準拠）
- [X] T060 [P] cmd/echo_test.goに無効フラグテスト追加（"-x" → エラー + ヘルプ表示、exit code 1）
- [X] T061 [P] cmd/echo_test.goに終了コードテスト追加（正常時0、エラー時1）

### Tasks - Verbose Logging

- [X] T062 cmd/echo.goにBoolフラグ"--verbose"を定義
- [X] T063 cmd/echo.goのRunEでverboseフラグ時にlog.SetOutput(cmd.ErrOrStderr())でデバッグ情報出力
- [X] T064 [P] cmd/echo_test.goにverboseフラグテスト追加（"--verbose" → stderrにデバッグ情報）

### Tasks - Performance & Edge Cases

- [X] T065 cmd/echo_test.goにパフォーマンステスト追加（起動時メモリ50MB以下 + 10,000引数100MB以下検証、runtime.MemStats使用、SC-004準拠）
- [X] T066 [P] cmd/echo_test.goにヘルプ表示パフォーマンステスト追加（"--help" → 50ms以内、SC-003）
- [X] T067 [P] cmd/echo_test.goに空文字列引数テスト追加（"", "test" → " test\n"）
- [X] T068 [P] cmd/echo_test.goに--引数区切りテスト追加（"-n", "--", "-e" → "-e\n" リテラル）

### Tasks - Documentation & Constitution Compliance

- [X] T069 cmd/echo.goにパッケージドキュメントコメントを追加
- [X] T070 [P] internal/echo/processor.goに関数ドキュメントコメントを追加
- [X] T071 [P] internal/echo/echo.goに関数ドキュメントコメントを追加
- [X] T072 README.mdにechoサブコマンドのクイックスタート例を追加

### Tasks - Final Validation

- [X] T073 make allを実行し、全品質ゲート（test → fmt → lint → build）をパス
- [X] T074 [P] bin/mycli echo --helpを実行し、ヘルプメッセージ品質を確認（FR-008）
- [X] T075 [P] 憲章チェックリスト（plan.md §Constitution Check）を再確認
- [X] T076 [P] UNIX標準echoとの互換性を手動テスト（SC-002検証）
- [X] T076a [P] cmd/echo_test.goにUTF-8テストケース追加（日本語文字列「こんにちは世界」、Emoji「🚀✨」の出力検証、FR-014準拠）

**Phase Validation**:
- ✅ 全エラーハンドリングテストがパス
- ✅ SC-003, SC-004達成（パフォーマンス目標）
- ✅ SC-006達成（TDDアプローチ、全テストパス）
- ✅ SC-007, SC-008達成（エラー/verboseフラグ）
- ✅ 憲章5原則すべて準拠確認済み

---

## Task Dependencies

### User Story Completion Order

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational) ← Must complete before any User Story
    ↓
Phase 3 (US1: P1) ← MVP, blocks nothing
    ↓
Phase 4 (US2: P2) ← Depends on US1 (基本出力機能)
    ↓
Phase 5 (US3: P3) ← Depends on US1 (基本出力機能)
    ↓
Phase 6 (US4: P4) ← Depends on US2 + US3 (両オプション実装)
    ↓
Phase 7 (Polish) ← Depends on US1-US4 (全機能実装)
```

### Critical Path

**Longest dependency chain**: Setup → Foundational → US1 → US2 → US4 → Polish  
**Estimated Duration**: 
- Setup: 0.5h
- Foundational: 1h
- US1: 3h (MVP)
- US2: 2h
- US3: 4h (並行可能)
- US4: 1.5h
- Polish: 2h
- **Total**: ~14h (シーケンシャル実行時)

### Parallel Execution Opportunities

#### Within US3 (Phase 5)
```bash
# 9種類のエスケープシーケンステスト（T025-T034）は並列実行可能
# 各エスケープの実装（T038-T047）も並列実行可能

# Example: 3つのエスケープを同時に実装
Terminal 1: T038 \n実装 → T039 \t実装
Terminal 2: T040 \\実装 → T041 \"実装
Terminal 3: T042 \a実装 → T043 \b実装
# → 所要時間を4hから2hに短縮可能
```

#### Within Polish (Phase 7)
```bash
# ドキュメント追加タスク（T069-T072）は並列実行可能
Terminal 1: T069 cmd/echo.goドキュメント
Terminal 2: T070 processor.goドキュメント
Terminal 3: T071 echo.goドキュメント
Terminal 4: T072 README.md更新
# → 所要時間を30分から10分に短縮可能
```

---

## Testing Strategy

### TDD Cycle per User Story

**Red Phase** (テスト先行):
- 各User Storyの最初のタスクで失敗するテストを作成
- テストが正しく失敗することを確認（`go test ./... -v`）

**Green Phase** (最小実装):
- テストをパスする最小限の実装を追加
- 各実装タスク完了後に`go test ./...`でパスを確認

**Refactor Phase** (リファクタリング):
- コードの可読性・保守性を向上
- テストが引き続きパスすることを確認
- `make all`で品質ゲートをパス

### Test Coverage Requirements

- ✅ **Unit Tests**: `internal/echo/processor_test.go`, `internal/echo/echo_test.go`
- ✅ **Integration Tests**: `cmd/echo_test.go`（Cobraコマンド統合）
- ✅ **Performance Tests**: 10,000引数、ヘルプ表示速度
- ✅ **Edge Case Tests**: 空文字列、無効エスケープ、--引数区切り

### Quality Gates

**Each Phase must pass**:
1. `go test ./...` - すべてのテストがパス
2. `gofmt -s -w .` - コードフォーマット
3. `golangci-lint run --enable=govet` - 静的解析
4. `go build -o bin/mycli` - ビルド成功

**Shortcut**: `make all` - 上記すべてを順次実行

---

## Progress Tracking

### Completed Tasks: 0/77

- Phase 1 (Setup): 0/2
- Phase 2 (Foundational): 0/3
- Phase 3 (US1 - P1): 0/10
- Phase 4 (US2 - P2): 0/9
- Phase 5 (US3 - P3): 0/28
- Phase 6 (US4 - P4): 0/5
- Phase 7 (Polish): 0/20

### Parallel Opportunities: 42 tasks marked with [P]

**Estimated Time Savings**: ~6h (シーケンシャル14h → 並列実行8h)

---

## Implementation Notes

### File Creation Order

1. **Tests First** (TDD必須): `*_test.go` → 実装ファイル
2. **Internal First**: `internal/echo/` → `cmd/echo.go`（依存方向）
3. **Incremental**: User Story単位で完全に実装・テスト完了

### Code Style

- Go標準の命名規約（PascalCase for exported, camelCase for unexported）
- `gofmt -s`で自動整形
- パッケージドキュメントコメント必須

### Commit Strategy

- User Story単位でコミット（Phase 3, 4, 5, 6完了時）
- コミットメッセージ例:
  - `feat(echo): implement basic text output (US1/P1)`
  - `feat(echo): add -n flag for newline suppression (US2/P2)`
  - `feat(echo): add -e flag for escape sequences (US3/P3)`
  - `feat(echo): support combined -n -e flags (US4/P4)`
  - `docs(echo): add help messages and README examples`

---

## Success Criteria Validation

### After Phase 7 Completion

- [ ] SC-001: 100ms以内の実行完了（パフォーマンステストT065で検証済み）
- [ ] SC-002: UNIX標準echoとの互換性（手動テストT076で検証済み）
- [ ] SC-003: 50ms以内のヘルプ表示（パフォーマンステストT066で検証済み）
- [ ] SC-004: 10,000引数で100MB以下（パフォーマンステストT065で検証済み）
- [ ] SC-005: 全エスケープシーケンス正しく解釈（ユニットテストT025-T034で検証済み）
- [ ] SC-006: TDDアプローチ、全テストパス（`make test`成功）
- [ ] SC-007: 無効オプション時のエラー処理（統合テストT060で検証済み）
- [ ] SC-008: verboseフラグでデバッグ情報出力（統合テストT064で検証済み）

### Constitution Compliance

- [ ] TDD必須: 全タスクでテストファースト実施済み
- [ ] パッケージ責務分離: `cmd/`（CLI）と`internal/echo/`（ロジック）明確に分離
- [ ] コード品質基準: `make all`成功（test + fmt + lint + build）
- [ ] 設定管理の一貫性: N/A（ステートレスコマンド）
- [ ] ユーザーエクスペリエンス: ヘルプメッセージ品質確認済み（T074）
- [ ] パフォーマンス要件: SC-001, SC-003, SC-004達成

---

## Next Steps

1. **Phase 1開始**: T001-T002を実行してSetup完了
2. **MVP実装**: Phase 2-3を完了してUser Story 1（P1）を動作可能に
3. **Incremental Delivery**: Phase 4-6で残りのUser Storyを順次実装
4. **Polish**: Phase 7で最終品質確認とドキュメント整備
5. **PR作成**: 全タスク完了後、`001-echo-subcommand` → `main`へのPR作成

**推定完了時間**: 8-14時間（並列実行度による）
