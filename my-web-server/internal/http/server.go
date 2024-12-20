package http

import (
    "database/sql"
    "encoding/json"
    "log"
    "my-web-server/internal/cache"
    "net/http"
    "html/template"
)

type Server struct {
    cache *cache.Cache
    db    *sql.DB
}

func NewServer(cache *cache.Cache, db *sql.DB) *Server {
    return &Server{
        cache: cache,
        db:    db,
    }
}

func (s *Server) Start(addr string) error {
    http.HandleFunc("/data", s.handleGetData)
    http.HandleFunc("/", s.handleIndex)
    return http.ListenAndServe(addr, nil)
}

func (s *Server) handleGetData(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "missing id", http.StatusBadRequest)
        return
    }

    if data, ok := s.cache.Get(id); ok {
        json.NewEncoder(w).Encode(data)
        return
    }

    var data string
    err := s.db.QueryRow("SELECT data FROM my_table WHERE id = $1", id).Scan(&data)
    if err != nil {
        http.Error(w, "data not found", http.StatusNotFound)
        return
    }

    s.cache.Set(id, data)
    json.NewEncoder(w).Encode(data)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
    tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Data Viewer</title>
    </head>
    <body>
        <h1>Data Viewer</h1>
        <form action="/data" method="get">
            <label for="id">ID:</label>
            <input type="text" id="id" name="id">
            <input type="submit" value="Get Data">
        </form>
    </body>
    </html>
    `
    t, err := template.New("index").Parse(tmpl)
    if err != nil {
        log.Printf("Ошибка создания шаблона: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }
    t.Execute(w, nil)
}