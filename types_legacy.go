package acp

import "context"

// TerminalId is kept as an exported compatibility alias for consumers pinned to
// fork builds that exposed a typed terminal identifier.
type TerminalId = string

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
