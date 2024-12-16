<h1>Todo app</h1>

Этот проект представляет собой простейшее веб-приложение для создания личного списка дел.

<h3>Выполненные задания повышенной трудности:</h3>
<ul>
<li>Определение пути к базе данных через переменную окружения.</li>
<li>Реализация обработки правил повторения по дням недели и по месяцам.</li>
<li>Реализация поиска в базе данных события по дате или подстроке</li>
<li>Аутентификация пользователя по паролю, определенному в .env файле, с помощью JWT</li>
<li>Создание Docker-образа с помощью Dockerfile</li>
</ul>

<h2>Требования для запуска.</h2>
<ul>
<li>Golang 1.23 и выше</li>
<li>Файл .env в корне проекта. Пример:</li>
<p></p>
<code>
# Переменные для работы сервера.<br />
SERVER_HOST=0.0.0.0<br />
SERVER_PORT=7540<br />
<br />
# Переменные для работы базы данных:<br />
# DB_NAME - название файла базы данных;<br />
# TODO_DBFILE - путь хранения фалйа базы данных;<br />
DB_NAME=scheduler<br />
TODO_DBFILE=<br />
<br />
# Пароль для доступа к api списка задач.<br />
TODO_PASSWORD=12345<br />
<br />
# Переменные для формирования токена:<br />
# JWT_SECRET - секрет для подписи токена;<br />
# JWT_TTL - время жизни токена;<br />
JWT_SECRET=scheduler<br />
JWT_TTL=8<br />
</code>
<li>Docker (для создания образа и запуска контейнера)</li>
</ul>

<h2>Инструкция по запуску.</h2>
<h4>Локально:</h4>
<ol>
<li>Клонируйте репозиторий:</li>
<p></p>
<code>git clone https://github.com/paudarco/todo.git</code>
<p></p>
<li>Запустите main.go:</li>
<p></p>
<code>go run cmd/main.go</code>
<p></p>
<li>Интерфейс для взаимодействия с API будет доступен на <code>localhost:7540</code></li>
</ol>

<h4>В Docker контейнере:</h4>
<ol>
<li>Клонируйте репозиторий:</li>
<p></p>
<code>git clone https://github.com/paudarco/todo.git</code>
<p></p>
<li>В bash (cmd/powershell) выполните команду для сборки образа (требуется запущенный Docker):</li>
<p></p>
<code>docker build -t todo .</code>
<p></p>
<li>Запустите контейнер:</li>
<p></p>
<code>docker run -d --name todo-container -p 7540:7540 --mount type=bind,source="$(pwd)/scheduler.db",target=/app/scheduler.db --env-file .env todo</code>
<p></p>
<li>Интерфейс для взаимодействия с API будет доступен на <code>localhost:7540</code></li>
</ol>

<h2>Тестирование</h2>
В файле <code>/tests/settings.go</code> следует использовать:
<p></p>
<code>
var Port = 7540<br />
var DBFile = "../scheduler.db"<br />
var FullNextDate = true<br />
var Search = true<br />
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOls4OSwxNDgsNzEsMjYsMTg3LDEsMTcsNDIsMjUyLDE5MywxMjksODksMjQ2LDIwNCwxMTYsMTgwLDI0NSwxNywxODUsMTUyLDYsMjE4LDg5LDE3OSwyMDIsMjQ1LDE2OSwxOTMsMTE1LDIwMiwyMDcsMTk3XX0.HBJ7UFemjWMn0LqwQcewKJXoWrPIYhkS6N3MmGN93WM`<br />
</code>
<p>Токен валиден при том же пароле и секрете, что объявлен в TODO_PASSWORD и JWT_SECRET в примере .env файла</p>
<p></p>
Для тестирования выполнить:
<p></p>
<code>go test ./tests</code>
<p></p>
