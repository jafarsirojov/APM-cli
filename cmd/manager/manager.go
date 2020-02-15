package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/jafarsirojov/APM-core/pkg/core"
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
			return false
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
		users, err := core.GetAllUsers(db)
		if err != nil {
			log.Printf("can't show list users to add card: %v", err)
			return true
		}
		listUsers(users)
		err = handleCard(db)
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
	case "5":
		operationsLoop(db, exportOperations, exportOperationsLoop)
	case "6":
		operationsLoop(db, importOperations, importOperationsLoop)
	case "7":
		usersShow, err := core.GetShowUsers(db)
		if err != nil {
			log.Printf("can't get show users: %v", err)
			return true
		}
		listShowUsers(usersShow)
		err = handleHide(db)
		if err != nil {
			log.Printf("can't hadle hide user: %v", err)
			return true
		}
	case "8":
		usersHide, err := core.GetHideUsers(db)
		if err != nil {
			log.Printf("can't get show users: %v", err)
			return true
		}
		listHideUsers(usersHide)
		err = handleShow(db)
		if err != nil {
			log.Printf("can't hadle show user: %v", err)
			return true
		}
	case "9":
		users, err := core.GetAllUsers(db)
		if err != nil {
			log.Printf("can't list users: %v", err)
			return true
		}
		listUsers(users)
	case "10":
		fmt.Println("Поиск пользователя")
		var phoneNumber int
		fmt.Print("Введите номер пользователя: ")
		fmt.Scan(&phoneNumber)
		users, err :=core.SearchUserByPhoneNumber(phoneNumber,db)
		if err != nil {
			return true
		}
		listUsers(users)
	case "11":
		var idUser int
		fmt.Print("Введите id пользователя, для просмотр журнал операции: ")
		fmt.Scan(&idUser)
		opLog, err := core.ViewOperationsLoggingToSearch(idUser,db)
		if err != nil {
			log.Printf("can't get operations logging: %v", err)
			return true
		}
		listUserOperationsLogging(opLog)
	case "12":
		fmt.Println("Просмотр событий журнала операций! ")
		opLog, err := core.ViewAllOperationsLogging(db)
		if err != nil {
			log.Printf("can't get operations logging: %v", err)
			return true
		}
		listUserOperationsLogging(opLog)

	case "13":
		fmt.Println(commandStatic)
		var cmdStatic string
		fmt.Scan(&cmdStatic)
		switch cmdStatic {
		case "a":
			fmt.Println("Количество пользователей в системе = ",core.StaticCountUsers(db))
		case "b":
			fmt.Println("Сколько денег у пользователей = ",core.StaticSumBalanceUsers(db))
		case "c":
			fmt.Println("На сколько оплачивано тех илииних услуг = ",core.StaticBalanceOfServices(db))
		case "d":
			fmt.Println("Сколько денег переведоно между счетами клиентов = ",core.StaticBalanceSumTransfer(db))
		}

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

func handleUser(db *sql.DB) (err error) {
	fmt.Println("Введите данные клиента:")
	var name string
	fmt.Print("Имя клиента: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
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
		return err
	}

	var phoneNumber int
	fmt.Print("Номер телефон: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}

	err = core.AddUser(name, login, password, passportSeries, phoneNumber, db)
	if err != nil {
		return err
	}

	fmt.Println("Клиент успешно добавлен!\n")

	return nil
}

func handleCard(db *sql.DB) (err error) {
	fmt.Println("Введите данные счёта:")
	var user_id int64
	fmt.Print("Введите id владелец счёта: ")
	_, err = fmt.Scan(&user_id)
	if err != nil {
		return err
	}
	var name string
	fmt.Print("Название счёта: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	var balance int64
	fmt.Print("Пополните счёт клиента: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	err = core.AddCard(name, balance, user_id, db)
	if err != nil {
		return err
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
		return err
	}

	err = core.AddService(name, db)
	if err != nil {
		return err
	}

	fmt.Println("Услуга успешно добавлено!")

	return nil
}

func handleAtm(db *sql.DB) (err error) {
	fmt.Println("Введите данные банкомата:")

	var name string
	fmt.Print("Название банкомата: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	var address string
	fmt.Print("Адресс банкомата: ")
	reader := bufio.NewReader(os.Stdin)
	address, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Can't read command: %v", err)
	}

	err = core.AddAtm(name, address, db)
	if err != nil {
		return err
	}

	fmt.Println("Банкомат успешно добавлен!")

	return nil
}

func listUsers(users []core.User) {
	fmt.Println("Список клиентов: ")
	for _, user := range users {
		fmt.Printf(
			"id: %d, name: %s, passportSeries: %s, numberPhone: %d\n",
			user.Id,
			user.Name,
			user.PassportSeries,
			user.NumberPhone,
		)
	}
}

func handleHide(db *sql.DB) (err error) {
	fmt.Println("Введите id клиента который хотите заблокировать:")
	var id int
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}

	err = core.UserHideManager(id, db)
	if err != nil {
		return err
	}

	fmt.Println("Клиент успешно заблокировань!!!\n")

	return nil
}
func handleShow(db *sql.DB) (err error) {
	fmt.Println("Введите id клиента который хотите разблокировать:")
	var id int
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}

	err = core.UserShowManager(id, db)
	if err != nil {
		return err
	}

	fmt.Println("Клиент успешно разблокировань!!!\n")

	return nil
}

func listHideUsers(users []core.UserHide) {
	fmt.Println("Список клиентов: ")
	for _, user := range users {
		fmt.Printf(
			"id: %d, name: %s, passportSeries: %s, numberPhone: %d\n",
			user.Id,
			user.Name,
			user.PassportSeries,
			user.NumberPhone,
		)
	}
}
func listShowUsers(users []core.UserShow) {
	fmt.Println("Список клиентов: ")
	for _, user := range users {
		fmt.Printf(
			"id: %d, name: %s, passportSeries: %s, numberPhone: %d\n",
			user.Id,
			user.Name,
			user.PassportSeries,
			user.NumberPhone,
		)
	}
}


func exportOperationsLoop(db *sql.DB, cmd string) bool {
	switch cmd {
	case "1":
		err := core.ExportAtmsToJSON(db)
		fmt.Println("Список банкоматов успешно экспортирован в JSON")
		if err != nil {
			log.Println(err)
		}
	case "2":
		err := core.ExportClientsToJSON(db)
		fmt.Println("Список клиентов успешно экспортирован в JSON")
		if err != nil {
			log.Println(err)
		}
	case "3":
		err := core.ExportAtmsToXML(db)
		fmt.Println("Список банкоматов успешно экспортирован в XML")
		if err != nil {
			log.Println(err)
		}
	case "4":
		err := core.ExportClientsToXML(db)
		fmt.Println("Список клиентов успешно экспортирован в XML")
		if err != nil {
			log.Println(err)
		}
	case "q":
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func importOperationsLoop(db *sql.DB, cmd string) bool {
	switch cmd {
	case "1":
		err := core.ImportAtmsFromJSON(db)
		fmt.Println("Список банкоматов успешно импортирован в JSON")
		if err != nil {
			log.Print(err)
		}
	case "2":
		err := core.ImportClientsFromJSON(db)
		fmt.Println("Список клиентов успешно импортирован в JSON")
		if err != nil {
			log.Println(err)
		}
	case "3":
		err := core.ImportAtmsFromXML(db)
		fmt.Println("Список банкоматов успешно импортирован в XML")
		if err != nil {
			log.Println(err)
		}
	case "4":
		err := core.ImportClientsFromXML(db)
		fmt.Println("Список клиентов успешно импортирован в XML")
		if err != nil {
			log.Println(err)
		}
	case "q":
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func listUserOperationsLogging(opLogs []core.OperationsLogging) {
	for _, opLog := range opLogs {
		fmt.Printf(
			"операция: %s, получатель-отправитель: %s, использовано денег в операции: %d, время операции: %s\n",
			opLog.Name,
			opLog.RecipientSender,
			opLog.Balance,
			opLog.Time,
		)
	}
}