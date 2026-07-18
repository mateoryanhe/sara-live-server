package migrate

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
)

// PostgresDSN converts a GoFrame pgsql link to a GORM postgres DSN.
// Example link: pgsql:user:pass@tcp(127.0.0.1:5432)/live_db?sslmode=disable
func PostgresDSN(link, extra string) string {
	link = strings.TrimPrefix(link, "pgsql:")
	at := strings.LastIndex(link, "@")
	if at <= 0 {
		panic("invalid pgsql link: missing @")
	}
	userPass := link[:at]
	rest := link[at+1:]

	colon := strings.Index(userPass, ":")
	if colon <= 0 {
		panic("invalid pgsql link: missing user password separator")
	}
	user := userPass[:colon]
	pass := userPass[colon+1:]

	rest = strings.TrimPrefix(rest, "tcp(")
	paren := strings.Index(rest, ")")
	if paren <= 0 || len(rest) < paren+2 || rest[paren+1] != '/' {
		panic("invalid pgsql link: missing tcp(host:port)/dbname")
	}
	hostPort := rest[:paren]
	after := rest[paren+2:]

	dbName := after
	query := extra
	if q := strings.Index(after, "?"); q >= 0 {
		dbName = after[:q]
		linkQuery := after[q+1:]
		if linkQuery != "" {
			if query != "" {
				query = linkQuery + "&" + query
			} else {
				query = linkQuery
			}
		}
	}
	if query == "" {
		query = "sslmode=disable"
	}

	host := hostPort
	port := "5432"
	if idx := strings.Index(hostPort, ":"); idx >= 0 {
		host = hostPort[:idx]
		port = hostPort[idx+1:]
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/" + dbName,
	}
	u.RawQuery = gstr.Replace(query, " ", "&")
	return u.String()
}
