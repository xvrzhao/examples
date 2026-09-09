"""
AsyncIterable 和 AsyncIterator

AsyncIterable/AsyncIterator 是异步版本的 Iterable/Iterator, 协议基本上是平行的。

AsyncIterable 需实现 __aiter__ 方法；
AsyncIterator 需实现 __aiter__, __anext__ 方法；

重点：__aiter__ 是同步方法，__anext__ 是异步方法；

__anext__ 方法没有新的值迭代时，抛出 StopAsyncIteration 异常；
循环语法使用: async for x in obj
"""
import asyncio

class AsyncCounter:
    def __init__(self, stop):
        self.stop = stop

    def __aiter__(self):
        return AsyncCounterIterator(self.stop)

class AsyncCounterIterator:
    def __init__(self, stop):
        self.stop = stop
        self.n = 0 # 内部维护推进状态

    def __aiter__(self):
        return self

    # 唯有 __anext__ 方法是异步的
    async def __anext__(self):
        if self.n >= self.stop:
            raise StopAsyncIteration
        await asyncio.sleep(0.1) # 模拟 I/O 操作
        self.n += 1
        return self.n

async def async_iterable_and_async_iterator():
    counter = AsyncCounter(5)
    async for v in counter:
        print(v)

    # 再次消费，依旧从 1 开始
    async for v in counter:
        print(v)

async def async_for_is_the_same_as():
    """async for 操作相当于以下逻辑"""
    counter = AsyncCounter(5)
    it = aiter(counter)
    while True:
        try:
            v = await anext(it)
        except StopAsyncIteration:
            break

        print(v)


async def main():
    await async_iterable_and_async_iterator()
    await async_for_is_the_same_as()


if __name__ == "__main__":
    asyncio.run(main())