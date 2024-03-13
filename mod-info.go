package main

import (
	"github.com/icholy/replace"
	"golang.org/x/text/transform"
	"io"
	"regexp"
	"strings"
)

type ModInfo map[string]string

func (m ModInfo) ReplaceInPath(s string) (string, error) {
	m2 := make(ModInfo)
	for k, v := range m {
		m2[k] = v
	}
	m2["modgroup"] = strings.ReplaceAll(m2["modgroup"], ".", "/")
	s2, _, err := transform.String(m2.transformer(), s)
	return s2, err
}

func (m ModInfo) ReplaceInStream(r io.Reader) io.Reader {
	return transform.NewReader(r, m.transformer())
}

func (m ModInfo) transformer() transform.Transformer {
	return replace.RegexpStringFunc(regexp.MustCompile("%%[a-z_]+%%"), func(s string) string {
		trim := strings.TrimFunc(s, func(r rune) bool {
			return r == '%'
		})
		if v, ok := m[trim]; ok {
			return v
		}
		return "%%unresolved:" + trim + "%%"
	})
}
