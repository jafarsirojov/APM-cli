package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/jafarsirojov/APM-cli/cmd/core"
	"log"
	"os"
	"strings"
)



func main() {
	file, err := os.OpenFile("logManager.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer func() {
		log.Print("close db")
		if err := db.Close(); err != nil {
			log.Fatalf("can't close db: %v", err)
		}
	}()
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init db: %v", err)
	}

	fmt.Fprintln(os.Stdout, `Добро пожаловать! ,`)
	log.Print("start operations loop")
	operationsLoop(db, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")
}



func operationsLoop(db *sql.DB, commands string, loop func(db *sql.DB, cmd string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func unauthorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		ok, err := handleLogin(db)
		if err != nil {
			log.Printf("can't handle login: %v", err)
			return true
		}
		if !ok {
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")
			return false
		}
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}



func authorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		err := handleUser(db)
		if err != nil {
			log.Printf("can't add user: %v", err)
			return true
		}
	case "2":
		err := handleCard(db)
		if err != nil {
			log.Printf("can't add card: %v", err)
			return true
		}
	case "3":
		err := handleService(db)
		if err != nil {
			log.Printf("can't add service: %v", err)
			return true
		}
	case "4":
		err := handleAtm(db)
		if err != nil {
			log.Printf("can't add atm: %v", err)
			return true
		}
	//case "5":
	//	err := handleExport(db)
	//	if err != nil {
	//		log.Printf("can't export: %v", err)
	//		return true
	//	}
	//case "6":
	//	err := handleImport(db)
	//	if err != nil {
	//		log.Printf("can't import: %v", err)
	//		return true
	//	}
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}


func handleLogin(db *sql.DB) (ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return false, err
	}
	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return false, err
	}

	ok, err = core.LoginManager(login, password, db)
	if err != nil {
		return false, err
	}

	return ok, err
}


func handleUser(db *sql.DB) ( err error) {
	fmt.Println("Введите данные клиента:")
	var name string
	fmt.Print("Имя клиента: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return  err
	}

	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return err
	}

	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return err
	}

	var passportSeries string
	fmt.Print("Серия пасспорта: ")
	_, err = fmt.Scan(&passportSeries)
	if err != nil {
		return  err
	}

	var phoneNumber int
	fmt.Print("Номер телефон: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}

	err = core.AddUser( name, login, password, passportSeries, phoneNumber, db)
	if err != nil {
		return  	err
	}

	fmt.Println("Клиент успешно добавлен!\n")

	return nil
}


func handleCard(db *sql.DB) ( err error) {
	fmt.Println("Введите данные счёта:")
	var name string
	fmt.Print("Название счёта: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return  err
	}

	var balance int64
	fmt.Print("Пополните счёт клиента: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return  err
	}

	var user_id int64
	fmt.Print("Введите id владелец счёта: ")
	_, err = fmt.Scan(&user_id)
	if err != nil {
		return  err
	}

	err = core.AddCard( name, balance, user_id, db)
	if err != nil {
		return  err
	}

	fmt.Println("Счёт успешно добавлен!")

	return nil
}


func handleService(db *sql.DB) (err error) {
	fmt.Println("Введите данные услуги:")
	var name string
	fmt.Print("Название услуги: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return  err
	}

	err = core.AddService( name, db)
	if err != nil {
		return  err
	}

	fmt.Println("Услуга успешно добавлено!")

	return nil
}

func handleAtm(db *sql.DB) ( err error) {
	fmt.Println("Введите данные банкомата:")

	var name string
	fmt.Print("Название банкомата: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return  err
	}

	var address string
	fmt.Print("Адресс банкомата: ")
	reader := bufio.NewReader(os.Stdin)
	address, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Can't read command: %v", err)
	}


	err = core.AddAtm( name, address, db)
	if err != nil {
		return  err
	}

	fmt.Println("Банкомат успешно добавлен!")

	return nil
}
