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
	tpl, err := template.New("name").Builtin().Parse(`
		{{print data}}
		{{getBaseDir}} -> {{setBaseDir "testdata"}}{{getBaseDir}}
		{{exists data "foo"}}
		{{import "file1.tmpl" "file2.tmpl" "tmpl1.tmpl" "tmpl2.tmpl"}}
		{{"\n"}}`)
	if err == nil {
		err = tpl.Execute(os.Stdout, "so easy")
	}
	if err != nil {
		T.Fatal(err)
	}
}
