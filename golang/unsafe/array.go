package unsafe

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ArrayOffsetElement 演示了获取数组元素地址的两种方式。
func ArrayOffsetElement() {
	arr := [5]int{}
	offset := 3

	fmt.Printf("%p \n", &arr[offset])                                                          // 0xc000014198
	fmt.Printf("0x%x \n", uintptr(unsafe.Pointer(&arr))+uintptr(offset)*unsafe.Sizeof(arr[0])) // 0xc000014198
}

// ArraySliceDiff 演示了 slice 与 array 本质上的区别。
// slice 是一个 reflect.SliceHeader 结构，而 array 是一段连续的内存结构。
func ArraySliceDiff() {
	s := []int{1, 2, 3}
	fmt.Printf("%p != %p == %x\n", &s, &(s[0]), (*reflect.SliceHeader)(unsafe.Pointer(&s)).Data)

	arr := [3]int{1, 2, 3}
	fmt.Printf("%p == %p\n", &arr, &arr[0])
}
