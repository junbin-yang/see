package see

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/junbin-yang/golib/bytesconv"
	"github.com/junbin-yang/golib/json"
	"github.com/junbin-yang/see/verify"
	"gopkg.in/yaml.v2"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type H map[string]interface{}

//使map类型支持xml解析
func (thisMap H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Space: "", Local: "map"}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range thisMap {
		elem := xml.StartElement{Name: xml.Name{Space: "", Local: key}, Attr: []xml.Attr{}}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (this Params) Get(key string) string {
	for _, param := range this {
		if param.Key == key {
			return param.Value
		}
	}
	return ""
}
func (this Params) GetMapByOneKey() map[string]string {
	out := map[string]string{}
	for _, param := range this {
		out[param.Key] = param.Value
	}
	return out
}

// 请求上下文信息
type Context struct {
	// 请求信息
	Path       string
	RequestURI string
	Method     string
	RemoteAddr string
	Params     Params
	// 响应信息
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	index    int
	// 处理过程中设置在上下文中的数据
	Keys sync.Map

	Req    *http.Request
	Writer http.ResponseWriter
	engine *Engine
}

// 初始化上下文实例
func (c *Context) SetContext(w http.ResponseWriter, r *http.Request) {
	c.Req = r
	c.Writer = w
	c.Path = r.URL.Path
	c.RequestURI = r.RequestURI
	c.Method = r.Method
	c.RemoteAddr = r.RemoteAddr
}

func (c *Context) Reset() {
	c.index = -1
}

func (c *Context) CopyRawData() ([]byte, error) {
	body, err := c.GetRawData()
	c.Req.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, err
}

func (c *Context) GetRawData() ([]byte, error) {
	if c.Req.Body != nil {
		return io.ReadAll(c.Req.Body)
	}
	return []byte{}, nil
}

func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

// 获取Url上的动态参数
// router.GET("/user/:id", func(c *see.Context) {
//		//a GET request to /user/john
//      id := c.Param("id") // id == "john"
// })
func (c *Context) Param(key string) string {
	return c.Params.Get(key)
}

// 添加上下文参数
func (c *Context) AddParam(key, value string) {
	c.Params = append(c.Params, Param{key, value})
}

// 获取Url上的参数,?x=y
func (c *Context) Query(name string) string {
	value, _ := c.GetQuery(name)
	return value
}

func (c *Context) DefaultQuery(name, defaultValue string) string {
	if value, ok := c.GetQuery(name); ok {
		return value
	}
	return defaultValue
}

func (c *Context) initQuery() url.Values {
	if c.Req != nil {
		return c.Req.URL.Query()
	}
	return url.Values{}
}

func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	data := c.initQuery()
	values, ok = data[key]
	return
}

func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// 返回磁盘上某个文件
func (c *Context) File(filepath string) {
	file, err := os.ReadFile(filepath)
	if err != nil || file == nil {
		c.StatusCode = http.StatusNotFound
		http.Error(c.Writer, "File Not Found", http.StatusNotFound)
		return
	}
	c.StatusCode = http.StatusOK
	http.ServeFile(c.Writer, c.Req, filepath)
}

// 获取表单参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value := c.PostForm(key); value != "" {
		return value
	}
	return defaultValue
}

// 单文件上传
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Req.MultipartForm == nil {
		if err := c.Req.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Req.ParseMultipartForm(c.engine.MaxMultipartMemory)
	return c.Req.MultipartForm, err
}

// 保存文件
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置头信息
// 注意在 WriteHeader() 后调用 Header().Set 是不会生效的
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 响应内容
// 字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain;charset=utf-8")
	c.Status(code)
	_, _ = c.Writer.Write(bytesconv.StringToBytes(fmt.Sprintf(format, values...)))
}

// yaml
func (c *Context) YAML(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/x-yaml; charset=utf-8")
	c.Status(code)
	bytes, err := yaml.Marshal(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	_, err = c.Writer.Write(bytes)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

// xml
func (c *Context) XML(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/xml; charset=utf-8")
	c.Status(code)
	encoder := xml.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

// json
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json;charset=utf-8")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

// Pure json
func (c *Context) PureJSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json;charset=utf-8")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

// Ascii json
func (c *Context) AsciiJSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	jstr, err := json.Marshal(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	var buffer bytes.Buffer
	for _, r := range bytesconv.BytesToString(jstr) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	_, _ = c.Writer.Write(buffer.Bytes())
}

// 返回第三方获取的数据
func (c *Context) DataFromReader(code int, contentLength int64, contentType string, body io.Reader, extraHeaders map[string]string) {
	for k, v := range extraHeaders {
		c.SetHeader(k, v)
	}
	if contentLength >= 0 {
		c.SetHeader("Content-Length", strconv.FormatInt(contentLength, 10))
	}
	if contentType != "" {
		c.SetHeader("Content-Type", contentType)
	}

	c.Status(code)
	_, _ = io.Copy(c.Writer, body)
}

// 返回字节流数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// 重定向
func (c *Context) Redirect(code int, location string) {
	if (code < http.StatusMultipleChoices || code > http.StatusPermanentRedirect) && code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot be redirected using status code %d", code))
	}
	c.StatusCode = code
	http.Redirect(c.Writer, c.Req, location, code)
}

// html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html;charset=utf-8")
	c.Status(code)
	_, _ = c.Writer.Write(bytesconv.StringToBytes(html))
}

// 执行下一个函数
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// 中间件的中断开关
func (c *Context) Abort() {
	// 直接让计数下角标等于中间件数组长度
	c.index = len(c.handlers)
}

// 中断并抛出信息给用户，内部使用
func (c *Context) fail(code int, err string) {
	c.Abort()
	c.String(code, err)
}

// 为这个上下文存储一个新的键/值对。
func (c *Context) Set(key string, value interface{}) {
	c.Keys.Store(key, value)
}

// 获取某个key的值，不存在返回(nil,false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys.Load(key)
	return
}

// 获取某个key的值,不存在则抛出异常
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// 获取其他具体类型的值
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

// 数据绑定
func (c *Context) Bind(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBind(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBind(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if c.Method == http.MethodGet {
		//return Form
		return c.ShouldBindForm(obj, validationfunc...)
	}

	switch c.ContentType() {
	case "application/json":
		return c.ShouldBindJSON(obj, validationfunc...)
	case "application/xml", "text/xml":
		return c.ShouldBindXML(obj, validationfunc...)
	case "application/x-protobuf":
		//return ProtoBuf
		return errors.New("未支持的类型")
	case "application/x-msgpack", "application/msgpack":
		//return MsgPack
		return errors.New("未支持的类型")
	case "application/x-yaml":
		return c.ShouldBindYAML(obj, validationfunc...)
	case "multipart/form-data":
		//return FormMultipart
		return c.ShouldBindForm(obj, validationfunc...)
	default: //"application/x-www-form-urlencoded"
		//return Form
		return c.ShouldBindForm(obj, validationfunc...)
	}
}

func (c *Context) BindForm(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBindForm(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBindForm(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if err := c.Req.ParseForm(); err != nil {
		return err
	}
	if _, err := c.MultipartForm(); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}

	form := map[string]string{}
	for key, _ := range c.Req.Form { //c.Req.Form 是一个 map[string][]string 类型
		form[key] = c.Req.Form.Get(key)
	}

	str, err := json.ObjectToJson(form)
	if err != nil {
		return err
	}
	json.JsonToObject(str, obj, "form")
	return c.validate(obj, validationfunc...)
}

func (c *Context) BindQuery(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBindQuery(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBindQuery(obj interface{}, validationfunc ...map[string]validator.Func) error {
	m := map[string]string{}
	for key, values := range c.initQuery() {
		if len(values) > 0 {
			m[key] = values[0]
		}
	}

	str, err := json.ObjectToJson(m)
	if err != nil {
		return err
	}
	json.JsonToObject(str, obj, "form")
	return c.validate(obj, validationfunc...)
}

func (c *Context) BindYAML(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBindYAML(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBindYAML(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if c.Req == nil || c.Req.Body == nil {
		return errors.New("invalid request")
	}

	decoder := yaml.NewDecoder(c.Req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return c.validate(obj, validationfunc...)
}

func (c *Context) BindXML(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBindXML(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBindXML(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if c.Req == nil || c.Req.Body == nil {
		return errors.New("invalid request")
	}

	decoder := xml.NewDecoder(c.Req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return c.validate(obj, validationfunc...)
}

func (c *Context) BindJSON(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBindJSON(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBindJSON(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if c.Req == nil || c.Req.Body == nil {
		return errors.New("invalid request")
	}

	decoder := json.NewDecoder(c.Req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return c.validate(obj, validationfunc...)
}

func (c *Context) BindUri(obj interface{}, validationfunc ...map[string]validator.Func) error {
	if e := c.ShouldBindUri(obj, validationfunc...); e != nil {
		c.fail(400, e.Error())
		return e
	}
	return nil
}

func (c *Context) ShouldBindUri(obj interface{}, validationfunc ...map[string]validator.Func) error {
	str, err := json.ObjectToJson(c.Params.GetMapByOneKey())
	if err != nil {
		return err
	}
	json.JsonToObject(str, obj, "uri")
	return c.validate(obj, validationfunc...)
}

func (c *Context) validate(obj interface{}, validationfunc ...map[string]validator.Func) error {
	v := verify.NewValidator(obj)
	if len(validationfunc) > 0 {
		for fn, f := range validationfunc[0] {
			v.RegisterValidation(fn, f)
		}
	}
	if v.Verify() == false {
		return v.GetErrors(verify.ZH)[0]
	}
	return nil
}

//获取请求的Content-Type
func (c *Context) ContentType() string {
	ct := c.GetHeader("Content-Type")
	for i, char := range ct {
		if char == ' ' || char == ';' {
			return ct[:i]
		}
	}
	return ct
}
