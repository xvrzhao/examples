"""
Generator 生成器
"""

from typing import Generator, Iterator, Iterable

def count(stop: int):
    """
    生成器函数，该函数返回一个生成器对象。
    凡是函数体中包含 yield 关键字的函数都是生成器函数。
    """
    n = 0
    while n < stop:
        yield n
        n += 1


def generator_is_iterator():
    """
    生成器内部实现了 Iterator 协议，包含 __iter__ 和 __next__ 方法，
    所以 Generator 也是一个 Iterator。
    """
    gen = count(5)                      # count 是生成器函数，调用函数返回一个生成器
    print(type(gen))                    # <class 'generator'>
    print(isinstance(gen, Generator))   # True

    # 重点：生成器内部实现了 Iterator 协议，包含 __iter__ 和 __next__ 方法
    print(hasattr(gen, "__iter__"))     # True
    print(hasattr(gen, "__next__"))     # True
    print(isinstance(gen, Iterator))    # True
    print(isinstance(gen, Iterable))    # True

    # 消费迭代器
    for i in gen:
        print(i, end=" ") # print: 0 1 2 3 4

    print()

    # 再次消费，无任何输出，记住：iterator 有状态，无法重复消费，可以回去看看第一节
    for i in gen:
        print(i, end="")

    print()

    # count(5) 可以被重复消费，因为每次调用都会重新返回一个 iterator
    for i in count(5):
        print(i, end=",") # print: 0,1,2,3,4,

    print()

    for i in count(5):
        print(i, end=".") # print: 0.1.2.3.4.

    print()


def gen() -> Iterator[int]:
    """如果 gen 函数只是用来被 for 消费，可以将函数返回值直接标注为 Iterator 而不是 Generator 更直观"""
    yield 1
    yield 2
    yield 3


def relationship_of_iterable_iterator_and_generator():
    """
    generator, iterator 和 iterable 之间的关系：
        generator ⊂ iterator ⊂ iterable
    """


def the_generator_version_of_CountDownIterator():
    """
    第一节中 CountDownIterator 迭代器的生成器版本。
    """

    # 可见使用 generator 的方式书写要简洁的多
    def count_down(n: int):
        while n > 0:
            n -= 1
            yield n

    for i in count_down(5):
        print(i)


def usage_of_yield_from():
    """yield from 关键字的用法"""

    def outer():
        yield from count(5) # 将一个 generator (具体来说是 iterable) 的值透传出来。等价于：for x from count(5): yield x
        yield 5
        yield 6

    for i in outer():
        print(i, end=" ") # print: 0 1 2 3 4 5 6


    print()


    def outer2():
        yield from [1, 2, 3] # 这次是从一个 iterable 中透传（第一节中有讲 list 是 iterable 而不是 iterator）
        yield 4
        yield 5

    for i in outer2():
        print(i, end=",") # print: 1,2,3,4,5,


def usage_of_send():
    """
    yield 不仅可以向外传值，还可以接受外部的传值。通过 send 可以向 generator 内部传值。
    (2026.09.10: 该应用场景目前还没怎么见过，先举例最简单的用法，以后见到实际使用案例再进一步讲解)
    """

    def receive_value_from_outer():
        while True:
            received_value = yield
            print("received:", received_value)

    g = receive_value_from_outer()
    print(type(g)) # <class 'generator'> 可见通过 yield 接收外部的值，函数返回的也是 generator

    next(g) # 需要先通过 next 启动，执行到第一次 yield 的位置。也可以使用 g.send(None)，第一次启动时必须传 None
    g.send("hello") # 在 generator 内部，从 yield 继续往下执行，直到下一个 yield 处暂停
    g.send("world")
    for i in range(1, 5):
        g.send(i)

    """
    输出:
    received: hello
    received: world
    received: 1
    received: 2
    received: 3
    received: 4
    """


if __name__ == "__main__":
    generator_is_iterator()
    relationship_of_iterable_iterator_and_generator()
    the_generator_version_of_CountDownIterator()
    usage_of_yield_from()
    usage_of_send()
    pass