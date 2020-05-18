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
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	. "github.com/toruta39/fosite"
	. "github.com/toruta39/fosite/internal"
)

func TestWriteAuthorizeResponse(t *testing.T) {
	oauth2 := &Fosite{}
	header := http.Header{}
	ctrl := gomock.NewController(t)
	rw := NewMockResponseWriter(ctrl)
	ar := NewMockAuthorizeRequester(ctrl)
	resp := NewMockAuthorizeResponder(ctrl)
	defer ctrl.Finish()

	for k, c := range []struct {
		setup  func()
		expect func()
	}{
		{
			setup: func() {
				redir, _ := url.Parse("https://foobar.com/?foo=bar")
				ar.EXPECT().GetRedirectURI().Return(redir)
				resp.EXPECT().GetFragment().Return(url.Values{})
				resp.EXPECT().GetHeader().Return(http.Header{})
				resp.EXPECT().GetQuery().Return(url.Values{})

				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			expect: func() {
				assert.Equal(t, http.Header{
					"Location": []string{"https://foobar.com/?foo=bar"},
				}, header)
			},
		},
		{
			setup: func() {
				redir, _ := url.Parse("https://foobar.com/?foo=bar")
				ar.EXPECT().GetRedirectURI().Return(redir)
				resp.EXPECT().GetFragment().Return(url.Values{"bar": {"baz"}})
				resp.EXPECT().GetHeader().Return(http.Header{})
				resp.EXPECT().GetQuery().Return(url.Values{})

				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			expect: func() {
				assert.Equal(t, http.Header{
					"Location": []string{"https://foobar.com/?foo=bar#bar=baz"},
				}, header)
			},
		},
		{
			setup: func() {
				redir, _ := url.Parse("https://foobar.com/?foo=bar")
				ar.EXPECT().GetRedirectURI().Return(redir)
				resp.EXPECT().GetFragment().Return(url.Values{"bar": {"baz"}})
				resp.EXPECT().GetHeader().Return(http.Header{})
				resp.EXPECT().GetQuery().Return(url.Values{"bar": {"baz"}})

				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			expect: func() {
				assert.Equal(t, http.Header{
					"Location": []string{"https://foobar.com/?bar=baz&foo=bar#bar=baz"},
				}, header)
			},
		},
		{
			setup: func() {
				redir, _ := url.Parse("https://foobar.com/?foo=bar")
				ar.EXPECT().GetRedirectURI().Return(redir)
				resp.EXPECT().GetFragment().Return(url.Values{"bar": {"baz"}, "scope": {"a b"}})
				resp.EXPECT().GetHeader().Return(http.Header{"X-Bar": {"baz"}})
				resp.EXPECT().GetQuery().Return(url.Values{"bar": {"b+az"}, "scope": {"a b"}})

				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			expect: func() {
				assert.Equal(t, http.Header{
					"X-Bar":    {"baz"},
					"Location": {"https://foobar.com/?bar=b%2Baz&foo=bar&scope=a%20b#bar=baz&scope=a%20b"},
				}, header)
			},
		},
		{
			setup: func() {
				redir, _ := url.Parse("https://foobar.com/?foo=bar")
				ar.EXPECT().GetRedirectURI().Return(redir)
				resp.EXPECT().GetFragment().Return(url.Values{"bar": {"baz"}, "scope": {"api:*"}})
				resp.EXPECT().GetHeader().Return(http.Header{"X-Bar": {"baz"}})
				resp.EXPECT().GetQuery().Return(url.Values{"bar": {"b+az"}, "scope": {"api:*"}})

				rw.EXPECT().Header().Return(header)
				rw.EXPECT().WriteHeader(http.StatusFound)
			},
			expect: func() {
				assert.Equal(t, http.Header{
					"X-Bar":    {"baz"},
					"Location": {"https://foobar.com/?bar=b%2Baz&foo=bar&scope=api%3A%2A#bar=baz&scope=api%3A%2A"},
				}, header)
			},
		},
	} {
		t.Logf("Starting test case %d", k)
		c.setup()
		oauth2.WriteAuthorizeResponse(rw, ar, resp)
		c.expect()
		header = http.Header{}
		t.Logf("Passed test case %d", k)
	}
}
