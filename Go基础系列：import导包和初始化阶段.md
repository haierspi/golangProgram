# Go基础系列：import导包和初始化阶段

## import导入包

### 搜索路径

import用于导入包：

```
import (
    "fmt"
    "net/http"
    "mypkg"
)

```

编译器会根据上面指定的相对路径去搜索包然后导入，这个相对路径是从GOROOT或GOPATH(workspace)下的src下开始搜索的。

假如go的安装目录为`/usr/local/go`，也就是说`GOROOT=/usr/local/go`，而GOPATH环境变量`GOPATH=~/mycode:~/mylib`，那么要搜索`net/http`包的时候，将按照**如下顺序**进行搜索：

```
/usr/local/go/srcnet/http
~/mycode/src/net/http
~/mylib/src/net/http

```

以下是`go install`搜索不到mypkg包时的一个报错信息：

```
can't load package: package mypkg: cannot find package "mypkg" in any of:
        /usr/lib/go-1.6/src/mypkg (from $GOROOT)
        /golang/src/mypkg (from $GOPATH)

```

也就是说，go总是先从`GOROOT`出先搜索，再从`GOPATH`列出的路径顺序中搜索，只要一搜索到合适的包就理解停止。当搜索完了仍搜索不到包时，将报错。

包导入后，就可以使用这个包中的属性。使用`包名.属性`的方式即可。例如，调用fmt包中的Println函数`fmt.Println`。

### 包导入的过程

![](https://img2018.cnblogs.com/blog/733013/201810/733013-20181023224911978-1960747966.png)

首先从main包开始，如果main包中有import语句，则会导入这些包，如果要导入的这些包又有要导入的包，则继续先导入所依赖的包。重复的包只会导入一次，就像很多包都要导入fmt包一样，但它只会导入一次。

每个被导入的包在导入之后，都会先将包的可导出函数(大写字母开头)、包变量、包常量等声明并初始化完成，然后如果这个包中定义了init()函数，则自动调用init()函数。init()函数调用完成后，才回到导入者所在的包。同理，这个导入者所在包也一样的处理逻辑，声明并初始化包变量、包常量等，再调用init()函数(如果有的话)，依次类推，直到回到main包，main包也将初始化包常量、包变量、函数，然后调用init()函数，调用完init()后，调用main函数，于是开始进入主程序的执行逻辑。

### 别名导入和特殊的导入方法

当要导入的包重名时会如何？例如`network/convert`包用于转换从网络上读取的数据，`file/convert`包用于转换从文件中读取的数据，如果要同时导入它们，当引用的时候指定`convert.FUNC()`，这个convert到底是哪个包？

可以为导入的包添加一个名称属性，为包设置一个别名。例如，除了导入标准库的fmt包外，自己还定义了一个mypkg/fmt包，那么可以如下导入：

```go
package main

import (
    "fmt"
    myfmt "mypkg/fmt"
)

func main() {
    fmt.Println()
    myfmt.myfunc()   // 使用别名进行访问
}

```

如果不想在访问包属性的时候加上包名，则import导入的时候，可以为其设置特殊的别名：点(.)。

```
import (
    . "fmt"
)

func main() {
    Println()    // 无需包名，直接访问Println
}

```

这时要访问fmt中的属性，**必须**不能使用包名fmt。

go要求import导入的包必须在后续中使用，否则会报错。如果想要避免这个错误，可以在包的前面加上下划线：

```
import (
    "fmt"
    _ "net/http"
    "mypkg"
)

```

这样在当前包中就无需使用`net/http`包。其实这也是为包进行命名，只不过命名为"\_"，而这个符号又正好表示丢弃赋值结果，使得这成为一个匿名包。

> **下划线(\_)**
> 在go中，下划线出现的频率非常高，它被称为blank identifier，可以用于赋值时丢弃值，可以用于保留import时的包，还可以用于丢弃函数的返回值。详细内容可参见官方手册：[https://golang.org/doc/effective\_go.html#blank](https://golang.org/doc/effective_go.html#blank)

导入而不使用看上去有点多此一举，但并非如此。因为导入匿名包仅仅表示无法再访问其内的属性。但导入这个匿名包的时候，会进行一些初始化操作(例如init()函数)，如果这个初始化操作会影响当前包，那么这个匿名导入就是有意义的。

## 远程包

现在通过分布式版本控制系统进行代码共享是一种大趋势。go集成了从gti上获取远程代码的能力。

例如：

```
$ go get github.com/golang/example

```

在import语句中也可以使用，首先从GOPATH中搜索路径，显然这是一个URL路径，于是调用go get进行fetch，然后导入。

```
import (
    "fmt"
    "github.com/golang/example"
)

```

当需要从git上获取代码的时候，将调用`go get`工具自动进行fetch、build、install。如果workspace中已经有这个包，那么将只进行最后的install阶段，如果没有这个包，将保存到GOPATH的第一个路径中，并build、install。

go get是递归的，所以可以直接fetch整个代码树。

## 常量和变量的初始化

Go中的常量在编译期间就会创建好，即使是那些定义为函数的本地常量也如此。常量只允许是数值、字符(runes)、字符串或布尔值。

由于编译期间的限制，定义它们的表达式必须是编译器可评估的常量表达式(constant expression)。例如，`1<<3`是一个常量表达式，而`math.Sin(math.Pi/4)`则不是常量表达式，因为涉及了函数math.Sin()的调用过程，而函数调用是在运行期间进行的。

变量的初始化和常量的初始化差不多，但初始化的变量允许是"需要在执行期间计算的一般表达式"。例如：

```
var (
    home   = os.Getenv("HOME")
    user   = os.Getenv("USER")
    gopath = os.Getenv("GOPATH")
)

```

## init()函数

Go中除了保留了main()函数，还保留了一个init()函数，这两个函数都不能有任何参数和返回值。它们都是在特定的时候自动调用的，无需我们手动去执行。

还是这张图：

![](https://img2018.cnblogs.com/blog/733013/201810/733013-20181023224911978-1960747966.png)

每个包中都可以定义init函数，甚至可以定义多个，但建议每个包只定义一个。每次导入包的时候，在导入完成后，且变量、常量等声明并初始化完成后，将会调用这个包中的init()函数。

对于main包，如果main包也定义了init()，那么它会在main()函数之前执行。当main包中的init()执行完之后，就会立即执行main()函数，然后进入主程序。

所以，init()经常用来初始化环境、安装包或其他需要在程序启动之前先执行的操作。如果import导入包的时候，发现前面命名为下划线`_`了，一般就说明所导入的这个包有init()函数，且导入的这个包除了init()函数外，没有其它作用。

**转载请注明出处：[https://www.cnblogs.com/f\-ck\-need\-u/p/9847554.html](https://www.cnblogs.com/f-ck-need-u/p/9847554.html)**