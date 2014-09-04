package template

import (
	"bytes"
	"github.com/achun/testing-want"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"text/template/parse"
)

var pwd = initPWD()

func initPWD() string {
	pwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	return filepath.ToSlash(pwd)
}

func walkTree(node parse.Node) {

	switch n := node.(type) {
	case *parse.ListNode:
		for _, node := range n.Nodes {
			walkTree(node)
		}
		return
	}
}

var uriTests = []string{
	".",
	"b",
	"../a",
	"../a/b",
	"../a/b/",
	"../a/b/./c/../../.././a",
	`$`,
	`$/.`,
	`C:\a/..\a/b`,
	`C:/a\\..\a/b/\`,
	`#/a/..\a/b`,
	`#/a/..\a/b/\`,
	`\\/\#/a/..\a/b/\`,
	`#/a\\/b///c\..\../\.././a`,
}

var fixtures = map[string]interface{}{
	"title": `>title`,
	"body":  `/admin/body.html`,
	"js":    `/admin/js.tmpl`,
	"href":  ">>>",
	"name":  "admin",
}

func TestCleanURI(T *testing.T) {
	wt := want.T(T)
	wd, _ := os.Getwd()
	wt.Equal(clearSlash(wd), pwd)
	for _, s := range uriTests {

		n := path.Clean(strings.Replace(s, "\\", "/", -1))
		if n[0] == '.' {
			n = ""
		}
		wt.Equal(cleanURI(s), n, s)
	}
}

func TestTemplate(T *testing.T) {
	var buf bytes.Buffer

	wt := want.T(T)
	t, err := New("./fixtures/base/layout.html")
	wt.Nil(err)
	wt.Equal(t.RootDir(), pwd+`/fixtures/base`)
	wt.Equal(t.Name(), pwd+"/fixtures/base/layout.html")
	t.Walk(pwd+`/fixtures/base`, ".html.tmpl")

	wt.Nil(t.Execute(&buf, fixtures))
	wt.Equal(buf.String(),
		`<html>
<head>
<meta charset="UTF-8">
<script>"admin"</script>
</head>
<body>
<h1><a href="%3e%3e%3e">admin</a></h1>
</body>
<script>"foot"</script>
</html>`)
}
