// internal/dtos/server_message.go

package dtos

type ServerMessage struct {
	Status                 string   `json:"status,omitempty"`
	Possibilities          []string `json:"possibilities,omitempty"`
	UserFlowName           string   `json:"user_flow_name,omitempty"`
	UserFlowID             uint     `json:"user_flow_id,omitempty"`
	Result                 string   `json:"result,omitempty"`
	QueryClearLevel        int      `json:"query_clear_level,omitempty"`
	QueryClearAnalysis     string   `json:"query_clear_analysis,omitempty"`
	UserFLowSwitchLevel    int      `json:"user_flow_switch_level,omitempty"`
	UserFlowAnalysis       string   `json:"user_flow_analysis,omitempty"`
	FormulatedSearxngQuery string   `json:"formulated_searxng_query,omitempty"`
	SelectedStyle          string   `json:"selected_style,omitempty"`
}
