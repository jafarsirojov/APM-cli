package core

import (
	"database/sql"
	"errors"
	"fmt"
	_"github.com/mattn/go-sqlite3"
)

var ErrInvalidPass = errors.New("invalid password")

type QueryError struct {
	Query string
	Err   error
}

type DbError struct {
	Err error
}

type DbTxError struct {
	Err         error
	RollbackErr error
}

type Atm struct {
	Id int64
	Name string
	Address string
}

type Service struct {
	Id int64
	Name string
	Balance int64
}

type Card struct {
	Id int64
	Name string
	Balance int64
	User_id int64
}

type User struct {
	Id int64
	Name string
	Login string
	Password string
	PassportSeries string
	NumberPhone int
}

func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", loginManagerSQL, receiver.Err.Error())
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}

func (receiver *DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

func (receiver *DbError) Unwrap() error {
	return receiver.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

//obnavlyon ------------------------
func Init(db *sql.DB) (err error) {
	ddls := []string{managerDDL, usersDDL, cardsDDL, atmDDL,servicesDDL}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return err
		}
	}

	initialData := []string{managerInitialData}
	for _, datum := range initialData {
		_, err = db.Exec(datum)
		if err != nil {
			return err
		}
	}

	return nil
}




//new-----------------------------------------------------------------------

func LoginManager(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginManagerSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginManagerSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}



func LoginUsers(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginUsersSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginUsersSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

func AddAtm( atmName string, atmAddress string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertAtmSQL,

		sql.Named("name", atmName),
		sql.Named("address", atmAddress),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllAtms(db *sql.DB) (atms []Atm, err error) {
	rows, err := db.Query(getAllAtmsSQL)
	if err != nil {
		return nil, queryError(getAllAtmsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		atm := Atm{}
		err = rows.Scan(&atm.Id, &atm.Name, &atm.Address)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, atm)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

func AddService( serviceName string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertServiceSQL,

		sql.Named("name", serviceName),
		sql.Named("balance", 0),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllServices(db *sql.DB) (services []Service, err error) {
	rows, err := db.Query(getAllServicesSQL)
	if err != nil {
		return nil, queryError(getAllServicesSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			services, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		service := Service{}
		err = rows.Scan(&service.Id, &service.Name, &service.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		services = append(services, service)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return services, nil
}

func AddCard( cardName string, cardBalance int64, cardUser_id int64, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertCardSQL,

		sql.Named("name", cardName),
		sql.Named("balance", cardBalance),
		sql.Named("user_id", cardUser_id),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllCards(db *sql.DB) (cards []Card, err error) {
	rows, err := db.Query(getAllCardsSQL)
	if err != nil {
		return nil, queryError(getAllCardsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			cards, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		card := Card{}
		err = rows.Scan(&card.Id, &card.Name, &card.Balance,&card.User_id)
		if err != nil {
			return nil, dbError(err)
		}
		cards = append(cards, card)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return cards, nil
}

func AddUser( userName string, userLogin string, userPassword string, userPassportSeries string, userPhoneNumber int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertUserSQL,

		sql.Named("name", userName),
		sql.Named("login", userLogin),
		sql.Named("password", userPassword),
		sql.Named("passportSeries", userPassportSeries),
		sql.Named("phoneNumber", userPhoneNumber),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsers(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(getAllUsersSQL)
	if err != nil {
		return nil, queryError(getAllUsersSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name,&user.Login, &user.Password,&user.PassportSeries,&user.NumberPhone)
		if err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return users, nil
}


func SelectUserCards(db *sql.DB) (cards []Card, err error) {
	rows, err := db.Query(selectUserCardsSQL)
	if err != nil {
		return nil, queryError(selectUserCardsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			cards, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		card := Card{}
		err = rows.Scan(&card.Id, &card.Name, &card.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		cards = append(cards, card)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return cards, nil
}