package domain

type TokenRequest struct {
	GrantType           string `json:"grant_type" validate:"required"`
	ClientAssertion     string `json:"client_assertion" validate:"required"`
	ClientAssertionType string `json:"client_assertion_type" validate:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type ParticipantInfo struct {
	PartyID     string `json:"party_id"`
	PartyName   string `json:"party_name"`
	Status      string `json:"status"`
	Certificate string `json:"certificate"`
}

type CertificateValidationResult struct {
	Valid       bool   `json:"valid"`
	ClientID    string `json:"client_id"`
	Certificate string `json:"certificate"`
	Error       string `json:"error,omitempty"`
}
