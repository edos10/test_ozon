package greetings

import "fmt"

// Hello возвращает приветствие для указанного человека.
func Hello(name string) string {
	// Возвращаем приветствие, включающее имя в сообщение.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
