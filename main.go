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
	dictionary := mydict.Dictionary{}
	baseWord := "hello"
	dictionary.Add(baseWord, "First")
	fmt.Println(dictionary)
	err := dictionary.Delete(baseWord)
	if err != nil {
		fmt.Println(err)
	}
	word, _ := dictionary.Search(baseWord)
	fmt.Println(word)
}
