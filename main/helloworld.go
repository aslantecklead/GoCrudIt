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

func main() {
    connStr := "host=localhost port=5433 dbname=postgres user=postgres password=123 sslmode=disable"
    db, err := sql.Open("postgres", connStr)
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
        fmt.Println("Выход")
        os.Exit(0)
    default:
        fmt.Println("Неправильный выбор, попробуйте снова.")
    }
}

func AddWorker(workerList *WorkerList) {
    
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
