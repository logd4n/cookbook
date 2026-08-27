package server

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

var (
	ClientIPErr = errors.New("Не удалось получить IP-адрес клиента...")
)

// Функция получения IP-адреса клиента
func getClientIP(ctx context.Context, r *http.Request) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctxErr
	default:
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ips := strings.Split(forwarded, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0]), nil
			}
		}

		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			return realIP, nil
		}

		return "unkwown IP", ClientIPErr
	}
}
