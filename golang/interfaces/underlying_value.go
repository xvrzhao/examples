package interfaces

import (
	"fmt"
)

// Equal 演示了接口值参与 == 比较的情况。
// 当一个接口值参与 == 比较时，无论另一方是接口值还是类型值，两者之间必须是实现与被实现关系，否则编译不通过
func Equal() {
	/* 接口值与类型值进行比较时，可以把接口值看成底层的类型值 */

	girl1 := new(Girl)
	girl2 := new(Girl)
	var human Human = girl1
	fmt.Println(human == girl1, human == girl2) // 相当于判断 girl1 == girl1 和 girl1 == girl2，故结果为 true false

	man1 := Man{"Gopher"}
	man2 := Man{"Gopher"}
	man3 := Man{"Xavier"}
	var walker Walker = man1
	fmt.Println(walker == man1, walker == man2, walker == man3) // 相当于判断 man1 == man1、man1 == man2、man1 == man3，故结果为 true true false

	/* 接口值与接口值进行比较时，同样把两个接口值都看成其底层的类型值 */

	var walker1, walker2 Walker
	fmt.Println(walker1 == walker2) // 相当于判断 nil == nil，故 true
	walker1, walker2 = man1, man2
	fmt.Println(walker1 == walker2) // 相当于判断 man1 == man2，故 true
	walker2 = man3
	fmt.Println(walker1 == walker2) // 相当于判断 man1 == man3，故 false

	var human1, human2 Human
	human1, human2 = girl1, girl2
	fmt.Println(human1 == human2) // 相当于判断 girl1 == girl2，故 false

	man4, man5 := new(Man), new(Man)
	var sw SpeakWalker = man4
	var s Speaker = man4
	fmt.Println(sw == s) // 相当于判断 man4 == man4，故 true
	s = man5
	fmt.Println(sw == s) // 相当于判断 man4 == man5，故 false
	s = sw
	fmt.Println(sw == s) // 相当于判断 man4 == man4，故 true
}
