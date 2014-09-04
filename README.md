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


License
=======
template is licensed under the BSD