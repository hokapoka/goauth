// Copyright (c) 2010 The GOAuth Authors. All rights reserved.
//
// Please note. 
// 
// The aim of this example is to show you how the GOAuth 
// Exported fields & methods are to be used.
//
// This _example_ will only work with a single twitter account.
// In order to use it with more than one account you will need 
// to store the AccessToken's that are associated with your 
// respective user/visitor accounts.  
//
// Additionally, you will need to replace the ConsumerKey, 
// ConsumerSecret & CallBackURL with you relevants values.
// 
// If you have any issues please feel free to contact : 
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
	"github.com/hoisie/web.go"
//	"http"
	"io/ioutil"
)

var goauthcon *oauth.OAuthConsumer
var AT *oauth.AccessToken
var RT *oauth.RequestToken

func main(){

	goauthcon = &oauth.OAuthConsumer{
		Service:"twitter",
		RequestTokenURL:"http://twitter.com/oauth/request_token",
		AccessTokenURL:"http://twitter.com/oauth/access_token",
		AuthorizationURL:"http://twitter.com/oauth/authorize",
		ConsumerKey:"change me",
		ConsumerSecret:"change me",
		CallBackURL:"oob",  // Twitter require the string "oob" be passed for Out Of Band Mode

	}

	web.Get("/signin/twitter(.*)", twitterSignIn)
	web.Post("/signin/twitter(.*)", twitterSignInVerify)
	web.Get("/callback/twitter(.*)", twitterCallback)

	web.Get("/twitter/hometimeline(.*)", twitterHomeTimeLine)
	web.Get("/twitter/updatestatus(.*)", twitterUpdateStatus)
	web.Get("/twitter/credentials(.*)", twitterVerifyCredentials)
	web.Get("/(.*)", noRespond)

	web.Run("0.0.0.0:7177")

}

func twitterSignIn(ctx *web.Context, name string) {
	ctx.WriteString("<h1>Twitter Out Of Band Verification</h1>")
	s, rt, err := goauthcon.GetRequestAuthorizationURL()

	if err != nil {
		ctx.WriteString(err.String())
	}

	if rt != nil {
		RT = rt
	}

	ctx.WriteString("<p>Visit this URL :<a href=\"" + s + "\" target=\"_new\">" + s + "</a> (which will open a new tab) and Allow access. Then, return to this tab and enter the PIN number into the box below.</p>")

	ctx.WriteString("<form method=\"POST\"><p><label>PIN</label><input type=\"text\" name=\"verifier\"/></p><p><input type=\"submit\" value=\"Verify\"/></form>")

}

func twitterSignInVerify(ctx *web.Context, name string) {

	ctx.WriteString("<h1>Twitter Out Of Band Verification</h1>")

	v := getParam(ctx, "verifier")
	if v == "" {
		ctx.WriteString("<p style=\"color:red;\">Please Enter the PIN Number!</p>")
		return
	}

	at := goauthcon.GetAccessToken(RT.Token, v)

	// Store at off to persistant data store for use later.
	AT = at

	ctx.WriteString("<h1>Access Token Received</h1>")
	defer func() { footer(ctx) }()

}
func twitterCallback(ctx *web.Context, name string) {
	if getParam(ctx, "denied") != "" {
		ctx.WriteString("<h1>OAuth Access Denied</h1>")
		return
	}

	oauth_token := getParam(ctx, "oauth_token")
	oauth_verifier := getParam(ctx, "oauth_verifier")
	at := goauthcon.GetAccessToken(oauth_token, oauth_verifier)

	// Store at off to persistant data store for use later.
	AT = at

	ctx.WriteString("<h1>Access Token Received</h1>")
	defer func() { footer(ctx) }()
}

func twitterVerifyCredentials(ctx *web.Context, name string){

	ctx.WriteString("<h1>Twitter Credentials</h1>")
	if AT == nil {
		ctx.WriteString("<p>Please <a href=\"/signin/twitter\">Sign in to Twitter</a></p>")
		return
	}
	defer func() { footer(ctx) }()

	ctx.WriteString("<p>Build &amp; Send request</p>")

	r, err := goauthcon.Get(
		"http://api.twitter.com/1/account/verify_credentials.json",
		nil,
		AT )

	if err != nil {
		ctx.WriteString("<p style=\"color:red\">Error : " + err.String() + "</p>")
		return
	}

	b, _ := ioutil.ReadAll( r.Body ) 

	ctx.WriteString("<h2>Twitter Response</h2>")
	ctx.WriteString("<textarea rows=\"20\" cols=\"60\">")
	ctx.Write(b)
	ctx.WriteString("</textarea>")



}

func twitterHomeTimeLine(ctx *web.Context, name string){

	ctx.WriteString("<h1>Twitter Home Time Line</h1>")
	if AT == nil {
		ctx.WriteString("<p>Please <a href=\"/signin/twitter\">Sign in to Twitter</a></p>")
		return
	}

	defer func() { footer(ctx) }()
	ctx.WriteString("<p>Build &amp; Send request</p>")

	r, err := goauthcon.Get(
		"http://api.twitter.com/1/statuses/home_timeline.json",
		nil,
		AT )

	if err != nil {
		ctx.WriteString("<p style=\"color:red\">Error : " + err.String() + "</p>")
		return
	}

	b, _ := ioutil.ReadAll( r.Body ) 

	ctx.WriteString("<h2>Twitter Response</h2>")
	ctx.WriteString("<textarea rows=\"20\" cols=\"60\">")
	ctx.Write(b)
	ctx.WriteString("</textarea>")


}

func twitterUpdateStatus(ctx *web.Context, name string){

	ctx.WriteString("<h1>Twiiter Status Update</h1>")
	if AT == nil {
		ctx.WriteString("<p>Please <a href=\"/signin/twitter\">Sign in to Twitter</a></p>")
		return
	}

	defer func() { footer(ctx) }()
	ctx.WriteString("<p>Build &amp; Send request</p>")

	r, err := goauthcon.Post(
		"http://api.twitter.com/1/statuses/update.json",
		oauth.Params{
			&oauth.Pair{Key:"status", Value:"Testing Status Update via GOAuth - OAuth consumer for #Golang"},
		},
		AT )

	if err != nil {
		ctx.WriteString("<p style=\"color:red\">Error : " + err.String() + "</p>")
		return
	}

	b, _ := ioutil.ReadAll( r.Body ) 

	ctx.WriteString("<h2>Twitter Response</h2>")
	ctx.WriteString("<textarea rows=\"5\" cols=\"30\">")
	ctx.Write(b)
	ctx.WriteString("</textarea>")


}


func noRespond(ctx *web.Context, name string) {
	ctx.WriteString("<h1>Testing OAuth With GoLang</h1>")
}

func footer(ctx *web.Context){
	ctx.WriteString("<p>Click Here <a href=\"/twitter/updatestatus\">to update your twitter status</a></p>")
	ctx.WriteString("<p>Click here <a href=\"/twitter/hometimeline\">to view home timeline</a></p>")
	ctx.WriteString("<p>Click here <a href=\"/twitter/credentials\">to verify your credentials</a></p>")
}


func getParam(ctx *web.Context, param string) (v string){

	c, ok := ctx.Request.Params[param]

	if !ok { return }
	
	v = c

	return
}




