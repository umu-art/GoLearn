package main

import (
	"fmt"
	"strings"
)

type Task struct {
	taskName     string
	taskCallable func()
}

var tasks = []Task{
	{
		taskName: "Привет, мир!",
		taskCallable: func() {
			println("Привет мир!")
		},
	},
	{
		taskName: "Сложение чисел",
		taskCallable: func() {
			var a, b int
			_, _ = fmt.Scan(&a, &b)
			println(a + b)
		},
	},
	{
		taskName: "Четное или нечетное",
		taskCallable: func() {
			var a int
			_, _ = fmt.Scan(&a)
			if a%2 == 0 {
				println("Четное")
			} else {
				println("Нечетное")
			}
		},
	},
	{
		taskName: "Максимум из трех чисел",
		taskCallable: func() {
			var arr [3]int
			for i := 0; i < 3; i++ {
				_, _ = fmt.Scan(&arr[i])
			}
			var result = arr[0]
			for i := 1; i < 3; i++ {
				if arr[i] > result {
					result = arr[i]
				}
			}
			println(result)
		},
	},
	{
		taskName: "Факториал числа",
		taskCallable: func() {
			var a int
			_, _ = fmt.Scan(&a)
			var result = 1
			for i := 1; i <= a; i++ {
				result *= i
			}
			println(result)
		},
	},
	{
		taskName: "Проверка символа",
		taskCallable: func() {
			var glasniye = "aeiouyаеёиоуыэюя"
			var input string
			_, _ = fmt.Scan(&input)

			var lowerInput = strings.ToLower(input)[:1]

			if strings.Contains(glasniye, lowerInput) {
				println("Гласная")
			} else {
				println("Согласная")
			}
		},
	},
	{
		taskName: "Простые числа",
		taskCallable: func() {
			var maxNumber int
			_, _ = fmt.Scan(&maxNumber)

			/*
			 Решето:
			*/
			var easyNumbers = make([]bool, maxNumber+1)
			for i := 2; i <= maxNumber; i++ {
				if easyNumbers[i] {
					continue
				}
				print(i, " ")
				for j := i * i; j <= maxNumber; j += i {
					easyNumbers[j] = true
				}
			}

			println()
		},
	},
	{
		taskName: "Строка и ее перевертыш",
		taskCallable: func() {
			var input string
			_, _ = fmt.Scan(&input)
			println(reverse(input))
		},
	},
	{
		taskName: "Массив и его сумма",
		taskCallable: func() {
			var n int
			_, _ = fmt.Scan(&n)
			var arr = make([]int, n)
			for i := 0; i < n; i++ {
				_, _ = fmt.Scan(&arr[i])
			}

			println(getSum(arr))
		},
	},
	{
		taskName: "Структуры и методы",
		taskCallable: func() {
			var a, b float32
			_, _ = fmt.Scan(&a, &b)

			var r = Rectangle{a, b}
			println(r.getSquare())
		},
	},
	{
		taskName: "Конвертер температур",
		taskCallable: func() {
			println("Давай в следующий спринт... :)")
		},
	},
}

func reverse(str string) string {
	var result = ""
	for i := len(str) - 1; i >= 0; i-- {
		result += string(str[i])
	}
	return result
}

func getSum(arr []int) int {
	var result = 0
	for i := 0; i < len(arr); i++ {
		result += arr[i]
	}
	return result
}

type Rectangle struct {
	weight float32
	height float32
}

func (r Rectangle) getSquare() float32 {
	return r.weight * r.height
}

func main() {
	println("Сделанные (наверное) задания:")

	println()
	for i := 0; i < len(tasks); i++ {
		println(i+1, ":", tasks[i].taskName)
	}
	println()

	print("Введите номер задания: ")
	var selectedTaskNum int
	_, _ = fmt.Scan(&selectedTaskNum)
	selectedTaskNum--

	if selectedTaskNum < 0 || selectedTaskNum >= len(tasks) {
		println("Неверный номер задания")
		return
	}

	println("Исполняю задание \"", tasks[selectedTaskNum].taskName, "\":")
	tasks[selectedTaskNum].taskCallable()
}
