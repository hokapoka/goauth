package oauth

import (
	"http"
	"os"
	"io"
	"net"
	"encoding/base64"
	"strings"
	"bufio"
	"crypto/tls"
)



// Get Taken from the golang source modifed to allow headers to be passed and no redirection allowed
func Get(url string,  headers map[string]string) (r *http.Response, err os.Error) {

	var req http.Request
	if err != nil { return }
	req.Header = headers
	req.URL, err = http.ParseURL(url)

	r, err = send(&req)
	if err != nil { return }
	return
}

// Post taken from Golang modified to allow Headers to be pased
func Post(url string, headers map[string]string, body io.Reader) (r *http.Response, err os.Error) {
    var req http.Request
    req.Method = "POST"
    req.ProtoMajor = 1
    req.ProtoMinor = 1
    req.Close = true
    req.Body = nopCloser{body}
	req.Header = headers
    req.TransferEncoding = []string{"chunked"}

    req.URL, err = http.ParseURL(url)
    if err != nil {
        return nil, err
    }

    return send(&req)
}

// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//

// From the http package - modified to allow Headers to be sent to the Post method
type nopCloser struct {
    io.Reader
}

func (nopCloser) Close() os.Error { return nil }

type readClose struct {
    io.Reader
    io.Closer
}

func send(req *http.Request) (resp *http.Response, err os.Error) {
    if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
        return nil, nil
    }

    addr := req.URL.Host
    if !hasPort(addr) {
        addr += ":" + req.URL.Scheme
    }
    info := req.URL.Userinfo
    if len(info) > 0 {
        enc := base64.URLEncoding
        encoded := make([]byte, enc.EncodedLen(len(info)))
        enc.Encode(encoded, []byte(info))
        if req.Header == nil {
            req.Header = make(map[string]string)
        }
        req.Header["Authorization"] = "Basic " + string(encoded)
    }

    var conn io.ReadWriteCloser
    if req.URL.Scheme == "http" {
        conn, err = net.Dial("tcp", "", addr)
    } else { // https
        conn, err = tls.Dial("tcp", "", addr)
    }
    if err != nil {
        return nil, err
    }

    err = req.Write(conn)
    if err != nil {
        conn.Close()
        return nil, err
    }

    reader := bufio.NewReader(conn)
    resp, err = http.ReadResponse(reader, req.Method)
    if err != nil {
        conn.Close()
        return nil, err
    }

    resp.Body = readClose{resp.Body, conn}

    return
}


func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }
