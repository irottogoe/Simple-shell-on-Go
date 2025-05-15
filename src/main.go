package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/peterh/liner" // библиотека для поддержки истории команд
)

const historyFile = ".shell_history" // файл для хранения истории команд

func main() {

	// создаем новый ввод с клавиатуры
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true) // Ctrl+C завершает ввод команды (не всю программу)

	// список встроенных команд
	commands := []string{"cd", "exit", "history", "help"}

	// автодополнение
	line.SetCompleter(func(line string) (c []string) {
		// список всех файлов
		files, err := os.ReadDir(".")
		if err == nil {
			for _, file := range files {
				name := file.Name()
				if strings.HasPrefix(name, line) { // если имя файла начинается с введенной строки
					c = append(c, name) // добавляем имя файла в список
				}
			}
		}

		// список всех команд
		for _, command := range commands {
			if strings.HasPrefix(command, line) { // если команда начинается с введенной строки
				c = append(c, command) // добавляем команду в список
			}
		}
		return
	})

	// загружаем историю команд из файла
	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f) // читаем историю из файла и добавляем в историю
		f.Close()
	}

	// список всех предыдущих команд
	var historyList []string

	// сама программа
	for {
		input, err := line.Prompt("> ") // ждем ввода от пользователя
		if err == liner.ErrPromptAborted || err == io.EOF {
			break // если Ctrl+C или Ctrl+D, то выходим из программы
		} else if err != nil {
			fmt.Println("Error:", err)
			continue // если ошибка, то продолжаем цикл
		}

		input = strings.TrimSpace(input) // удаляем пробелы
		if input == "" {
			continue // если пустая строка, то продолжаем цикл
		}

		// обработка встроенной команды history
		if input == "history" {
			for i, cmd := range historyList {
				fmt.Printf("%d: %s\n", i+1, cmd) // выводим все предыдущие команды
			}
			continue
		}

		historyList = append(historyList, input) // добавляем команду в историю
		line.AppendHistory(input)                // сохраняем команду в историю

		// обработка встроенной команды help
		if input == "help" {
			fmt.Println("Встроенные команды:")
			fmt.Println("  cd <путь>     — сменить директорию")
			fmt.Println("  exit          — выйти из оболочки")
			fmt.Println("  history       — показать историю команд")
			fmt.Println("  help          — показать список команд")
			line.AppendHistory(input) // сохраняем команду в историю
			continue
		}

		line.AppendHistory(input) // добавляем команду в историю

		if err := execInput(input); err != nil { // выполняем команду
			fmt.Fprintln(os.Stderr, "Error:", err) // если ошибка, то выводим ее
		}
	}

	// сохраняем историю команд в файл
	if f, err := os.Create(historyFile); err == nil {
		line.WriteHistory(f) // записываем историю в файл
		f.Close()
	} else {
		fmt.Fprintln(os.Stderr, "Error:", err) // если ошибка, то выводим ее
	}
}

func execInput(input string) error {
	// удаляем /r и /n
	input = strings.TrimRight(input, "\r\n")

	// разбиваем строку на аргументы
	args := strings.Fields(input)

	if len(args) == 0 {
		return nil // если пустая строка, то ничего не делаем
	}

	// проверяем наличие аргументов
	switch args[0] {
	case "cd":
		// если cd без аргументов, то возвращаем ошибку
		if len(args) < 2 {
			return errors.New("path required")
		}
		// если cd с аргументами, то меняем директорию
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0) // если exit, то выходим из программы
	}

	// выполняем команду и передаем аргументы
	// args[0] - имя программы, args[1:] - аргументы программы
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin   // ввод из терминала
	cmd.Stderr = os.Stderr // вывод ошибок в терминал
	cmd.Stdout = os.Stdout // вывод результата в терминал

	// выполняем программу и возвращаем ошибку
	return cmd.Run()
}
