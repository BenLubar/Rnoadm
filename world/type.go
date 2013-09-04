package world

import (
	"encoding/gob"
	"math/big"
	"reflect"
	"time"
)

var registeredIdentifierObject = make(map[string]reflect.Type)
var registeredObjectIdentifier = make(map[reflect.Type]string)

func Register(identifier string, obj ObjectLike) {
	if identifier == "" {
		panic("attempt to register an empty identifier")
	}
	if obj == nil {
		panic("attempt to register a nil object type")
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if _, ok := registeredIdentifierObject[identifier]; ok {
		panic("duplicate registration for object identifier " + identifier)
	}
	if _, ok := registeredObjectIdentifier[t]; ok {
		panic("duplicate registration for object type " + t.Name())
	}
	registeredObjectIdentifier[t] = identifier
	registeredIdentifierObject[identifier] = t
}

func getObjectTypeIdentifier(obj ObjectLike) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	identifier := registeredObjectIdentifier[t]
	if identifier == "" {
		panic("unregistered object type: " + t.Name())
	}
	return identifier
}

func getObjectByIdentifier(identifier string) ObjectLike {
	t := registeredIdentifierObject[identifier]
	if t == nil {
		panic("unregistered object identifier: " + identifier)
	}
	return reflect.New(t).Interface().(ObjectLike)
}

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(time.Time{})
	gob.Register(&big.Int{})
}
