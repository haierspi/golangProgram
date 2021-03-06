# 一. Golint介绍

Golint is a linter for Go source code.

1.  Golint 是一个源码检测工具用于检测代码规范
2.  Golint 不同于gofmt, Gofmt用于代码格式化

Golint会对代码做以下几个方面检查

1.  package注释 必须按照 “Package xxx 开头”
2.  package命名 不能有大写字母、下划线等特殊字符
3.  struct、interface等注释 必须按照指定格式开头
4.  struct、interface等命名
5.  变量注释、命名
6.  函数注释、命名
7.  各种语法规范校验等

# 二. Golint安装

1.  go get \-u [github.com/golang/lint/golint](http://github.com/golang/lint/golint)
2.  ls $GOPATH/bin (可以发现已经有golint可执行文件)

# 三. Golint使用

golint检测代码有2种方式

1.  golint file
2.  golint directory

golint校验常见的问题如下所示

1.  `don't use ALL_CAPS in Go names; use CamelCase`
    不能使用下划线命名法，使用驼峰命名法
2.  `exported function Xxx should have comment or be unexported`
    外部可见程序结构体、变量、函数都需要注释
3.  `var statJsonByte should be statJSONByte`
    `var taskId should be taskID`
    通用名词要求大写
    iD/Id \-> ID
    Http \-> HTTP
    Json \-> JSON
    Url \-> URL
    Ip \-> IP
    Sql \-> SQL
4.  `don't use an underscore in package name`
    `don't use MixedCaps in package name; xxXxx should be xxxxx`
    包命名统一小写不使用驼峰和下划线
5.  `comment on exported type Repo should be of the form "Repo ..." (with optional leading article)`
    注释第一个单词要求是注释程序主体的名称，注释可选不是必须的
6.  `type name will be used as user.UserModel by other packages, and that stutters; consider calling this Model`
    外部可见程序实体不建议再加包名前缀
7.  `if block ends with a return statement, so drop this else and outdent its block`
    if语句包含return时，后续代码不能包含在else里面
8.  `should replace errors.New(fmt.Sprintf(...)) with fmt.Errorf(...)`
    errors.New(fmt.Sprintf(…)) 建议写成 fmt.Errorf(…)
9.  `receiver name should be a reflection of its identity; don't use generic names such as "this" or "self"`
    receiver名称不能为this或self
10.  `error var SampleError should have name of the form ErrSample`
    错误变量命名需以 Err/err 开头
11.  `should replace num += 1 with num++`
    `should replace num -= 1 with num--`
    a+=1应该改成a++，a\-=1应该改成a–

# 四. Goland配置golint

1.  新增tool: Goland \-> Tools \-> External Tools 新建一个tool 配置如下
    ![](https://github.com/chenguolin/chenguolin.github.io/blob/master/data/image/goland-add-golint-tool.png?raw=true)

2.  新增快捷键: Goland \-> Tools \-> Keymap \-> External Tools \-> External Tools \-> golint 右键新增快捷键
    ![](https://github.com/chenguolin/chenguolin.github.io/blob/master/data/image/goland-add-golint-shortcut.png?raw=true)

3.  使用: 打开Go文件，然后使用快捷键就可以进行代码检测了
    有golint的输出日志说明存在代码不规范的地方，需要进行修改
    ![](https://github.com/chenguolin/chenguolin.github.io/blob/master/data/image/goland-add-golint-check.png?raw=true)

# 五. gitlab提交限制

为了保证项目代码规范，我们可以在gitlab上做一层约束限制，当代码提交到gitlab的时候先做golint校验，校验不通过则不让提交代码。

我们可以为Go项目创建gitlab CI流程，通过`.gitlab-ci.yml`配置CI流程会自动使用`govet`进行代码静态检查、`gofmt`进行代码格式化检查、`golint`进行代码规范检查、`gotest`进行单元测试

例如[go\-common项目](https://github.com/chenguolin/golang/blob/master/go-common) .gitlab.yml文件如下，相关的脚本可以查看[scripts](https://github.com/chenguolin/golang/tree/master/go-common/scripts)

```
# gitlab CI/CD pipeline配置文件
# 默认使用本地定制过的带有golint的golang镜像
image: golang:custom

stages:
    - test

before_script:
    - mkdir -p /go/src/gitlab.local.com/golang
    - ln -s `pwd` /go/src/gitlab.local.com/golang/go-common && cd /go/src/gitlab.local.com/golang/go-common

# test stage
# job 1 test go vet
job_govet:
    stage: test
    script:
        - bash ./scripts/ci-govet-check.sh
    tags:
        - dev
# job 2 test go fmt
job_gofmt:
    stage: test
    script:
        - bash ./scripts/ci-gofmt-check.sh
    tags:
        - dev
# job 3 test go lint
job_golint:
    stage: test
    script:
        - bash ./scripts/ci-golint-check.sh
    tags:
        - dev
# job 4 test go unit test
job_unit:
    stage: test
    script:
        - bash ./scripts/ci-gotest-check.sh
    tags:
        - dev
```