package storage

// Структура данных для хранения новостей
type News struct {
	Id      int    //Идентификатор новости
	Title   string // Заголовок новости
	Content string // Текст новости
	PubTime int64  // Время публикации новости в источнике
	Link    string // Сыылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	AddNews(News) error           // Добавлеет новость
	LastNews(int) ([]News, error) // Возвращает заданное количество последних новостей
}
