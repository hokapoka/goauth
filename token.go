package oauth

import (
	"fmt"
)

type RequestToken struct{
	Token string
	Secret string
	Verifier string
}

func (t *RequestToken)  test(){
	fmt.Println("hi")
}

type AccessToken struct{
	Token string
	Secret string
	UserRef string
	Verifier string
}




