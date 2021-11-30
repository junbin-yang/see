package verify

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

const (
	ZH = "zh"
	EN = "en"
)

type bind struct {
	data           interface{}
	e              map[string][]error
	validationfunc map[string]validator.Func
}

func NewValidator(obj interface{}) *bind {
	return &bind{
		data:           obj,
		e:              make(map[string][]error),
		validationfunc: make(map[string]validator.Func),
	}
}

func (this *bind) Verify() bool {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()

	//注册自定义函数
	if len(this.validationfunc) > 0 {
		for name, fun := range this.validationfunc {
			validate.RegisterValidation(name, fun)
		}
	}

	//验证器注册翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		this.e["en"] = append(this.e["en"], err)
		return false
	}

	err = validate.Struct(this.data)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			this.e["en"] = append(this.e["en"], e)
			this.e["zh"] = append(this.e["zh"], errors.New(e.Translate(trans)))
		}
		return false
	}
	return true
}

func (this *bind) GetErrors(k string) []error {
	return this.e[k]
}

func (this *bind) RegisterValidation(name string, h validator.Func) {
	this.validationfunc[name] = h
}
