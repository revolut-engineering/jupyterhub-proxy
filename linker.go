package main

import (
	"regexp"
	"strings"
)

// A Linker is used to correct the links found in the response body
// It will add the prefix when a regex match is found
type Linker struct {
	quoted_reg   *regexp.Regexp
	hue_base_url *regexp.Regexp
	prefix_link  *regexp.Regexp
}

func new_linker(paths string) Linker {
	quoted := `([='"])(` + paths + ")"
	hue_base := `window.HUE_BASE_URL\s*[+]\s*'/hue'`
	prefix_link := `var prefix = '/hue';`

	linker := Linker{
		quoted_reg:   regexp.MustCompile(quoted),
		hue_base_url: regexp.MustCompile(hue_base),
		prefix_link:  regexp.MustCompile(prefix_link),
	}

	return linker
}

func (l *Linker) replace(file []byte, service_prefix string) []byte {
	prefix := "$1" + strings.TrimSuffix(service_prefix, "/") + "$2"
	replaced := l.quoted_reg.ReplaceAll(file, []byte(prefix))
	replaced = l.hue_base_url.ReplaceAll(replaced, []byte("window.HUE_BASE_URL"))
	replaced = l.prefix_link.ReplaceAll(replaced, []byte("var prefix = '';"))

	return replaced
}
