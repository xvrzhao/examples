package unsafe

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Person struct {
	Name string
	Age  int
}

// StructOffsetField 演示了获取结构体字段的三种不同方式。
func StructOffsetField() {
	p := Person{}
	fmt.Printf("%p \n", &p.Age)                                               // 0xc00000c090
	fmt.Printf("0x%x \n", uintptr(unsafe.Pointer(&p.Age)))                    // 0xc00000c090
	fmt.Printf("0x%x \n", uintptr(unsafe.Pointer(&p))+unsafe.Offsetof(p.Age)) // 0xc00000c090
}

// Person2Bytes 将 Person 结构体转为 bytes。
// 这里或许不能称为转，应该为获取结构体的内存引用。
// 需要注意 GC 问题，如果结构体被回收了，根据 bytes 的 Data(起点) 和 Len(长度) 获取到的内存可能已经不再是之前的结构体了。
//
// TODO: Translate to English.
func Person2Bytes(p *Person) []byte {
	size := unsafe.Sizeof(*p)

	s := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(p)),
		Len:  int(size),
		Cap:  int(size),
	}

	return *(*[]byte)(unsafe.Pointer(&s))
}

// Bytes2Person 将 bytes 再转回结构体。
// 结合 Person2Bytes 函数进行来回互转见单元测试：TestPerson2Bytes。项目根目录下执行：
//   go test -v -run=Person2Bytes ./unsafe
//
// TODO: Translate to English.
func Bytes2Person(b []byte) *Person {

	// 以下为错误的方式，执行 go vet 将提示 possible misuse of unsafe.Pointer
	// 不能以中间变量的形式存储 uintptr 值，防止其值表示的内存地址被 GC
	//
	//personAddress := (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data // 此步拿到字节数组地址，但是 b 不再被使用，很可能被 GC，此后 personAddress 表示的地址不再有意义
	//return (*Person)(unsafe.Pointer(personAddress))

	return (*Person)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data))
}
