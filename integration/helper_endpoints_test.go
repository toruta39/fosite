/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/toruta39/fosite"
	"github.com/toruta39/fosite/handler/oauth2"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func tokenRevocationHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		err := oauth2.NewRevocationRequest(ctx, req)
		if err != nil {
			t.Logf("Revoke request failed because %+v", err)
		}
		oauth2.WriteRevocationResponse(rw, err)
	}
}

func tokenIntrospectionHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		ar, err := oauth2.NewIntrospectionRequest(ctx, req, session)
		if err != nil {
			t.Logf("Introspection request failed because: %+v", err)
			oauth2.WriteIntrospectionError(rw, err)
			return
		}

		oauth2.WriteIntrospectionResponse(rw, ar)
	}
}

func tokenInfoHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		_, resp, err := oauth2.IntrospectToken(ctx, fosite.AccessTokenFromRequest(req), fosite.AccessToken, session)
		if err != nil {
			t.Logf("Info request failed because: %+v", err)
			http.Error(rw, errors.Cause(err).(*fosite.RFC6749Error).Description, errors.Cause(err).(*fosite.RFC6749Error).Code)
			return
		}

		t.Logf("Introspecting caused: %+v", resp)

		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func authEndpointHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()

		ar, err := oauth2.NewAuthorizeRequest(ctx, req)
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", ar)
			oauth2.WriteAuthorizeError(rw, ar, err)
			return
		}

		if ar.GetRequestedScopes().Has("fosite") {
			ar.GrantScope("fosite")
		}

		if ar.GetRequestedScopes().Has("offline") {
			ar.GrantScope("offline")
		}

		if ar.GetRequestedScopes().Has("openid") {
			ar.GrantScope("openid")
		}

		for _, a := range ar.GetRequestedAudience() {
			ar.GrantAudience(a)
		}

		// Normally, this would be the place where you would check if the user is logged in and gives his consent.
		// For this test, let's assume that the user exists, is logged in, and gives his consent...

		response, err := oauth2.NewAuthorizeResponse(ctx, ar, session)
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", ar)
			oauth2.WriteAuthorizeError(rw, ar, err)
			return
		}

		oauth2.WriteAuthorizeResponse(rw, ar, response)
	}
}

func authCallbackHandler(t *testing.T) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		if q.Get("code") == "" && q.Get("error") == "" {
			assert.NotEmpty(t, q.Get("code"))
			assert.NotEmpty(t, q.Get("error"))
		}

		if q.Get("code") != "" {
			rw.Write([]byte("code: ok"))
		}
		if q.Get("error") != "" {
			rw.WriteHeader(http.StatusNotAcceptable)
			rw.Write([]byte("error: " + q.Get("error")))
		}

	}
}

func tokenEndpointHandler(t *testing.T, provider fosite.OAuth2Provider) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseMultipartForm(1 << 20)
		ctx := fosite.NewContext()

		accessRequest, err := provider.NewAccessRequest(ctx, req, &oauth2.JWTSession{})
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", accessRequest)
			provider.WriteAccessError(rw, accessRequest, err)
			return
		}

		if accessRequest.GetRequestedScopes().Has("fosite") {
			accessRequest.GrantScope("fosite")
		}

		response, err := provider.NewAccessResponse(ctx, accessRequest)
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", accessRequest)
			provider.WriteAccessError(rw, accessRequest, err)
			return
		}

		provider.WriteAccessResponse(rw, accessRequest, response)
	}
}
