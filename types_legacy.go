package acp

import "context"

// LegacyModels is the pre-v0.13.5 top-level "models" payload that some agents
// still emit on session/new, session/load, and session/fork responses.
type LegacyModels struct {
	AvailableModels []LegacyModelInfo `json:"availableModels"`
	CurrentModelId  string            `json:"currentModelId"`
}

// LegacyModelInfo describes one model entry in the legacy models payload.
type LegacyModelInfo struct {
	ModelId     string         `json:"modelId"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Meta        map[string]any `json:"_meta,omitempty"`
}

// AuthMethodId is kept as an exported compatibility alias for consumers pinned
// to fork builds that exposed a typed authentication method identifier.
type AuthMethodId = string

// TerminalId is kept as an exported compatibility alias for consumers pinned to
// fork builds that exposed a typed terminal identifier.
type TerminalId = string

// DeleteSessionRequest is kept as a compatibility alias to the current
// unstable delete-session request type.
type DeleteSessionRequest = UnstableDeleteSessionRequest

// DeleteSessionResponse is kept as a compatibility alias to the current
// unstable delete-session response type.
type DeleteSessionResponse = UnstableDeleteSessionResponse

// LegacyAgentMethodSessionSetModel is the JSON-RPC method name for the
// pre-v0.13.5 unstable session model API (session/set_model).
const LegacyAgentMethodSessionSetModel = "session/set_model"

// UnstableSetSessionModelRequest is the pre-v0.13.5 request for
// session/set_model. It is kept here as a client-side compatibility shim for
// consumers that still need to talk to agents on the older surface.
type UnstableSetSessionModelRequest struct {
	SessionId SessionId      `json:"sessionId"`
	ModelId   string         `json:"modelId"`
	Meta      map[string]any `json:"_meta,omitempty"`
}

// UnstableSetSessionModelResponse is the pre-v0.13.5 empty response for
// session/set_model.
type UnstableSetSessionModelResponse struct {
	Meta map[string]any `json:"_meta,omitempty"`
}

// UnstableSetSessionModel calls the legacy session/set_model JSON-RPC method
// on the connected agent.
func (c *ClientSideConnection) UnstableSetSessionModel(
	ctx context.Context,
	params UnstableSetSessionModelRequest,
) (UnstableSetSessionModelResponse, error) {
	resp, err := SendRequest[UnstableSetSessionModelResponse](c.conn, ctx, LegacyAgentMethodSessionSetModel, params)
	return resp, err
}
