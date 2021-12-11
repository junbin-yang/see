





## See

# 简介
参考gin开发的高性能轻量级web框架。作为一个练手项目。以尽量精简的代码实现，最大程度兼容和优化gin使用习惯，添加一些新特性。经过不断优化，效率已略优于gin。

# 新特性
🚩 访问日志，类似nginx的access.log，支持rotate和过期自动删除。

🚩 更直接的自定义参数验证器，在数据绑定时传入作为可选参数传入即可。

🚩 新增CopyRawData()方法支持*http.Request读出后重新写入。

🚩 新增ShouldBindForm()和BindForm()方法，绑定form表单数据。

# Benchmarks

性能对比：

<table>
   <tr>
      <td>Benchmark name</td>
      <td>（1）</td>
      <td>（2）</td>
      <td>（3）</td>
      <td>（4）</td>
   </tr>
   <tr>
      <td>BenchmarkGin_Param        </td>
      <td>16790644</td>
      <td>        70.58 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_Param        </td>
      <td>18993596</td>
      <td>        64.59 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_Param5       </td>
      <td>8565904</td>
      <td>       141.1 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_Param5       </td>
      <td>10001650</td>
      <td>       118.3 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_Param20      </td>
      <td>3168889</td>
      <td>       386.2 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_Param20      </td>
      <td>4002708</td>
      <td>       313.8 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_ParamWrite   </td>
      <td>8736094</td>
      <td>       130.0 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_ParamWrite   </td>
      <td>12166742</td>
      <td>       100.2 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GithubStatic </td>
      <td>13365705</td>
      <td>        92.69 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GithubStatic </td>
      <td>13938786</td>
      <td>        87.14 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GithubParam  </td>
      <td>7792669</td>
      <td>       153.0 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GithubParam  </td>
      <td>8141613</td>
      <td>       147.5 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GithubAll    </td>
      <td>36870</td>
      <td>     33976 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GithubAll    </td>
      <td>42343</td>
      <td>     28180 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GPlusStatic  </td>
      <td>19503213</td>
      <td>        61.31 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GPlusStatic  </td>
      <td>18894066</td>
      <td>        64.43 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GPlusParam   </td>
      <td>12364990</td>
      <td>       101.1 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GPlusParam   </td>
      <td>12786169</td>
      <td>        94.64 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GPlus2Params </td>
      <td>9736572</td>
      <td>       124.7 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GPlus2Params </td>
      <td>9554286</td>
      <td>       124.9 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_GPlusAll     </td>
      <td>866647</td>
      <td>      1489 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_GPlusAll     </td>
      <td>929016</td>
      <td>      1335 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_ParseStatic  </td>
      <td>18856722</td>
      <td>        66.26 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_ParseStatic  </td>
      <td>17481632</td>
      <td>        66.95 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_ParseParam   </td>
      <td>16241710</td>
      <td>        75.11 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_ParseParam   </td>
      <td>17228764</td>
      <td>        70.17 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_Parse2Params </td>
      <td>12969364</td>
      <td>        94.29 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_Parse2Params </td>
      <td>12308853</td>
      <td>        90.03 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_ParseAll     </td>
      <td>514867</td>
      <td>      2368 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_ParseAll     </td>
      <td>552721</td>
      <td>      2151 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkGin_StaticAll    </td>
      <td>54718</td>
      <td>     22569 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
   <tr>
      <td>BenchmarkSee_StaticAll    </td>
      <td>56215</td>
      <td>     20705 ns/op</td>
      <td>       0 B/op</td>
      <td>       0 allocs/op</td>
   </tr>
</table>

- (1): Total Repetitions achieved in constant time, higher means more confident result
- (2): Single Repetition Duration (ns/op), lower is better
- (3): Heap Memory (B/op), lower is better
- (4): Average Allocations per Repetition (allocs/op), lower is better

# 快速入门
运行这段代码并在浏览器中访问 [http://localhost:8080](http://localhost:8080/)

```go
package main

import "github.com/junbin-yang/see"

func main() {
	r := see.Default()
	r.GET("/ping", func(c *see.Context) {
		c.JSON(200, see.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
```

# 使用GET, POST, PUT等

```go
package main

import "github.com/junbin-yang/see"

func main() {
	// 使用默认中间件（logger and recovery）创建一个路由器
	router := see.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)

	// 默认启动的是 8080端口，也可以自己定义启动端口
	router.Run()
	// router.Run(":3000") for a hard coded port
}
```

# 获取路径中的参数

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()
	
	router.GET("/user/:name", func(c *see.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	
	router.GET("/user/name/zhangsan", func(c *see.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.Run(":8080")
}
```

# 获取Get参数

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()

	// 匹配的url格式:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *see.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

# 获取Post参数

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()

	router.POST("/form_post", func(c *see.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous") // 此方法可以设置默认值

		c.JSON(200, see.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}
```

# Get + Post 混合

```
示例：
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
```

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()

	router.POST("/post", func(c *see.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		
		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	router.Run(":8080")
}
```

```
结果：id: 1234; page: 1; name: manu; message: this_is_great
```

# 上传文件

单文件上传

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()
	// 给表单限制上传大小 (默认 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *see.Context) {
		// 单文件
		file, _ := c.FormFile("file")

		// 上传文件到指定的路径
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	router.Run(":8080")
}
```

测试：

```
curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
  -H "Content-Type: multipart/form-data"
```

多文件上传

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()
	// 给表单限制上传大小 (默认 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *see.Context) {
		// 多文件
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			
			// 上传文件到指定的路径
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":8080")
}
```

测试：

```
curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
```

# 路由分组

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}
```

# 无中间件启动

使用

```go
r := see.New()
```

代替

```go
// 默认启动方式，包含 Logger、Recovery 中间件（Logger信息在stdout输出）
r := see.Default()
```

# 日志模式启动 🟢

```go
// 传入日志文件前缀、日志保存目录、是否rotate、日志保存天数
// 启动包含Logger、Recovery 中间件
r := see.Enable("seeAccess", "/var/log", true, 7)
```

# 使用中间件 🟢

```go
package main

import "github.com/junbin-yang/see"

func main() {
	// 创建一个不包含中间件的路由器
	r := see.New()

	// 全局中间件
	// 设置log参数
	//see.SetLoggerConfig("seeAccess", "/var/log/see", true, 7)
	
	// 使用 Logger 中间件
	r.Use(see.Logger())

	// 使用 Recovery 中间件
	r.Use(see.Recovery())

	// 路由添加中间件，可以添加任意多个
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// 路由组中添加中间件
	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

# 模型绑定和验证 🟢

🔴注意：`BindForm`和`ShouldBindForm`是See新增的方法。

若要将请求主体绑定到结构体中，请使用模型绑定，目前支持JSON、XML、YAML和标准表单值(foo=bar&boo=baz)的绑定。

See使用 [go-playground/validator.v10](https://github.com/go-playground/validator 验证参数。

需要在绑定的字段上设置tag，比如，绑定格式为json，需要这样设置 `json:"fieldname"` 。

此外，还提供了两套绑定方法：

- Must bind
- - Methods - `Bind`,`BindUri`, `BindJSON`, `BindXML`, `BindQuery`, `BindYAML`,`BindForm`
- - Behavior - 这些方法如果存在绑定错误，响应状态代码会被设置为400，请求头`Content-Type`被设置为`text/plain; charset=utf-8`。注意，如果你试图在此之后设置响应代码，将会发出一个警告，如果你希望更好地控制行为，请使用`ShouldBind`相关的方法
- Should bind
- - Methods - `ShouldBind`,`ShouldBindUri`, `ShouldBindJSON`, `ShouldBindXML`, `ShouldBindQuery`, `ShouldBindYAML`,`ShouldBindForm`
- - Behavior - 这些方法如果存在绑定错误，则返回错误，开发人员可以正确处理请求和错误。

你还可以给字段指定特定规则的修饰符，如果一个字段用`validate:"required"`修饰，并且在绑定时该字段的值为空，那么将返回一个错误。参数验证这一部分可以直接使用validator.v10。

```go
package main

import "github.com/junbin-yang/see"

// 绑定为json
type Login struct {
	User     string `json:"user" xml:"user"  validate:"required"`
	Password string `json:"password" xml:"password" validate:"required,max=20,min=6"`
	Code     string `json:"code" xml:"code" validate:"required,len=6"`
}

func main() {
	router := see.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *see.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, see.H{"error": err.Error()})
			return
		}
		
		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, see.H{"status": "unauthorized"})
			return
		} 
		
		c.JSON(http.StatusOK, see.H{"status": "you are logged in"})
	})

	// Example for binding XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>user</user>
	//		<password>123</password>
	//	</root>)
	router.POST("/loginXML", func(c *see.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, see.H{"error": err.Error()})
			return
		}
		
		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, see.H{"status": "unauthorized"})
			return
		} 
		
		c.JSON(http.StatusOK, see.H{"status": "you are logged in"})
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
```

**请求示例：**

```
$ curl -v -X POST \
  http://localhost:8080/loginJSON \
  -d '{ "user": "manu" }'
> POST /loginJSON HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.51.0
> Accept: */*
> content-type: application/json
> Content-Length: 18
>
* upload completely sent off: 18 out of 18 bytes
< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
< Date: Fri, 04 Aug 2017 03:51:31 GMT
< Content-Length: 100
<
{"error":"Key: 'Login.Password' Error:Field validation for 'Password' failed on the 'required' tag"}
```

**跳过验证：**

当使用上面的curl命令运行上面的示例时，返回错误，因为示例中`Password`字段使用了`validate:"required"`，如果我们使用`validate:"-"`，那么它就不会报错。

**验证规则：**

```
required ：必填
email：验证字符串是email格式；例：“email”
url：这将验证字符串值包含有效的网址;例：“url”
max：字符串最大长度；例：“max=20”
min:字符串最小长度；例：“min=6”
excludesall:不能包含特殊字符；例：“excludesall=0x2C”//注意这里用十六进制表示。
len：字符长度必须等于n，或者数组、切片、map的len值为n，即包含的项目数；例：“len=6”
eq：数字等于n，或者或者数组、切片、map的len值为n，即包含的项目数；例：“eq=6”
ne：数字不等于n，或者或者数组、切片、map的len值不等于为n，即包含的项目数不为n，其和eq相反；例：“ne=6”
gt：数字大于n，或者或者数组、切片、map的len值大于n，即包含的项目数大于n；例：“gt=6”
gte：数字大于或等于n，或者或者数组、切片、map的len值大于或等于n，即包含的项目数大于或等于n；例：“gte=6”
lt：数字小于n，或者或者数组、切片、map的len值小于n，即包含的项目数小于n；例：“lt=6”
lte：数字小于或等于n，或者或者数组、切片、map的len值小于或等于n，即包含的项目数小于或等于n；例：“lte=6”
```

**跨字段验证：**

如想实现比较输入密码和确认密码是否一致等类似场景

```
eqfield=Field: 必须等于 Field 的值；
nefield=Field: 必须不等于 Field 的值；
gtfield=Field: 必须大于 Field 的值；
gtefield=Field: 必须大于等于 Field 的值；
ltfield=Field: 必须小于 Field 的值；
ltefield=Field: 必须小于等于 Field 的值；
eqcsfield=Other.Field: 必须等于 struct Other 中 Field 的值；
necsfield=Other.Field: 必须不等于 struct Other 中 Field 的值；
gtcsfield=Other.Field: 必须大于 struct Other 中 Field 的值；
gtecsfield=Other.Field: 必须大于等于 struct Other 中 Field 的值；
ltcsfield=Other.Field: 必须小于 struct Other 中 Field 的值；
ltecsfield=Other.Field: 必须小于等于 struct Other 中 Field 的值；
```

示例：

验证Passwd和Repasswd值是否相等

```go
type UserReg struct {
	Passwd 		string `json:"passwd" 	validate:"required,max=20,min=6"`
 	Repasswd 	string `json:"repasswd" validate:"required,max=20,min=6,eqfield=Passwd"`
}
```

# 自定义验证器 🟢

简化了这一部分的使用方式，直接在绑定模型时传入自定义的验证方法即可。

```go
package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/junbin-yang/see"
)

type User struct {
	MyName string `json:"name" validate:"required,CK"`
	Phone  string `json:"phone" validate:"required"`
}

func Cust(fl validator.FieldLevel) bool {
	return fl.Field().String() == "admin"
}

func main() {
	r := see.Default()
	r.POST("/post", func(c *see.Context) {
		var userinfo User
		err := c.BindJSON(&userinfo,map[string]validator.Func{"CK": Cust})
		if err != nil {
			return
		}
		c.JSON(200, see.H{"myName": userinfo.MyName})
	})
	
	route.Run(":8085")
}
```

# 只绑定Get参数

`ShouldBindQuery` 函数只绑定Get参数，不绑定post数据。

```go
package main

import "github.com/junbin-yang/see"

type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	route := see.Default()
	route.Get("/", startPage)
	route.Run(":8085")
}

func startPage(c *see.Context) {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}
```

# 绑定Get参数或者Post参数

```go
package main

import "github.com/junbin-yang/see"

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	route := see.Default()
	route.GET("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *see.Context) {
	var person Person
	// 如果是Get，那么接收不到请求中的Post的数据
	// 如果是Post, 首先判断 `content-type` 的类型, 然后使用对应的绑定器获取数据.
	if c.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
	}
	c.String(200, "Success")
}
```

# 绑定uri

```go
package main

import "github.com/junbin-yang/see"

type Person struct {
	ID string `uri:"id" validate:"required,uuid"`
	Name string `uri:"name" validate:"required"`
}

func main() {
	route := see.Default()
	route.GET("/:name/:id", func(c *see.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, see.H{"msg": err})
			return
		}
		c.JSON(200, see.H{"name": person.Name, "uuid": person.ID})
	})
	route.Run(":8088")
}
```

测试用例：

```
$ curl -v localhost:8088/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3
$ curl -v localhost:8088/thinkerou/not-uuid
```

# 绑定Post参数 🟢

```go
package main

import "github.com/junbin-yang/see"

type LoginForm struct {
	User     string `form:"user" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func main() {
	router := see.Default()
	router.POST("/login", func(c *see.Context) {
		var form LoginForm
		// c.ShouldBind()
		if c.ShouldBindForm(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(200, see.H{"status": "you are logged in"})
			} else {
				c.JSON(401, see.H{"status": "unauthorized"})
			}
		}
	})
	router.Run(":8080")
}
```

测试用例：

```
$ curl -v localhost:8088/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3
$ curl -v localhost:8088/thinkerou/not-uuid
```

# 输出格式XML、JSON、YAML 🟢

```go
package main

import "github.com/junbin-yang/see"

func main() {
	r := see.Default()

	r.GET("/someJSON", func(c *see.Context) {
		c.JSON(http.StatusOK, see.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someXML", func(c *see.Context) {
		c.XML(http.StatusOK, see.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *see.Context) {
		c.YAML(http.StatusOK, see.H{"message": "hey", "status": http.StatusOK})
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

**AsciiJSON**

使用AsciiJSON将使特殊字符编码

```go
package main

import "github.com/junbin-yang/see"

func main() {
	r := see.Default()

	r.GET("/someJSON", func(c *see.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 将输出: {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

**PureJSON**

通常情况下，JSON会将特殊的HTML字符替换为对应的unicode字符，比如`<`替换为`\u003c`，如果想原样输出html，则使用PureJSON，这个特性在Go 1.6及以下版本中无法使用。

```go
package main

import "github.com/junbin-yang/see"

func main() {
	r := see.Default()
	
	// Serves unicode entities
	r.GET("/json", func(c *see.Context) {
		c.JSON(200, see.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	
	// Serves literal characters
	r.GET("/purejson", func(c *see.Context) {
		c.PureJSON(200, see.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	
	// listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

# 设置静态文件路径

访问静态文件需要先设置路径

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()
	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.Run(":8080")
}
```

# 返回第三方获取的数据

```go
package main

import "github.com/junbin-yang/see"

func main() {
	router := see.Default()
	router.GET("/someDataFromReader", func(c *see.Context) {
		response, err := http.Get("https://wx4.sinaimg.cn/large/008aq1Apgy1gwo3onis8rj30mh0cn74z.jpg")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.jpg"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	router.Run(":8080")
}
```

# 重定向

发布HTTP重定向很容易，支持内部和外部链接

```go
r.GET("/test", func(c *see.Context) {
	//c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com/")
	c.Redirect(http.StatusMovedPermanently, "/json")
})
```

# 自定义中间件

```go
package main

import "github.com/junbin-yang/see"

func Logger() see.HandlerFunc {
	return func(c *see.Context) {
		t := time.Now()
		
		// Set example variable
		c.Set("example", "12345")

		// before request
		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.StatusCode
		log.Println(status)
	}
}

func main() {
	r := see.New()
	r.Use(Logger())
	r.GET("/test", func(c *see.Context) {
		example := c.MustGet("example").(string)
		// it would print: "12345"
		log.Println(example)
	})
	r.Run(":8080")
}
```

# 中间件中使用Goroutines 🟢

在中间件或处理程序中启动新的Goroutines时，gin的做法是c.Copy()拷贝一个完整的上下文只读副本。see不支持Copy()函数。常用字段已经存储到上下文中，直接使用即可。

多次读取Body数据的问题：gin使用GetRawData()方法读取*http.Request.Body数据，后续的处理流程里将无法通过参数绑定解析到数据。（一般是在写访问日志中间件时记录请求的数据使用），see新增CopyRawData()方法，将数据读出后重新写回上下文。

```go
package main

import "github.com/junbin-yang/see"

func Logger() HandlerFunc {
	return func(c *Context) {
		// 开始时间
		startTime := time.Now()
		// 请求数据
		body, _ := c.CopyRawData()	

		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Method
		// 请求路由
		reqUri := c.RequestURI
		// 请求参数
		userAgent := c.Req.Header["User-Agent"]
		// 状态码
		statusCode := c.StatusCode
		// 请求IP
		clientIP := c.RemoteAddr
		
		// ...
	}
}

func main() {
	r := see.Default()
	r.Use(Logger())

	r.GET("/long_async", func(c *see.Context) {
		go func() {
			time.Sleep(5 * time.Second)
			log.Println("Done! in path " + c.Path)
		}()
	})

	r.Run(":8080")
}
```

# 自定义HTTP配置

直接像这样使用`http.ListenAndServe()`

```go
func main() {
	router := see.Default()
	http.ListenAndServe(":8080", router)
}
```

或者

```go
func main() {
	router := see.Default()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
```

# 优雅重启或停止

想要优雅地重启或停止你的Web服务器，使用http.Server内置的[Shutdown()](https://golang.org/pkg/net/http/#Server.Shutdown)方法进行优雅关闭

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/junbin-yang/see"
)

func main() {
	router := see.Default()
	router.GET("/", func(c *see.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome See Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
```

# 自定义路由日志的格式

默认的路由日志是这样的：

```
2021-11-30 10:13:09,514 Register Route: POST /post/18
2021-11-30 10:13:09,514 Register Route: GET /v1/index
2021-11-30 10:13:09,514 Register Route: GET /v2/index
2021-11-30 10:13:09,514 Register Route: GET /json
```

如果你想以给定的格式记录这些信息（例如 JSON，键值对或其他格式），你可以使用`see.DebugPrintRouteFunc`来定义格式，在下面的示例中，我们使用标准日志包记录路由日志，你可以使用其他适合你需求的日志工具

```go
package main

import "github.com/junbin-yang/see"

func main() {
	r := see.Default()
	see.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string) {
		log.Printf("endpoint %v %v %v\n", httpMethod, absolutePath, handlerName)
	}

	r.POST("/foo", func(c *see.Context) {
		c.JSON(http.StatusOK, "foo")
	})

	r.GET("/bar", func(c *see.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	r.GET("/status", func(c *see.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// Listen and Server in http://0.0.0.0:8080
	r.Run()
}
```
