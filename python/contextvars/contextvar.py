"""
上下文变量 ContextVar
"""

from contextvars import ContextVar

# 声明一个上下文变量
# 类型为 ContextVar[str]
# 名称为 user
# 默认值为 空字符串
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

