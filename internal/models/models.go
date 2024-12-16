package models

// TodoItem представляет один элемент Todo-списка в формате хранения в базе данных.
type TodoItem struct {
	Id      int    `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// TodoItemResponse представляет один элемент Todo-списка в формате, который ожидает frontend.
type TodoItemResponse struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// User представляет пользователя.
type User struct {
	Password string `json:"password"`
}

// AuthResponse представляет ответ при успешной аутентификации.
type AuthResponse struct {
	Token string `json:"token"`
}

// Format - формат времени, в котором передается и хранится дата.
const Format = "20060102"
