// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"html"
	"net/url"
	"strings"
)

var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&quot;",
)

func escapeHTML(html string) string {
	return htmlEscaper.Replace(html)
}

func unescapeString(str string) string {
	str = html.UnescapeString(str) // FIXME: 此处应该用内部的实体转义方式
	runes := []rune(str)

	var retRunes []rune
	for i := 0; i < len(runes); i++ {
		if isBackslashEscape(runes, i) {
			retRunes = retRunes[:len(retRunes)-1]
		}
		retRunes = append(retRunes, runes[i])
	}

	return string(retRunes)
}

func isBackslashEscape(runes []rune, pos int) bool {
	if !isASCIIPunct(runes[pos]) {
		return false
	}

	backslashes := 0
	for i := pos - 1; 0 <= i; i-- {
		if '\\' != runes[i] {
			break
		}

		backslashes++
	}

	return 0 != backslashes%2
}

func isASCIIPunct(c rune) bool {
	return (0x21 <= c && 0x2F >= c) || (0x3A <= c && 0x40 >= c) || (0x5B <= c && 0x60 >= c) || (0x7B <= c && 0x7E >= c)
}

func encodeDestination(destination string) (ret string) {
	destination = decodeDestination(destination)

	parts := strings.SplitN(destination, ":", 2)
	var scheme string
	remains := destination
	if 1 < len(parts) {
		scheme = parts[0]
		remains = parts[1]
	}

	index := strings.Index(remains, "?")
	var query string
	path := remains
	if 0 <= index {
		query = remains[index+1:]
		queries := strings.Split(query, "&")
		query = ""
		length := len(queries)
		for i, q := range queries {
			kv := strings.Split(q, "=")
			if 1 < len(kv) {
				query += url.QueryEscape(kv[0]) + "=" + url.QueryEscape(kv[1])
			} else {
				query += url.QueryEscape(kv[0])
			}
			if i < length-1 {
				query += "&"
			}
		}

		path = remains[:index]
	}

	parts = strings.Split(path, "/")
	path = ""
	length := len(parts)
	for i, part := range parts {
		unescaped := url.PathEscape(part)
		path += unescaped
		if i < length-1 {
			path += "/"
		}
	}

	if "" == scheme {
		ret = path
	} else {
		ret = scheme + ":" + path
	}
	if "" != query {
		ret += "?" + query
	}

	ret = strings.ReplaceAll(ret, "%2A", "*")
	ret = strings.ReplaceAll(ret, "%29", ")")
	ret = strings.ReplaceAll(ret, "%28", "(")
	ret = strings.ReplaceAll(ret, "%23", "#")
	ret = strings.ReplaceAll(ret, "%2C", ",")

	return
}

func decodeDestination(destination string) (ret string) {
	parts := strings.SplitN(destination, ":", 2)
	var scheme string
	remains := destination
	if 1 < len(parts) {
		scheme = parts[0]
		remains = parts[1]
	}

	parts = strings.Split(remains, "/")
	remains = ""
	length := len(parts)
	for i, part := range parts {
		unescaped, err := url.QueryUnescape(part)
		if nil != err {
			unescaped = part
		}
		remains += unescaped
		if i < length-1 {
			remains += "/"
		}
	}

	if "" == scheme {
		ret = remains
	} else {
		ret = scheme + ":" + remains
	}

	return
}
