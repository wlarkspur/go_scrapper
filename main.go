package main

import (
	"fmt"
	"learngo/mydict"
)

/* type person struct {
	name    string
	age     int
	favFood []string
} */

/* func main() {
	account := accounts.NewAccount("jack")
	account.Deposit(10)

	fmt.Println(account)
}
*/

func main() {
	dictionary := mydict.Dictionary{"first": "First word"}
	word := "Hello"
	definition := "Greeting"
	err := dictionary.Add(word, definition)
	if err != nil {
		fmt.Println(err)
	}
	hello, _ := dictionary.Search(word)
	fmt.Println("Found", word, "definition:", hello)
	err2 := dictionary.Add(word, definition)
	if err2 != nil {
		fmt.Println(err2)
	}
}
