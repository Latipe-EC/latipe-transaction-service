package dto

type AuthorizationHeader struct {
	BearerToken string `reqHeader:"Authorization"`
}
