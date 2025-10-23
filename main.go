package main

import (
	"encoding/json"
	"fmt"
	color "github.com/fatih/color"
	"math/rand"
	"os"
	"time"
)

var (
	win  = color.New(color.FgGreen, color.Bold).PrintfFunc()
	lose = color.New(color.FgRed, color.Bold).PrintfFunc()
	warn = color.New(color.FgYellow).PrintfFunc()
	attS = color.New(color.FgYellow).SprintFunc()
	info = color.New(color.FgCyan).PrintfFunc()
)

type Level struct {
	name       string
	maxNumber  int
	maxAttempt int
}

type Result struct {
	Date     string `json:"data"`
	Result   string `json:"result"`
	Attempts int    `json:"attempts"`
	Level    string `json:"level"`
}

func inArray(arr []int, x int) bool {
	for _, n := range arr {
		if n == x {
			return true
		}
	}
	return false
}

func enterLevel() int {
	fmt.Printf("Введите уровень сложности\n1 - Easy: 1–50, 15 попыток\n2 - Medium: 1–100, 10 попыток\n3 - Hard: 1–200, 5 попыток\nТвой выбор: ")

	var choice int

	if _, err := fmt.Scan(&choice); err != nil || choice > 3 || choice < 1 {
		fmt.Println("Введите число от 1 до 3")
		var discard string
		fmt.Scanln(&discard)
		return enterLevel()
	}

	return choice
}

func switchLevel(choice int) (Level, bool) {
	switch choice {
	case 1:
		return Level{"Easy", 50, 15}, true
	case 2:
		return Level{"Medium", 100, 10}, true
	case 3:
		return Level{"Hard", 200, 5}, true
	default:
		return Level{}, false
	}
}

func readGuess() int {
	for {
		var x int
		if _, err := fmt.Scan(&x); err != nil {
			fmt.Printf("Введите число!\n")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		return x
	}
}

func playGuess(lvl Level) {
	number := rand.Intn(lvl.maxNumber) + 1
	var user_number int
	var user_attemps []int
	var count int

	info("Введите число от 1 до %d, которое я загадал\n", lvl.maxNumber)

	for count = 1; count <= lvl.maxAttempt; {
		if user_attemps != nil {
			fmt.Println("Попытки:", attS(fmt.Sprint(user_attemps)))
		}

		fmt.Printf("%d попытка: ", count)
		user_number = readGuess()

		if inArray(user_attemps, user_number) {
			fmt.Println("Вы уже вводили это число")
			continue
		}
		if user_number <= 0 || user_number > lvl.maxNumber {
			fmt.Printf("Введите число в диапазоне!\n")
			continue
		}

		count++
		diff := absInt(user_number - number)
		switch {
		case user_number > number:
			fmt.Printf("Число %s загаданного\n", attS("больше"))
			switchDiff(diff)
			user_attemps = append(user_attemps, user_number)
			continue
		case user_number < number:
			fmt.Printf("Число %s загаданного\n", attS("меньше"))
			switchDiff(diff)
			user_attemps = append(user_attemps, user_number)
			continue
		default:
			number = -1
			win("Вы угадали\n")
			saveResult(Result{Date: time.Now().Format("15:04:05 02-01-2006"), Result: "WIN", Attempts: count, Level: lvl.name})
			return
		}

	}
	if count == lvl.maxAttempt+1 && number > 0 {
		lose("Вы проиграли\nЗагаданное число: %d\n", number)
		saveResult(Result{Date: time.Now().Format("15:04:05 02-01-2006"), Result: "LOSE", Attempts: count - 1, Level: lvl.name})
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func switchDiff(diff int) {
	switch {
	case diff <= 5:
		warn(" 🔥 Горячо \n")
	case diff <= 15:
		info("🙂 Тепло \n")
	case diff > 15:
		fmt.Printf("❄️ Холодно \n")
	}
	return
}

func saveResult(r Result) {
	const filename = "result.json"

	var results []Result

	if data, err := os.ReadFile(filename); err == nil && len(data) > 0 {
		json.Unmarshal(data, &results)
	}

	results = append(results, r)

	jsonData, _ := json.MarshalIndent(results, "", "	")
	os.WriteFile(filename, jsonData, 0665)
}

func nextAction() int {
	fmt.Printf("\nСыграть заново?\n1 - Сыграть на том же уровне сложности\n2 - Выбрать уровень сложности\n3 - Выход\n")
	var choice int
	if _, err := fmt.Scan(&choice); err != nil || choice > 3 || choice < 1 {
		fmt.Printf("Введите число от 1 до 3!\n")
		var discard string
		fmt.Scanln(&discard)
		return nextAction()
	}
	return choice
}

func main() {
	var lvl Level
	var ok bool
	var needLevel bool = true
	for {
		if needLevel {
			lvl, ok = switchLevel(enterLevel())
		}
		if ok {
			playGuess(lvl)
		}

		action := nextAction()
		switch action {
		case 1:
			needLevel = false
			continue
		case 2:
			needLevel = true
			continue
		case 3:
			fmt.Println("Пока!")
			return
		}
	}
}
