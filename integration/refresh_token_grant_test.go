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
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/toruta39/fosite"
	"github.com/toruta39/fosite/compose"
	"github.com/toruta39/fosite/handler/openid"
	"github.com/toruta39/fosite/internal"
	"github.com/toruta39/fosite/token/jwt"
)

type introspectionResponse struct {
	Active    bool     `json:"active"`
	ClientID  string   `json:"client_id,omitempty"`
	Scope     string   `json:"scope,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Username  string   `json:"username,omitempty"`
}

func TestRefreshTokenFlow(t *testing.T) {
	session := &defaultSession{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Subject: "peter",
			},
			Headers:  &jwt.Headers{},
			Subject:  "peter",
			Username: "peteru",
		},
	}
	fc := new(compose.Config)
	fc.RefreshTokenLifespan = -1
	f := compose.ComposeAllEnabled(fc, fositeStore, []byte("some-secret-thats-random-some-secret-thats-random-"), internal.MustRSAKey())
	ts := mockServer(t, f, session)
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	state := "1234567890"
	fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"

	refreshCheckClient := &fosite.DefaultClient{
		ID:            "refresh-client",
		Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
		RedirectURIs:  []string{ts.URL + "/callback"},
		ResponseTypes: []string{"id_token", "code", "token", "token code", "id_token code", "token id_token", "token code id_token"},
		GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
		Scopes:        []string{"fosite", "offline", "openid"},
		Audience:      []string{"https://www.ory.sh/api"},
	}
	fositeStore.Clients["refresh-client"] = refreshCheckClient

	fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"
	for _, c := range []struct {
		description   string
		setup         func(t *testing.T)
		pass          bool
		params        []oauth2.AuthCodeOption
		check         func(t *testing.T, original, refreshed *oauth2.Token, or, rr *introspectionResponse)
		beforeRefresh func(t *testing.T)
		mockServer    func(t *testing.T) *httptest.Server
	}{
		{
			description: "should fail because refresh scope missing",
			setup: func(t *testing.T) {
				oauthClient.Scopes = []string{"fosite"}
			},
			pass: false,
		},
		{
			description: "should pass but not yield id token",
			setup: func(t *testing.T) {
				oauthClient.Scopes = []string{"offline"}
			},
			pass: true,
			check: func(t *testing.T, original, refreshed *oauth2.Token, or, rr *introspectionResponse) {
				assert.NotEqual(t, original.RefreshToken, refreshed.RefreshToken)
				assert.NotEqual(t, original.AccessToken, refreshed.AccessToken)
				assert.Nil(t, refreshed.Extra("id_token"))
			},
		},
		{
			description: "should pass and yield id token",
			params:      []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("audience", "https://www.ory.sh/api")},
			setup: func(t *testing.T) {
				oauthClient.Scopes = []string{"fosite", "offline", "openid"}
			},
			pass: true,
			check: func(t *testing.T, original, refreshed *oauth2.Token, or, rr *introspectionResponse) {
				assert.NotEqual(t, original.RefreshToken, refreshed.RefreshToken)
				assert.NotEqual(t, original.AccessToken, refreshed.AccessToken)
				assert.NotEqual(t, original.Extra("id_token"), refreshed.Extra("id_token"))
				assert.NotNil(t, refreshed.Extra("id_token"))

				assert.NotEmpty(t, or.Audience)
				assert.NotEmpty(t, or.ClientID)
				assert.NotEmpty(t, or.Scope)
				assert.NotEmpty(t, or.ExpiresAt)
				assert.NotEmpty(t, or.IssuedAt)
				assert.True(t, or.Active)
				assert.EqualValues(t, "peter", or.Subject)
				assert.EqualValues(t, "peteru", or.Username)

				assert.EqualValues(t, or.Audience, rr.Audience)
				assert.EqualValues(t, or.ClientID, rr.ClientID)
				assert.EqualValues(t, or.Scope, rr.Scope)
				assert.NotEqual(t, or.ExpiresAt, rr.ExpiresAt)
				assert.True(t, or.ExpiresAt < rr.ExpiresAt)
				assert.NotEqual(t, or.IssuedAt, rr.IssuedAt)
				assert.True(t, or.IssuedAt < rr.IssuedAt)
				assert.EqualValues(t, or.Active, rr.Active)
				assert.EqualValues(t, or.Subject, rr.Subject)
				assert.EqualValues(t, or.Username, rr.Username)
			},
		},
		{
			description: "should fail because scope is no longer allowed",
			setup: func(t *testing.T) {
				oauthClient.ClientID = refreshCheckClient.ID
				oauthClient.Scopes = []string{"fosite", "offline", "openid"}
			},
			beforeRefresh: func(t *testing.T) {
				refreshCheckClient.Scopes = []string{"offline", "openid"}
			},
			pass: false,
		},
		{
			description: "should fail because audience is no longer allowed",
			params:      []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("audience", "https://www.ory.sh/api")},
			setup: func(t *testing.T) {
				oauthClient.ClientID = refreshCheckClient.ID
				oauthClient.Scopes = []string{"fosite", "offline", "openid"}
				refreshCheckClient.Scopes = []string{"fosite", "offline", "openid"}
			},
			beforeRefresh: func(t *testing.T) {
				refreshCheckClient.Audience = []string{"https://www.not-ory.sh/api"}
			},
			pass: false,
		},
		{
			description: "should fail with expired refresh token",
			setup: func(t *testing.T) {
				fc = new(compose.Config)
				fc.RefreshTokenLifespan = time.Nanosecond
				f = compose.ComposeAllEnabled(fc, fositeStore, []byte("some-secret-thats-random-some-secret-thats-random-"), internal.MustRSAKey())
				ts = mockServer(t, f, session)

				oauthClient = newOAuth2Client(ts)
				oauthClient.Scopes = []string{"fosite", "offline", "openid"}
				fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"
			},
			pass: false,
		},
		{
			description: "should pass with limited but not expired refresh token",
			setup: func(t *testing.T) {
				fc = new(compose.Config)
				fc.RefreshTokenLifespan = time.Minute
				f = compose.ComposeAllEnabled(fc, fositeStore, []byte("some-secret-thats-random-some-secret-thats-random-"), internal.MustRSAKey())
				ts = mockServer(t, f, session)

				oauthClient = newOAuth2Client(ts)
				oauthClient.Scopes = []string{"fosite", "offline", "openid"}
				fositeStore.Clients["my-client"].(*fosite.DefaultClient).RedirectURIs[0] = ts.URL + "/callback"
			},
			beforeRefresh: func(t *testing.T) {
				refreshCheckClient.Audience = []string{}
			},
			pass:  true,
			check: func(t *testing.T, original, refreshed *oauth2.Token, or, rr *introspectionResponse) {},
		},
	} {
		t.Run("case="+c.description, func(t *testing.T) {
			c.setup(t)

			var intro = func(token string, p interface{}) {
				req, err := http.NewRequest("POST", ts.URL+"/introspect", strings.NewReader(url.Values{"token": {token}}.Encode()))
				require.NoError(t, err)
				req.SetBasicAuth("refresh-client", "foobar")
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				r, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, r.StatusCode)

				dec := json.NewDecoder(r.Body)
				dec.DisallowUnknownFields()
				require.NoError(t, dec.Decode(p))
			}

			resp, err := http.Get(oauthClient.AuthCodeURL(state, c.params...))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			if resp.StatusCode != http.StatusOK {
				return
			}

			token, err := oauthClient.Exchange(oauth2.NoContext, resp.Request.URL.Query().Get("code"))
			require.NoError(t, err)
			require.NotEmpty(t, token.AccessToken)

			var ob introspectionResponse
			intro(token.AccessToken, &ob)

			t.Logf("Token %s\n", token)
			token.Expiry = token.Expiry.Add(-time.Hour * 24)

			if c.beforeRefresh != nil {
				c.beforeRefresh(t)
			}

			tokenSource := oauthClient.TokenSource(oauth2.NoContext, token)

			// This sleep guarantees time difference in exp/iat
			time.Sleep(time.Second * 2)

			refreshed, err := tokenSource.Token()
			if c.pass {
				require.NoError(t, err)

				var rb introspectionResponse
				intro(refreshed.AccessToken, &rb)
				c.check(t, token, refreshed, &ob, &rb)
			} else {
				require.Error(t, err)
			}
		})
	}
}
