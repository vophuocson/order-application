package action

// VerificationResponse represents the response from a service verification
type VerificationResponse struct {
	ServiceName string
	Accepted    bool
	Message     string
	Error       error
}
