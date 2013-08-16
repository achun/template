package template_test

import (
	"github.com/achun/template"
	"os"
	"testing"
)

func TestErrEmpty(T *testing.T) {
	err := template.New("").Execute(os.Stdout, "")
	if err == nil || err.Error() != template.ErrEmpty {
		T.Fatal("wanted:", template.ErrEmpty, "but:", err)
	}
	err = template.NewHtml("").Execute(os.Stdout, "")
	if err == nil || err.Error() != template.ErrEmpty {
		T.Fatal("wanted:", template.ErrEmpty, "but:", err)
	}
}

func TestBuiltin(T *testing.T) {
	tpl, err := template.New("").Builtin().Parse(`
		{{print data}}
		{{getBaseDir}} -> {{setBaseDir "testdata"}}{{getBaseDir}}
		{{if exists data "foo"|not}}not exists "foo"{{end}}
		{{import "tmpl1.tmpl" "tmpl2.tmpl"}}
		{{template "tmpl1.tmpl"}}{{"\n"}}`)
	if err == nil {
		err = tpl.Execute(os.Stdout, "import so easy")
	}
	if err != nil {
		T.Fatal(err)
	}
}
