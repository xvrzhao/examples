"""
usage of the contextmanager decorator
"""

from contextlib import contextmanager
import os
from typing import Generator, ContextManager

@contextmanager
def open_file(path, flags):
    """粗略模拟标准库的 open 上下文管理器"""
    class File:
        def __init__(self, path, flags):
            fd = os.open(path, flags)
            self.fd = fd
        def write(self, content: str):
            os.write(self.fd, content.encode("utf-8"))
        def read(self, length: int):
            return os.read(self.fd, length).decode("utf-8")
        def close(self):
            os.close(self.fd)

    print("before user code, initializing file instance...")
    file = File(path, flags)
    print("before user code, file instance initialized")

    # 一定要用 try ... finally ..., 否则用户代码如果抛出异常则影响 yield 之后逻辑的执行
    try:
        # yield 之前 = __enter__ 的逻辑；yield 的值 = as 接住的值
        yield file
    finally:
        # yield 之后（在 finally 里）= __exit__ 的逻辑
        print("after user code, closing file descriptor...")
        file.close()
        print("after user code, file descriptor closed")

if __name__ == "__main__":
    ctx = open_file(..., ...) # 这一步，open_file 函数中的逻辑并不会被执行
    print(isinstance(ctx, Generator)) # False, 被装饰过的 open_file 函数的返回值不再是一个 Generator

    # 注意：
    #   经过 contextmanager 装饰过的 Generator 函数返回的是一个 ContextManager 实现！
    #   这和 “实现了 ContextManager 协议的 class” 的实例化过程是一个意思。
    #   两者在这一步都只是实例化的过程，函数内部逻辑并不会被执行！
    #   with 后面跟的是一个 ContextManager 实例
    print(isinstance(ctx, ContextManager)) # True
    print(hasattr(ctx, "__enter__")) # True
    print(hasattr(ctx, "__exit__"))  # True
    print("-"*20)

    """
    猜想内部大致原理：
        1. contextmanager 装饰器大概是先执行 open_file 函数返回一个 Generator 实例 gen
        2. 返回一个 ContextManager 实现（实例），__enter__ 方法中调用 next(gen)，__exit__ 方法再次调用 next(gen) 并返回 False
    """
    
    with open_file("./hello.txt", os.O_CREAT|os.O_TRUNC|os.O_RDWR) as file:
        file.write("hello, world")

    print("-"*20)

    file_ctx = open_file("./hello.txt", os.O_CREAT|os.O_TRUNC|os.O_RDWR)
    with file_ctx as file:
        print("example 1")

    print("-"*20)

    # 再次执行会失败：@contextmanager 返回的 ContextManager 实例内部的 generator 执行 (next) 到头不能重来
    # 这和第一节 (./1_protocol.py) 中的使用类实现的方式不同，那边可以重复执行
    with file_ctx as file:
        print("example 2")