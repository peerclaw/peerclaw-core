package contacts

import "time"

// ContactRequestStatus represents the status of a contact request.
type ContactRequestStatus string

const (
	StatusPending  ContactRequestStatus = "pending"
	StatusApproved ContactRequestStatus = "approved"
	StatusRejected ContactRequestStatus = "rejected"
)

// ContactRequest represents a request from one agent to add another as a contact.
type ContactRequest struct {
	ID           string               `json:"id"`
	FromAgentID  string               `json:"from_agent_id"`
	ToAgentID    string               `json:"to_agent_id"`
	Message      string               `json:"message,omitempty"`
	Status       ContactRequestStatus `json:"status"`
	RejectReason string               `json:"reject_reason,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}
