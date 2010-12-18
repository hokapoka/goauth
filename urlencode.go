// Copyright 2010 Calvin McAnarney <calvin@mcanarney.org>
//
// Use of this file is governed by the ISC license. See the LICENSE file for
// details.

package oauth

import (
	"bytes"
	"encoding/base64"
)

// base64encode returns a string representation of Base64-encoded src.
func base64encode(src []byte) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(src)
	encoder.Close()
	return buf.String()
}

// Encode percent-encodes a string as defined in RFC 3986.
func Encode(s string) string {
	var enc string
	for _, c := range []byte(s) {
		if isEncodable(c) {
			enc += "%"
			enc += string("0123456789ABCDEF"[c>>4])
			enc += string("0123456789ABCDEF"[c&15])
		} else {
			enc += string(c)
		}
	}
	return enc
}

// isEncodable returns true if a given character should be percent-encoded
// according to RFC 3986.
func isEncodable(c byte) bool {
	// return false if c is an unreserved character (see RFC 3986 section 2.3)
	switch {
	case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z'):
		return false
	case c >= '0' && c <= '9':
		return false
	case c == '-' || c == '.' || c == '_' || c == '~':
		return false
	}
	return true
}

