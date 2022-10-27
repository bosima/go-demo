package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Title  string
		Port   int
		Target target
	}

	target struct {
		PrivateKey string
		Host       string
		Uris       []string
	}
)

var confs config

func main() {

	// 读取配置
	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err := toml.DecodeFile(f, &confs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 读取支持的Uri列表，遍历注册http处理函数
	if len(confs.Target.Uris) > 0 {
		for _, uri := range confs.Target.Uris {
			http.HandleFunc(uri, proxyHandle)
		}
	} else {
		http.HandleFunc("/", proxyHandle)
	}

	// 读取配置的端口号
	port := strconv.Itoa(confs.Port)

	// 启动http服务
	fmt.Println("服务启动，监听端口:" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func proxyHandle(w http.ResponseWriter, r *http.Request) {
	// 读取代理服务的基地址
	host := confs.Target.Host

	// 从请求中读取请求Uri和参数，拼接之
	requestUri := host + r.RequestURI

	// 拼接代理服务的key
	if strings.Trim(confs.Target.PrivateKey, " ") != "" {
		requestUri = requestUri + "&key=" + confs.Target.PrivateKey
	}

	// 读请求体
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 创建请求
	req, err := http.NewRequest(r.Method, requestUri, strings.NewReader(string(reqBody)))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 写请求头
	for k, v := range r.Header {
		w.Header().Set(k, v[0])
	}

	cli := &http.Client{}
	rsp, err := cli.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if rsp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rsp.Body.Close()

	// HTTP响应头
	if rsp.Header != nil {
		for k, v := range rsp.Header {
			w.Header().Set(k, v[0])
		}
	}

	// HTTP状态码
	w.WriteHeader(rsp.StatusCode)

	// HTTP Body
	body, _ := ioutil.ReadAll(rsp.Body)
	w.Write(body)
}
