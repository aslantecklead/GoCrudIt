package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Worker struct {
    ID_WORKER   int
    WORKER_NAME string
    AGE         int
    SALARY      int
    ADRESS      string
}

type WorkerList struct {
    Workers []Worker
}
var db *sql.DB 

func main() {
    connStr := "host=localhost port=5433 dbname=postgres user=postgres password=123 sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr) 
    if err != nil {
        panic(err)
    }
    defer db.Close()


    rows, err := db.Query("select * from Worker")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    Workers := []Worker{}
    for rows.Next() {
        p := Worker{}
        err := rows.Scan(&p.ID_WORKER, &p.WORKER_NAME, &p.AGE, &p.SALARY, &p.ADRESS)
        if err != nil {
            fmt.Println(err)
            continue
        }
        Workers = append(Workers, p)
    }

    workerList := WorkerList{Workers: Workers}

    for {
        fmt.Println("Действие:")
        fmt.Println("1 - Добавить работника")
        fmt.Println("2 - Вывести среднюю зарплату работников")
        fmt.Println("4 - Всех")
        fmt.Println("3 - Выход")

        var num string
        fmt.Scanln(&num)
        ChooseOption(num, &workerList)
    }
}

func ChooseOption(num string, workerList *WorkerList) {
    switch num {
    case "1":
        AddWorker(workerList)
    case "2":
        CalculateAverageSalary(workerList)
    case "3":
        ShowAll(workerList);
    case "4":
        fmt.Println("Выход")
        os.Exit(0)
    default:
        fmt.Println("Неправильный выбор, попробуйте снова.")
    }
}

func AddWorker(workerList *WorkerList) {
    var newWorker Worker
    fmt.Print("Введите имя работника: ")
    fmt.Scanln(&newWorker.WORKER_NAME)
    
    fmt.Print("Введите возраст работника: ")
    fmt.Scanln(&newWorker.AGE)
    
    fmt.Print("Введите зарплату работника: ")
    fmt.Scanln(&newWorker.SALARY)
    
    fmt.Print("Введите адрес работника: ")
    fmt.Scanln(&newWorker.ADRESS)

    _, err := db.Exec("INSERT INTO Worker (WORKER_NAME, AGE, SALARY, ADRESS) VALUES ($1, $2, $3, $4)",
        newWorker.WORKER_NAME, newWorker.AGE, newWorker.SALARY, newWorker.ADRESS)
    if err != nil {
        fmt.Println("Ошибка при добавлении работника:", err)
    } else {
        fmt.Println("Работник успешно добавлен!")
        workerList.Workers = append(workerList.Workers, newWorker)
    }
}

func ShowAll(workerList *WorkerList) {
    if len(workerList.Workers) == 0 {
        fmt.Println("Список работников пуст.")
    } else {
        fmt.Println("Список всех работников:")
        for _, worker := range workerList.Workers {
            fmt.Printf("ID: %d, Имя: %s, Возраст: %d, Зарплата: %d, Адрес: %s\n", worker.ID_WORKER, worker.WORKER_NAME, worker.AGE, worker.SALARY, worker.ADRESS)
        }
    }
}

func CalculateAverageSalary(workerList *WorkerList) {
    totalSalary := 0
    for _, w := range workerList.Workers {
        totalSalary += w.SALARY
    }
    if len(workerList.Workers) > 0 {
        averageSalary := float64(totalSalary) / float64(len(workerList.Workers))
        fmt.Printf("Средняя зарплата работников: %.2f\n", averageSalary)
    } else {
        fmt.Println("Список работников пуст.")
    }
}
