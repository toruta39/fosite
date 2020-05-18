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

package openid

import (
	"testing"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/toruta39/fosite"
	"github.com/toruta39/fosite/token/jwt"
)

func TestOpenIDConnectRefreshHandler_HandleTokenEndpointRequest(t *testing.T) {
	h := &OpenIDConnectRefreshHandler{}
	for _, c := range []struct {
		areq        *fosite.AccessRequest
		expectedErr error
		description string
	}{
		{
			description: "should not pass because grant_type is wrong",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"foo"},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should not pass because grant_type is right but scope is missing",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"something"},
				},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should not pass because client may not execute this grant type",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client:       &fosite.DefaultClient{},
				},
			},
			expectedErr: fosite.ErrInvalidGrant,
		},
		{
			description: "should pass",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
						//ResponseTypes: []string{"id_token"},
					},
					Session: &DefaultSession{},
				},
			},
		},
	} {
		t.Run("case="+c.description, func(t *testing.T) {
			err := h.HandleTokenEndpointRequest(nil, c.areq)
			if c.expectedErr != nil {
				require.EqualError(t, errors.Cause(err), c.expectedErr.Error(), "%v", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOpenIDConnectRefreshHandler_PopulateTokenEndpointResponse(t *testing.T) {
	h := &OpenIDConnectRefreshHandler{
		IDTokenHandleHelper: &IDTokenHandleHelper{
			IDTokenStrategy: j,
		},
	}
	for _, c := range []struct {
		areq        *fosite.AccessRequest
		expectedErr error
		check       func(t *testing.T, aresp *fosite.AccessResponse)
		description string
	}{
		{
			description: "should not pass because grant_type is wrong",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"foo"},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		{
			description: "should not pass because grant_type is right but scope is missing",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"something"},
				},
			},
			expectedErr: fosite.ErrUnknownRequest,
		},
		// Disabled because this is already handled at the authorize_request_handler
		//{
		//	description: "should not pass because client may not ask for id_token",
		//	areq: &fosite.AccessRequest{
		//		GrantTypes: []string{"refresh_token"},
		//		Request: fosite.Request{
		//			GrantedScope: []string{"openid"},
		//			Client: &fosite.DefaultClient{
		//				GrantTypes: []string{"refresh_token"},
		//			},
		//		},
		//	},
		//	expectedErr: fosite.ErrUnknownRequest,
		//},
		{
			description: "should pass",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
						//ResponseTypes: []string{"id_token"},
					},
					Session: &DefaultSession{
						Subject: "foo",
						Claims: &jwt.IDTokenClaims{
							Subject: "foo",
						},
					},
				},
			},
			check: func(t *testing.T, aresp *fosite.AccessResponse) {
				assert.NotEmpty(t, aresp.GetExtra("id_token"))
				idToken, _ := aresp.GetExtra("id_token").(string)
				decodedIdToken, _ := jwtgo.Parse(idToken, func(token *jwtgo.Token) (interface{}, error) {
					return key.PublicKey, nil
				})
				claims, _ := decodedIdToken.Claims.(jwtgo.MapClaims)
				assert.NotEmpty(t, claims["at_hash"])
			},
		},
		{
			description: "should fail because missing subject claim",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
						//ResponseTypes: []string{"id_token"},
					},
					Session: &DefaultSession{
						Subject: "foo",
						Claims:  &jwt.IDTokenClaims{},
					},
				},
			},
			expectedErr: fosite.ErrServerError,
		},
		{
			description: "should fail because missing session",
			areq: &fosite.AccessRequest{
				GrantTypes: []string{"refresh_token"},
				Request: fosite.Request{
					GrantedScope: []string{"openid"},
					Client: &fosite.DefaultClient{
						GrantTypes: []string{"refresh_token"},
					},
				},
			},
			expectedErr: fosite.ErrServerError,
		},
	} {
		t.Run("case="+c.description, func(t *testing.T) {
			aresp := fosite.NewAccessResponse()
			err := h.PopulateTokenEndpointResponse(nil, c.areq, aresp)
			if c.expectedErr != nil {
				require.EqualError(t, errors.Cause(err), c.expectedErr.Error(), "%v", err)
			} else {
				require.NoError(t, err)
			}

			if c.check != nil {
				c.check(t, aresp)
			}
		})
	}
}
