package env

import "os"

func Env(name string, def string) string {
	if e, ok := os.LookupEnv(name); ok {
		return e
	}

	return def
}
