package core

const managerDDL = `
CREATE TABLE IF NOT EXISTS manager
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL
);`

const usersDDL = `
CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
	passportSeries TEXT NOT NULL UNIQUE,
	phoneNumber INTEGER NOT NULL
	
);`

const atmDDL = `
CREATE TABLE IF NOT EXISTS atm
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL UNIQUE,
    address TEXT NOT NULL
);`

const cardsDDL = `
CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL,
   user_id INTEGER REFERENCES users(id)
);`

const servicesDDL = `
CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL
);`

const managerInitialData = `INSERT INTO manager(name, login, password)
VALUES ('IBank', 'admin', 'boss')
       ON CONFLICT DO NOTHING;`

const loginManagerSQL = `SELECT login, password FROM manager WHERE login = ?`
const loginUsersSQL = `SELECT login, password FROM users WHERE login = ?`

const getAllAtmsSQL = `SELECT id, name, address FROM atm;`
const getAllServicesSQL = `SELECT id, name, balance FROM services;`
const getAllCardsSQL = `SELECT id, name, balance, user_id FROM cards;`

const addBalanceToUser = `SELECT u.name,(SELECT c.name from cards c WHERE u.id=c.user_id);`



const getAllUsersSQL = `SELECT name , login, password, passportSeries, phoneNumber,  FROM users;`

const insertAtmSQL = `INSERT INTO atm(name, address) VALUES ( :name, :address);`
const insertServiceSQL = `INSERT INTO services( name , balance) VALUES( :name, :balance);`
const insertCardSQL = `INSERT INTO cards( name , balance, user_id) VALUES ( :name, :balance, :user_id);`
const insertUserSQL = `INSERT INTO users( name , login, password, passportSeries, phoneNumber) VALUES (:name , :login, :password, :passportSeries, :phoneNumber);`

const selectUserCardsSQL = `SELECT id, name , balance  FROM cards WHERE users.id=user_id;`

