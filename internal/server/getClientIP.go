package server

import (
	"context"
	"errors"
	"net"
	"net/http"
)

var (
	ctxErr      = errors.New("Пользователь отключился! getClientIP остановлен!")
	ClientIPErr = errors.New("Не удалось получить IP-адрес клиента...")
)

// Функция получения IP-адреса клиента
func getClientIP(ctx context.Context, r *http.Request) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctxErr
	default:
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "unknown client"
		}
		return ip, err
	}
}
