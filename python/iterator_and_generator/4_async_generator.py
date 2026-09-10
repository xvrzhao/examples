"""
AsyncGenerator 异步生成器
"""

import asyncio
from typing import AsyncIterable, AsyncIterator, AsyncGenerator

async def count(stop: int) -> AsyncIterator[int]:
    """
    异步生成器函数，该函数返回一个异步生成器对象。

    凡是函数体中包含 yield 关键字的函数都是生成器函数;
    如果函数体中既有 yield 关键字，又有 await 异步 io 操作，则该函数则是异步生成器函数，返回一个异步生成器对象。
    """
    n = 0
    while n < stop:
        await asyncio.sleep(0.1) # 模拟异步io操作
        yield n
        n += 1

async def async_generator_is_async_iterator():
    """
    异步生成器内部实现了 AsyncIterator 协议，包含 __aiter__ 和 __anext__ 方法，
    所以 AsyncGenerator 也是一个 AsyncIterator。
    """

    # 注意: 
    # 要注意分辨异步生成器函数和普通的异步函数，两者都由 async def 关键字定义，
    # 但是！调用异步生成器函数返回的是一个 AsyncGenerator 对象，是直接返回，不需要 await；
    # 而普通的 async 函数返回的是一个 coroutine，执行时需要 await。
    # 
    # 所以，为了让代码更易读，最好标注一下异步生成器函数的返回值类型为 AsyncGenerator 或者直接标注为 AsyncIterator

    ag = count(5) # 不需要 await
    print(type(ag)) # <class 'async_generator'>
    print(isinstance(ag, AsyncGenerator)) # True

    print(hasattr(ag, "__aiter__")) # True
    print(hasattr(ag, "__anext__")) # True

    print(isinstance(ag, AsyncIterator)) # True
    print(isinstance(ag, AsyncIterable)) # True

    # 消费
    async for i in ag:
        print(i, end="|") # 0|1|2|3|4|

    print()

    # 再次消费，无任何输出，AsyncIterator 中的推进状态已到头
    async for i in ag:
        print(i) # 无输出

    print()

    # 但是每次调用生成器函数，都会返回一个全新的 iterator，内部状态都是全新的，可以重复消费
    async for i in count(5):
        print(i, end=",") # 0,1,2,3,4,

    print()

    async for i in count(5):
        print(i, end=".") # 0.1.2.3.4.

    print()


async def gen() -> AsyncIterator[int]:
    """如果异步生成器函数只是用来被 for 消费，可以将函数返回值直接标注为 AsyncIterator 而不是 AsyncGenerator 会更直观"""
    await asyncio.sleep(0.1)
    yield 1
    await asyncio.sleep(0.1)
    yield 2
    await asyncio.sleep(0.1)
    yield 3


def relationship_of_AsyncIterable_AsyncIterator_and_AsyncGenerator():
    """
    AsyncGenerator, AsyncIterator 和 AsyncIterable 之间的关系：
        AsyncGenerator ⊂ AsyncIterator ⊂ AsyncIterable
    """


async def yield_from_a_async_iterable():
    """
    yield from 可否用于 AsyncIterable ？
    答案是：不可以，yield from 只能用于同步的 Iterable；并且也没有 `async yield from` 这样的语法
    """

    async def outer():
        """尝试 yield from count(5)"""

        # yield from count(5)       # SyntaxError: 'yield from' inside async function 编译器直接不允许
        # 而且即使编译器允许 count(5) 也无法被 yield from 识别，yield from 底层实际是调用 __iter__ 和 __next__ 方法，但是 count(5) 返回的 generator 只有 __aiter__ 和 __anext__ 方法

        # async yield from count(5) # SyntaxError: invalid syntax 没有 async yield from 这样的语法

    def outer2():
        """同步迭代器中使用 yield from"""

        yield from count(5) # 这样看似编译器没有报错，但是运行时报错：TypeError: 'async_generator' object is not iterable

        # 原因很简单：
        # yield from 只用于同步 Iterable，
        # 底层实际是先调用对象的 __iter__ 方法，但是 count(5) 返回的 generator 只有 __aiter__ 和 __anext__ 方法，
        # 并不满足同步 Iterable 协议，所以报：'async_generator' object is not iterable

    try:
        for i in outer2():
            pass
    except Exception as e:
        print("catched exception:", e)


    async def outer3():
        """只能用 async for 取值然后手动 yield"""

        async for i in count(5):
            yield i

        yield 5
        yield 6

    async for i in outer3():
        print(i) # print: 0 ~ 6


async def usage_of_asend():
    """
    yield 接受外部的传值 在异步生成器函数中的应用
    (2026.09.10: 该应用场景目前还没怎么见过，先举例最简单的用法，以后见到实际使用案例再进一步讲解)
    """

    async def receive_value_from_outer():
        while True:
            await asyncio.sleep(0.1) # 模拟io等待
            received_value = yield
            print("received:", received_value)

    ag = receive_value_from_outer()
    print(type(ag)) # <class 'async_generator'>

    await anext(ag) # 需要先通过 anext 启动，执行到第一次 yield 的位置。也可以使用 g.asend(None)，第一次启动时必须传 None
    await ag.asend("hello") # 在 generator 内部，从 yield 继续往下执行，直到下一个 yield 处暂停
    await ag.asend("world")
    for i in range(1, 5):
        await ag.asend(i)


async def main():
    await async_generator_is_async_iterator()
    await yield_from_a_async_iterable()
    await usage_of_asend()


if __name__ == "__main__":
    asyncio.run(main())