Context Manager 跟迭代器协议是同一种设计思路: **用一对"魔法方法"定义一套协议，`with` 语句负责按协议调用它们**。

## 1. 协议本身：`__enter__` + `__exit__`

一个对象只要同时实现了这两个方法，就是一个上下文管理器，可以放进 `with` 语句：

```python
class MyContext:
    def __enter__(self):
        print("进入上下文")
        return self          # 返回值会被 as 后面的变量接住

    def __exit__(self, exc_type, exc_value, traceback):
        print("退出上下文")
        return False          # 返回值决定异常是否被"吞掉"

with MyContext() as ctx:
    print("上下文内部")
```

输出：

```
进入上下文
上下文内部
退出上下文
```

## 2. `with` 语句的展开逻辑

```python
with MyContext() as ctx:
    BODY
```

大致等价于：

```python
mgr = MyContext()
ctx = mgr.__enter__()
try:
    BODY
except BaseException:
    if not mgr.__exit__(*sys.exc_info()):
        raise
else:
    mgr.__exit__(None, None, None)
```

关键点：**不管 `BODY` 有没有抛异常，`__exit__` 都保证会被调用**——这正是 `with` 语句存在的核心价值：把"资源获取"和"资源释放"绑在一起，释放这一步是**无条件执行**的，不需要你手写 `try/finally`。

```python
# with 语句出现之前，等价的写法必须手动 try/finally
mgr = MyContext()
ctx = mgr.__enter__()
try:
    BODY
finally:
    mgr.__exit__(None, None, None)   # 简化版，没处理异常传递给 __exit__ 的细节
```

## 3. `__exit__` 的三个参数和返回值：异常处理的核心

`__exit__(self, exc_type, exc_value, traceback)` 这三个参数就是 `with` 块里发生的异常信息（如果没异常，全是 `None`）。**返回值决定这个异常要不要被"吞掉"**：

```python
class SuppressError:
    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        if exc_type is ValueError:
            print(f"忽略了一个 ValueError: {exc_value}")
            return True     # 返回 True，异常被吞掉，不再向外传播
        return False         # 返回 False（或不写返回值，默认 None 也是 falsy），异常正常抛出

with SuppressError():
    raise ValueError("出错了")
print("程序继续运行")   # 能执行到这里，因为异常被吞掉了

with SuppressError():
    raise TypeError("类型错误")
# 这个异常不会被吞，会正常往外抛出
```

这也是标准库 `contextlib.suppress` 的实现原理：

```python
from contextlib import suppress

with suppress(FileNotFoundError):
    os.remove("不存在的文件.txt")   # 不会报错，因为 __exit__ 判断异常类型匹配就返回 True
```

**这里有个常见的坑**：`__exit__` 里不写 `return` 语句，默认返回 `None`，`None` 是 falsy，等价于 `False`，所以**不会**意外吞掉异常——这是个安全的默认行为。

## 4. 最常见的应用：资源管理（文件、锁、连接）

```python
with open("file.txt") as f:
    data = f.read()
# 离开 with 块，f.close() 自动被调用，即便 read() 抛异常也一样会 close
```

`open()` 返回的文件对象内部实现了：

```python
class TextIOWrapper:
    def __enter__(self):
        return self
    def __exit__(self, *args):
        self.close()
        return False
```

同理，线程锁：

```python
import threading
lock = threading.Lock()

with lock:            # 等价于 lock.acquire() ... lock.release()
    critical_section()
# 即便 critical_section() 抛异常，lock 也一定会被 release，避免死锁
```

## 5. `contextlib.contextmanager`：用生成器函数简化写法

跟之前 `yield` 简化迭代器写法思路一样，Python 提供了 `@contextmanager` 装饰器，让你**用一个生成器函数**就能写出上下文管理器，而不用手写一整个类的 `__enter__`/`__exit__`：

```python
from contextlib import contextmanager

@contextmanager
def my_context():
    print("进入上下文")
    try:
        yield "some_value"    # yield 之前 = __enter__ 的逻辑；yield 的值 = as 接住的值
    finally:
        print("退出上下文")     # yield 之后（在 finally 里）= __exit__ 的逻辑

with my_context() as val:
    print(val)
```

输出：

```
进入上下文
some_value
退出上下文
```

**这里为什么必须用 `try/finally`？** 因为如果 `with` 块内部抛了异常，这个异常会在生成器暂停的 `yield` 那一行被**重新抛出**（本质上是通过 `gen.throw()` 机制），如果不用 `try/finally` 包裹，异常会直接从 `yield` 处往外冒，跳过后面的清理代码：

```python
@contextmanager
def bad_context():
    print("进入")
    yield
    print("退出")   # 如果 with 块内抛异常，这行永远不会执行！

with bad_context():
    raise ValueError("出错")
# 输出只有: 进入   然后异常直接抛出，"退出" 没打印
```

如果想在 `@contextmanager` 版本里实现"吞掉异常"的效果，对应的是捕获异常但不重新抛出：

```python
@contextmanager
def suppress_value_error():
    try:
        yield
    except ValueError as e:
        print(f"忽略了: {e}")
    # 不 re-raise，等价于 __exit__ 返回 True
```

## 6. `ExitStack`：动态管理多个上下文管理器

有时候你不知道运行时要打开多少个资源（比如不定数量的文件），`ExitStack` 可以动态地把多个上下文管理器"压栈"，统一在退出时按**后进先出**顺序清理：

```python
from contextlib import ExitStack

filenames = ["a.txt", "b.txt", "c.txt"]

with ExitStack() as stack:
    files = [stack.enter_context(open(fname)) for fname in filenames]
    # 所有文件都已打开，且离开 with 块时会自动全部 close，顺序是后开的先关
```

## 7. 异步版本：`__aenter__` + `__aexit__`

跟之前讲的 `AsyncIterator` 是 `Iterator` 的异步平行版一样，`with` 也有异步版本 `async with`，对应协议是 `__aenter__`/`__aexit__`，两者都必须是 `async def`（这点和 `__aiter__` 不同，`__aiter__` 必须同步，但 `__aenter__`/`__aexit__` 是真的需要异步等待，所以必须是协程方法）：

```python
class AsyncResource:
    async def __aenter__(self):
        await asyncio.sleep(0.1)   # 模拟异步建立连接
        print("异步资源就绪")
        return self

    async def __aexit__(self, exc_type, exc_value, traceback):
        await asyncio.sleep(0.1)   # 模拟异步关闭连接
        print("异步资源已释放")
        return False

async def main():
    async with AsyncResource() as res:
        print("使用资源")

asyncio.run(main())
```

对应也有 `@asynccontextmanager` 简化写法：

```python
from contextlib import asynccontextmanager

@asynccontextmanager
async def async_context():
    print("进入")
    try:
        yield "value"
    finally:
        print("退出")

async def main():
    async with async_context() as val:
        print(val)

asyncio.run(main())
```

这也回应了之前你问的"`__init__` 不能异步"的问题——**`__aenter__` 正是官方推荐的"异步初始化"落地方式**之一：

```python
class Database:
    def __init__(self, dsn):
        self.dsn = dsn
        self.connection = None

    async def __aenter__(self):
        self.connection = await establish_connection(self.dsn)   # 真正的异步初始化在这里
        return self

    async def __aexit__(self, *args):
        await self.connection.close()

async def main():
    async with Database("postgres://...") as db:
        ...
```

## 8. 一次 `with` 支持多个上下文管理器

```python
with open("a.txt") as fa, open("b.txt") as fb:
    ...
# 等价于嵌套的 with open("a.txt") as fa:
#            with open("b.txt") as fb:
```

Python 3.10+ 还支持圆括号换行写法，更清晰：

```python
with (
    open("a.txt") as fa,
    open("b.txt") as fb,
):
    ...
```

## 9. 协议对照小结

| | 同步 | 异步 |
|---|---|---|
| 协议方法 | `__enter__` / `__exit__` | `__aenter__` / `__aexit__` |
| 语句 | `with` | `async with` |
| 生成器简化写法 | `@contextlib.contextmanager` | `@contextlib.asynccontextmanager` |
| 方法是否必须是协程 | 都是普通方法 | 都必须是 `async def` |

## 一句话总结

上下文管理器的核心价值就一句话：**把"进入前要做的事"和"不管发生什么、离开时都必须做的事"绑定在一起**，靠 `__enter__`/`__exit__`（或异步版 `__aenter__`/`__aexit__`）这套协议，让 `with` 语句自动帮你处理"正常退出"和"异常退出"两种情况下都要执行清理逻辑这件事，本质上是对 `try/finally` 的一层更简洁、更难被误用（比如忘记写 finally）的封装。