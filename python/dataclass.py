import json
from dataclasses import dataclass, asdict
from typing import TypedDict


@dataclass
class Person:
    name: str
    age: int


@dataclass
class Persons:
    persons: list[Person]


class Fruit(TypedDict):
    name: str
    color: str
    is_sweet: bool = True


if __name__ == "__main__":

    # 1. dataclass 类是一个数据集合，每个字段注释类型，但在实例化时并不校验类型

    person = Person(name=123, age=456)
    print(person) # 通过，打印 Person(name=123, age=456)


    # 2. dataclass 装饰器会往类中注入构造方法，将在类中声明的字段作为构造方法的参数，
    #    所以需要注意字段声明顺序，包含默认值的字段应放在后面

    @dataclass
    class Animal:
        name: str
        age: int = 0
    
    animal = Animal("dog")
    print(animal)


    # 3. dataclass 与 TypedDict 的区别：
    #
    #    是否要求包含所有已声明的字段：
    #    - dataclass 实例化时要求所有字段都传入初始值，除非设有默认值；
    #    - TypedDict 本质是一个注释类型的 dict，内部字段是动态的，初始化时可以不传某个字段；
    #
    #    字段默认值：
    #    - dataclass 可以设置字段默认值；
    #    - TypedDict 不需要给字段设置默认值，设置了也没有任何效果；
    #
    #    JSON 序列化：
    #    - TypedDict 因为本质为 dict，所以可以直接序列化；
    #    - dataclass 需要先通过 asdict 转为 dict，才能进行序列化；

    fruit = Fruit(color="red")
    print(fruit) # {'color': 'red'}

    print(json.dumps(fruit))

    # print(json.dumps(person)) # TypeError: Object of type Person is not JSON serializable
    print(json.dumps(asdict(person))) # 先通过 asdict 将 classdata 实例转为 dict，再进行 json 序列化

    persons = Persons(persons=[Person(name="Lateef", age=20)])
    print(json.dumps(asdict(persons))) # 嵌套结构的 classdata 仍然可以转为 dict (classdata 嵌套 pydantic，或者 pydantic 嵌套 classdata 的复杂结构仍需验证如何序列化)
