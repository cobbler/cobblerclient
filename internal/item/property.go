package item

import "reflect"

type Property[T interface{}] interface {
	Get() T
	Set(value T)
}

type InheritableProperty[T interface{}] interface {
	Property[T]
	GetCurrentRawType() reflect.Type
	IsInherited() (bool, error)
	GetRaw() interface{}
	SetRaw(interface{})
}
