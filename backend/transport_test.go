// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
)

func TestFirebaseClient(t *testing.T) {
	ch := make(chan bool, 1)
	// Check that the Firebase client correctly adds the auth parameter
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		if request.URL.Query().Get("auth") != config.Firebase.Secret {
			t.Errorf("Expected auth query parameter to be the FB secret, URL was %v", request.URL)
		}

		ch <- true
	}))
	defer ts.Close()

	oldTransport := httpTransport
	httpTransport = func(_ context.Context) http.RoundTripper {
		return &http.Transport{}
	}
	defer func() {
		httpTransport = oldTransport
	}()

	c := newTestContext()

	if _, err := firebaseClient(c).Get(ts.URL); err != nil {
		t.Fatal(err)
	}

	select {
	case <-ch:
		// passed
	default:
		t.Fatalf("firebase request never happened")
	}
}
