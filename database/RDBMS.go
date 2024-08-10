package database

import (
	"encoding/json"
	"os"
)

// Note структура записи в бд
type Note struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Priority string `json:"priority"`
	Time     string `json:"time"`
	Date     string `json:"date"`
	Done     bool   `json:"done"`
}

// CRUD интерфейс для выполнения операций CRUD
type CRUD interface {
	Create(note Note) error
	Read() ([]Note, error)
	Update(id int, newNote Note) error
	Delete(id int) error
}

// JSONDatabase это одна из реализаций CRUD через json файла как бд
type JSONDatabase struct {
	FileName string // Имя файла базы данных
	nextID   int    // Следующий ID для новой записи
}

// NewJSONDatabase инициализирует экземпляр принимая путь к файлу
func NewJSONDatabase(fileName string) (*JSONDatabase, error) {
	db := &JSONDatabase{
		FileName: fileName,
	}

	// Чтение существующих данных, чтобы определить следующий ID
	notes, err := db.Read()
	if err != nil && !os.IsNotExist(err) { // Проверка, существует ли файл
		return nil, err
	}

	// Определяем следующий ID на основе существующих данных
	for _, note := range notes {
		if note.ID >= db.nextID {
			db.nextID = note.ID + 1
		}
	}

	return db, nil
}

// Create создает новую запись в файле JSONDatabase
func (db *JSONDatabase) Create(note Note) error {
	notes, err := db.Read()
	if err != nil && !os.IsNotExist(err) { // Если файла нет, это не ошибка
		return err
	}

	// Присвоение нового ID и инкрементация nextID
	note.ID = db.nextID
	db.nextID++

	// Добавляем новую запись к существующим
	notes = append(notes, note)

	// Записываем обновленный слайс в файл
	file, err := os.Create(db.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(notes); err != nil {
		return err
	}

	return nil
}

// Read читает все записи из файла JSONDatabase
func (db *JSONDatabase) Read() ([]Note, error) {
	var notes []Note

	file, err := os.Open(db.FileName)
	if err != nil {
		if os.IsNotExist(err) { // Если файла нет, возвращаем пустой список и nil
			return notes, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&notes); err != nil {
		return nil, err
	}

	return notes, nil
}

// Update обновляет запись в файле JSONDatabase по ID
func (db *JSONDatabase) Update(id int, newNote Note) error {
	notes, err := db.Read()
	if err != nil {
		return err
	}

	for i, v := range notes {
		if v.ID == id {
			newNote.ID = id
			notes[i] = newNote
			break
		}
	}

	file, err := os.Create(db.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(notes); err != nil {
		return err
	}

	return nil
}

// Delete удаляет запись из файла JSONDatabase по ID
func (db *JSONDatabase) Delete(id int) error {
	notes, err := db.Read()
	if err != nil {
		return err
	}

	var updatedNotes []Note
	for _, v := range notes {
		if v.ID != id {
			updatedNotes = append(updatedNotes, v)
		}
	}

	file, err := os.Create(db.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(updatedNotes); err != nil {
		return err
	}

	return nil
}
