package template

import (
	"errors"
	_ "fmt"
	"io"
	"path/filepath"
	"reflect"

	html "html/template"
	text "text/template"
	"text/template/parse"
)

const (
	ErrEmpty  = "template: Templates() is empty"
	ErrNever  = "template: Never"
	ErrChRoot = "template: Can not chroot"
)

// FuncsMap is global variable for Funcs(FuncsMap) on New/NewHtml function
var FuncsMap = map[string]interface{}{}

type Template struct {
	tpl       interface{}
	typ       string
	parsed    bool
	w         io.Writer
	data      interface{}
	baseDir   string
	firstName string
	chroot    string
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

func New(name string) *Template {
	tmpl := &Template{typ: "text"}
	tmpl.tpl = text.New(name).Funcs(FuncsMap)
	return tmpl
}

func NewHtml(name string) *Template {
	tmpl := &Template{typ: "html"}
	tmpl.tpl = html.New(name).Funcs(FuncsMap)
	return tmpl
}
func (t *Template) clone() *Template {
	return &Template{
		tpl:       t.tpl,
		typ:       t.typ,
		parsed:    t.parsed,
		w:         t.w,
		data:      t.data,
		baseDir:   t.baseDir,
		firstName: t.firstName,
	}
}
func (t *Template) setFirstName(name string) {
	if t.firstName == "" {
		t.firstName = filepath.Base(name)
	}
}

func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error) {
	var err error
	switch t.typ {
	case "text":
		_, err = t.tpl.(*text.Template).AddParseTree(name, tree)
	case "html":
		_, err = t.tpl.(*html.Template).AddParseTree(name, tree)
	}
	if err == nil {
		t.setFirstName(name)
		return t, nil
	}
	return nil, err
}
func (t *Template) Clone() (*Template, error) {
	var err error
	tmpl := t.clone()
	switch tmpl.typ {
	case "text":
		_, err = tmpl.tpl.(*text.Template).Clone()
	case "html":
		_, err = tmpl.tpl.(*html.Template).Clone()
	}
	if err == nil {
		return tmpl, nil
	}
	return nil, err
}
func (t *Template) Delims(left, right string) *Template {
	switch t.typ {
	case "text":
		t.tpl.(*text.Template).Delims(left, right)
	case "html":
		t.tpl.(*html.Template).Delims(left, right)
	}
	return t
}

func (t *Template) Execute(wr io.Writer, data interface{}) (err error) {
	t.data = data

	switch t.typ {
	case "text":
		tpl := t.tpl.(*text.Template)
		if t.parsed {
			err = tpl.Execute(wr, data)
		} else {
			if t.firstName == "" {
				err = errors.New(ErrEmpty)
			} else {
				err = tpl.ExecuteTemplate(wr, t.firstName, data)
			}
		}
	case "html":
		tpl := t.tpl.(*html.Template)
		if t.parsed {
			err = tpl.Execute(wr, data)
		} else {
			if t.firstName == "" {
				err = errors.New(ErrEmpty)
			} else {
				err = tpl.ExecuteTemplate(wr, t.firstName, data)
			}
		}
	}
	return
}

func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) (err error) {
	t.data = data
	switch t.typ {
	case "text":
		err = t.tpl.(*text.Template).ExecuteTemplate(wr, name, data)
	case "html":
		err = t.tpl.(*html.Template).ExecuteTemplate(wr, name, data)
	}
	return
}

func (t *Template) Funcs(funcMap map[string]interface{}) *Template {
	switch t.typ {
	case "text":
		t.tpl.(*text.Template).Funcs(funcMap)
	case "html":
		t.tpl.(*html.Template).Funcs(funcMap)
	}
	return t
}
func (t *Template) Lookup(name string) *Template {
	tmpl := &Template{typ: t.typ}
	switch t.typ {
	case "text":
		tmpl.tpl = t.tpl.(*text.Template).Lookup(name)
	case "html":
		tmpl.tpl = t.tpl.(*html.Template).Lookup(name)
	}
	if tmpl.tpl == nil {
		return nil
	}
	return tmpl
}
func (t *Template) Name() string {
	switch t.typ {
	case "text":
		return t.tpl.(*text.Template).Name()
	case "html":
		return t.tpl.(*html.Template).Name()
	}
	return ""
}
func (t *Template) New(name string) *Template {
	tmpl := &Template{typ: t.typ}
	switch t.typ {
	case "text":
		tmpl.tpl = t.tpl.(*text.Template).New(name)
	case "html":
		tmpl.tpl = t.tpl.(*html.Template).New(name)
	}
	return tmpl
}
func (t *Template) Parse(src string) (*Template, error) {
	var err error
	switch t.typ {
	case "text":
		_, err = t.tpl.(*text.Template).Parse(src)
	case "html":
		_, err = t.tpl.(*html.Template).Parse(src)
	}
	if err != nil {
		return nil, err
	}
	t.parsed = true
	return t, nil
}

func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	var err error
	Abs(t.baseDir, filenames)
	if len(filenames) == 0 {
		return nil, errors.New(ErrEmpty)
	}
	if !t.canChRoot(filenames) {
		return nil, errors.New(ErrChRoot)
	}
	switch t.typ {
	case "text":
		_, err = t.tpl.(*text.Template).ParseFiles(filenames...)
	case "html":
		_, err = t.tpl.(*html.Template).ParseFiles(filenames...)
	}
	if err != nil {
		return nil, err
	}
	t.setFirstName(filenames[0])
	return t, nil
}

func (t *Template) ParseGlob(pattern string) (*Template, error) {
	tmpl := &Template{typ: t.typ}
	var err error
	switch t.typ {
	case "text":
		tmpl.tpl, err = t.tpl.(*text.Template).ParseGlob(pattern)
	case "html":
		tmpl.tpl, err = t.tpl.(*html.Template).ParseGlob(pattern)
	}
	if err == nil {
		return tmpl, nil
	}
	return nil, err
}

func (t *Template) Templates() (tmpls []*Template) {
	switch t.typ {
	case "text":
		mp := t.tpl.(*text.Template).Templates()
		tmpls = make([]*Template, len(mp))
		for _, tpl := range mp {
			tmpl := t.clone()
			tmpl.tpl = tpl
			tmpls = append(tmpls, tmpl)
		}
	case "html":
		mp := t.tpl.(*html.Template).Templates()
		tmpls = make([]*Template, len(mp))
		for _, tpl := range mp {
			tmpl := t.clone()
			tmpl.tpl = tpl
			tmpls = append(tmpls, tmpl)
		}
	}
	return
}

// TextTemplate returns raw *text/template.Template ify
func (t *Template) TextTemplate() *text.Template {
	if t.typ == "text" {
		return t.tpl.(*text.Template)
	}
	return nil
}

// HtmlTemplate returns raw *html/template.Template ify
func (t *Template) HtmlTemplate() *html.Template {
	if t.typ == "html" {
		return t.tpl.(*html.Template)
	}
	return nil
}

// Builtin enable builtin funcs
func (t *Template) Builtin() *Template {
	return t.Funcs(t.BuiltinFuncs())
}

// BuiltinFuncs returns builtin funcs
func (t *Template) BuiltinFuncs() map[string]interface{} {
	// builtin funcMap
	mp := map[string]interface{}{
		// data returns origin data
		"data": func() interface{} {
			return t.data
		},
		// getBaseDir returns the base directory path for template files
		"getBaseDir": func() string {
			return t.baseDir
		},
		// setBaseDir setting the base directory path for template files
		"setBaseDir": func(baseDir string) string {
			t.BaseDir(baseDir)
			return ""
		},
		// exists to determine whether the key of data exists
		"exists": func(data interface{}, key string) bool {
			v := reflect.Indirect(reflect.ValueOf(data))
			switch v.Kind() {
			case reflect.Map:
				return v.MapIndex(reflect.ValueOf(key)).IsValid()
			case reflect.Struct:
				return v.FieldByName(key).IsValid() || v.MethodByName(key).IsValid()
			}
			return false
		},
	}
	switch t.typ {
	case "text":
		// tpl return origin Template
		mp["tpl"] = func() *text.Template {
			return t.tpl.(*text.Template)
		}
		// import wrapper for ParseFiles
		mp["import"] = func(filenames ...string) (string, error) {
			Abs(t.baseDir, filenames)
			if len(filenames) == 0 {
				return "", errors.New(ErrEmpty)
			}
			if !t.canChRoot(filenames) {
				return "", errors.New(ErrChRoot)
			}
			_, err := t.tpl.(*text.Template).ParseFiles(filenames...)
			return "", err
		}
	case "html":
		mp["tpl"] = func() *html.Template {
			return t.tpl.(*html.Template)
		}
		mp["import"] = func(filenames ...string) (string, error) {
			Abs(t.baseDir, filenames)
			if len(filenames) == 0 {
				return "", errors.New(ErrEmpty)
			}
			if !t.canChRoot(filenames) {
				return "", errors.New(ErrChRoot)
			}
			_, err := t.tpl.(*html.Template).ParseFiles(filenames...)
			return "", err
		}
	}
	return mp
}
func (t *Template) canChRoot(filenames []string) bool {
	if t.chroot != "" {
		for _, filename := range filenames {
			if len(filename) < len(t.chroot) || filename[0:len(t.chroot)] != t.chroot {
				return false
			}
		}
	}
	return true
}

// BaseDir setting the base directory path for template files
func (t *Template) BaseDir(dir string) *Template {
	if dir == "" {
		t.baseDir = dir
		return t
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return t
	}
	t.baseDir = dir + string(filepath.Separator)
	return t
}

//ChRoot setting the root directory path for template files
func (t *Template) ChRoot(dir string) *Template {
	if dir == "" {
		t.chroot = dir
		return t
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return t
	}
	t.chroot = dir + string(filepath.Separator)
	return t
}

// Abs transform filenames to absolute path base basedir
func Abs(basedir string, filenames []string) {
	for i, filename := range filenames {
		if !filepath.IsAbs(filename) {
			filenames[i] = basedir + filename
		}
	}
	return
}
