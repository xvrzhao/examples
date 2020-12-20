package unsafe

import (
	"fmt"
	"unsafe"
)

// Note: All annotations are results of execution under 64-bit platform.

// Sizeof prints memory size of all go types under 64-bit platform.
func Sizeof() {
	// int
	fmt.Println(
		unsafe.Sizeof(int(0)),     // 8
		unsafe.Sizeof(uint(0)),    // 8
		unsafe.Sizeof(int8(0)),    // 1
		unsafe.Sizeof(uint8(0)),   // 1
		unsafe.Sizeof(int16(0)),   // 2
		unsafe.Sizeof(uint16(0)),  // 2
		unsafe.Sizeof(int32(0)),   // 4
		unsafe.Sizeof(uint32(0)),  // 4
		unsafe.Sizeof(int64(0)),   // 8
		unsafe.Sizeof(uint64(0)),  // 8
		unsafe.Sizeof(uintptr(0)), // 8
	)

	// float
	fmt.Println(unsafe.Sizeof(float32(0)), unsafe.Sizeof(float64(0))) // 4 8

	// bool
	fmt.Println(unsafe.Sizeof(false), unsafe.Sizeof(true)) // 1 1

	// string is underlying type reflect.StringHeader
	str := "xvr"
	fmt.Println(
		unsafe.Sizeof(str),       // 16 uintptr + int
		unsafe.Sizeof("xvrzhao"), // 16
	)

	// slice is underlying type reflect.SliceHeader, though s is not allocated memory (s == nil)
	var (
		s0 []int8
		s1 []int16
	)
	fmt.Println(
		unsafe.Sizeof(s0),                 // 24 uintptr + int + int
		unsafe.Sizeof(s1),                 // 24
		unsafe.Sizeof([]int8{}),           // 24
		unsafe.Sizeof([]int8{0, 1, 2, 3}), // 24
	)

	// array
	var (
		arr1 [3]uint32 // 4 * 3
		arr2 [5][]bool // 24 * 5
	)
	fmt.Println(unsafe.Sizeof(arr1), unsafe.Sizeof(arr2)) // 12 120

	// unsafe.Pointer
	fmt.Println(unsafe.Sizeof(unsafe.Pointer(&arr1))) // 8

	// map is a pointer to underlying type runtime.hmap, so a map variable(a pointer) is 8 bytes
	var m1 map[string]int
	m2 := make(map[string]int, 100)
	fmt.Println(unsafe.Sizeof(m1), unsafe.Sizeof(m2)) // 8 8

	// interface, a interface variable is a underlying type runtime.iface while
	// a empty interface is a underlying type runtime.eface, iface and eface are
	// structs with 2 pointer-type fields.
	//
	//   type iface struct {
	//   	tab  *itab
	//   	data unsafe.Pointer
	//   }
	//
	//   type eface struct {
	//   	_type *_type
	//   	data  unsafe.Pointer
	//   }
	type myInterface interface {
		DoSth()
	}
	var (
		mi myInterface
		ei interface{}
	)
	fmt.Println(unsafe.Sizeof(mi), unsafe.Sizeof(ei)) // 16 16
}

// StructCompare 示例说明了：具有相同字段的结构体，不同的字段排序方式，所产生的内存占用不同，
// 原因是由于结构体字段和结构体整体的对齐方式。
//
// 知识点：
//   1. 64 位平台的机器字长（CPU 一次读入的字节数）为 8 字节，32 位平台为 4 字节。
//   2. 一个数据类型在对齐时，要保证其起始地址为其对齐值的整数倍。
//   3. 所有类型的对齐值可以通过 unsafe.Alignof 计算，array 类型的对齐值为其元素的对齐值，
//      struct 类型的对齐值为其最长字节的字段的对齐值。
//   4. 结构体所占的字节长度需为其对齐值的整数倍（据此来推算结构体最后需要多少个字节的 padding）。
//
// 参考链接：
//   - https://www.jianshu.com/p/49f7e6f56568
//   - https://www.bilibili.com/video/BV1iZ4y1j7TT
//   - https://github.com/talk-go/night/issues/588
func StructCompare() {
	type S1 struct {
		a int8  // 1 0 0 0 0 0 0 0
		b int64 // 1 1 1 1 1 1 1 1
		c int16 // 1 1 0 0 0 0 0 0
	}
	type S2 struct {
		a int8 // 1 0 1 1 0 0 0 0
		c int16
		b int64 // 1 1 1 1 1 1 1 1
	}
	fmt.Println(unsafe.Sizeof(S1{}), unsafe.Sizeof(S2{})) // 24 16

	type S3 struct {
		a int8 // 1 0 0 0 1 1 1 1
		b int32
		c int16 // 1 1 0 0
	}
	s3 := S3{}
	s1 := S1{}
	// 结构体大小是结构体对齐值的整数倍，结构体对齐值是结构体中最大字段的对齐值
	fmt.Println(unsafe.Alignof(s3), unsafe.Sizeof(s3), unsafe.Offsetof(s3.b), unsafe.Offsetof(s3.c)) // 4 12 4 8
	fmt.Println(unsafe.Alignof(s1), unsafe.Sizeof(s1))                                               // 8 24
}
