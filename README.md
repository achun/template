# 重构

建立这个项目最初的原因是想要扩展一些官方`template`包所不支持的功能.
2013-08-14 日作者发现, 种种想法, 原来在官方包之下居然全部都能实现. 这完全是使用方法问题.这里有一个和闭包结合的[例子](http://play.golang.org/p/Fil_Vi2ZhU)可以说明这种技巧.

因此此项目原先fork官方包进行源码修改的方法是错误的, 项目被重构.

重构借鉴了官方 `html/template` 包的方法.


# Template

`Template` 是对官方 `text/template` 和 `html/template` 包进行的再次包装.
除了重新包装所有的接口函数外, 新增了一些函数, 还通过闭包增加了一些模板函数.


新增函数
--------

* (*Template) TextTemplate() // 返回原 `*text/template.Template` 对象, 如果是的话.
* (*Template) HtmlTemplate() // 返回原 `*html/template.Template` 对象, 如果是的话.
* (*Template) BuiltinFuncs() // 返回内建模板函数.
* (*Template) Builtin() // 使内建模板函数在模板执行期生效.
* (*Template) BaseDir(dir string) // 为要解析的模板文件设置 base directory.
* Abs(basedir string, filenames []string) // 以 basedir 为基础, 转变 filenames 到绝对路径.

内建模板函数
------------

* data() interface{} 返回执行模板时所传入的原始 data.
* getBaseDir() string 返回 base directory.
* setBaseDir(baseDir string) string 对 `BaseDir` 的包装, 总是返回"".
* exists(data interface{}, key string) bool 返回 data 是否包含 key 属性, 支持map和struct
* import(filenames ...string) string 动态 ParseFiles, 总是返回"", 注意这不是执行

# 使用

```go
import "github.com/achun/template"
```

您应该发现`Template`的接口尽量和官方包一致, 所以使用方法和官方包是完全一致的.

在 `Parse` 或 `ParseFiles` 之前调用 `Builtin()` 内建模板函数才会生效.

几句代码演示使用方法

```go
tpl := template.New("") // 是的不需要名字, tpl 其实是个集合, 官方的也是集合. 自动匹配到第一个有效的模板很容易实现.
tpl.Bultin() // 先执行这一句那些内建模板函数才会生效, 如果您需要那些内建模板函数的话.
content := "get.tmpl" // 这里示意定义一个 content 模板文件的名字
// 用闭包的方法加入您自己的FuncMap
tpl.Funcs(map[string]interface{}{
	"request": func() *http.Request { //
		return r // 这个就是是客户端发起 *http.Request, 提供这个只是一个演示, 说明这种方法完全可用
	},
	"content": func() string { // 用函数的方法运行输出 get.tmpl
		err := tpl.ExecuteTemplate(w, content, dat) // w 就是http.ResponseWriter了
		if err != nil {
			return err.Error()
		}
		return ""
	},
})
// 事实上我们只需要这两个文件既可, 其余的文件可以在模板中用import加载
tpl.ParseFiles("layout.html", "install/get.tmpl")
// 万事俱备, 执行
tpl.Execute(w, "your data")
```

layout.html的样子

```html
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <!-- 这样写只是表示 request 函数是可用的 -->
    <title>{{request.URL.Path}}</title>
    <link rel="stylesheet" href="/css/pure.css">
    <link rel="stylesheet" href="/css/default.css">
</head>

<body>
<!-- 哈, content 用函数的方法完成是不是更优雅呢 -->
{{content}} 
</body>
<script src="/js/jquery-2.0.3.min.js"></script>
<script src="/js/md5.js"></script>
<script src="/js/common.js"></script>

</html>
```

get.tmpl

```html
<h1>这里您随意吧<h1>
<!-- import 另一个模板 -->
{{import "foo.tmpl"}}
<!--
执行仍然是标准语法
data是个内建函数, 返回的是原始传入的模板数据
模板嵌套多次时,可以用data获取最原始的传入值
-->
{{template "foo.tmpl" data}}
```

License
=======
template is licensed under the BSD