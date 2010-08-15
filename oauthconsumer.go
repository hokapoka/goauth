package oauth

import(
	"fmt"
	"http"
	"rand"
	"strconv"
	"time"
	"sort"
	"strings"
	"bytes"
	"crypto/hmac"
	"encoding/base64"
	"io/ioutil"
	"os"
)


type OAuthConsumer struct{
	Service string
	RequestTokenURL string
	AccessTokenURL string
	AuthorizationURL string
	ConsumerKey string
	ConsumerSecrect string
	CallBackURL string
	requestTokens []*RequestToken
}



func (oc *OAuthConsumer) GetRequestAuthorizationURL() (string, os.Error){
	// Gather the params
	p := Params{}

	// Add required OAuth params
	p.Add( &Pair{ Key:"oauth_callback", Value:oc.CallBackURL } )
	p.Add( &Pair{ Key:"oauth_consumer_key", Value:oc.ConsumerKey } )
	p.Add( &Pair{ Key:"oauth_nonce", Value:strconv.Itoa64(rand.Int63()) } )
	p.Add( &Pair{ Key:"oauth_signature_method", Value:"HMAC-SHA1" } )
	p.Add( &Pair{ Key:"oauth_timestamp", Value:strconv.Itoa64(time.Seconds()) } )
	p.Add( &Pair{ Key:"oauth_version", Value:"1.0" } )

	// Sort the collection
	sort.Sort(p)

	// Generate string of sorted params
	sigBaseCol := make([]string, len(p))
	for i := range p {
		sigBaseCol[i] = EncodeURL(p[i].Key) + "=" + EncodeURL( p[i].Value )
	}

	sigBaseStr :=	"POST&" +
					EncodeURL(oc.RequestTokenURL) + "&" +
					EncodeURL(strings.Join(sigBaseCol, "&"))

	// Generate Composite Signing key
	key := oc.ConsumerSecrect + "&" + "" // token secrect is blank on the Resquest Token

	// Generate Signature
	d := oc.digest(key, sigBaseStr)

	// Build Auth Header
	authHeader := "OAuth "
	for i := range p {
		authHeader +=  p[i].Key + "=\"" + EncodeURL(p[i].Value ) + "\", "
	}

	// Add the signature
	authHeader += "oauth_signature=\"" + EncodeURL(d) + "\""

	// Add Header & Buffer for params
	buf := &bytes.Buffer{}
	headers := map[string]string{
		"Content-Type":"application/x-www-form-urlencoded",
		"Authorization":authHeader,
	}

	r, err := Post(oc.RequestTokenURL, headers, buf)

	if err != nil {
		return "", err
	}

	b, _ := ioutil.ReadAll( r.Body ) 
	s := string(b)

	rt := &RequestToken{}

	if strings.Index(s, "&") > -1 {
		vals := strings.Split(s, "&", 10)

		for i := range vals {
			if strings.Index(vals[i], "=") > -1 {
				kv := strings.Split(vals[i], "=", 2)
				if len(kv) > 0 { // Adds the key even if there's no value. 
					switch kv[0]{
						case "oauth_token":					if len(kv) > 1 { rt.Token = kv[1] }; break
						case "oauth_token_secret":			if len(kv) > 1 { rt.Secret = kv[1] }; break
					}
				}
			}
		}
	}

	oc.appendRequestToken(rt)

	return oc.AuthorizationURL + "?oauth_token=" + rt.Token, nil

}

// GetAccessToken gets the access token
func (oc *OAuthConsumer) GetAccessToken(token string, verifier string, ) *AccessToken{

	fmt.Println("***************************** GET ACCESS TOKEN **********************")
	var rt *RequestToken

	// Match the RequestToken by Token
	for i := range oc.requestTokens {
		if oc.requestTokens[i].Token == token{
			rt = oc.requestTokens[i]
		}
	}

	rt.Verifier = verifier

	// Gather the params
	p := Params{}

	// Add required OAuth params
	p.Add( &Pair{ Key:"oauth_consumer_key", Value:oc.ConsumerKey } )
	p.Add( &Pair{ Key:"oauth_nonce", Value:strconv.Itoa64(rand.Int63()) } )
	p.Add( &Pair{ Key:"oauth_signature_method", Value:"HMAC-SHA1" } )
	p.Add( &Pair{ Key:"oauth_timestamp", Value:strconv.Itoa64(time.Seconds()) } )
	p.Add( &Pair{ Key:"oauth_version", Value:"1.0" } )
	p.Add( &Pair{ Key:"oauth_token", Value:rt.Token })
	p.Add( &Pair{ Key:"oauth_verifier", Value:rt.Verifier })

	// Sort the collection
	sort.Sort(p)

	// Generate string of sorted params
	sigBaseCol := make([]string, len(p))
	for i := range p {
		sigBaseCol[i] = EncodeURL(p[i].Key) + "=" + EncodeURL( p[i].Value )
	}

	sigBaseStr :=	"POST&" +
					EncodeURL(oc.AccessTokenURL) + "&" +
					EncodeURL(strings.Join(sigBaseCol, "&"))

	// Generate Composite Signing key
	key := oc.ConsumerSecrect + "&" + rt.Secret

	// Generate Signature
	d := oc.digest(key, sigBaseStr)

	// Build Auth Header
	authHeader := "OAuth "
	for i := range p {
		authHeader +=  p[i].Key + "=\"" + EncodeURL(p[i].Value ) + "\", "
	}

	// Add the signature
	authHeader += "oauth_signature=\"" + EncodeURL(d) + "\""

	// Add Header & Buffer for params
	buf := &bytes.Buffer{}
	headers := map[string]string{
		"Content-Type":"application/x-www-form-urlencoded",
		"Authorization":authHeader,
	}

	// Action the POST to get the AccessToken
	r, err :=  Post(oc.AccessTokenURL, headers, buf)
	if err != nil {
		fmt.Println(err.String())
		return nil
	}

	// Read response Body & Create AccessToken
	b, _ := ioutil.ReadAll( r.Body ) 
	s := string(b)
	at := &AccessToken{}

	if strings.Index(s, "&") > -1 {
		vals := strings.Split(s, "&", 10)

		for i := range vals {
			if strings.Index(vals[i], "=") > -1 {
				kv := strings.Split(vals[i], "=", 2)
				if len(kv) > 0 { // Adds the key even if there's no value. 
					switch kv[0]{
						case "oauth_token":					if len(kv) > 1 { at.Token = kv[1] };  break
						case "oauth_token_secret":			if len(kv) > 1 { at.Secret = kv[1] }; break
					}
				}
			}
		}
	}

	// Return the AccessToken
	return at

}

// OAuthRequestGet return the response via a GET for the url with the AccessToken passed
func (oc *OAuthConsumer) Get( url string, at *AccessToken) (r *http.Response, err os.Error) {
	return oc.oAuthRequest(url, Params{}, at, "GET")
}

// OAuthRequest returns the response via a POST for the url with the AccessToken passed & the Form params passsed in fparams
func (oc *OAuthConsumer) Post( url string, fparams Params, at *AccessToken) (r *http.Response, err os.Error) {
	return oc.oAuthRequest( url, fparams, at, "POST")
}

func (oc *OAuthConsumer) oAuthRequest( url string, fparams Params, at *AccessToken, method string) (r *http.Response, err os.Error) {

	fmt.Println("***************************** DO REQUEST **********************")
	// Gather the params
	p := Params{}

	// Add required OAuth params
	//p.Add( &Pair{ Key:"oauth_callback", Value:oc.CallBackURL } )
	p.Add( &Pair{ Key:"oauth_consumer_key", Value:oc.ConsumerKey } )
	p.Add( &Pair{ Key:"oauth_nonce", Value:strconv.Itoa64(rand.Int63()) } )
	p.Add( &Pair{ Key:"oauth_signature_method", Value:"HMAC-SHA1" } )
	p.Add( &Pair{ Key:"oauth_timestamp", Value:strconv.Itoa64(time.Seconds()) } )
	p.Add( &Pair{ Key:"oauth_version", Value:"1.0" } )


	// Add the token params if passed
	if at != nil {
		p.Add( &Pair{ Key:"oauth_token", Value:at.Token })
	}

	fparamsStr := ""
	// Add any additional params passed
	for i := range fparams{
		k, v := fparams[i].Key, fparams[i].Value
		p.Add( &Pair{ Key:k, Value:v } )
		fparamsStr += k + "=" + v + "&"
	}

	// Sort the collection
	sort.Sort(p)

	// Generate string of sorted params
	sigBaseCol := make([]string, len(p))
	for i := range p {
		sigBaseCol[i] = EncodeURL(p[i].Key) + "=" + EncodeURL( p[i].Value )
	}

	sigBaseStr :=	method + "&" +
					EncodeURL(url) + "&" +
					EncodeURL(strings.Join(sigBaseCol, "&"))

	// Generate Composite Signing key
	key := oc.ConsumerSecrect + "&" + at.Secret

	// Generate Signature
	d := oc.digest(key, sigBaseStr)

	// Build Auth Header
	authHeader := "OAuth "
	for i := range p {
		if strings.Index(p[i].Key, "oauth") == 0 {
			//Add it to the authHeader
			authHeader +=  p[i].Key + "=\"" + EncodeURL(p[i].Value ) + "\", "
		}
	}

	// Add the signature
	authHeader += "oauth_signature=\"" + EncodeURL(d) + "\""

	// Add Header & Buffer for params
	buf := bytes.NewBufferString(fparamsStr)
	headers := map[string]string{
		"Content-Type":"application/x-www-form-urlencoded",
		"Authorization":authHeader,
	}

	fmt.Println(authHeader)
	if method == "GET" {
		// return Get response
		return Get(url, headers)
	}

	// return POSTs response
	return Post(url, headers, buf)

}


// digest Generates a HMAC-1234 for the signature
func (oc *OAuthConsumer) digest(key string, m string) string {
	h := hmac.NewSHA1([]byte(key))
	h.Write([]byte(m))
	s := bytes.TrimSpace(h.Sum())
	d := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(d, s)
	ds := strings.TrimSpace(bytes.NewBuffer(d).String())
	return ds
}

// appendRequestToken adds the Request Tokens to a localy temp collection
func (oc *OAuthConsumer) appendRequestToken(token *RequestToken){

	if oc.requestTokens == nil { oc.requestTokens = make([]*RequestToken, 0, 4) }

	n := len(oc.requestTokens)

	if n+1 > cap(oc.requestTokens) {
		s := make([]*RequestToken, n, 2*n+1)
		copy(s, oc.requestTokens)
		oc.requestTokens = s
	}
	oc.requestTokens = oc.requestTokens[0 : n+1]
	oc.requestTokens[n] = token

}


