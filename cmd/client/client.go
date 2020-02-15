package main

import (
	"database/sql"
	"fmt"
	"github.com/jafarsirojov/APM-core/pkg/core"
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
			return false
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
			return true
		}
		listAtms(atms)

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
		userOnlineCards, err := core.GetUserCards(db)
		if err != nil {
			log.Printf("can't list of count client: %v", err)
			return true
		}
		cardsOnlineUser(userOnlineCards)

	case "2":
		err := transferMoney(db, transferToPhoneNumberOrCardNumber)
		if err != nil {
			log.Printf("can't transfer money: %v", err)
			return true
		}
	case "3":
		services, err := core.GetAllServices(db)
		if err != nil {
			log.Printf("can't get all services: %v", err)
			return true
		}
		listServices(services)
		err = handleTransferServices(db)
		if err != nil {
			fmt.Print("Упс неверный ход!!!")
			log.Printf("Can't pay to service: %v",err)
			return true
		}
	case "4":
		atms, err := core.GetAllAtms(db)
		if err != nil {
			log.Printf("can't get all atms: %v", err)
			return true
		}
		listAtms(atms)
	case "5":
		opLog, err := core.ViewOperationsLogging(db)
		if err != nil {
			log.Printf("can't get all atms: %v", err)
			return true
		}
		listUserOperationsLogging(opLog)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

func listAtms(atms []core.Atm) {
	for _, atm := range atms {
		fmt.Printf(
			"id: %d, name: %s, address: %s\n",
			atm.Id,
			atm.Name,
			atm.Address,
		)
	}
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

func cardsOnlineUser(cards []core.Card) {
	for _, card := range cards {
		fmt.Printf(
			"number count: %d, name: %s, balance: %d, number Count: %d\n",
			card.Id,
			card.Name,
			card.Balance,
			card.NumberCard,
		)
	}
}

func transferMoney(db *sql.DB, commands string) (err error) {
	fmt.Println(commands)
	var cmd string
	_, err = fmt.Scan(&cmd)
	if err != nil {
		log.Fatalf("Can't read input: %v", err)
	}
	switch cmd {
	case "1":
		transferMoneyForCountNumber(db)
	case "2":
		transferMoneyForPhoneNumber(db)
	case "q":
		return nil
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return nil
}

func transferMoneyForCountNumber(db *sql.DB) {
	fmt.Print("Введите номер счёта:")
	var countNumber string
	fmt.Scan(&countNumber)

	fmt.Print("Введите сумму перевода:")
	var currency int
	fmt.Scan(&currency)
	if currency < 0{
		fmt.Println("Упс)) Вы зломали банк!!! \n Хитрый ход!!!")
		return
	}

	err := core.TransferMoneyCardNumber(countNumber, db)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = core.TransferMoney(currency, db)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Перевод успешно выполнен!!!")
}

func transferMoneyForPhoneNumber(db *sql.DB) {
	fmt.Print("Введите номер телфона:")
	var phoneNumber int
	fmt.Scan(&phoneNumber)

	fmt.Print("Введите сумму перевода:")
	var currency int
	fmt.Scan(&currency)
	if currency < 0{
		fmt.Println("Упс)) Вы зломали банк!!! \n Хитрый ход!!!")
		return
	}

	err := core.TransferMoneyForPhoneNumber(phoneNumber, db)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = core.TransferMoney(currency, db)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Перевод успешно выполнен!!!")
}

func handleTransferServices(db *sql.DB) (err error) {
	fmt.Println("Оплачивайте услуг))")
	fmt.Print("Введите имя услуги: ")
	var nameService string
	fmt.Scan(&nameService)

	fmt.Print("Введите сумму оплачиваемого услуги: ")
	var currency int
	fmt.Scan(&currency)
	err = core.TransferServices(currency, nameService, db)
	if err != nil {
		return err
	}
	fmt.Print("Улуга успешно оплачено!!!")
	return nil
}

func listServices(services []core.Service) {
	fmt.Println("Список услуг: ")
	for _, service := range services {
		fmt.Printf(
			"%s\n",
			service.Name,
		)
	}
}
