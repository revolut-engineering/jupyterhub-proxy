package main

import "testing"

type TestCase struct {
	file string
	want string
}

func TestLinker(t *testing.T) {
	paths := `/static/|/debug/`

	linker := new_linker(paths)
	prefix := "/user/user@comp/"

	cases := []TestCase{
		{"'/static/file/something.css'", "'/user/user@comp/static/file/something.css'"},
		{"'/debug/something.js'", "'/user/user@comp/debug/something.js'"},
		{
			`sourceMappingURL=/static/desktop/js/bundles/hue/defaultVendors~hue-cdd60c62-bundle-5ac0133f57e0`,
			`sourceMappingURL=/user/user@comp/static/desktop/js/bundles/hue/defaultVendors~hue-cdd60c62-bundle-5ac0133f57e0`,
		},
		{
			`<link href="/static/desktop/css/roboto.895233d7bf84.css" rel="stylesheet">
			<link href="/static/desktop/ext/css/font-awesome.min.bf0c425cdb73.css" rel="stylesheet">`,
			`<link href="/user/user@comp/static/desktop/css/roboto.895233d7bf84.css" rel="stylesheet">
			<link href="/user/user@comp/static/desktop/ext/css/font-awesome.min.bf0c425cdb73.css" rel="stylesheet">`,
		},
	}

	for _, tc := range cases {
		res := linker.replace([]byte(tc.file), prefix)
		if string(res) != tc.want {
			t.Fatalf("replace(%s, %s) = %s want %s", tc.file, prefix, res, tc.want)
		}
	}
}
