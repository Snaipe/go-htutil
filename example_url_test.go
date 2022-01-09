// Copyright 2022 Franklin "Snaipe" Mathieu.
//
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

package htutil_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"snai.pe/go-htutil"
)

type T struct {
	URL htutil.URL `json:"url"`
}

func ExampleURL_UnmarshalText() {
	jsonstr := []byte(`{"url":"https://google.com"}`)

	var value T
	if err := json.Unmarshal(jsonstr, &value); err != nil {
		log.Fatal(err)
	}

	fmt.Println(value.URL)
	// Output: https://google.com
}

func ExampleURL_MarshalText() {
	u, err := url.Parse("https://google.com")
	if err != nil {
		log.Fatal(err)
	}

	value := T{
		URL: htutil.URL{u},
	}

	txt, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(txt))
	// Output: {"url":"https://google.com"}
}
