
// Copyright (c) 2010 The GOAuth Authors. All rights reserved.
//
//     email - hoka@hokapoka.com
//       web - http://go.hokapoka.com
//      buzz - hokapoka.com@gmail.com 
//   twitter - @hokapokadotcom
//    github - github.com/hokapoka/goauth
//   
package main

import (
	oauth "github.com/hokapoka/goauth"
	"fmt"
)

var goauthcon *oauth.OAuthConsumer
var AT *oauth.AccessToken

func main(){

	goauthcon = &oauth.OAuthConsumer{
		Service:"twitter",
		RequestTokenURL:"http://twitter.com/oauth/request_token",
		AccessTokenURL:"http://twitter.com/oauth/access_token",
		AuthorizationURL:"http://twitter.com/oauth/authorize",
		ConsumerKey:"OdFyxuGBBcBx4edyWGvsQ",
		ConsumerSecret:"tyhEcdpaJoKUNsQju0PYjKAxOAQUNwAnjEjOb3tRYTs",
		CallBackURL:"oob",

	}

	s, rt, err := goauthcon.GetRequestAuthorizationURL()
	if err != nil {
		fmt.Println(err.String())
		return
	}
	var pin string

	fmt.Printf("Open %s In your browser.\n Allow access and then enter the PIN number\n", s);
	fmt.Printf("PIN Number: ")
	fmt.Scanln(&pin)

	at := goauthcon.GetAccessToken(rt.Token, pin)

	_, err = goauthcon.Post(
		"http://api.twitter.com/1/statuses/update.json",
		oauth.Params{
			&oauth.Pair{Key:"status", Value:"Testing Status Update via GOAuth - OAuth consumer for #Golang"},
		},
		at )

	if err != nil {
		fmt.Println(err.String())
		return
	}

	fmt.Println("Twitter Status is updated")
	return

}

