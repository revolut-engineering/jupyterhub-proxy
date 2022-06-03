package main

import (
	"regexp"
	"strings"
)

// A Linker is used to correct the links found in the response body
// It will add the prefix when a regex match is found
type Linker struct {
	quoted_reg *regexp.Regexp
}

func new_linker(paths string) Linker {
	quoted := `([='"])(` + paths + ")"

	linker := Linker{
		quoted_reg: regexp.MustCompile(quoted),
	}

	return linker
}

func (l *Linker) replace(file []byte, service_prefix string) []byte {
	prefix := "$1" + strings.TrimSuffix(service_prefix, "/") + "$2"
	replaced := l.quoted_reg.ReplaceAll(file, []byte(prefix))

	return replaced
}
