# Context Manager 协议包含两个方法：__enter__ 和 __exit__
class MyContext:
    def __enter__(self):
        print("enter", end=", ")
        return self

    def __exit__(self, exc_type, exc, tb):
        """
        注意：
            1. 函数后面三个参数为用户代码的异常信息，若无异常则均为 None
            2. 函数返回真值（truthy）代表捕获异常；若忘记返回，即返回 None，相当于是返回了 False
        """
        print("exit")
        return False # 若用户代码块中出现异常，是否选择捕获，False 代表不捕获，即继续向外层抛出


if __name__ == "__main__":

    from typing import ContextManager
    print(isinstance(MyContext(), ContextManager)) # True

    # 用户层使用方法：
    #   1. as 后面的变量为 __enter__ 方法的返回值，若不需要使用返回值，可以省略 as 语句；
    #   2. 代码中的 MyContext() 是实例化过程，调用的依然是类的 __init__ 方法，并不是 __enter__ 方法；
    #   3. with 语句代码的执行顺序是先执行 __enter__，再执行用户代码块，最后执行 __exit__ 方法；
    #   4. 无论用户代码块会不会抛出异常，都不影响 __exit__ 方法的执行。若用户代码有异常产生，异常信息会传递到 __exit__ 方法中，context manager 可以选择是否向外传递（即传递到 with 外层）；
    with MyContext() as ctx:
        print("executing-user-code", end=", ") # print: enter, executing-user-code, exit

    # 实例化出来的 ContextManager 可重复使用
    ctx = MyContext()
    with ctx:
        print("hello", end=", ") # print: enter, hello, exit
    with ctx:
        print("world", end=", ") # print: enter, world, exit
