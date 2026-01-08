package server

import (
	"net"
	"net/http"
)

// Функция получения IP-адреса клиента
func getClientIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = "unknown client"
	}
	return ip, err
}
