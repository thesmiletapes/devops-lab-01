package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// AppInfoResponse — структура для информативного ответа в формате JSON
type AppInfoResponse struct {
	Service     string    `json:"service"`
	Environment string    `json:"environment"`
	Status      string    `json:"status"`
	HostName    string    `json:"hostname"`
	ClientIP    string    `json:"client_ip"`
	ServerTime  time.Time `json:"server_time"`
}

// HealthResponse — структура для пробы здоровья приложения (Healthcheck)
type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

var startTime = time.Now()

func main() {
	// Считываем порт из переменной окружения PORT (по умолчанию 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Считываем среду окружения
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "virtualbox-debian-dev"
	}

	// Имя хоста (в Docker это ID контейнера)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	// Основной роут: http://localhost:8080/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] Запрос: %s %s с адреса %s", r.Method, r.URL.Path, r.RemoteAddr)

		resp := AppInfoResponse{
			Service:     "devops-lab-microservice",
			Environment: appEnv,
			Status:      "RUNNING",
			HostName:    hostname,
			ClientIP:    r.RemoteAddr,
			ServerTime:  time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	// Эндпоинт проверки здоровья: http://localhost:8080/healthz
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(startTime).Round(time.Second).String()
		resp := HealthResponse{
			Status: "healthy",
			Uptime: uptime,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("[START] Веб-сервер запущен на порту :%s (Окружение: %s)", port, appEnv)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("[FATAL] Ошибка запуска сервера: %s", err)
	}
}

