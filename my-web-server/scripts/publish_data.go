package main

import (
    "log"
    "github.com/nats-io/stan.go"
)

const (
    DATABASE_URL = "postgres://user:password@localhost:5432/mydb?sslmode=disable"
    NATS_URL = "nats://localhost:4222"
)

func main() {
    sc, err := stan.Connect("test-cluster", "publisher")
    if err != nil {
        log.Fatalf("Ошибка подключения к NATS: %v", err)
    }
    defer sc.Close()

    data := `{"id": "1", "data": "example data"}`
    err = sc.Publish("my-channel", []byte(data))
    if err != nil {
        log.Fatalf("Ошибка публикации данных: %v", err)
    }

    log.Println("Данные успешно опубликованы")
}