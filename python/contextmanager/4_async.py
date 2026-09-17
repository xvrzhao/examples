"""
AsyncContextManager 异步上下文管理器

AsyncContextManager 必须实现 __aenter__ 和 __aexit__ 方法，两者都是 async 方法。
这点和 AsyncIterator 不一样，AsyncIterator 的 __aiter__ 方法是同步方法，__anext__ 是异步方法。
"""
import os
import asyncio
from typing import AsyncContextManager
from contextlib import asynccontextmanager

class MyFile:
    """
    模拟异步版文件上下文管理器

    一般，类中即实现了上下文管理器协议，同时也具备了一些供用户使用的方法（如以下的 write），
    这样，在 __aenter__ 方法中就可以返回自身实例 self 来供用户调用方法。
    """
    def __init__(self, path: str, flags: int):
        self.path = path
        self.flags = flags

    async def __aenter__(self):
        # 将文件 open 操作放在线程池中去处理，外部挂起转而去执行其他协程，等线程执行完毕再回来继续往下执行
        fd = await asyncio.to_thread(lambda path, flags: os.open(path, flags), self.path, self.flags)
        # 或直接写成:
        # fd = await asyncio.to_thread(os.open, self.path, self.flags)

        self.fd = fd
        print("aenter finished, fd:", self.fd)

        return self

    async def __aexit__(self, exc_type, exc, tb):
        await asyncio.to_thread(lambda fd: os.close(fd), self.fd)
        print("aexit finished, fd closed")

    async def write(self, content: str):
        data = content.encode()
        await asyncio.to_thread(lambda fd, data: os.write(fd, data), self.fd, data)

"""
使用 @asynccontextmanager 装饰器 + yield 创建的 AsyncContextManager，
被装饰器包装之后，执行函数返回的是一个 AsyncContextManager 实例，实例包含 __aenter__ 和 __aexit__ 方法，也就是说，实现了 AsyncContextManager 协议。
"""
@asynccontextmanager
async def my_file(path: str, flags: int):
    fd = await asyncio.to_thread(lambda path, flags: os.open(path, flags), path, flags)
    print("file opened, fd:", fd)
    try:
        yield fd
    finally:
        await asyncio.to_thread(lambda fd: os.close(fd), fd)
        print("file closed, fd:", fd)


async def main():
    my_file_ctx = MyFile("./hello.txt", os.O_CREAT|os.O_TRUNC|os.O_RDWR)
    print(isinstance(my_file_ctx, AsyncContextManager)) # True
    
    async with my_file_ctx as file:
        print("first")

    # 通过类实现的 AsyncContextManager 可以多次作为上下文执行
    async with my_file_ctx as file:
        await file.write("hello, async")
        print("second")

    print("-" * 20)

    my_file_ctx = my_file("./hello.txt", os.O_CREAT|os.O_TRUNC|os.O_RDWR)

    print(isinstance(my_file_ctx, AsyncContextManager)) # True
    print(hasattr(my_file_ctx, "__aenter__")) # True
    print(hasattr(my_file_ctx, "__aexit__")) # True

    async with my_file_ctx as fd:
        await asyncio.to_thread(lambda data: os.write(fd, data), b"hello, decorator")
        print("first")

    try:
        # 报错，使用装饰器生成的上下文管理器无法多次被调用
        async with my_file_ctx as fd:
            print("second")
    except Exception as e:
        print("caught error:", e)

    print("-" * 20)

    # with 后可以跟多个上下文管理器
    async with (
        my_file("./a.txt", os.O_CREAT|os.O_WRONLY) as a_fd,
        my_file("./b.txt", os.O_CREAT|os.O_WRONLY) as b_fd,
    ):
        print("user code")
        await asyncio.to_thread(os.write, a_fd, b"aaa")
        await asyncio.to_thread(os.write, b_fd, b"bbb")

    """
    以上代码输出：
        file opened, fd: 6
        file opened, fd: 7
        user code
        file closed, fd: 7
        file closed, fd: 6

    可以看到，多个上下文时，遵从先进后出的原则
    """


if __name__ == "__main__":
    asyncio.run(main())