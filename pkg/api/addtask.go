package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alonsoF100/13_14-sprint/pkg/db"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := processTaskDate(&task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func processTaskDate(task *db.Task) error {
	now := time.Now()
	nowDate := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = nowDate
	}

	taskTime, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if AfterNow(now, taskTime) {
			task.Date = nextDate
		}
	} else {
		if AfterNow(now, taskTime) {
			task.Date = nowDate
		}
	}

	return nil
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSONResponse(w http.ResponseWriter, data map[string]string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}