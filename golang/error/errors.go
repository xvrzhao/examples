package error

import (
	"errors"
	"fmt"
)

var (
	// implement an error chain: err3(top) -> err2 -> err1(root)
	err1 = errors.New("err1")
	err2 = fmt.Errorf("err2: %w", err1)
	err3 = fmt.Errorf("err3: %w", err2)

	// 在一个 error 中包含多个 error 将不被允许，go vet 静态检查会检测出该问题，例如：
	//   err4 = fmt.Errorf("err4: %w, %w", err2, err3)
)

func ExampleOfUnwrap() {
	fmt.Println(err3) // err3: err2: err1

	_err2 := errors.Unwrap(err3)
	fmt.Println(_err2, err2 == err2) // err2: err1  true

	_err1 := errors.Unwrap(_err2)
	fmt.Println(_err1, _err1 == err1) // err1  true
}

// ExampleOfIs 演示了 errors.Is(err, target error) 的用法。
//   1. err 或 组成 err 的链中，只要有一个等于 target，则返回 true。
//   2. 若 err 底层类型包含 Is(target error) 方法，则优先使用该方法的返回结果。
func ExampleOfIs() {
	fmt.Println(
		// true
		errors.Is(err3, err3),
		errors.Is(err3, err2),
		errors.Is(err3, err1),

		errors.Is(err2, err2),
		errors.Is(err2, err1),

		errors.Is(err1, err1),

		// false
		errors.Is(err2, err3),
		errors.Is(err1, err2),
		errors.Is(err1, err3),
	)
}
