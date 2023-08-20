package bind

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

// copy from gin
// go get github.com/go-playground/validator/v10
// 集成第三方验证
type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *defaultValidator) Validate(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.Validate(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.Validate(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj any) error {
	v.lazyInit()
	return v.validate.Struct(obj)
}

/*
once 只初始化一次
SetTagName 注意设置了什么名字就要求tag就是什么名字
*/
func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("bind")
	})
}
