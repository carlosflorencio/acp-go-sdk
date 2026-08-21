# Fork changes vs upstream

This repository is a fork of [`coder/acp-go-sdk`](https://github.com/coder/acp-go-sdk),
maintained for [Kandev](https://github.com/kdlbs/kandev). It tracks upstream
`v0.13.5` and adds the changes below. Everything here exists because Kandev needs
behavior that upstream does not provide; keep this file current when the fork's
delta changes.

## Why the fork exists

Kandev talks to several ACP agents (Claude Code, Gemini, Cursor, and others) that
do not all track the latest ACP schema. Two of them, plus one Cursor-specific
protocol quirk, forced changes the upstream SDK rejects by design:

1. Cursor emits a non-standard `cursor/task` client request that upstream refuses
   because it is not underscore-prefixed.
2. Older agents still send pre-`v0.13.5` payloads (top-level `models`, the legacy
   `session/set_model` method) that upstream dropped.
3. Kandev needs a bounded inbound notification queue with a fail-fast overflow
   signal, which upstream does not expose.

## Changes

### 1. Vendor client-request extension seam (the primary reason for the fork)

Upstream only routes underscore-prefixed (`_`) methods to
`ExtensionMethodHandler`. Cursor sends `cursor/task`, so upstream returns
`method not found` and Cursor's subagent turn ends with no feedback.

The fork widens the seam **for inbound client requests only**:
`ClientSideConnection.handleWithExtensions` now falls through to the client's
`HandleExtensionMethod` for any method that is not a known stable method, as long
as it is a real inbound request. This is scoped deliberately:

- Only the **client** side is widened; the agent side still requires `_`.
- Only **requests** are delegated; vendor **notifications** are not routed to the
  extension handler.
- **Outbound** `CallExtension` / `NotifyExtension` still reject non-underscore
  names via `validateExtensionMethodName`.

Supporting pieces:

- `ClientSideConnection.isKnownMethod(method)` is generated into `client_gen.go`
  by `cmd/generate` (see `cmd/generate/internal/emit/dispatch.go`). It is
  reproducible by `make version` and must not be hand-edited.
- `isInboundRequest(ctx)` gates the fallback to genuine inbound requests.

Files: `extensions.go`, `connection.go`, `client_gen.go`,
`cmd/generate/internal/emit/dispatch.go`. Tests:
`TestExtensionMethods_*` in `acp_test.go`.

### 2. Legacy agent compatibility shims

Hand-written in `types_legacy.go` (with `types_legacy_test.go`):

- `LegacyModels` / `LegacyModelInfo`: read-only parsing of the pre-`v0.13.5`
  top-level `models` payload.
- `LegacyAgentMethodSessionSetModel` (`"session/set_model"`),
  `UnstableSetSessionModelRequest` / `UnstableSetSessionModelResponse`, and
  `ClientSideConnection.UnstableSetSessionModel(...)`: legacy model-selection
  wire method for unmigrated agents.
- Compatibility type aliases: `AuthMethodId`, `MessageId`, `TerminalId`,
  `DeleteSessionRequest`, `DeleteSessionResponse`.

### 3. Bounded inbound notification queue

In `connection.go`:

- `WithMaxQueuedNotifications(n)`: `ConnectionOption` to cap the inbound
  notification queue on both `NewClientSideConnection` and
  `NewAgentSideConnection`.
- `ErrNotificationQueueOverflow`: exported sentinel returned when the queue
  overflows, so the connection fails fast instead of growing unbounded.

Tests: `TestConnectionFailsFastOnNotificationQueueOverflow*` in `acp_test.go`.

## Maintenance hazard: `LegacyModels` fields in `types_gen.go`

`types_gen.go` is a **generated** file, but it carries three manually injected
fields that the generator does not know about:

- `NewSessionResponse.LegacyModels *LegacyModels`
- `LoadSessionResponse.LegacyModels *LegacyModels`
- `UnstableForkSessionResponse.LegacyModels *LegacyModels`

Running `make version` (or `go run ./cmd/generate`) **silently deletes these
three fields**, which breaks legacy `models` parsing. This is confirmed: a clean
regenerate drops all three and leaves no other functional drift.

**When you regenerate `types_gen.go`, re-add the three `LegacyModels` fields by
hand** (or port them into the generator). The `LegacyModels` *type* lives in the
hand-written `types_legacy.go` and is safe; only the *fields* are at risk.

## Upgrading from upstream

1. Rebase the fork branch onto the new upstream tag.
2. Run `make version` to regenerate bindings.
3. Re-apply the `LegacyModels` fields in `types_gen.go` (see hazard above).
4. Re-run `go test ./...` and confirm the `TestExtensionMethods_*`,
   `TestLegacy*`, and `TestConnectionFailsFastOnNotificationQueueOverflow*`
   tests still pass.
5. Re-point Kandev's `apps/backend/go.mod` `replace` at the new commit.

## Related open work

- `fix/large-jsonrpc-lines` (PR #1): raises the 10 MiB per-line Scanner cap to
  64 MiB with a chunked reader, and adds `ErrPeerDisconnected`. Not in this
  branch.
- `feat/acp-schema-1.20-sdk-0.14` (PR #2): schema 1.20 regen and independent
  `v0.14` versioning, including PR #1's fixes. Not in this branch.

Neither of those branches carries the vendor extension seam, so a future merge
must preserve the changes documented above.
