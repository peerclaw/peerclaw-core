package agentcard

// Reputation score thresholds used across PeerClaw modules.
const (
	ReputationLow    = 0.3 // Below this score an agent is considered untrusted.
	ReputationMedium = 0.7 // Minimum score for moderate trust.
	ReputationHigh   = 0.8 // Score at or above this level is considered Trusted.
)
