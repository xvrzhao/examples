package data_structure

import "fmt"

func RunMapPanic() {
	// m := make(map[string]int) // declared and initialized
	var m map[string]int  // declared but not initialized
	fmt.Println(m == nil) // true
	m["xvrzhao"] = 23     // panic: assignment to entry in nil map
}

func RunMapRefer() {
	m1 := make(map[string]string, 2)
	m1["apple"] = "red"
	m1["banana"] = "yellow"

	m2 := m1 // map is reference type, m1 and m2 point to the same block of memory
	m2["apple"] = "green"

	fmt.Println(m1) // map[apple:green banana:yellow]
	fmt.Println(m2) // map[apple:green banana:yellow]
}

func RunMapExceededCap() {
	m1 := make(map[string]string, 2)
	m1["apple"] = "red"
	m1["banana"] = "yellow"

	m2 := m1 // m2 points to m1, or m2 and m1 point to the same block of memory

	m1["orange"] = "orange" // m2 also changes, even though m1 has exceeded capacity limit
	fmt.Println(m1)         // map[apple:red banana:yellow orange:orange]
	fmt.Println(m2)         // map[apple:red banana:yellow orange:orange]
}

type Age struct {
	value int
}

// RunMapNotAddressable demonstrates that a map's v(value) are not addressable.
//
// When a v stores the struct type or array type, you can not change fields/elements of that value directly.
// Therefore, usually a map's v just stores a reference type like struct pointer, array pointer, slice, or map.
func RunMapNotAddressable() {
	m := map[string]Age{
		"xavier": Age{value: 23},
		"john":   Age{value: 13},
	}

	//m["xavier"].value = 10 // compile error: cannot assign to struct field m["xavier"].value in map

	// change fields of struct value in a map
	xavier := m["xavier"]
	xavier.value = 10
	m["xavier"] = xavier

	// or the following function
}

func RunMapAddressable() {
	// when the type of value in a map is a reference type
	m := map[string]*Age{
		"xavier": &Age{value: 23},
	}

	m["xavier"].value = 10
	fmt.Println(m["xavier"])
}

func RunMapIterate() {
	var m map[string]string

	initM := func() {
		m = map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
			"d": "4",
			"e": "5",
			"f": "6",
		}
	}

	for i := 0; i < 10; i++ {
		initM()

		// do both add and delete when iterate a map, the result will be unknown
		for k, v := range m {
			m[v] = k
			delete(m, k)
		}

		fmt.Println(m)
	}
	// stdout:
	// map[2:b 3:c 4:d 5:e 6:f a:1]
	// map[1:a 2:b 3:c 4:d 5:e 6:f]
	// map[1:a 3:c 4:d 5:e 6:f b:2]
	// map[2:b 3:c 4:d 5:e 6:f a:1]
	// map[1:a 2:b 3:c 5:e 6:f d:4]
	// map[1:a 2:b 3:c 4:d 5:e 6:f]
	// map[1:a 2:b 4:d 5:e 6:f c:3]
	// map[1:a 2:b 3:c 4:d 5:e f:6]
	// map[1:a 3:c 4:d 5:e 6:f b:2]
	// map[1:a 2:b 3:c 4:d 5:e 6:f]
}
