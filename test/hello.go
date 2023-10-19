package main

import "fmt"

var helloPrefixes = map[string]string{
	"English":"Hello, ",
	"Spanish":"Hola, ",
	"French":"Bonjour, ",
}


func main() {
	fmt.Println(Hello("James", "English"))
}

func Hello(name string, language string) string {
	if name =="" {
		name = "world" 
	}
	return fmt.Sprint(getPrefix(language), name)
}

func getPrefix(language string) string {
	prefix:= helloPrefixes[language]
	if prefix==""{
		return helloPrefixes["English"]
	}
	return prefix
}
