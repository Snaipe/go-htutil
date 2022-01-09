// Copyright 2022 Franklin "Snaipe" Mathieu.
//
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

package htutil

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"testing"
)

func TestSortAcceptables(t *testing.T) {
	t.Parallel()

	tcases := []struct {
		In  []Acceptable
		Out []Acceptable
	}{
		// Test cases from RFC2616 14.1
		{
			In: []Acceptable{
				{ Value: "text/*",    Quality: 1.0 },
				{ Value: "text/html", Quality: 1.0 },
				{ Value: "text/html", Quality: 1.0, Params: map[string]string{"level": "1"} },
				{ Value: "*/*",       Quality: 1.0 },
			},
			Out: []Acceptable{
				{ Value: "text/html", Quality: 1.0, Params: map[string]string{"level": "1"} },
				{ Value: "text/html", Quality: 1.0 },
				{ Value: "text/*",    Quality: 1.0 },
				{ Value: "*/*",       Quality: 1.0 },
			},
		},
		{
			In: []Acceptable{
				{ Value: "text/*",    Quality: 0.3 },
				{ Value: "text/html", Quality: 0.7 },
				{ Value: "text/html", Quality: 1.0, Params: map[string]string{"level": "1"} },
				{ Value: "text/html", Quality: 0.4, Params: map[string]string{"level": "2"} },
				{ Value: "*/*",       Quality: 0.5 },
			},
			Out: []Acceptable{
				{ Value: "text/html", Quality: 1.0, Params: map[string]string{"level": "1"} },
				{ Value: "text/html", Quality: 0.7 },
				{ Value: "*/*",       Quality: 0.5 },
				{ Value: "text/html", Quality: 0.4, Params: map[string]string{"level": "2"} },
				{ Value: "text/*",    Quality: 0.3 },
			},
		},
	}

	for i, tcase := range tcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			sort.Slice(tcase.In, func(i, j int) bool { return Acceptable.Less(tcase.In[i], tcase.In[j]) })
			if !reflect.DeepEqual(tcase.In, tcase.Out) {
				t.Fatalf("expected %v, got %v", tcase.Out, tcase.In)
			}
		})
	}
}

func TestNegotiateContent(t *testing.T) {
	t.Parallel()

	tcases := []struct {
		Accept string
		Offers []string
		Expect string
	}{
		{
			Accept: "text/plain; q=0.1, application/json",
			Offers: []string{"text/plain", "application/json"},
			Expect: "application/json",
		},
		{
			Accept: "text/plain; q=0.1, application/json",
			Offers: []string{"text/plain"},
			Expect: "text/plain",
		},
		{
			Accept: "text/plain; q=0.1, application/json",
			Offers: []string{},
			Expect: "",
		},
	}

	for i, tcase := range tcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			hdr := http.Header{}
			hdr.Set("Accept", tcase.Accept)

			actual, _ := NegotiateContent(hdr, "Accept", tcase.Offers...)
			if actual != tcase.Expect {
				t.Fatalf("expected %v, got %v", tcase.Expect, actual)
			}
		})
	}
}
