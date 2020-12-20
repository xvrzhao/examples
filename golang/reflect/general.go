package reflect

import (
	"fmt"
	"reflect"
)

func RunSetExample() {
	f := 3.14

	fv := reflect.ValueOf(&f)
	if fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}
	if fv.CanSet() && fv.Kind() == reflect.Float64 {
		fv.SetFloat(3.1415)
	}

	fmt.Println(f) // 3.1415
}

type person struct {
	Name   string
	Age    uint8
	gender uint8
}

func (p person) tellGender() {
	fmt.Println(p.gender)
}

func (p person) TellName() {
	fmt.Println(p.Name)
}

func (p *person) Grow(years uint8) {
	p.Age += years
}

func RunSetStructExample() {
	v := reflect.ValueOf(person{})
	fmt.Println(v.CanSet(), v.NumField(), v.NumMethod()) // false 3 1(methods that their receiver is `person` and are exported)

	v = reflect.ValueOf(&person{})
	fmt.Println(v.CanSet(), v.NumMethod()) // false 2(methods that their receiver could be `person` or `*person` and are exported)
	// fmt.Println(v.NumField()) // panic: reflect: call of reflect.Value.NumField on ptr Value

	v = reflect.ValueOf(&person{}).Elem()
	fmt.Println(v.CanSet(), v.NumField(), v.NumMethod()) // true 3 1(methods that their receiver is `person` and are exported)

	p := person{}
	v = reflect.ValueOf(&p).Elem()
	nv := v.FieldByName("Name")
	fmt.Println(nv.Kind(), nv.CanSet()) // string true
	nv.SetString("Xavier")
	fmt.Println(p.Name) // Xavier

	gv := v.FieldByName("gender")
	fmt.Println(gv.Kind(), gv.Uint(), gv.CanSet()) // uint8 0(value is accessible even though gender is unexported) false(gender is unexported)

	zv := v.FieldByName("notExistField") // sv is zero Value
	// fmt.Println(sv.IsZero()) // panic: reflect: call of reflect.Value.IsZero on zero Value
	fmt.Println(zv.IsValid(), zv.CanSet()) // false false

	tellName := v.MethodByName("TellName")
	if tellName.IsValid() {
		tellName.Call(nil) // Xavier
	}

	tellGender := v.MethodByName("tellGender")
	grow := v.MethodByName("Grow")
	fmt.Println(tellGender.IsValid(), grow.IsValid()) // false false
}
