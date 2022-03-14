package crawler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type Request struct {
	url          string
	Values       url.Values
	ResponseData []byte
	HtmlNode     *html.Node
}

func (req *Request) NewRequest() {
	req.Values = url.Values{}
	req.ResponseData = make([]byte, 0)
}

//setURL 設定請求網址
func (req *Request) SetURL(url string) {
	req.url = url
}

//sendRequest 發送請求
func (req *Request) SendPostRequest() (err error) {

	//送出請求並將返回結果放入res
	res, err := http.Post(req.url, "application/x-www-form-urlencoded", bytes.NewBufferString(req.Values.Encode()))
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	req.ResponseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

//sendGetRequest 發送GET請求
func (req *Request) SendGetRequest() (err error) {

	//送出請求並將返回結果放入res
	res, err := http.Get(req.url)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	req.ResponseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

//parseHTML 剖析HTML字串放入req.htmlNode中
func (req *Request) ParseHTML() (err error) {
	node, err := html.Parse(strings.NewReader(string(req.ResponseData)))
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	req.HtmlNode = node
	return nil
}
