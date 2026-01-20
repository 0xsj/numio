package rpc

import (
	"encoding/json"
	"testing"
)

// ════════════════════════════════════════════════════════════════
// REQUEST TESTS
// ════════════════════════════════════════════════════════════════

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		params  any
		wantErr bool
	}{
		{
			name:   "simple request without params",
			method: MethodInitialize,
			params: nil,
		},
		{
			name:   "request with struct params",
			method: MethodKeypress,
			params: KeypressParams{Event: KeyEvent{Key: "a"}},
		},
		{
			name:   "request with map params",
			method: MethodSetOption,
			params: map[string]any{"option": "precision", "value": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest(tt.method, tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if req.JSONRPC != JSONRPCVersion {
				t.Errorf("JSONRPC = %q, want %q", req.JSONRPC, JSONRPCVersion)
			}

			if req.Method != tt.method {
				t.Errorf("Method = %q, want %q", req.Method, tt.method)
			}

			if req.ID == 0 {
				t.Error("ID should not be 0")
			}

			if tt.params == nil && req.Params != nil {
				t.Errorf("Params = %v, want nil", req.Params)
			}

			if tt.params != nil && req.Params == nil {
				t.Error("Params should not be nil")
			}
		})
	}
}

func TestNewRequestWithID(t *testing.T) {
	req, err := NewRequestWithID(42, MethodInitialize, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.ID != 42 {
		t.Errorf("ID = %d, want 42", req.ID)
	}
}

func TestRequest_ParseParams(t *testing.T) {
	params := KeypressParams{
		Event: KeyEvent{
			Key:       "Enter",
			Modifiers: []string{"ctrl", "shift"},
		},
	}

	req, err := NewRequest(MethodKeypress, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed KeypressParams
	if err := req.ParseParams(&parsed); err != nil {
		t.Fatalf("ParseParams error: %v", err)
	}

	if parsed.Event.Key != "Enter" {
		t.Errorf("Key = %q, want %q", parsed.Event.Key, "Enter")
	}

	if len(parsed.Event.Modifiers) != 2 {
		t.Errorf("Modifiers len = %d, want 2", len(parsed.Event.Modifiers))
	}
}

func TestRequest_ParseParams_NilParams(t *testing.T) {
	req, _ := NewRequest(MethodShutdown, nil)

	var parsed map[string]any
	if err := req.ParseParams(&parsed); err != nil {
		t.Errorf("ParseParams with nil should not error: %v", err)
	}
}

func TestRequest_IDsAreUnique(t *testing.T) {
	seen := make(map[int64]bool)

	for i := 0; i < 1000; i++ {
		req, _ := NewRequest(MethodInitialize, nil)
		if seen[req.ID] {
			t.Errorf("duplicate ID: %d", req.ID)
		}
		seen[req.ID] = true
	}
}

// ════════════════════════════════════════════════════════════════
// RESPONSE TESTS
// ════════════════════════════════════════════════════════════════

func TestNewResponse(t *testing.T) {
	tests := []struct {
		name   string
		id     int64
		result any
	}{
		{
			name:   "response with nil result",
			id:     1,
			result: nil,
		},
		{
			name:   "response with struct result",
			id:     2,
			result: InitializeResult{CoreVersion: "1.0.0", ProtocolVersion: "2.0"},
		},
		{
			name:   "response with primitive result",
			id:     3,
			result: "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := NewResponse(tt.id, tt.result)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.JSONRPC != JSONRPCVersion {
				t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, JSONRPCVersion)
			}

			if resp.ID != tt.id {
				t.Errorf("ID = %d, want %d", resp.ID, tt.id)
			}

			if resp.Error != nil {
				t.Error("Error should be nil for success response")
			}

			if resp.IsError() {
				t.Error("IsError() should return false")
			}
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	rpcErr := NewError(ErrCodeInvalidParams, "missing field")
	resp := NewErrorResponse(42, rpcErr)

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, JSONRPCVersion)
	}

	if resp.ID != 42 {
		t.Errorf("ID = %d, want 42", resp.ID)
	}

	if resp.Result != nil {
		t.Error("Result should be nil for error response")
	}

	if !resp.IsError() {
		t.Error("IsError() should return true")
	}

	if resp.Error.Code != ErrCodeInvalidParams {
		t.Errorf("Error.Code = %d, want %d", resp.Error.Code, ErrCodeInvalidParams)
	}
}

func TestResponse_ParseResult(t *testing.T) {
	original := InitializeResult{
		CoreVersion:     "1.0.0",
		ProtocolVersion: "2.0",
		Capabilities: Capabilities{
			SupportsCurrency: true,
			SupportsCrypto:   true,
		},
	}

	resp, _ := NewResponse(1, original)

	var parsed InitializeResult
	if err := resp.ParseResult(&parsed); err != nil {
		t.Fatalf("ParseResult error: %v", err)
	}

	if parsed.CoreVersion != "1.0.0" {
		t.Errorf("CoreVersion = %q, want %q", parsed.CoreVersion, "1.0.0")
	}

	if !parsed.Capabilities.SupportsCurrency {
		t.Error("Capabilities.SupportsCurrency should be true")
	}
}

// ════════════════════════════════════════════════════════════════
// NOTIFICATION TESTS
// ════════════════════════════════════════════════════════════════

func TestNewNotification(t *testing.T) {
	tests := []struct {
		name   string
		method string
		params any
	}{
		{
			name:   "notification without params",
			method: MethodRender,
			params: nil,
		},
		{
			name:   "notification with params",
			method: MethodRateStatus,
			params: RateStatusParams{Status: RateStatus{State: RateStateFetching}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notif, err := NewNotification(tt.method, tt.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if notif.JSONRPC != JSONRPCVersion {
				t.Errorf("JSONRPC = %q, want %q", notif.JSONRPC, JSONRPCVersion)
			}

			if notif.Method != tt.method {
				t.Errorf("Method = %q, want %q", notif.Method, tt.method)
			}
		})
	}
}

func TestNotification_ParseParams(t *testing.T) {
	params := RateStatusParams{
		Status: RateStatus{
			State:   RateStateSuccess,
			Message: "Rates updated",
		},
	}

	notif, _ := NewNotification(MethodRateStatus, params)

	var parsed RateStatusParams
	if err := notif.ParseParams(&parsed); err != nil {
		t.Fatalf("ParseParams error: %v", err)
	}

	if parsed.Status.State != RateStateSuccess {
		t.Errorf("Status.State = %q, want %q", parsed.Status.State, RateStateSuccess)
	}
}

// ════════════════════════════════════════════════════════════════
// ERROR TESTS
// ════════════════════════════════════════════════════════════════

func TestNewError(t *testing.T) {
	err := NewError(ErrCodeMethodNotFound, "unknown method")

	if err.Code != ErrCodeMethodNotFound {
		t.Errorf("Code = %d, want %d", err.Code, ErrCodeMethodNotFound)
	}

	if err.Message != "unknown method" {
		t.Errorf("Message = %q, want %q", err.Message, "unknown method")
	}

	if err.Data != nil {
		t.Error("Data should be nil")
	}
}

func TestNewErrorWithData(t *testing.T) {
	err := NewErrorWithData(ErrCodeInvalidParams, "validation failed", map[string]string{
		"field": "viewport",
		"issue": "width must be positive",
	})

	if err.Code != ErrCodeInvalidParams {
		t.Errorf("Code = %d, want %d", err.Code, ErrCodeInvalidParams)
	}

	if err.Data == nil {
		t.Fatal("Data should not be nil")
	}

	data, ok := err.Data.(map[string]string)
	if !ok {
		t.Fatal("Data should be map[string]string")
	}

	if data["field"] != "viewport" {
		t.Errorf("Data[field] = %q, want %q", data["field"], "viewport")
	}
}

func TestError_ErrorInterface(t *testing.T) {
	err := NewError(ErrCodeInternal, "something went wrong")

	// Should implement error interface
	var _ error = err

	if err.Error() != "something went wrong" {
		t.Errorf("Error() = %q, want %q", err.Error(), "something went wrong")
	}
}

func TestStandardErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		wantCode int
	}{
		{"ErrParse", ErrParse("invalid json"), ErrCodeParse},
		{"ErrInvalidRequest", ErrInvalidRequest("missing jsonrpc"), ErrCodeInvalidRequest},
		{"ErrMethodNotFound", ErrMethodNotFound("foo"), ErrCodeMethodNotFound},
		{"ErrInvalidParams", ErrInvalidParams("bad type"), ErrCodeInvalidParams},
		{"ErrInternal", ErrInternal("panic"), ErrCodeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %d, want %d", tt.err.Code, tt.wantCode)
			}

			if tt.err.Data == nil {
				t.Error("Data should contain detail")
			}
		})
	}
}

// ════════════════════════════════════════════════════════════════
// MESSAGE DETECTION TESTS
// ════════════════════════════════════════════════════════════════

func TestDetectMessageType(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantType MessageType
		wantErr  bool
	}{
		{
			name:     "request",
			json:     `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
			wantType: MessageRequest,
		},
		{
			name:     "response with result",
			json:     `{"jsonrpc":"2.0","id":1,"result":{"version":"1.0"}}`,
			wantType: MessageResponse,
		},
		{
			name:     "response with error",
			json:     `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid"}}`,
			wantType: MessageResponse,
		},
		{
			name:     "notification",
			json:     `{"jsonrpc":"2.0","method":"render","params":{}}`,
			wantType: MessageNotification,
		},
		{
			name:     "invalid json",
			json:     `{not valid json`,
			wantType: MessageUnknown,
			wantErr:  true,
		},
		{
			name:     "empty object",
			json:     `{}`,
			wantType: MessageUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgType, err := DetectMessageType([]byte(tt.json))

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if msgType != tt.wantType {
				t.Errorf("MessageType = %v, want %v", msgType, tt.wantType)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	t.Run("parse request", func(t *testing.T) {
		data := []byte(`{"jsonrpc":"2.0","id":42,"method":"keypress","params":{"event":{"key":"a"}}}`)

		msg, msgType, err := ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if msgType != MessageRequest {
			t.Errorf("MessageType = %v, want %v", msgType, MessageRequest)
		}

		req, ok := msg.(*Request)
		if !ok {
			t.Fatal("message should be *Request")
		}

		if req.ID != 42 {
			t.Errorf("ID = %d, want 42", req.ID)
		}

		if req.Method != "keypress" {
			t.Errorf("Method = %q, want %q", req.Method, "keypress")
		}
	})

	t.Run("parse response", func(t *testing.T) {
		data := []byte(`{"jsonrpc":"2.0","id":42,"result":"ok"}`)

		msg, msgType, err := ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if msgType != MessageResponse {
			t.Errorf("MessageType = %v, want %v", msgType, MessageResponse)
		}

		resp, ok := msg.(*Response)
		if !ok {
			t.Fatal("message should be *Response")
		}

		if resp.ID != 42 {
			t.Errorf("ID = %d, want 42", resp.ID)
		}
	})

	t.Run("parse notification", func(t *testing.T) {
		data := []byte(`{"jsonrpc":"2.0","method":"render","params":{"lines":[]}}`)

		msg, msgType, err := ParseMessage(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if msgType != MessageNotification {
			t.Errorf("MessageType = %v, want %v", msgType, MessageNotification)
		}

		notif, ok := msg.(*Notification)
		if !ok {
			t.Fatal("message should be *Notification")
		}

		if notif.Method != "render" {
			t.Errorf("Method = %q, want %q", notif.Method, "render")
		}
	})
}

// ════════════════════════════════════════════════════════════════
// JSON SERIALIZATION TESTS
// ════════════════════════════════════════════════════════════════

func TestRequest_JSONRoundtrip(t *testing.T) {
	original, _ := NewRequestWithID(99, MethodKeypress, KeypressParams{
		Event: KeyEvent{Key: "Escape", Modifiers: []string{"ctrl"}},
	})

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded Request
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.JSONRPC != original.JSONRPC {
		t.Errorf("JSONRPC = %q, want %q", decoded.JSONRPC, original.JSONRPC)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}

	if decoded.Method != original.Method {
		t.Errorf("Method = %q, want %q", decoded.Method, original.Method)
	}

	var params KeypressParams
	if err := decoded.ParseParams(&params); err != nil {
		t.Fatalf("ParseParams error: %v", err)
	}

	if params.Event.Key != "Escape" {
		t.Errorf("Event.Key = %q, want %q", params.Event.Key, "Escape")
	}
}

func TestResponse_JSONRoundtrip(t *testing.T) {
	original, _ := NewResponse(123, InitializeResult{
		CoreVersion:     "1.2.3",
		ProtocolVersion: JSONRPCVersion,
		Capabilities: Capabilities{
			SupportsCurrency: true,
			SupportsExplain:  true,
		},
	})

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}

	if decoded.IsError() {
		t.Error("should not be error response")
	}

	var result InitializeResult
	if err := decoded.ParseResult(&result); err != nil {
		t.Fatalf("ParseResult error: %v", err)
	}

	if result.CoreVersion != "1.2.3" {
		t.Errorf("CoreVersion = %q, want %q", result.CoreVersion, "1.2.3")
	}
}

func TestErrorResponse_JSONRoundtrip(t *testing.T) {
	original := NewErrorResponse(456, NewErrorWithData(
		ErrCodeEvaluation,
		"division by zero",
		map[string]int{"line": 5, "col": 10},
	))

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !decoded.IsError() {
		t.Error("should be error response")
	}

	if decoded.Error.Code != ErrCodeEvaluation {
		t.Errorf("Error.Code = %d, want %d", decoded.Error.Code, ErrCodeEvaluation)
	}

	if decoded.Error.Message != "division by zero" {
		t.Errorf("Error.Message = %q, want %q", decoded.Error.Message, "division by zero")
	}
}

func TestNotification_JSONRoundtrip(t *testing.T) {
	original, _ := NewNotification(MethodRateStatus, RateStatusParams{
		Status: RateStatus{
			State:     RateStateSuccess,
			Message:   "Updated 150 rates",
			UpdatedAt: 1700000000,
		},
	})

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded Notification
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.Method != MethodRateStatus {
		t.Errorf("Method = %q, want %q", decoded.Method, MethodRateStatus)
	}

	var params RateStatusParams
	if err := decoded.ParseParams(&params); err != nil {
		t.Fatalf("ParseParams error: %v", err)
	}

	if params.Status.State != RateStateSuccess {
		t.Errorf("Status.State = %q, want %q", params.Status.State, RateStateSuccess)
	}

	if params.Status.UpdatedAt != 1700000000 {
		t.Errorf("Status.UpdatedAt = %d, want %d", params.Status.UpdatedAt, 1700000000)
	}
}

// ════════════════════════════════════════════════════════════════
// JSON FORMAT COMPLIANCE TESTS
// ════════════════════════════════════════════════════════════════

func TestRequest_OmitsNullParams(t *testing.T) {
	req, _ := NewRequest(MethodShutdown, nil)

	data, _ := json.Marshal(req)
	jsonStr := string(data)

	// params should be omitted when nil, not "params":null
	if contains(jsonStr, `"params":null`) {
		t.Error("nil params should be omitted, not serialized as null")
	}
}

func TestResponse_OmitsNullFields(t *testing.T) {
	// Success response should omit error
	success, _ := NewResponse(1, "ok")
	successData, _ := json.Marshal(success)

	if contains(string(successData), `"error"`) {
		t.Error("success response should omit error field")
	}

	// Error response should omit result
	errResp := NewErrorResponse(2, NewError(ErrCodeInternal, "fail"))
	errData, _ := json.Marshal(errResp)

	if contains(string(errData), `"result"`) {
		t.Error("error response should omit result field")
	}
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
