package main

import (
	"regexp"
	"strings"
)

// A Linker is used to correct the links found in the response body
// It will add the prefix when a regex match is found
type Linker struct {
	quoted_reg *regexp.Regexp
	js_reg     *regexp.Regexp
}

func new_linker(paths string) Linker {
	quoted := `(['"])(` + paths + ")"
	js_string := `(\\)(` + paths + ")"

	linker := Linker{
		quoted_reg: regexp.MustCompile(quoted),
		js_reg:     regexp.MustCompile(js_string),
	}

	return linker
}

func (l *Linker) replace(file []byte, service_prefix string) []byte {
	prefix := "$1" + strings.TrimSuffix(service_prefix, "/") + "$2"
	replaced := l.quoted_reg.ReplaceAll(file, []byte(prefix))

	new_prefix := strings.Replace(service_prefix, "/", `\/`, -1)
	prefix = strings.TrimSuffix(new_prefix, "/") + "$2"
	replaced = l.js_reg.ReplaceAll(replaced, []byte(prefix))

	return replaced
}
