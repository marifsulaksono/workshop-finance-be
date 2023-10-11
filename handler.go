package main

import (
	"fmt"
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
		err := rows.Scan(&obj.Id, &obj.UserId, &obj.Date, &obj.Status, &obj.Amount)
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
func GetTransactionById(id int) (transaction, error) {
	// declaration new object
	var obj transaction

	db := connect()
	defer db.Close()

	// statement query for get row match and scan
	err := db.QueryRow("select * from transaction where id = ?", id).
		Scan(&obj.Id, &obj.UserId, &obj.Date, &obj.Status, &obj.Amount)
	if err != nil {
		return transaction{}, err
	}

	return obj, nil
}

func InsertNewTransaction() {
	// declaration new object
	var obj transaction

	// form input
	fmt.Print("Masukkan User ID : ")
	fmt.Scan(&obj.UserId)
	fmt.Print("Masukkan Status (in/out) : ")
	fmt.Scan(&obj.Status)
	fmt.Print("Masukkan Jumlah : Rp. ")
	fmt.Scan(&obj.Amount)

	db := connect()
	defer db.Close()

	// initial date value
	obj.Date = time.Now()

	// call validate total expenses function
	if obj.Status == "out" {
		ValidateTotalExpenses(obj.UserId, obj.Amount, 100000)
	}

	// statement execution query for inserting data
	_, err := db.Exec("insert into transaction values (?,?,?,?,?)",
		nil, &obj.UserId, &obj.Date, &obj.Status, &obj.Amount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Insert Success")
}

func UpdateTransaction() {
	// declaration obj and date variable
	var obj transaction
	var date string

	// form ID input
	fmt.Print("Masukkan ID : ")
	fmt.Scan(&obj.Id)

	// check the ID provided on database
	result, err := GetTransactionById(obj.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)

	fmt.Print("Masukkan User ID : ")
	fmt.Scan(&obj.UserId)
	fmt.Print("Masukkan Tanggal (2023-10-01) : ")
	fmt.Scan(&date)
	fmt.Scanln()
	fmt.Print("Masukkan Status (in/out) : ")
	fmt.Scan(&obj.Status)
	fmt.Scanln()
	fmt.Print("Masukkan Jumlah : Rp. ")
	fmt.Scan(&obj.Amount)

	obj.Date, err = time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db := connect()
	defer db.Close()

	// call validate total expenses function
	if obj.Status == "out" {
		ValidateTotalExpenses(obj.UserId, obj.Amount, 100000)
	}

	// statement execution query for updating data
	_, err = db.Exec("update transaction set user_id=?, date=?, status=?, amount=? where id=?", &obj.UserId, &obj.Date, &obj.Status, &obj.Amount, obj.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
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
	_, err := GetTransactionById(id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db := connect()
	defer db.Close()

	// statement execution query for deleting data
	_, err = db.Exec("delete from transaction where id=?", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Delete Success")
}

func ValidateTotalExpenses(userId, current, limit int) {
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

	// suplement total plust current expenses
	total += current
	fmt.Println(total)

	// limit validation
	if total > limit {
		fmt.Println("Pengeluaran Anda melebihi limit")
		return
	}
}
