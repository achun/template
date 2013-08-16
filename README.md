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


License
=======
template is licensed under the BSD