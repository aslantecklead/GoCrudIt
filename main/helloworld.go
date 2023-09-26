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

	workerList := loadWorkersFromDatabase(db)

		for {
			fmt.Println("Действие:")
			fmt.Println("1 - Добавить работника")
			fmt.Println("2 - Вывести среднюю зарплату работников")
			fmt.Println("3 - Вывести список всех работников")
			fmt.Println("4 - Выход")
			fmt.Println("5 - Изменить сотрудника")
			fmt.Println("6 - Удалить сотрудника")

			var num string
			fmt.Scanln(&num)
			ChooseOption(num, &workerList, db)
		}
}

func loadWorkersFromDatabase(db *sql.DB) WorkerList {
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
	return WorkerList{Workers: Workers}
}

func printWorkers(workerList *WorkerList) {
	for _, p := range workerList.Workers {
		fmt.Println("id: ", p.ID_WORKER, "Name: ", p.WORKER_NAME, "Age: ", p.AGE, "Salary: ", p.SALARY, "Adress: ", p.ADRESS)
	}
}

func ChooseOption(num string, workerList *WorkerList, db *sql.DB) {
	switch num {
	case "1":
			var newWorker Worker
			fmt.Print("Имя работника: ")
			fmt.Scanln(&newWorker.WORKER_NAME)
			fmt.Print("Возраст работника: ")
			fmt.Scanln(&newWorker.AGE)
			fmt.Print("Зарплата работника: ")
			fmt.Scanln(&newWorker.SALARY)
			fmt.Print("Адрес работника: ")
			fmt.Scanln(&newWorker.ADRESS)
			AddWorker(workerList, newWorker)
	case "2":
			CalculateAverageSalary(workerList)
	case "3":
			printWorkers(workerList)
	case "4":
			fmt.Println("Выход")
			os.Exit(0)
	case "5":
			UpdateOption(workerList, db)
	case "6":
			DeleteOption(workerList, db)
	default:
			fmt.Println("Неправильный выбор, попробуйте снова.")
	}
}

func UpdateOption(workerList *WorkerList, db *sql.DB) {
	var workerID int
	fmt.Print("Введите ID работника для обновления: ")
	fmt.Scanln(&workerID)

	// Найдем сотрудника по ID в списке
	index := -1
	for i, worker := range workerList.Workers {
			if worker.ID_WORKER == workerID {
					index = i
					break
			}
	}

	if index == -1 {
			fmt.Println("Работник с указанным ID не найден.")
			return
	}

	var updatedWorker Worker
	fmt.Print("Имя работника: ")
	fmt.Scanln(&updatedWorker.WORKER_NAME)
	fmt.Print("Возраст работника: ")
	fmt.Scanln(&updatedWorker.AGE)
	fmt.Print("Зарплата работника: ")
	fmt.Scanln(&updatedWorker.SALARY)
	fmt.Print("Адрес работника: ")
	fmt.Scanln(&updatedWorker.ADRESS)

	// Обновляем сотрудника в списке
	workerList.Workers[index] = updatedWorker

	// Обновляем сотрудника в базе данных
	_, err := db.Exec("UPDATE Worker SET WORKER_NAME = $1, AGE = $2, SALARY = $3, ADRESS = $4 WHERE ID_WORKER = $5",
			updatedWorker.WORKER_NAME, updatedWorker.AGE, updatedWorker.SALARY, updatedWorker.ADRESS, workerID)
	if err != nil {
			fmt.Println("Ошибка при обновлении работника в базе данных:", err)
			// Восстанавливаем старые данные сотрудника в списке, так как операция в базе данных не удалась
			workerList.Workers[index] = workerList.Workers[index]
			return
	}

	fmt.Println("Работник успешно обновлен.")
}

func DeleteOption(workerList *WorkerList, db *sql.DB) {
	var workerID int
	fmt.Print("Введите ID работника для удаления: ")
	fmt.Scanln(&workerID)

	// Найдем сотрудника по ID в списке
	index := -1
	for i, worker := range workerList.Workers {
			if worker.ID_WORKER == workerID {
					index = i
					break
			}
	}

	if index == -1 {
			fmt.Println("Работник с указанным ID не найден.")
			return
	}

	// Удаляем сотрудника из списка
	deletedWorker := workerList.Workers[index]
	workerList.Workers = append(workerList.Workers[:index], workerList.Workers[index+1:]...)

	// Удаляем сотрудника из базы данных
	_, err := db.Exec("DELETE FROM Worker WHERE ID_WORKER = $1", workerID)
	if err != nil {
			fmt.Println("Ошибка при удалении работника из базе данных:", err)
			// Восстанавливаем сотрудника в списке, так как операция в базе данных не удалась
			workerList.Workers = append(workerList.Workers, deletedWorker)
			return
	}

	fmt.Println("Работник успешно удален.")
}


func CalculateAverageSalary(workerList *WorkerList) {
	totalSalary := 0
	numWorkers := len(workerList.Workers)

	if numWorkers == 0 {
		fmt.Println("Нет работников для вычисления средней зарплаты.")
		return
	}

	for _, worker := range workerList.Workers {
		totalSalary += worker.SALARY
	}

	averageSalary := totalSalary / numWorkers
	fmt.Printf("Средняя зарплата работников: %d\n", averageSalary)
}

func AddWorker(workerList *WorkerList, worker Worker) {
	workerList.Workers = append(workerList.Workers, worker)

	connStr := "host=localhost port=5433 dbname=postgres user=postgres password=123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO Worker (WORKER_NAME, AGE, SALARY, ADRESS) VALUES ($1, $2, $3, $4)",
		worker.WORKER_NAME, worker.AGE, worker.SALARY, worker.ADRESS)
	if err != nil {
		fmt.Println("Ошибка при добавлении работника в базу данных:", err)
		return
	}

	fmt.Println("Работник успешно добавлен.")
}
