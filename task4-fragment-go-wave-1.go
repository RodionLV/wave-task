package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

// ВНИМАНИЕ: в этом фрагменте есть несколько ошибок и плохих практик.
// Кандидату нужно:
// 1) Найти и описать проблемы.
// 2) Предложить, как переписать код лучше.

var db *sql.DB

func initDB() {
	// Потенциальная проблема: ошибка игнорируется
	db, _ = sql.Open("postgres", "host=localhost user=app dbname=devices sslmode=disable")
	// Нет проверки доступности соединения и таймаута
}

// Device простая модель устройства
type Device struct {
	ID       int64
	Hostname string
	IP       string
}

// handler получает устройство по id и пишет в лог таблицу audit_log
func deviceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Потенциальная проблема: контекст без таймаута, возможная утечка
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("long debug operation finished")
	}()

	// Потенциальная проблема: строка запроса строится конкатенацией
	query := "SELECT id, hostname, ip FROM devices WHERE id = " + idStr
	row := db.QueryRowContext(ctx, query)

	var d Device
	err := row.Scan(&d.ID, &d.Hostname, &d.IP)
	if err != nil {
		// Ошибка обрабатывается одинаково, не различаем NotFound и т.п.
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// Потенциальная проблема: игнорируется ошибка вставки в audit_log
	db.ExecContext(ctx, "INSERT INTO audit_log(device_id, ts, action) VALUES ($1, now(), 'view')", d.ID)

	// Потенциальная проблема: нет установки Content-Type
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Device: %s (%s)", d.Hostname, d.IP)))
}

func main() {
	initDB()
	http.HandleFunc("/device", deviceHandler)
	// Потенциальная проблема: сервер никогда не завершится, ошибки ListenAndServe игнорируются
	http.ListenAndServe(":8080", nil)
}

