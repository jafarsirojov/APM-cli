package main

const unauthorizedOperations = `Список доступных операций:
1. Авторизация
q. Выйти из приложения

Введите команду`

const authorizedOperations = `Список доступных операций:
1. Добавление пользователя
2. Добавить счёт пользователю
3. Добавить услуги
4. Добавление банкомата
5. Экспорт
6. Импорт
--------- 4 балла
7. Блокировать пользователя
8. Разблокировать пользователя
9. Вывод списка пользователей
10. Поиск пользователя
--------- 5 балла
11. Просмотр событий журнала операций по пользователям
12. Просмотр событий журнала операций всех пользователей
13. Статистика
q. Выйти (разлогиниться)

Введите команду`

const commandStatic =` Список доступных операций по статистике:
	a. Количество пользователей в системе
	b. Сколько денег у пользователей
	c. На сколько оплачивано тех илииних услуг
	d. Сколько денег переведоно между счетами клиентов
	Введите команду`

const exportOperations  = `Список доступных операций:
  1. Экспортировать список банокматов в JSON
  2. Экспортировать список клиентов в JSON
  3. Экспортировать список банкоматов в XML
  4. Экспортировать список клиентов в XML

  q. Назад

Введите команду`

const importOperations  = `Список доступных операций:
  1. Импортировать список банкоматов в JSON
  2. Импортировать список клиентов в JSON
  3. Импортировать список банкоматов в XML
  4. Импортировать список клиентов в XML

  q. Назад

Введите команду`