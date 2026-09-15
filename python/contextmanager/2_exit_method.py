"""
Context Manager 中 __exit__ 方法的使用
"""

class ExceptionSuppress:
    """
    用户实例化时传入需要抑制的异常类型，当 with 代码块中抛出对应的异常时会被抑制，不会中断代码的执行。
    通过该示例展示 __exit__ 方法中参数和返回值的使用。
    """

    def __init__(self, *exc_types):
        """传入需要抑制的异常类型"""
        self.exc_types = exc_types

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc, tb):
        if exc_type in self.exc_types:
            print(f"WARNING: exception suppressed: type: {exc_type}, value: {exc}, traceback: {tb}")
            return True
        else:
            return False

if __name__ == "__main__":

    # 用户代码中异常被抑制的案例

    with ExceptionSuppress(ValueError):
        print("before error")
        raise ValueError("something wrong")
        print("after error, won't be executed")

    print("ValueError is suppressed, so here will be executed.")


    print("-" * 20)


    # 用户代码中异常未被抑制的案例

    try:

        with ExceptionSuppress(ValueError):
            print("before error")
            raise NameError("something wrong") # this error will not be suppressed
            print("after error, won't be executed.")

        print("after error, won't be executed.")

    except Exception as e:
        print(f"catched exception: {e}")


    print("-" * 20)


    # 以上是在模拟标准库中的 suppress 类，无需自己实现

    from contextlib import suppress
    from typing import ContextManager

    print(isinstance(suppress(), ContextManager))

    with suppress(BaseException):
        raise Exception("my exception") # Exception 是 BaseException 的子类，所以也会被抑制

    print("after suppress") # 会被执行
