package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	name := os.Getenv("APP_NAME")
	port := os.Getenv("APP_PORT")
	bkeHost := os.Getenv("BKE_HOST")
	bkePort := os.Getenv("BKE_PORT")
	bkePath := os.Getenv("BKE_PATH")

	// ハンドラー関数を定義する
	proxy := func(w http.ResponseWriter, _ *http.Request) {

		fmt.Println(fmt.Sprintf("app_name: %s", name))

		cli, err := newHttpClient(fmt.Sprintf("http://%s:%s/%s", bkeHost, bkePort, bkePath), "GET")
		if err != nil {
			fmt.Printf("http client error: %v\n", err)
			w.WriteHeader(500)
			return
		}
		defer cli.Close()

		cli.SendRequest()

		apiRes, err := cli.RespToString()
		if err != nil {
			fmt.Printf("http client error: %v\n", err)
			w.WriteHeader(500)
			return
		}

		io.WriteString(w, fmt.Sprintf("route: %s -> %s", name, apiRes))
	}

	backend := func(w http.ResponseWriter, _ *http.Request) {
		fmt.Println(fmt.Sprintf("app_name: %s", name))
		io.WriteString(w, name)

	}

	// パスとハンドラー関数を結びつける
	http.HandleFunc("/", proxy)
	http.HandleFunc("/backend", backend)

	// Web サーバーを起動する
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

type HttpClient struct {
	req    *http.Request
	client *http.Client
	method string
	url    string
	resp   *http.Response
}

func newHttpClient(url string, method string) (*HttpClient, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return &HttpClient{}, err
	}
	client := new(http.Client)
	return &HttpClient{
		req:    req,
		client: client,
		method: method,
		url:    url,
	}, nil
}

func (h *HttpClient) setHeader(header map[string]string) {
	for k, v := range header {
		h.req.Header.Set(k, v)
	}
}

func (h *HttpClient) SendRequest() error {
	r, err := h.client.Do(h.req)
	if err != nil {
		return err
	}
	h.resp = r
	return nil
}

func (h *HttpClient) Close() {
	h.resp.Body.Close()
}

func (h *HttpClient) RespToString() (string, error) {
	b, err := ioutil.ReadAll(h.resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
