// Copyright 2022 Franklin "Snaipe" Mathieu.
//
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

package htutil

import (
	"fmt"
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func qualityEq(lhs, rhs float32) bool {
	// The quality factor of a MIME type has up to 3 precision digits
	const ε = 0.001
	return rhs - lhs <= ε && lhs - rhs <= ε
}

// Acceptable represents an acceptable value for a response; typical use is
// representing media (or MIME) type value as per RFC2616 §14.1 in an Accept
// header, as well as encoding (Accept-Encoding), language (Accept-Language),
// and character set (Accept-Charset).
type Acceptable struct {
	// The value that is acceptable for a response. May contain wildcard
	// ('*') characters.
	Value string

	// The quality, between 0 and 1, of the accepted value.
	Quality float32

	// Params contains any extra optional parameters for this value.
	Params map[string]string
}

// ParseAcceptable parses a single acceptable value, as laid out in an
// Accept{,-*} or Content-* header as per RFC2616 §14.1
func ParseAcceptable(v string) (Acceptable, error) {
	// mime.ParseMediaType actually understands other things than pure media
	// types, like encoding, language, and charsets. It also ensures that
	// 
	value, params, err := mime.ParseMediaType(v)
	if err != nil {
		return Acceptable{}, err
	}

	quality := 1.0
	if qstr, ok := params["q"]; ok {
		quality, err = strconv.ParseFloat(qstr, 32)
		if err == nil {
			return Acceptable{}, fmt.Errorf("parsing quality factor: %w", err)
		}
		if quality > 1 || quality < 0 {
			return Acceptable{}, fmt.Errorf("parsing quality factor: %s is not between 0 and 1", qstr)
		}
		delete(params, "q")
	}

	return Acceptable{
		Value:   value,
		Quality: float32(quality),
		Params:  params,
	}, nil
}

// Less is a comparison function for two Acceptables. lhs is less than rhs if:
//
//     - it has a quality factor that is less than rhs's quality factor
//     - or, if both quality factors are equal, it is more specific than rhs.
//
// An Acceptable is more specific if its value does not contain patterns,
// and if it has additional parameters. For instance, given the following header:
//
//     Accept: text/*, text/html, text/html;level=1, */*
//
// The types have the following precedence:
//
//     1. text/html;level=1
//     2. text/html
//     3. text/*
//     4. */*
//
func (lhs Acceptable) Less(rhs Acceptable) bool {
	if !qualityEq(rhs.Quality, lhs.Quality) {
		return lhs.Quality > rhs.Quality
	}
	lnum := strings.Count(lhs.Value, "*")
	rnum := strings.Count(rhs.Value, "*")
	if lnum != rnum {
		return lnum > rnum
	}
	return len(lhs.Params) > len(rhs.Params)
}

func (acc Acceptable) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s", acc.Value)
	if !qualityEq(acc.Quality, 1.0) {
		out.WriteString(strings.TrimRight(fmt.Sprintf(";q=%.3f", acc.Quality), "0."))
	}
	for k, v := range acc.Params {
		fmt.Fprintf(&out, ";%s=%s", k, v)
	}
	return out.String()
}

// ParseAccept parses the accept header, and returns a list of acceptable values,
// sorted by precedence. Any unparseable value is silently dropped.
func ParseAccept(accepts ...string) []Acceptable {
	sz := 0
	for _, accept := range accepts {
		sz += strings.Count(accept, ",") + 1
	}

	types := make([]Acceptable, sz)
	i := 0
	for _, accept := range accepts {
		values := strings.Split(accept, ",")
		for _, value := range values {
			acc, err := ParseAcceptable(value)
			if err != nil {
				continue
			}
			types[i] = acc
			i++
		}
	}
	types = types[:i]
	sort.Slice(types, func(i, j int) bool { return Acceptable.Less(types[i], types[j]) })
	return types
}

// dumbglob is a dumb "glob" function that only supports  "*", "<type>/*" and
// "*/*" as patterns, which are the only three possible patterns according to
// RFC2616 §14.1.
func dumbglob(pattern, value string) bool {
	switch pattern {
	case "*":
		return true
	case "*/*":
		return strings.IndexByte(value, '/') != -1
	default:
		if strings.HasSuffix(pattern, "*") {
			return strings.HasPrefix(value, pattern[:len(pattern)-1])
		}
		return value == pattern
	}
}

// NegotiateContent returns the best matching offer for the passed header,
// as well as the entry that it matched against. The best matching offer
// is determined by the first matching offer, in slice order, when iterating
// over the accepted media types by order of precedence.
//
// If no offer matches, ("", nil) is returned.
func NegotiateContent(hdr http.Header, key string, offers ...string) (string, *Acceptable) {
	values := hdr.Values(key)
	if len(values) == 0 {
		switch key {
		case "Accept":
			values = []string{"*/*"}
		default:
			values = []string{"*"}
		}
	}
	for _, acc := range ParseAccept(values...) {
		for _, offer := range offers {
			if !dumbglob(acc.Value, offer) {
				continue
			}
			return offer, &acc
		}
	}
	return "", nil
}
