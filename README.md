Template
========

基于官方 `text/template` 和 `html/template` 的模板引擎.
Template 通过几种惯用方式组合, 为模板提供简洁的使用方式.

特性
====

* 模板名仿效 URI 格式, 使用全路径名称命名.
* 依据命名扩展名自动判断 TEXT 或 HTML 风格.
* 模板源码可使用相对路径名指示目标模板.
* 引入 RootDir 限制模板文件根目录.
* 内置 import 函数支持变量名表示模板名.

使用
====

以源码 fixtures/base 目录下的文件为例:

```
\---base
    |   foot.html     <script>"foot"</script>
    |   layout.html   单独列出
    |
    \---admin
            body.html <h1><a href="{{.href}}">{{.name}}</a></h1>
            js.tmpl   <script>"admin"</script>
```

layout.html 内容, 注意 import 支持变量, 支持目标模板名采用相对路径:

```html
<html>
<head>
<meta charset="UTF-8">
{{import .js}}
</head>
<body>
{{import .body .}}
</body>
{{template "foot.html"}}
</html>
```

GoLang 代码:

```go
package main

import (
    "github.com/achun/template"
    "os"
)

var data = map[string]interface{}{
    "title": `>title`,
    "body":  `/admin/body.html`,
    "js":    `/admin/js.tmpl`,
    "href":  ">>>",
    "name":  "admin",
}

func main() {
    pwd, _ := os.Getwd()
    t, err := template.New("./fixtures/base/layout.html")
    t.Walk(pwd+`/fixtures/base`, ".html.tmpl")
    t.Execute(os.Stdout, data)
}
```

输出:

```html
<html>
<head>
<meta charset="UTF-8">
<script>"admin"</script>
</head>
<body>
<h1><a href="%3e%3e%3e">admin</a></h1>
</body>
<script>"foot"</script>
</html>
```

内部实现
========

通过重写 *parse.Tree 中的 `template` 用 `import` 函数替换. 并且重新计算目标路径为绝对路径. 最终形成的 `import` 定义为:

```go
func(from, name string, data ...interface{}) (template.HTML, error)
```

* from 为发起调用的模板名称, form 是自动计算出的, 使用者不能定义.
* name 为模板模板名称
* data 为用户数据

name 支持变量, 此变量有可能采用相对路径, form 为计算绝对路径提供了参照.
而 rootdir 保证所有的相对路径都可以计算出绝对路径.

License
=======
template is licensed under the BSD