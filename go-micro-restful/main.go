package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	httpServer "github.com/go-micro/plugins/v4/server/http"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

type HttpResult struct {
	Content    string
	StatusCode int
}

type Hello struct {
	Language string
}

func (h *Hello) getHelloString() string {
	if h.Language == "cn" {
		return "你好"
	} else {
		return "Hello"
	}
}

func (h *Hello) Say(r *http.Request) (result *HttpResult) {
	result = &HttpResult{
		StatusCode: http.StatusOK,
	}

	if r.Method != http.MethodGet {
		result.StatusCode = http.StatusMethodNotAllowed
		return
	}

	paras := r.URL.Query()
	name, ok := paras["name"]
	if !ok {
		result.StatusCode = http.StatusBadRequest
		return
	}

	result.Content = fmt.Sprintf(`{"Message":"%s %s"}`, h.getHelloString(), name[0])
	return
}

func (h *Hello) Set(r *http.Request) (result *HttpResult) {
	result = &HttpResult{
		StatusCode: http.StatusOK,
	}

	if r.Method != http.MethodPost {
		result.StatusCode = http.StatusMethodNotAllowed
		return
	}

	if r.Body == http.NoBody {
		result.StatusCode = http.StatusBadRequest
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result.StatusCode = http.StatusBadRequest
		result.Content = err.Error()
		return
	}

	rMap := make(map[string]interface{})
	err = json.Unmarshal(b, &rMap)
	if err != nil {
		result.StatusCode = http.StatusBadRequest
		return
	}

	if lang, ok := rMap["Language"]; ok {
		h.Language = lang.(string)
		result.Content = `{"errcode":0}`
		return
	} else {
		result.StatusCode = http.StatusBadRequest
		return
	}
}

func main() {
	srv := httpServer.NewServer(
		server.Name("helloworld"),
		server.Address("0.0.0.0:8002"),
	)

	hello := &Hello{}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello/say", func(w http.ResponseWriter, r *http.Request) {
		result := hello.Say(r)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(result.StatusCode)
		w.Write([]byte(result.Content))
	})
	mux.HandleFunc("/hello/set", func(w http.ResponseWriter, r *http.Request) {
		result := hello.Set(r)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(result.StatusCode)
		w.Write([]byte(result.Content))
	})

	hd := srv.NewHandler(mux)

	srv.Handle(hd)

	service := micro.NewService(
		micro.Server(srv),
	)
	service.Init()
	service.Run()
}
