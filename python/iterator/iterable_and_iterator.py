"""
Iterable:
1. Iterable 是可迭代对象
2. 只需实现 __iter__ 方法, 该方法返回一个 Iterator 迭代器

Iterator:
1. Iterator 是迭代器，需要实现两个方法，__iter__ 和 __next__
2. 其中 __iter__ 方法返回自身 self, 所以 Iterator 也是 Iterable
3. __next__ 方法则是返回下一个迭代值, 若没有新的迭代值, 则 raise StopIteration 异常

for 循环原理:
1. `for ... in obj` 中 obj 只需要是 Iterable 对象即可，也就是需包含 __iter__ 方法, Iterator 因为也是 Iterable 所以也可以被 for 循环
2. for 循环会先调用 iter(obj), 内部执行 obj.__iter__(), 拿到一个 Iterator 对象，我们命名为 it
3. 之后会循环执行 next(it), 也就是调用 it.__next__(), 直到遇到 StopIteration 异常停止

```
for v in obj:
    print(v)

# 等价于

it = iter(obj) # 内部调用 obj.__iter__() 拿到一个 Iterator
while True:
    try:
        v = next(it) # 内部调用 it.__next__()
        print(v)
    except StopIteration:
        break
```
"""

from typing import Iterable, Iterator

def list_is_iterable():
    """
    list 对象只是一个数据容器，但是内部包含 __iter__ 方法，所以只是一个 Iterable 并不是 Iterator。
    这样的好处是，在对 list 对象进行 for 循环时，list 对象每次只需要返回一个新的 Iterator 对象即可，
    Iterator 对象内部维护一个索引的推进状态，list 本身不需要维护推进状态，这样重复多次对 list 进行遍历，或即使是并发对 list 进行遍历，
    都会从 list 的第一个值开始。

    这也就是为什么要区分 Iterable 和 Iterator 的原因，即 Iterable 不需要维护迭代推进状态，只需要返回一个全新的 Iterator 对象，而 Iterator 对象内部维护一个推进状态。
    一句话概括：Iterable 无状态，Iterator 有状态。
    """
    lst = [1, 2, 3]

    is_iterable = isinstance(lst, Iterable)
    print(is_iterable) # True

    is_not_iterator = isinstance(lst, Iterator)
    print(is_not_iterator) # False

    print(hasattr(lst, '__iter__')) # True
    print(hasattr(lst, '__next__')) # False

    lst_iterator = iter(lst)
    print(isinstance(lst_iterator, Iterator)) # True
    print(hasattr(lst_iterator, '__iter__')) # True
    print(hasattr(lst_iterator, '__next__')) # True


class CountDown:
    """CountDown 实现了 Iterable"""
    def __init__(self, step):
        self.step = step

    def __iter__(self):
        return CountDownIterator(self.step)
    
class CountDownIterator:
    """CountDownIterator 实现了 Iterator"""
    def __init__(self, step):
        self.step = step

    def __iter__(self):
        return self

    def __next__(self):
        if self.step <= 0:
            raise StopIteration
        self.step -= 1
        return self.step


def iterable_and_iterator():
    """
    实际案例区分 Iterable 和 Iterator
    先说结论：如果只被迭代消费一次可以选择 Iterator，如果需要多次从头开始迭代消费则选择 Iterable
    """

    cd = CountDown(5)
    print(isinstance(cd, Iterable)) # True
    print(isinstance(cd, Iterator)) # Fasle

    # 输出: 4 3 2 1 0
    for i in cd: 
        print(i, end=" ")

    print()

    # 再次遍历仍然输出：4 3 2 1 0
    # 因为每次 iter(cd) 都会返回一个新的 Iterator
    for i in cd:
        print(i, end=" ")

    print()

    """ 直接对 Iterator 进行遍历会是什么后果 """

    cdi = CountDownIterator(5)
    print(isinstance(cdi, Iterable)) # True
    print(isinstance(cdi, Iterator)) # True

    # 首次遍历输出：4 3 2 1 0
    for i in cdi:
        print(i, end=" ") 

    # 再次遍历，无任何输出，因为 cdi 内部推进状态已经到头了
    for i in cdi:
        print("nothing")

    print()


if __name__ == "__main__":
    # list_is_iterable()
    iterable_and_iterator()