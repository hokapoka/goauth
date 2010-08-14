package oauth

import (
	"fmt"
)

type PersistantMognoDB struct{

}

func ( p *PersistantMognoDB) Save(at *AccessToken){
	fmt.Println("Saving to mongoDB")
}

func (p *PersistantMognoDB) Get( userRef string) *AccessToken {
	// Search monogDB for Users Access Token
	return nil
}


