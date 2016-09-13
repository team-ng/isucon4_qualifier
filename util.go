package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"github.com/gin-gonic/gin"
)

func getEnv(key string, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}

	return v
}

func getFlash(c *gin.Context, key string) string {
	//session, _ := store.Get(c.Request, key)
	//if value := session.Values[key]; value== nil {
	//	return ""
	//} else {
	//	session.Values[key] = ""
	//	session.Save(c.Request, c.Writer)
	//	return value.(string)
	//}
	return ""
}

func calcPassHash(password, hash string) string {
	h := sha256.New()
	io.WriteString(h, password)
	io.WriteString(h, ":")
	io.WriteString(h, hash)

	return fmt.Sprintf("%x", h.Sum(nil))
}
