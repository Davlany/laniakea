package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Открываем файл .ans
	file, err := os.Open("2.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Создаем сканер для чтения файла построчно
	scanner := bufio.NewScanner(file)

	// Считываем и выводим каждую строку файла
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Проверяем наличие ошибок при сканировании файла
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при сканировании файла:", err)
	}
}
