package oauth

import (
	"fmt"
)

type RequestToken struct{
	Token string
	Secret string
	Verifier string
}

type AccessToken struct{
	Token string
	Secret string
	UserRef string
	Verifier string
}




