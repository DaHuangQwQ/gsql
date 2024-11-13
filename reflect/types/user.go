package types

import "fmt"

type User struct {
	Name string
	// 因为同属一个包，所以 age 还可以被测试访问到
	// 如果是不同包，就访问不到了
	age int
}

func NewUser(name string, age int) User {
	return User{
		Name: name,
		age:  age,
	}
}

func NewUserPtr(name string, age int) *User {
	return &User{
		Name: name,
		age:  age,
	}
}

// func GetAge(u User) int {
//
// }
func (u User) GetAge() int {
	return u.age
}

func (u *User) ChangeName(newName string) {
	u.Name = newName
}

func (u User) private() {
	fmt.Println("private")
}
