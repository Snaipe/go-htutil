// Copyright 2022 Franklin "Snaipe" Mathieu.
//
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

package htutil

import "net/url"

// URL embeds *url.URL, but implements encoding.TextMarshaler and
// encoding.TextUnmarshaler to simply call MarshalBinary and UnmarshalBinary
// respectively.
type URL struct {
	*url.URL
}

func (u *URL) UnmarshalText(data []byte) error {
	if u.URL == nil {
		u.URL = new(url.URL)
	}
	return u.UnmarshalBinary(data)
}

func (u URL) MarshalText() ([]byte, error) {
	return u.MarshalBinary()
}
