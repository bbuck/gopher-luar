package luar

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

func getStructMetaTable(L *lua.LState) lua.LValue {
	key := registryPrefix + "struct"
	table := L.G.Registry.RawGetH(lua.LString(key))
	if table != lua.LNil {
		return table
	}
	newTable := L.NewTable()
	newTable.RawSetH(lua.LString("__index"), L.NewFunction(structIndex))
	newTable.RawSetH(lua.LString("__newindex"), L.NewFunction(structNewIndex))
	L.G.Registry.RawSetH(lua.LString(key), newTable)
	return newTable
}

func structIndex(L *lua.LState) int {
	ud := L.CheckUserData(1)
	name := L.CheckString(2)

	value := reflect.ValueOf(ud.Value)
	if value.Kind() == reflect.Ptr {
		if method := value.MethodByName(name); method.IsValid() {
			L.Push(getLuaFuncWrapper(L, method))
			return 1
		}
		value = value.Elem()
	}

	if method := value.MethodByName(name); method.IsValid() {
		L.Push(getLuaFuncWrapper(L, method))
		return 1
	}

	field := value.FieldByName(name)
	if field.IsValid() {
		if val := New(L, field.Interface()); val != nil {
			L.Push(val)
			return 1
		}
	}

	return 0
}

func structNewIndex(L *lua.LState) int {
	ud := L.CheckUserData(1)
	name := L.CheckString(2)
	lValue := L.Get(3)

	value := reflect.ValueOf(ud.Value)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	field := value.FieldByName(name)
	field.Set(lValueToReflect(lValue, field.Type()))
	return 0
}