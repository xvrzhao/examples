package unsafe

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// StringHeader 演示了 string 的内部数据结构。
// string 内部的数据结构为 reflect.StringHeader：
//
//   type StringHeader struct {
//   	// 底层字符串真实的内存地址
//   	Data uintptr
//   	// 字符串长度
//   	Len  int
//   }
//
// 直接打印字符串变量的地址其实是打印的 StringHeader 结构体的地址，字符串内存的真实地址为 Data 字段的值。
//
// TODO: Translate to English.
func StringHeader() {
	s := "xavier"

	fmt.Printf("%p\n", &s)                                                 // 0xc0000764e0，StringHeader 结构体地址
	fmt.Printf("0x%x\n", uintptr(unsafe.Pointer(&s)))                      // 0xc0000764e0，StringHeader 结构体地址
	fmt.Printf("%p\n", &(*reflect.StringHeader)(unsafe.Pointer(&s)).Data)  // 0xc0000764e0，StringHeader 结构体第一个字段的地址，也是 StringHeader 结构体地址
	fmt.Printf("0x%x\n", (*reflect.StringHeader)(unsafe.Pointer(&s)).Data) // 0x1148654，底层字符串真实的内存地址
}

// String2Bytes 演示了使用 unsafe 的方法将字符串转为字节切片。
// 因为 字符串内部结构 (reflect.StringHeader) 和切片的内部结构 (reflect.SliceHeader) 存在不同，直接转会出现问题。
//
// TODO: Translate to English.
func String2Bytes() {
	s := "xavier"

	// 直接转，将导致 sliceHeader 缺少 Cap 字段，读取到的 cap 值是 Len 字段后的内存中的值，存在不确定性。
	b1 := *(*[]byte)(unsafe.Pointer(&s))
	fmt.Println(len(b1), cap(b1)) // 6 17740064

	// 提取出 stringHeader 的字段，来构造 sliceHeader，再转为 []byte。
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	b2 := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
	fmt.Println(len(b2), cap(b2)) // 6 6

	// 构造一个匿名结构体，继承 string 的两个字段后再添加一个 Cap 属性，来模拟 sliceHeader。
	b3 := *(*[]byte)(unsafe.Pointer(&struct {
		string // 包含 Data 和 Len
		Cap    int
	}{s, len(s)}))
	fmt.Println(len(b3), cap(b3)) // 6 6
}

// ReadOnlyBytes 演示了只读的 bytes 切片。
// 通过字面量初始化的字符串，编译时会将内存设为只读，即使转换为 bytes 类型也不可
// 操控这部分内存，否则会抛出致命错误，无法通过 recover 捕获。
//
// TODO: Translate to English.
func ReadOnlyBytes() {
	s := "xavier"
	b := *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
	b[0] = 0x61 // throw fatal error
	fmt.Println(s, b)
}

func unsafe2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}

// StringIsReferenceType 演示了 string 其实也是引用类型。
// string 类型的内部结构为 reflect.StringHeader，其中 Data 字段代表了
// 字符串的真实内存地址，而 string/stringHeader 只是一个引用值。
//
// Go 中函数参数都是按值传递的，代码中将字符串 s 传递给函数 unsafe2Bytes 相当于 copy 了一个字符串给函数，
// 因此按照普通思维，你可能会认为在函数内部的行为不会反应到外部的 s 上，函数返回的 bytes 也只是新的一块内存罢了，
// 但是结果却是，返回的 bytes 和 s 指向同一块内存区域。
//
// 这是因为 string 只是一个引用 (reflect.StringHeader)，类似于 slice，函数传递的 s 的确是传递了一个 s 的 copy，
// 但是这个 copy 中的 Data 字段依然是底层字符串的地址，所以返回的 bytes 也是根据这个地址生成的一个 slice，他们共享
// 一块内存区域。
//
// TODO: Translate to English.
func StringIsReferenceType() {
	s := strings.Repeat("a", 3) // s: aaa
	b := unsafe2Bytes(s)
	b[1] = 98
	b[2] = 99
	fmt.Println(s, b) // abc [97 98 99]
}

// StringIsReferenceType1 同样演示了 string 是引用类型。
// 将一个 string 类型的变量赋值给另一个变量时，两者指向同一块内存区域，修改其中一个，另外一个同样会被修改。
//
// 你可能会疑问，函数中第一个代码段里修改了 str 变量，为什么 anotherStr 没被修改。这是因为 str 变量是被
// 重新赋值，内部引用已经被修改了，而 anotherStr 还是指向原来的内存区域。
// 另外需要注意，我们为什么不直接以字符串字面量的形式赋值 str，这是因为字面量形式的字符串是只读的，没法被修
// 改，也就没法验证 "修改其一，两者均改"。
func StringIsReferenceType1() {
	{
		str := strings.Repeat("a", 1)
		anotherStr := str
		fmt.Println(str, anotherStr) // a a
		str = "b"
		fmt.Println(str, anotherStr) // b a
	}
	fmt.Println("---")
	{
		str := strings.Repeat("a", 1)
		anotherStr := str
		fmt.Println(str, anotherStr) // a a

		bytesOfStr := unsafe2Bytes(str)
		bytesOfStr[0] = 98
		fmt.Println(str, anotherStr) // b b
	}
}

// Bytes2String 演示了使用 unsafe 方法将 bytes 转为 string。
// 转换之后的 bytes 和 string 指向同一块内存区域。
//
// 因为 reflect.StringHeader 的两个字段 (Data 和 Len) 与 reflect.SliceHeader 是对齐的，
// 所以转换时不需要另外增加字段。
//
// TODO: Translate to English.
func Bytes2String() {
	b := []byte{0x61, 0x62, 0x63}

	s := *(*string)(unsafe.Pointer(&b))
	fmt.Println(s, len(s)) // abc 3

	b[0] = 0x64
	fmt.Println(s) // dbc
}
