package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	appName string
)

func init() {
	appName = os.Getenv("APP_NAME")
	if len(appName) == 0 {
		fmt.Println("Must set APP_NAME.")
		os.Exit(1)
	}
}

func EnvUint64(name string, defaultVal uint64) uint64 {
	tmp := os.Getenv(name)
	if len(tmp) != 0 {
		n, err := strconv.ParseUint(tmp, 10, 64)
		if err == nil {
			return n
		}
	}
	return defaultVal
}

func EnvString(name string, defaultVal string) string {
	tmp := os.Getenv(name)
	if len(tmp) != 0 {
		return tmp
	}
	return defaultVal
}

func EnvInt(name string, defaultVal int) int {
	tmp := os.Getenv(name)
	if len(tmp) != 0 {
		n, err := strconv.Atoi(tmp)
		if err == nil {
			return n
		}
	}
	return defaultVal
}

func EnvDuration(name string, defaultVal int) time.Duration {
	return time.Duration(EnvInt(name, defaultVal))
}

func MeasureI(measurement string, val int64) {
	m := fmt.Sprintf("%s.%s", appName, measurement)
	fmt.Printf("measure=%q val=%d\n", m, val)
}

func MeasureT(measurement string, t time.Time) {
	m := fmt.Sprintf("%s.%s", appName, measurement)
	fmt.Printf("measure=%q val=%d\n", m, time.Since(t)/time.Millisecond)
}

func WriteJsonBytes(w http.ResponseWriter, status int, b []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(b)
	w.Write([]byte("\n"))
}

// Convenience
func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("at=error error=%s\n", err)
		WriteJson(w, 500, map[string]string{"error": "Internal Server Error"})
	}
	WriteJsonBytes(w, status, b)
}

func RoundTime(t time.Time, d time.Duration) time.Time {
	return time.Unix(0, int64((time.Duration(t.UnixNano())/d)*d))
}

func ParseRedisUrl() (string, string, error) {
	u, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		return "", "", errors.New("utils: Missing REDIS_URL")
	}
	var password string
	if u.User != nil {
		password, _ = u.User.Password()
	}
	return u.Host, password, nil
}

func ParseAuth(r *http.Request) (user, pass string, err error) {
	header, ok := r.Header["Authorization"]
	if !ok {
		err = errors.New("Authorization header not set.")
		return
	}

	auth := strings.SplitN(header[0], " ", 2)
	if len(auth) != 2 {
		err = errors.New("Malformed header.")
		return
	}

	userPass, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		err = errors.New("Malformed encoding.")
		return
	}

	parts := strings.Split(string(userPass), ":")

	//In this case, the request requires db-backed authentication.
	//the token.id is expected to be in parts[1] (password field).
	if parts[0] == "l2met" && len(parts[1]) > 0 {
		user = parts[0]
		pass = parts[1]
		return
	}

	//If the user is not == "l2met" then dbless-auth is requested.
	//ATM we assume the first part (user field) contains a base64 encoded
	//representation of the outlet credentials.
	if len(parts[0]) > 0 && len(parts[1]) == 0 {
		var decodedPart []byte
		decodedPart, err = base64.StdEncoding.DecodeString(parts[0])
		outletCreds := strings.Split(string(decodedPart), ":")
		//If the : is absent in parts[0], outletCreds[0] will contain the entire string in parts[0].
		user = outletCreds[0]
		//It is not required for the outletCreds to contain a pass.
		//The empty string that is returned will be handled by the outlet.
		if len(outletCreds) > 1 {
			pass = outletCreds[1]
		}
	}
	return
}
