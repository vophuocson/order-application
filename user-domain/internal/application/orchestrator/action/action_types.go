package action

// Action name constants
const (
	USER_UPDATE                 = "UserUpdate"
	USER_UPDATE_COMPENSATE      = "UserUpdateCompensate"
	USER_UPDATE_APPROVE         = "UserUpdateApprove"
	PAYMENT_UPDATE_EXECUTE      = "PaymentUpdateExecute"
	PAYMENT_UPDATE_COMPENSATE   = "PaymentUpdateCompensate"
	PAYMENT_UPDATE_VERIFICATION = "PaymentUpdateVerification"
	PAYMENT_UPDATE_APPROVE      = "PaymentUpdateApprove"
)

// VerificationResponse represents the response from a service verification
type VerificationResponse struct {
	ServiceName string
	Accepted    bool
	Message     string
	Error       error
}
