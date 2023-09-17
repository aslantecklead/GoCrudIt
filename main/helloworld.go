package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Worker struct{
	ID_WORKER int
	WORKER_NAME string
	AGE int
	SALARY int
	ADRESS string
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
	for rows.Next(){
		p := Worker{}
		err := rows.Scan(&p.ID_WORKER, &p.WORKER_NAME, &p.AGE, &p.SALARY, &p.ADRESS)
		if err != nil{
				fmt.Println(err)
				continue
		}
		Workers = append(Workers, p)
	}
	for _, p := range Workers {
		fmt.Println("id: ", p.ID_WORKER, "Name: ",p.WORKER_NAME, "Age: ",p.AGE, "Salary: ",p.SALARY, "Adress: ",p.ADRESS);
	}
}