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

package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/toruta39/fosite/token/jwt"
)

var idTokenClaims = &IDTokenClaims{
	JTI:                                 "foo-id",
	Subject:                             "peter",
	IssuedAt:                            time.Now().UTC().Round(time.Second),
	Issuer:                              "fosite",
	Audience:                            []string{"tests"},
	ExpiresAt:                           time.Now().UTC().Add(time.Hour).Round(time.Second),
	AuthTime:                            time.Now().UTC(),
	RequestedAt:                         time.Now().UTC(),
	AccessTokenHash:                     "foobar",
	CodeHash:                            "barfoo",
	AuthenticationContextClassReference: "acr",
	AuthenticationMethodsReference:      "amr",
	Extra: map[string]interface{}{
		"foo": "bar",
		"baz": "bar",
	},
}

func TestIDTokenAssert(t *testing.T) {
	assert.NoError(t, (&IDTokenClaims{ExpiresAt: time.Now().UTC().Add(time.Hour)}).
		ToMapClaims().Valid())
	assert.Error(t, (&IDTokenClaims{ExpiresAt: time.Now().UTC().Add(-time.Hour)}).
		ToMapClaims().Valid())

	assert.NotEmpty(t, (new(IDTokenClaims)).ToMapClaims()["jti"])
}

func TestIDTokenClaimsToMap(t *testing.T) {
	assert.Equal(t, map[string]interface{}{
		"jti":       idTokenClaims.JTI,
		"sub":       idTokenClaims.Subject,
		"iat":       float64(idTokenClaims.IssuedAt.Unix()),
		"rat":       float64(idTokenClaims.RequestedAt.Unix()),
		"iss":       idTokenClaims.Issuer,
		"aud":       idTokenClaims.Audience,
		"nonce":     idTokenClaims.Nonce,
		"exp":       float64(idTokenClaims.ExpiresAt.Unix()),
		"foo":       idTokenClaims.Extra["foo"],
		"baz":       idTokenClaims.Extra["baz"],
		"at_hash":   idTokenClaims.AccessTokenHash,
		"c_hash":    idTokenClaims.CodeHash,
		"auth_time": idTokenClaims.AuthTime.Unix(),
		"acr":       idTokenClaims.AuthenticationContextClassReference,
		"amr":       idTokenClaims.AuthenticationMethodsReference,
	}, idTokenClaims.ToMap())
}
