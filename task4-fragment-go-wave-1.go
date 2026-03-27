package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// ВНИМАНИЕ: в этом фрагменте есть несколько ошибок и плохих практик.
// Кандидату нужно:
// 1) Найти и описать проблемы.
// 2) Предложить, как переписать код лучше.

var db *sql.DB

func initDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost user=app dbname=devices sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Device простая модель устройства
type Device struct {
	ID       int64
	Hostname string
	IP       string
}

type DeviceRepo struct {
	db *sql.DB
}

func NewDeviceRepo(db *sql.DB) *DeviceRepo {
	return &DeviceRepo{db: db}
}

func (r *DeviceRepo) GetById(ctx context.Context, id int) (*Device, error) {
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("long debug operation finished")
	}()

	query := "SELECT id, hostname, ip FROM devices WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var d Device
	err := row.Scan(&d.ID, &d.Hostname, &d.IP)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Not found row")
		}
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO audit_log(device_id, ts, action) VALUES ($1, now(), 'view')", d.ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

type IDeviceRepo interface {
	GetById(ctx context.Context, id int) (*Device, error)
}

type DeviceHandler struct {
	repo IDeviceRepo
}

func NewDeviceHandler(repo IDeviceRepo) *DeviceHandler {
	return &DeviceHandler{
		repo: repo,
	}
}

// handler получает устройство по id и пишет в лог таблицу audit_log
func (h *DeviceHandler) DeviceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id should be number", http.StatusBadRequest)
		return
	}

	d, err := h.repo.GetById(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Device: %s (%s)", d.Hostname, d.IP)))
}

func main() {
	db := initDB()
	repo := NewDeviceRepo(db)
	handler := NewDeviceHandler(repo)

	http.HandleFunc("/device", handler.DeviceHandler)
	// Потенциальная проблема: сервер никогда не завершится, ошибки ListenAndServe игнорируются
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
