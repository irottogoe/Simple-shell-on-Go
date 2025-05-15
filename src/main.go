package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		// Читаем ввод с клавиатуры
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Обработка ввода
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	// Удаляем /n
	input = strings.TrimRight(input, "\r\n")

	// Разделяем ввод на аргументы
	args := strings.Split(input, " ")

	// Проверяем наличие аргументов
	switch args[0] {
	case "cd":
		// Если cd без аргументов, то возвращаем ошибку
		if len(args) < 2 {
			return errors.New("path required")
		}
		// Если cd с аргументами, то меняем директорию
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	// Передаем программу и аргументы отдельно
	cmd := exec.Command(args[0], args[1:]...)

	// Правильное устройство вывода
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Выполняем программу и возвращаем ошибку
	return cmd.Run()
}
