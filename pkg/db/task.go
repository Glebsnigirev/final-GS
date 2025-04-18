package db

import "time"

// Структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Функция для добавления задачи
func AddTask(task *Task) (int64, error) {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (?, ?, ?, ?)
    `
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	args := []interface{}{}
	where := ""

	// Обработка параметра поиска
	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// Поиск по дате
			where = " WHERE date = ?"
			args = append(args, t.Format("20060102"))
		} else {
			// Поиск по подстроке в заголовке или комментарии
			where = " WHERE title LIKE ? OR comment LIKE ?"
			s := "%" + search + "%"
			args = append(args, s, s)
		}
	}

	query = query + where + " ORDER BY date LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}
