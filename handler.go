package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Get all transaction data from database
func GetAllTransactions() {
	// declaration new array
	var result []transaction

	// declaration new sql.DB variable
	db := connect()
	defer db.Close()

	// statement query for getting all data
	rows, err := db.Query("select * from transaction")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	// process scanning data row by row on table
	for rows.Next() {
		var obj = transaction{}
		err := rows.Scan(&obj.Id, &obj.UserId, &obj.Date, &obj.Status, &obj.Amount, &obj.Detail)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// adding the matching data to be an object
		result = append(result, obj)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
	}

	// print result
	fmt.Println(result)
}

// Get transaction by id
func GetTransactionById(id int) {
	// declaration new object
	var obj transaction

	db := connect()
	defer db.Close()

	// statement query for get row match and scan
	err := db.QueryRow("select * from transaction where id = ?", id).
		Scan(&obj.Id, &obj.UserId, &obj.Date, &obj.Status, &obj.Amount, &obj.Detail)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(obj)
}

func InsertNewTransaction() {
	// declaration new object
	var obj transaction
	sc := bufio.NewScanner(os.Stdin)

	// form input
	fmt.Print("Masukkan User ID : ")
	fmt.Scan(&obj.UserId)
	fmt.Print("Masukkan Status (in/out) : ")
	fmt.Scan(&obj.Status)
	fmt.Print("Masukkan Jumlah : Rp. ")
	fmt.Scan(&obj.Amount)
	fmt.Print("Masukkan keterangan : ")
	sc.Scan()
	obj.Detail = sc.Text()

	db := connect()
	defer db.Close()

	// initial date value
	obj.Date = time.Now()

	// statement execution query for inserting data
	_, err := db.Exec("insert into transaction values (?,?,?,?,?,?)",
		nil, &obj.UserId, &obj.Date, &obj.Status, &obj.Amount, &obj.Detail)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// call validate total expenses function
	if obj.Status == "out" {
		ValidateTotalExpenses(obj.UserId, 100000)
	}

	fmt.Println("Insert Success")
}

func UpdateTransaction() {
	// declaration obj and date variable
	var obj transaction
	var date string
	var err error
	sc := bufio.NewScanner(os.Stdin)

	// form ID input
	fmt.Print("Masukkan ID : ")
	fmt.Scan(&obj.Id)

	// check the ID provided on database
	GetTransactionById(obj.Id)

	fmt.Print("Masukkan User ID : ")
	fmt.Scan(&obj.UserId)
	fmt.Print("Masukkan Tanggal (2023-10-01) : ")
	sc.Scan()
	date = sc.Text()
	fmt.Print("Masukkan Status (in/out) : ")
	sc.Scan()
	obj.Status = sc.Text()
	fmt.Print("Masukkan Jumlah : Rp. ")
	fmt.Scan(&obj.Amount)
	fmt.Print("Masukkan Keterangan : ")
	sc.Scan()
	obj.Detail = sc.Text()

	obj.Date, err = time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db := connect()
	defer db.Close()

	// statement execution query for updating data
	_, err = db.Exec("update transaction set user_id=?, date=?, status=?, amount=?, detail=? where id=?",
		&obj.UserId, &obj.Date, &obj.Status, &obj.Amount, &obj.Detail, obj.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// call validate total expenses function
	if obj.Status == "out" {
		ValidateTotalExpenses(obj.UserId, 100000)
	}

	fmt.Println("Update Success")
}

func DeleteTransaction() {
	// declaration id variable
	var id int

	// form ID input
	fmt.Print("Masukkan ID : ")
	fmt.Scan(&id)

	// check the ID provided on database
	GetTransactionById(id)

	db := connect()
	defer db.Close()

	// statement execution query for deleting data
	_, err := db.Exec("delete from transaction where id=?", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Delete Success")
}

func ValidateTotalExpenses(userId, limit int) {
	// declaration total and date variable
	var total int
	var date = time.Now().Format("2006-01-02")

	db := connect()
	defer db.Close()

	// statement query for get row match and scan
	// coalesce is return the first non-null value in a list
	err := db.QueryRow("select coalesce(sum(amount), 0) from transaction where date=? and status='out' and user_id=?", date, userId).Scan(&total)
	if err != nil {
		fmt.Printf("Error count : %v", err)
	}

	// limit validation
	if total > limit {
		fmt.Println("Pengeluaran Anda melebihi limit")
		return
	}
}
