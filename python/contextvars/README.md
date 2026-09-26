# 上下文变量 contextvars

这个模块是专门为异步/并发场景设计的，解决的问题是：**在协程之间传递"上下文相关"的数据，同时保证不同并发任务之间互不干扰**。

## 1. 先看它要解决的问题：线程局部变量在异步场景下失效

在同步多线程编程里，Python 有 `threading.local()` 来实现"每个线程有自己独立的一份变量"：

```python
import threading

local_data = threading.local()

def worker(name):
    local_data.user = name
    print(local_data.user)

threading.Thread(target=worker, args=("Alice",)).start()
threading.Thread(target=worker, args=("Bob",)).start()
# 各自打印各自的 name，互不干扰，因为每个线程有独立的存储
```

但异步编程里，**所有协程跑在同一个线程里**，`threading.local()` 完全没用——因为它是按"线程"隔离的，而一个线程内可能同时交替执行成百上千个协程：

```python
import asyncio, threading

local_data = threading.local()

async def handle_request(user):
    local_data.user = user
    await asyncio.sleep(0.1)   # 让出控制权，另一个协程可能在这期间抢先跑
    print(local_data.user)     # 可能打印出别的协程设置的值！

async def main():
    await asyncio.gather(
        handle_request("Alice"),
        handle_request("Bob"),
    )

asyncio.run(main())
# 可能输出: Bob Bob （两次都是 Bob，Alice 的值被覆盖了！）
```

问题很明确：两个协程实际上跑在**同一个线程**里，`local_data` 是共享的，`handle_request("Alice")` 在 `await asyncio.sleep(0.1)` 让出控制权后，`handle_request("Bob")` 抢着把 `local_data.user` 改成了 "Bob"，等 Alice 的协程恢复执行时，读到的已经是被污染的值。

## 2. `ContextVar` 怎么解决这个问题

`contextvars.ContextVar` 提供的是**按"逻辑执行流"隔离**，而不是按线程隔离——每个协程 task 都有自己独立的一份"上下文"（Context），`ContextVar` 的值在这份上下文里读写，不会互相污染：

```python
import asyncio
from contextvars import ContextVar

user_var = ContextVar("user")   # 声明一个上下文变量，名字仅用于调试显示

async def handle_request(user):
    user_var.set(user)
    await asyncio.sleep(0.1)
    print(user_var.get())   # 每个协程读到的都是自己设置的值

async def main():
    await asyncio.gather(
        handle_request("Alice"),
        handle_request("Bob"),
    )

asyncio.run(main())
# 输出: Alice Bob  —— 正确隔离，互不干扰
```

## 3. 核心 API

```python
from contextvars import ContextVar

# 创建，可以给默认值
user_var = ContextVar("user", default=None)

user_var.set("Alice")       # 设置当前上下文里的值
user_var.get()               # 'Alice'，读取当前上下文的值
user_var.get("默认值")       # 如果没设置过，也没给 default，可以在 get() 时传兜底值

token = user_var.set("Bob")  # set() 返回一个 token，记录了"改之前"的状态
user_var.reset(token)         # 用 token 精确地恢复到 set 之前的值（比重新 set 更安全）
```

`reset(token)` 比直接再 `set()` 回旧值更严谨，因为它能正确处理"这个变量之前到底有没有被设置过"的边界情况（比如原来是"未设置"状态，`reset` 能恢复成"未设置"，而不是错误地设成某个默认值）。

## 4. 背后的机制：`Context` 对象和"隐式传递"

**关键概念：`Context`**。每次调用 `asyncio.create_task()`（或者 `loop.call_soon()` 等）创建一个新的并发执行单元时，Python 会**自动复制当前的 Context**，作为新任务的初始上下文：

```python
import asyncio
from contextvars import ContextVar

var = ContextVar("var", default="根节点")

async def child():
    print("子任务读到:", var.get())   # 继承了父任务当时的值
    var.set("子任务改的值")
    print("子任务改后:", var.get())

async def main():
    var.set("主任务设置的值")
    task = asyncio.create_task(child())
    await task
    print("主任务读到:", var.get())   # 不受子任务影响，仍是主任务自己设置的值

asyncio.run(main())
```

输出：

```
子任务读到: 主任务设置的值
子任务改后: 子任务改的值
主任务读到: 主任务设置的值
```

这里体现了两个关键特性：

1. **子任务能读到父任务设置的值**（创建 task 时复制了当时的 context 快照）
2. **子任务修改后不会影响父任务**（复制出来的是独立的一份，子任务在自己的副本上改）

这个模型跟"每个协程 task 拥有一份从创建时刻继承下来、之后独立演化的上下文快照"很像——有点类似 fork 出一个子进程，子进程改自己的内存不影响父进程。

## 5. 手动操作 Context：`copy_context()` + `run()`

`ContextVar` 底层依赖的 `Context` 对象也可以手动拿到、手动运行：

```python
import contextvars

var = contextvars.ContextVar("var")
var.set("原始值")

ctx = contextvars.copy_context()   # 拷贝当前上下文

def modify():
    var.set("在拷贝的上下文里改的值")
    print("内部:", var.get())

ctx.run(modify)   # 在拷贝出的上下文里执行这个函数
print("外部:", var.get())   # 不受影响
```

输出：

```
内部: 在拷贝的上下文里改的值
外部: 原始值
```

`asyncio.create_task()` 内部本质上就是做了类似 `copy_context()` 的操作，给每个 task 一份独立的上下文快照。

## 6. 一个实际场景：Web 服务里的 request_id / 用户身份追踪

这是 `contextvars` 最经典的应用场景，跟你在做的多用户 agent backend 高度相关——**不用每层函数调用都手动传参，就能让底层的日志、数据库查询等代码拿到当前请求的上下文信息**：

```python
from contextvars import ContextVar
import uuid

request_id_var = ContextVar("request_id", default=None)

def log(message):
    rid = request_id_var.get()
    print(f"[{rid}] {message}")

async def handle_request(user):
    request_id_var.set(str(uuid.uuid4())[:8])
    log(f"开始处理 {user} 的请求")
    await process(user)

async def process(user):
    # 注意：这个函数完全不需要接收 request_id 参数，
    # 但内部调用 log() 依然能拿到正确的、当前这个请求的 request_id
    log(f"正在处理 {user}")

async def main():
    await asyncio.gather(
        handle_request("Alice"),
        handle_request("Bob"),
    )
```

如果没有 `contextvars`，想做到"每条日志都带上当前请求的 request_id"，通常得把 `request_id` 当参数一层一层往下传，或者用 `threading.local()`（但正如前面说的，在异步场景下这是错的）。

**FastAPI/Starlette 生态里的 middleware 常用这个模式**：在最外层中间件里 `set()` 一个 `request_id` 或当前用户信息，后续所有的路由处理函数、依赖注入、日志记录，都可以在任意深度的函数里直接 `get()` 到，不需要显式传参，天然还带着"每个请求互不干扰"的隔离性。

## 7. 和 `threading.local()` 的关键区别总结

| | `threading.local()` | `contextvars.ContextVar` |
|---|---|---|
| 隔离粒度 | 按线程 | 按"逻辑执行流"（每个 asyncio task / 协程） |
| 异步场景下是否安全 | ❌ 不安全，同线程内所有协程共享 | ✅ 安全，专门为此设计 |
| 子任务是否自动继承父级值 | 不适用（是新线程，默认不继承） | ✅ 创建 task 时自动复制当前 context |
| 典型场景 | 多线程程序里的线程私有状态 | asyncio 协程里的请求级上下文（request_id、当前用户、trace 信息） |

## 8. 一个容易踩的坑：在协程里 `set()` 之后再 `create_task`，顺序有讲究

```python
async def main():
    task = asyncio.create_task(child())   # 创建 task 时复制的是"此刻"的 context
    var.set("之后才设置的值")               # 这个 set 发生在 create_task 之后
    await task

async def child():
    print(var.get())   # 读不到"之后才设置的值"！因为复制发生在 set 之前
```

因为 context 的复制是在 `create_task()` **调用那一刻**发生的，之后主协程再修改自己的 context，不会回溯影响已经创建好的子任务副本。这个坑在实际开发中比较隐蔽，尤其是中间件里"先创建后台任务，再设置上下文"这种顺序错误的写法。

## 一句话总结

`contextvars` 解决的是异步并发场景下"每个协程需要有自己独立、互不干扰的上下文数据，同时又要能被创建时自动继承父级快照"这个需求，是 `threading.local()` 在异步世界里的正确替代品。核心机制是：每次 `create_task()`（或类似的调度点）都会**复制当前的 `Context` 快照**给新任务，之后各自在自己的副本上独立读写，`ContextVar.get()`/`set()` 操作的都是"当前所在的那份 Context"里的值。这套机制是现代异步 Web 框架里做请求级追踪（request_id、当前用户、trace context）的标准工具。


# 相关问题解答

## 1. ContextVar 是 Task 级别隔离吗？

**严格说，`contextvars` 隔离的单位是 `Context` 对象，不是 Task 本身——但由于"复制 Context"这个动作,只在创建新 Task（或类似调度点）时发生，所以在实践中，隔离的边界确实和 Task 的边界重合。**

细化一下机制：

```python
class Task:
    def __init__(self, coro):
        self._context = contextvars.copy_context()   # Task 创建时，复制一份当前 Context
        ...

    def __step(self):
        # Task 每次被事件循环唤醒、推进执行时，
        # 都是在这份"自己专属的 Context"里跑的
        self._context.run(self._coro.send, None)
```

这解释了两件事：

- **一个 Task 内部，不管套多少层 `await coroutine`，大家共用的都是同一个 Context**（因为压根没有创建新 Task，也就没有新的 `copy_context()` 发生），所以 `ContextVar.set()` 在内层协程改的值，外层协程能读到——这跟你上一个问题里验证的"`current_task()` 全程都是同一个 Task"是同一件事的两个侧面。
- **只有创建新 Task 那一刻，才会触发一次 Context 复制**，产生一份新的、独立的快照，之后这个新 Task 在自己的 Context 里独立读写，不回流影响父级。

所以更精确的表述是：**ContextVar 隔离的最小单位是 Context，而"Context 何时被复制、产生隔离边界"这件事，是跟着 Task 的创建走的**。二者在 asyncio 的默认使用方式下几乎总是重合的，所以说"Task 级别隔离"作为一种实用理解是对的。

**一个能体现"本质是 Context 而非 Task"的边界情况**：你可以完全不创建 Task，只手动用 `contextvars.copy_context()` + `ctx.run()` 就能拿到独立的隔离效果：

```python
import contextvars

var = contextvars.ContextVar("var", default="root")

def modify():
    var.set("独立上下文里的值")
    print(var.get())

ctx = contextvars.copy_context()
ctx.run(modify)
print(var.get())   # 'root'，不受影响 —— 隔离发生了，但全程没有涉及任何 asyncio Task
```

这说明"隔离"这件事的根本载体是 Context，Task 只是 asyncio 里最常见的、自动帮你做 Context 复制的那个触发点。另外像 `loop.call_soon(callback, context=...)` 这种更底层的调度 API，也可以手动指定用哪个 Context 去跑一个回调，同样不涉及创建 Task。

## 2. 每个用户请求分配独立 Task，是 ASGI 服务器做的还是 FastAPI 做的？

**是 ASGI 服务器（比如 uvicorn）做的，不是 FastAPI 或者 Starlette 应用层代码做的。**

搜索确认了这一点：由于 uvicorn 会为每一个具体的请求创建异步 task，所以预期处理过程中设置的任何 context var 都应该在请求范围内被隔离——这条来自 uvicorn 仓库的 issue 讨论，明确说明了"每个请求一个 task"是 **uvicorn 这一层**做的事情。

具体链路大致是这样：

```
TCP 连接进来
  → uvicorn 的 Protocol 实现（基于 asyncio.Protocol / httptools）解析出一个 HTTP 请求
  → uvicorn 为这次请求 asyncio.create_task(...)，把请求交给 ASGI 应用去处理
      → 这个 Task 内部才真正去调用你的 ASGI app（也就是 FastAPI/Starlette 实例）的 __call__
          → FastAPI/Starlette 在这个 Task 内部，一路 await 路由匹配、依赖注入、
            你写的 async def 视图函数、中间件……全部都是同一个 Task 内的顺序执行
```

也就是说：

| 层次 | 职责 |
|---|---|
| uvicorn（ASGI 服务器） | 负责网络 I/O、HTTP 协议解析，并**为每个请求创建一个 Task**，在这个 Task 里调用 ASGI 应用 |
| FastAPI / Starlette（ASGI 应用） | 在 uvicorn 已经分配好的这个 Task **内部**运行——路由分发、依赖注入、中间件链、你的视图函数，都只是这一个 Task 内的层层 `await` 调用，不会自己额外创建新 Task |

这也解释了为什么"每个请求之间互不干扰"（包括你可能在中间件里用 `ContextVar` 设置 `request_id`）能天然成立——**根本原因是 uvicorn 在最外层已经把每个请求安排到了独立的 Task 里**，FastAPI 层面写的代码只是在这个既定的隔离边界内部运行，它本身不需要（也没有）额外做"为每个请求单独开 Task"这件事。

**需要留意的例外**：FastAPI/Starlette 自己在某些具体功能点上，确实会主动创建**额外的**子 Task，但这些不是"为整个请求分配主 Task"，而是请求处理过程中的局部并发需求，比如：

- `BackgroundTasks`（响应发送后在后台继续跑的任务）
- 某些中间件实现里对"读请求体"和"处理超时"做的并发控制（常见于用 `anyio` 的 task group）

这些都是在 uvicorn 已经建立好的"请求主 Task"内部，视图逻辑执行过程中按需临时开的子任务，跟"一个请求对应一个顶层 Task"这件事是两回事，后者始终是 ASGI 服务器负责的。