include $(GOROOT)/src/Make.inc

TARG=github.com/hokapoka/goauth
GOFILES=\
	urlencode.go\
	token.go\
	pairs.go\
	http.go\
	oauthconsumer.go\
	oauth.go\

include $(GOROOT)/src/Make.pkg

