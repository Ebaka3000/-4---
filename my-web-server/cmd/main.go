package main

import (
    "encoding/json"
    "log"
    "my-web-server/config"
    "my-web-server/internal/cache"
    "my-web-server/internal/database"
    "my-web-server/internal/http"
    "my-web-server/internal/messaging"
    "my-web-server/internal/models"
    "github.com/nats-io/stan.go"
)

func main() {
    // Загрузка конфигурации
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Ошибка загрузки конфигурации: %v", err)
    }

    // Инициализация базы данных
    db, err := database.NewPostgresDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Ошибка инициализации базы данных: %v", err)
    }
    defer db.Close()

    // Инициализация NATS Streaming
    natsClient, err := messaging.NewNATSClient("test-cluster", "client-1")
    if err != nil {
        log.Fatalf("Ошибка инициализации NATS: %v", err)
    }
    defer natsClient.Close()

    // Инициализация кэша
    cache := cache.NewCache()

    // Восстановление кэша из базы данных
    rows, err := db.Query("SELECT id, data FROM my_table")
    if err != nil {
        log.Fatalf("Ошибка восстановления кэша из базы данных: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var id string
        var data string
        if err := rows.Scan(&id, &data); err != nil {
            log.Fatalf("Ошибка сканирования строки: %v", err)
        }
        cache.Set(id, data)
    }

    // Подписка на канал в NATS Streaming
    _, err = natsClient.Subscribe("my-channel", func(msg *stan.Msg) {
        var data models.Data
        if err := json.Unmarshal(msg.Data, &data); err != nil {
            log.Printf("Ошибка десериализации сообщения: %v", err)
            return
        }

        // Сохранение данных в базе данных
        _, err := db.Exec("INSERT INTO my_table (id, data) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data", data.ID, data.Data)
        if err != nil {
            log.Printf("Ошибка сохранения данных в базе данных: %v", err)
            return
        }

        // Сохранение данных в кэше
        cache.Set(data.ID, data.Data)
    })
    if err != nil {
        log.Fatalf("Ошибка подписки на канал: %v", err)
    }

    // Запуск HTTP сервера
    server := http.NewServer(cache, db)
    log.Println("HTTP сервер запущен на :8080")
    if err := server.Start(":8080"); err != nil {
        log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
    }
}