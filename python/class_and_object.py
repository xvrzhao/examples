"""
Python 中类与实例的关系
"""

class Person:

    # 直接在类中声明的属性为类属性
    name: str
    age: int = 30
    job: str = "saxophonist"

    def __init__(self, name, age, gender = "male"):
        # 只有通过 self 赋值的属性才是对象属性
        self.name = name
        self.age = age
        self.gender = gender

    @classmethod
    def grow(cls):
        cls.age += 1

    @staticmethod
    def greet():
        print(f"I'm {Person.name}.")


if __name__ == "__main__":

    # 结论 1: 类属性和对象属性互相隔离、互不干扰，即使同名也没关系
    Person.name = "John"
    person = Person("Coltrane", 40)
    print("class:", Person.name, Person.age) # class: John 30
    print("object:", person.name, person.age) # object: Coltrane 40

    # 结论 2: 实例可以直接通过自身访问到类属性、类方法
    print("class property", person.job) # class property job: saxophone
    Person.grow()
    person.grow()
    print(Person.age, person.age) # 32, 40

    # 结论 3: 类和实例都可以直接通过自身访问到静态方法
    Person.greet()
    person.greet()