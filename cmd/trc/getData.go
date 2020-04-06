package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type request struct {
	url          string
	values       url.Values
	responseData []byte
	htmlNode     *html.Node
}

//getTeacherDataRequest
type getTDRequest struct {
	request
	teacherInfo []teacherJSON
}

//getUnitDataRequest
type getUDRequest struct {
	request
	unitTitle map[string]unitTitleJSON
}

//getContainUnitDataRequest
type getCUDRequest struct {
	request
	unitInfo []collegeJSON
}

//getSyllabusVideoLinkRequest
type getSVLRequest struct {
	request
	academicYear string //academicYear
	semester     string //semester
	svXi         map[string][]syllabusVideo
}

//getYoutubeVideoDurationRequest
type getYTVDRequest struct {
	request
	title         string
	duration      string
	seconds       int
	youtubeAPIKey string
	videoInfo     ytVideoInfo
}

func (req *request) newRequest() {
	req.values = url.Values{}
	req.responseData = make([]byte, 0)
}

//setURL 設定請求網址
func (req *request) setURL(url string) {
	req.url = url
}

//sendRequest 發送請求
func (req *request) sendPostRequest() (err error) {

	//送出請求並將返回結果放入res
	res, err := http.Post(req.url, "application/x-www-form-urlencoded", bytes.NewBufferString(req.values.Encode()))
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	req.responseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

//sendGetRequest 發送GET請求
func (req *request) sendGetRequest() (err error) {

	//送出請求並將返回結果放入res
	res, err := http.Get(req.url)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	req.responseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

//parseHTML 剖析HTML字串放入req.htmlNode中
func (req *request) parseHTML() (err error) {
	node, err := html.Parse(strings.NewReader(string(req.responseData)))
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	req.htmlNode = node
	return nil
}
