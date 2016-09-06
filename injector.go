package sector

import "reflect"

//-----------------------------------------------------------------------------

// NewInjector creates an Injector
func NewInjector(fac Factory) Injector {
	res := new(injector)
	res.fac = fac
	return res
}

//-----------------------------------------------------------------------------

type injector struct {
	fac Factory
}

func (x *injector) Inject(ptr interface{}) {
	valOf := reflect.ValueOf(ptr)
	elem := valOf.Elem()
	if elem.Kind() == reflect.Interface {
		return
	}

	fieldCount := elem.NumField()

	typeOf := valOf.Type().Elem()

	for i := 0; i < fieldCount; i++ {
		i := i
		field := elem.Field(i)

		if !field.IsValid() || !field.CanSet() {
			continue
		}

		typeField := typeOf.Field(i)

		tag := typeField.Tag.Get("inject")
		recurse := false
		if tag == "*" {
			recurse = true
		} else if tag != "+" {
			continue
		}
		// if !(typeField.Tag == "inject" || typeField.Tag.Get("inject") != "") {
		// 	continue
		// }

		var v interface{}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				back := reflect.New(field.Type().Elem())
				v = back.Interface()
				field.Set(back)
			} else {
				v = field.Interface()
			}
		} else {
			v = field.Addr().Interface()
		}

		if x.fac.Fill(v) {
			if !recurse {
				continue
			}
		}

		x.Inject(v)
	}
}

//-----------------------------------------------------------------------------

// Injector /
type Injector interface {
	// Inject , must accept only pointer
	Inject(interface{})
}
