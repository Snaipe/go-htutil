// Copyright 2022 Franklin "Snaipe" Mathieu.
//
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

package htutil_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"snai.pe/go-htutil"
)

func ExampleNegotiateContent() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctype, _ := htutil.NegotiateContent(req.Header, "Accept",
			"text/plain",
			"application/json",
		)
		if ctype == "" {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		w.Header().Set("Content-Type", ctype)
		w.WriteHeader(http.StatusOK)
		switch ctype {
		case "application/json":
			fmt.Fprint(w, `{"message":"OK"}`)
		case "text/plain":
			fmt.Fprint(w, `OK`)
		}
	})

	server := http.Server{Addr: ":8080"}
	defer server.Shutdown(context.Background())

	go server.ListenAndServe()

	get := func(addr, accept string) string {
		req, err := http.NewRequest("GET", addr, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Accept", accept)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		var out strings.Builder
		io.Copy(&out, resp.Body)
		return out.String()
	}

	fmt.Println(get("http://localhost:8080", "*/*"))
	fmt.Println(get("http://localhost:8080", "text/plain"))
	fmt.Println(get("http://localhost:8080", "application/json"))
	fmt.Println(get("http://localhost:8080", "text/html"))
	fmt.Println(get("http://localhost:8080", "application/json, text/*;q=0.5, */*;q=0.1"))
	fmt.Println(get("http://localhost:8080", ""))
	// Output: OK
	// OK
	// {"message":"OK"}
	// 406 Not Acceptable
	// {"message":"OK"}
	// 406 Not Acceptable
}
