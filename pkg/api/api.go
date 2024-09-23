package api

import (
	"encoding/json"
	"goNews/pkg/storage/postgres"
	"net/http"
	"strconv"
)

// Структура API
type API struct {
	db     *postgres.Store // Хранилище данных
	router *http.ServeMux  //Маршрутизатор запросов
}

// Конструктор API
func NewAPI(db *postgres.Store) *API {
	api := API{}
	api.db = db
	api.router = http.NewServeMux()
	api.endpoints()
	return &api
}

// Возвращает маршрутизатор API
func (api *API) Router() *http.ServeMux {
	return api.router
}

// Регистрируем обработчики
func (api *API) endpoints() {
	// Основной метод, возвращающий заданное количество последних новостей
	api.router.HandleFunc("GET /news/{quantity}", api.lastNews)
	// Обработчик запроса к статическим фалам веб-приложения, предоставленого в задании
	api.router.Handle("GET /", http.FileServer(http.Dir("./webapp")))
}

// lastNews возвращает заданное количество последних новостей
func (api *API) lastNews(w http.ResponseWriter, r *http.Request) {
	// Заголовки указывают тип данных ответа и разрешают кросс-доменные запросы
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Читаем количество новостей к выдаче из параметров запроса
	quantity, err := strconv.Atoi(r.PathValue("quantity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Запрашиваем в БД заданное количество последних новостей
	lastNews, err := api.db.LastNews(quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Возвращаем список новостей в JSON формате
	json.NewEncoder(w).Encode(lastNews)
}
