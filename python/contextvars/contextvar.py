"""
上下文变量 ContextVar
"""

from contextvars import ContextVar, copy_context
import asyncio


# 声明一个上下文变量:
#   类型为 ContextVar[str]
#   名称为 user
#   默认值为 空字符串
user_var: ContextVar[str] = ContextVar("user", default="")


def context_var_without_default_value():
    """
    最好给 contextvar 设置默认值。
    没有默认值的 contextvar 在获取时 (get 方法) 会抛 LookupError。
    给 get 方法传入默认值时不会报错，但这样不能确认 contextvar 原本就是这个值，还是因从未设置才返回默认值。
    """
    foo = ContextVar("foo")
    # print(foo.get()) # LookupError: <ContextVar name='foo' at 0x1024bc450>
    print(foo.get(None))


def core_apis():
    """
    ContextVar 对象的三个方法：set / get / reset
    """
    bar_var: ContextVar[int] = ContextVar("bar")

    # 1. get with default value
    val = bar_var.get(0)
    print(val) # print: 0

    # 2. reset
    token = bar_var.set(1) # return previous state which is not set
    print(bar_var.get()) # print: 1
    bar_var.reset(token)
    # bar_var.get() # raise LookupError, as previous state is not set
    print(bar_var.get(-1)) # print: -1

    # 3. set
    bar_var.set(3)

    # get
    val = bar_var.get()
    print(val) # print: 3


async def context_isolation_between_tasks():
    """
    同一个 ContextVar 对象，在不同的 asyncio.Task 中具有独立的访问空间，互不干扰。
    
    关于 Task 和 Coroutine:
        1. Task 为事件循环唯一认识、可以独立调度切换的执行实体，Task 内部的 await coroutine 是串行执行的。
        2. 代码里并发的数量取决于正在运行的 Task 的数量，Task 相当于逻辑并发的入口。
        3. 事件循环负责在 Tasks 之间切换交替执行，遇到 await 则让出当前 Task 的控制权，转去执行其他 Task。
        4. 可以把 Task 理解为异步世界里的线程，Coroutine 只是线程中的一个函数。
        5. 同一个 Task 中，即使 coroutine 存在嵌套，他们之间也都是共享同一个 ContextVar 上下文，改动可彼此看见。

    该函数输出：
        worker a set: John
        worker b set: Coltrane
        worker b get: Coltrane
        worker a get: John
    """

    async def worker_a(user: str):
        print("worker a set:", user)
        user_var.set(user)
        await asyncio.sleep(1)
        print("worker a get:", user_var.get())

    async def worker_b(user: str):
        await asyncio.sleep(0.5)
        print("worker b set:", user)
        user_var.set(user)
        print("worker b get:", user_var.get())

    await asyncio.gather(worker_a("John"), worker_b("Coltrane")) # asyncio.gather 会将两个 coroutine 包装成两个 task 分别独立运行


async def context_sharing_within_one_task():
    """
    同一个 Task 中，即使 coroutine 存在嵌套，他们之间也都是共享同一个 ContextVar 上下文，改动可彼此看见。
    """
    async def task():
        async def inner():
            user_var.set("inner")

        user_var.set("outer")
        await inner()

        # 子协程中对 contextvar 的改动，外层也可以看到
        print(user_var.get()) # print: inner

    await asyncio.create_task(task())


async def copy_and_inherit():
    """
    关键概念：Context

    每次调用 asyncio.create_task()（或者 loop.call_soon() 等）创建一个新的并发执行单元时，
    Python 会自动复制当前的 Context，作为新任务的初始上下文。

    程序输出：
        child read: value father set
        child set: value child set
        father read: value father set
    """

    async def child():
        """
        记住：child task 复制出一份独立的上下文，同时初始值继承自 father task。
        """
        print("child read:", user_var.get())
        user_var.set("value child set")
        print("child set:", user_var.get())

    user_var.set("value father set")
    task = asyncio.create_task(child()) # spawn a task
    await task
    print("father read:", user_var.get())


def copy_context_manually():
    """
    手动拷贝上下文，上一个例子的等效操作，实际上在创建新的 asyncio.Task 时就是执行的类似操作

    输出：
        child read: value outter set
        child set: value child set
        outter read: value outter set
    """

    def child():
        """
        该函数具有了一份独立的上下文，同时初始值继承自外部
        """
        print("child read:", user_var.get())
        user_var.set("value child set")
        print("child set:", user_var.get())

    user_var.set("value outter set")
    ctx = copy_context()
    ctx.run(child)
    print("outter read:", user_var.get())


async def main():
    # await context_sharing_within_one_task()
    # await copy_and_inherit()
    copy_context_manually()


if __name__ == "__main__":
    asyncio.run(main())