package domain

type TokenRequest struct {
	GrantType           string
	ClientAssertion     string
	ClientAssertionType string
}

type TokenResponse struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
}

type ParticipantInfo struct {
	PartyID     string
	PartyName   string
	Status      string
	Certificate string
}

type CertificateValidationResult struct {
	Valid       bool
	ClientID    string
	Certificate string
	Error       string
}
