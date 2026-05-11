package models

type Eat struct {
	ID           int      `json:"ID"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
}

type RecipeShort struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type LogLevel string

const (
	Error LogLevel = "error"
	Info  LogLevel = "info"
	Panic LogLevel = "panic"
	Warn  LogLevel = "warning"
	Fatal LogLevel = "fatal"
	Debug LogLevel = "debug"
)

type LogMessage struct {
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
}
