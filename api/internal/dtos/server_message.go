// dtos/server_message.go

package dtos

type ServerMessage struct {
	Status  string `json:"status,omitempty"`
	Result  string `json:"result,omitempty"`
	Details string `json:"details,omitempty"`
}
