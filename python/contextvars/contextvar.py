"""
ContextVar 上下文变量

同一个上下文变量，在不同协程中
"""
from contextvars import ContextVar

# 声明一个上下文变量，类型为 ContextVar[str]，名称为 user，默认值为空字符串
user_var: ContextVar[str] = ContextVar("user", default="")

print(user_var.get())