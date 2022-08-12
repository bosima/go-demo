package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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
	for _, uri := range confs.Target.Uris {
		http.HandleFunc(uri, proxyHandle)
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
	//fmt.Println(requestUri)

	// 拼接代理服务的key
	requestUri = requestUri + "&key=" + confs.Target.PrivateKey

	// 发起http调用
	rsp, err := http.Get(requestUri)
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = rsp.Body.Close() }()

	// HTTP响应头
	for k, v := range rsp.Header {
		w.Header().Set(k, v[0])
	}

	// HTTP状态码
	w.WriteHeader(rsp.StatusCode)

	// HTTP Body
	body, _ := ioutil.ReadAll(rsp.Body)
	w.Write(body)
}
