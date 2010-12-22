package main

import (
	oauth "github.com/hokapoka/goauth"
	"github.com/hoisie/web.go"
	"http"
	"io/ioutil"
	"gobuzz" // github.com/hokapoka/gobuzz
	"json"
	"time"
)

var goauthcon *oauth.OAuthConsumer
var googleConn *oauth.OAuthConsumer
var diggConn *oauth.OAuthConsumer
var AT *oauth.AccessToken
var diggAT *oauth.AccessToken
var googleAT *oauth.AccessToken

//51.8872, -1.7575 - bourton

func main(){

	googleConn = &oauth.OAuthConsumer{
		Service:"google",
		RequestTokenURL:"https://www.google.com/accounts/OAuthGetRequestToken",
		AccessTokenURL:"https://www.google.com/accounts/OAuthGetAccessToken",
		AuthorizationURL:"https://www.google.com/accounts/OAuthAuthorizeToken",
		ConsumerKey:"change me",
		ConsumerSecret:"change me",
		CallBackURL:"http://oauth.hokapoka.com/callback/google",
		AdditionalParams:oauth.Params{
			&oauth.Pair{ Key:"scope", Value:"https://www.googleapis.com/auth/buzz"},
		},
	}

	web.Get("/signin/google(.*)", googleSignin)
    web.Get("/callback/google(.*)", googleRespond)
    web.Get("/google/getpublic(.*)", googleGetPublic)
    web.Get("/google/search(.*)", googleGetSearch)
    web.Post("/google/search(.*)", googleSearch)

    web.Get("/(.*)", noRespond)
	web.Run("0.0.0.0:7177")

}

func googleSignin(ctx *web.Context, name string) {
	s, err := googleConn.GetRequestAuthorizationURL()
	if err != nil {
		ctx.WriteString(err.String())
	}
	ctx.Redirect(http.StatusFound, s)
}



func googleRespond(ctx *web.Context, name string) {
	if GetParam(ctx, "denied") != "" {
		ctx.WriteString("<h1>OAuth Access Denied</h1>")
		return
	}

	oauth_token := GetParam(ctx, "oauth_token")
	oauth_verifier := GetParam(ctx, "oauth_verifier")

	at := googleConn.GetAccessToken(oauth_token, oauth_verifier)

	// Store at off to persistant data store for use later.
	googleAT = at

	googleAT.UserRef = "foo"

	ctx.WriteString("<h1>Access Token Received</h1>")
	defer func() { googleFooter(ctx) }()
}

func googleFooter(ctx *web.Context){
	ctx.WriteString("<p>Click here <a href=\"/google/getpublic\">to get latest public</a></p>")
	ctx.WriteString("<p>Click here <a href=\"/google/search\">to search buzz</a></p>")
}

func googleGetPublic(ctx *web.Context, name string){

	ctx.WriteString("<h1>Google public</h1>")
	if googleAT == nil {
		ctx.WriteString("<p>Please <a href=\"/signin/google\">Sign in to Google</a></p>")
		return
	}

	defer func() { googleFooter(ctx) }()
	ctx.WriteString("<p>Build &amp; Send request</p>")

	r, err := googleConn.Get(
		"https://www.googleapis.com/buzz/v1/activities/googlebuzz/@public",
		nil,
		googleAT )

	if err != nil {
		ctx.WriteString("<p style=\"color:red\">Error : " + err.String() + "</p>")
	}else{
		ctx.WriteString("<p style=\"color:green\">Sent Request Sent :-</p>")
	}

	b, _ := ioutil.ReadAll( r.Body ) 

	ctx.WriteString("<h2>Googles Response</h2>")
	ctx.WriteString("<textarea rows=\"20\" cols=\"60\">")
	ctx.Write(b)
	ctx.WriteString("</textarea>")

}


func googleGetSearch(ctx *web.Context, name string){

	ctx.WriteString("<h1>Google Buzz Search</h1>")

	if googleAT == nil {
		ctx.WriteString("<p>Please <a href=\"/signin/google\">Sign in to Google</a></p>")
		return
	}

	ctx.WriteString("<form method=\"post\"><p>Search : <input type=\"text\" name=\"q\"/></p><p>Lat :<input type=\"test\" name=\"lat\"/> Lon : <input type\"text\" name=\"lon\"/>  Radius : <input type=\"text\" name=\"radius\"/></p><p><input type=\"submit\" value=\"Search\"/></p></form>")

	defer func() { googleFooter(ctx) }()
}

func googleSearch(ctx *web.Context, name string){

	if googleAT == nil {
		ctx.WriteString("<h1>Google Search</h1>")
		ctx.WriteString("<p>Please <a href=\"/signin/google\">Sign in to Google</a></p>")
		return
	}
	googleGetSearch(ctx, name)

	q := GetParam( ctx, "q")
	lat := GetParam( ctx, "lat")
	lon := GetParam( ctx, "lon")
	radius := GetParam( ctx, "radius")

	// Defaults
	params := oauth.Params{
		&oauth.Pair{Key:"q", Value:q},
		&oauth.Pair{Key:"key", Value:"Change to your app API KEY"},
		&oauth.Pair{Key:"alt", Value:"json"},
	}

	//q=query&lat=latitude&lon=longitude&radius=radius
	if lat != "" && lon != "" && radius != "" {
		params.Add( &oauth.Pair{ Key:"lat", Value:lat } )
		params.Add( &oauth.Pair{ Key:"lon", Value:lon } )
		params.Add( &oauth.Pair{ Key:"radius", Value:radius } )
	}


	r, err := googleConn.Get(
		"https://www.googleapis.com/buzz/v1/activities/search",
		params,
		googleAT)

	b, _ := ioutil.ReadAll( r.Body ) 

	var m map[string]gobuzz.ActivityFeed

	err = json.Unmarshal(b, &m)
	if err != nil {
		ctx.WriteString(err.String())
		return
	}

	feed := m["data"]

	ctx.WriteString("<h2>" + feed.Title + "</h2>");

	for i := range feed.Items {
		activity := feed.Items[i]

		t := (*time.Time)(activity.Published)

		c := t.Format("1 January 2006 @ 3:04 pm")

		ctx.WriteString("<div>");

		ctx.WriteString("<h2>" + activity.Title + "</h2>")
		ctx.WriteString("<div>" + activity.Object.Content + "<p>Published : " + c + "</p></div>")

		ctx.WriteString("</div>");

	}
}



func noRespond(ctx *web.Context, name string) {
	ctx.WriteString("<h1>Testing OAuth With GoLang</h1>")

}

func GetParam(ctx *web.Context, param string) (v string){

	c, ok := ctx.Request.Params[param]

	if !ok { return }

	v = c
	return
}


