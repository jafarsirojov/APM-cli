package main

import (
	"database/sql"
	"fmt"
	"github.com/jafarsirojov/APM-cli/cmd/core"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.OpenFile("logClient.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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

	fmt.Fprintln(os.Stdout, `Добро пожаловать в наше приложение! ,`)
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
	case "2":
		atms, err := core.GetAllAtms(db)
		if err != nil {
			log.Printf("can't get all atms: %v", err)
			return true // TODO: may be log fatal
		}
		handleListAtm(atms)

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

	ok, err = core.LoginUsers(login, password, db)
	if err != nil {
		return false, err
	}

	return ok, err
}


func authorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		err := handleCardsUser(db)
		if err != nil {
			log.Printf("can't list of cards client: %v", err)
			return true
		}
	//case "2":
	//	err := handleTransferMoney(db)
	//	if err != nil {
	//		log.Printf("can't transfer money: %v", err)
	//		return true
	//	}
	//case "3":
	//	err := handlePayForTheService(db)
	//	if err != nil {
	//		log.Printf("can't pay for the service: %v", err)
	//		return true
	//	}
	case "4":
		atms, err := core.GetAllAtms(db)
		if err != nil {
			log.Printf("can't get all atms: %v", err)
			return true // TODO: may be log fatal
		}
		handleListAtm(atms)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}


//-------------handle

func handleCardsUser(db *sql.DB) ( err error) {
	fmt.Println("Список ваших счётов:")
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

	fmt.Println("Клиент успешно добавлен!")

	return nil
}


func handleCard(db *sql.DB) ( err error) {
	fmt.Println("Введите данные счёта:")
	var name string
	fmt.Print("Имя карту: ")
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
	fmt.Print("Имя услуги: ")
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
	fmt.Print("Адрес банкомата: ")
	_, err = fmt.Scan(&address)
	if err != nil {
		return  err
	}

	err = core.AddAtm( name, address, db)
	if err != nil {
		return  err
	}

	fmt.Println("Банкомат успешно добавлен!")

	return nil
}
//new
func handleListAtm(atms []core.Atm) {
	for _, atm := range atms {
		fmt.Printf(
			"id: %d, name: %s, address: %s\n",
			atm.Id,
			atm.Name,
			atm.Address,
		)
	}
}
