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

package fosite_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/toruta39/fosite"
	. "github.com/toruta39/fosite"
	"github.com/toruta39/fosite/compose"
	"github.com/toruta39/fosite/internal"
	"github.com/toruta39/fosite/storage"
)

func TestIntrospectionResponse(t *testing.T) {
	r := &fosite.IntrospectionResponse{
		AccessRequester: fosite.NewAccessRequest(nil),
		Active:          true,
	}

	assert.Equal(t, r.AccessRequester, r.GetAccessRequester())
	assert.Equal(t, r.Active, r.IsActive())
}

func TestNewIntrospectionRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	validator := internal.NewMockTokenIntrospector(ctrl)
	defer ctrl.Finish()

	f := compose.ComposeAllEnabled(new(compose.Config), storage.NewExampleStore(), []byte{}, nil).(*Fosite)
	httpreq := &http.Request{
		Method: "POST",
		Header: http.Header{},
		Form:   url.Values{},
	}
	newErr := errors.New("asdf")

	for k, c := range []struct {
		description string
		setup       func()
		expectErr   error
		isActive    bool
	}{
		{
			description: "should fail",
			setup: func() {
			},
			expectErr: ErrInvalidRequest,
		},
		{
			description: "should fail",
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						"Authorization": []string{"bearer some-token"},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(context.TODO(), "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenType(""), nil)
				validator.EXPECT().IntrospectToken(context.TODO(), "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenType(""), newErr)
			},
			isActive:  false,
			expectErr: ErrInactiveToken,
		},
		{
			description: "should pass",
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						"Authorization": []string{"bearer some-token"},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(context.TODO(), "some-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenType(""), nil)
				validator.EXPECT().IntrospectToken(context.TODO(), "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenType(""), nil)
			},
			isActive: true,
		},
		{
			description: "should pass with basic auth if username and password encoded",
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						//Basic Authorization with username=encoded:client and password=encoded&password
						"Authorization": []string{"Basic ZW5jb2RlZCUzQWNsaWVudDplbmNvZGVkJTI2cGFzc3dvcmQ="},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(context.TODO(), "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenType(""), nil)
			},
			isActive: true,
		},
		{
			description: "should pass with basic auth if username and password not encoded",
			setup: func() {
				f.TokenIntrospectionHandlers = TokenIntrospectionHandlers{validator}
				httpreq = &http.Request{
					Method: "POST",
					Header: http.Header{
						//Basic Authorization with username=my-client and password=foobar
						"Authorization": []string{"Basic bXktY2xpZW50OmZvb2Jhcg=="},
					},
					PostForm: url.Values{
						"token": []string{"introspect-token"},
					},
				}
				validator.EXPECT().IntrospectToken(context.TODO(), "introspect-token", gomock.Any(), gomock.Any(), gomock.Any()).Return(TokenType(""), nil)
			},
			isActive: true,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.setup()
			res, err := f.NewIntrospectionRequest(context.TODO(), httpreq, &DefaultSession{})

			if c.expectErr != nil {
				assert.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.isActive, res.IsActive())
			}
		})
	}
}
