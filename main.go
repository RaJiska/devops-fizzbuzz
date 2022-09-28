package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var globalDataStore = map[string]string{}

func init() {
	rand.NewSource(time.Now().Unix())
}

func writeResult(w http.ResponseWriter, status int, result interface{}) {
	w.WriteHeader(status)
	responseBody := map[string]interface{}{
		"status": status,
		"result": result,
	}
	response, _ := json.Marshal(responseBody)
	w.Write(response)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	headers, _ := json.Marshal(r.Header)
	writeResult(w, http.StatusOK, map[string]interface{}{"headers": string(headers)})
}

func handleStubbedProcess1(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(100) > 20 {
		writeResult(w, http.StatusOK, "ok")
		return
	}
	writeResult(w, http.StatusInternalServerError, "not ok")
}

func handleStubbedProcess2(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(100) > 40 {
		writeResult(w, http.StatusOK, "ok")
		return
	}
	writeResult(w, http.StatusInternalServerError, "not ok")
}

func handleStubbedProcess3(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(100) > 60 {
		writeResult(w, http.StatusOK, "ok")
		return
	}
	writeResult(w, http.StatusInternalServerError, "not ok")
}

func handleStubbedProcess4(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(100) > 80 {
		writeResult(w, http.StatusOK, "ok")
		return
	}
	writeResult(w, http.StatusInternalServerError, "not ok")
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeResult(w, http.StatusMethodNotAllowed, "use post")
		return
	}
	bodyData, _ := ioutil.ReadAll(r.Body)
	requestId := uuid.New().String()
	globalDataStore[requestId] = string(bodyData)
	writeResult(w, http.StatusOK, requestId)
}

func handleLoad(w http.ResponseWriter, r *http.Request) {
	// get the part after /load/, which should be the uuid returned from /save-me
	dataId := strings.TrimPrefix(r.URL.Path, "/load/")
	if globalDataInstance, ok := globalDataStore[dataId]; ok {
		writeResult(w, http.StatusOK, globalDataInstance)
		return
	}
	writeResult(w, http.StatusNotFound, "not found")
}

func requestReceived(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"remoteAddr": r.RemoteAddr,
			"proto":      r.Proto,
			"method":     r.Method,
			"url":        r.URL,
			"requestId":  r.Header.Get("X-Request-Id"),
		}).Info()
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// endpoint routing configuration
	handler := http.NewServeMux()

	handler.HandleFunc("/", handleRoot)
	handler.HandleFunc("/stubbed-process-1", handleStubbedProcess1)
	handler.HandleFunc("/stubbed-process-2", handleStubbedProcess2)
	handler.HandleFunc("/stubbed-process-3", handleStubbedProcess3)
	handler.HandleFunc("/stubbed-process-4", handleStubbedProcess4)
	handler.HandleFunc("/save", handleSave)
	handler.HandleFunc("/load/", handleLoad)

	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "3000"))
	healthCheckTimeout, _ := strconv.Atoi(getEnv("HEALTHCHECK_TIMEOUT", "1"))

	// healthcheck function
	go func() {
		for {
			// default to pinging ourselves
			targetService := fmt.Sprintf("http://127.0.0.1:%d", serverPort)
			log.Info(fmt.Sprintf("pinging %s...", targetService))
			req, _ := http.NewRequest("GET", targetService, nil)
			res, _ := http.DefaultClient.Do(req)
			if res.StatusCode != http.StatusOK {
				log.Error(fmt.Sprintf("failed to complete healthcheck (status code: %v)", res.StatusCode))
			}
			responseBody, _ := ioutil.ReadAll(res.Body)
			log.Info(fmt.Sprintf("body data: %s", string(responseBody)))
			<-time.After(time.Duration(healthCheckTimeout) * time.Second)
		}
	}()

	// start the server
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", serverPort), requestReceived(handler)); err != nil {
		log.Fatal(fmt.Sprintf("failed to start server: %s\n", err))
	}
}
