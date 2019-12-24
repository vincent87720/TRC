package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Luxurioust/excelize"
	"golang.org/x/net/html"
)

type searchRequest struct {
	smye     string //academicYear
	smty     string //semester
	day      string
	lesson   string
	htmlNode *html.Node
}

type course struct {
	courseXi []courseDetail
}

type courseDetail struct {
	gradeNClass         string
	creditNChooseSelect string
	courseID            string
	courseName          string
	teacher             string
	courseInfo          string
	remark              string
	courseURL           string
}

//searchRequest的建構函式
func newCourseList() *searchRequest {
	return &searchRequest{
		smye:   academicYear,
		smty:   semester,
		day:    "'1','2','3','4','5','6','7'",
		lesson: "'1','2','3','4','N','5','6','7','8','9','A','B','C','D','E'",
	}
}

//設定學年度
func (sr *searchRequest) setAcademicYear(year string) {
	i, err := strconv.Atoi(year)
	if err != nil || i < 0 {
		fmt.Println("請輸入合法的年分!")
		panic("Download Fail")
	} else {
		sr.smye = year
	}
}

//設定學期
func (sr *searchRequest) setSemester(semester string) {
	i, err := strconv.Atoi(semester)
	if err != nil || i < 0 || i > 2 {
		fmt.Println("請輸入合法的學期!")
		panic("Download Fail")
	} else {
		sr.smty = semester
	}
}

//設定查詢星期
func (sr *searchRequest) setDay(day string) {
	sr.day = day
}

//設定查詢節次
func (sr *searchRequest) setLesson(lesson string) {
	sr.lesson = lesson
}

//發送查詢請求並儲存返回結果
func (sr *searchRequest) getCourseData() {
	v := url.Values{}
	//post form data
	v.Add("smye", sr.smye)
	v.Add("smty", sr.smty)
	v.Add("str_time", sr.day+"sec"+sr.lesson)

	//送出請求並將返回結果放入res
	res, err := http.Post("http://syl.dyu.edu.tw/sl_cour_time.php?itimestamp="+string(int32(time.Now().Unix())), "application/x-www-form-urlencoded", bytes.NewBufferString(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(strings.NewReader(string(sitemap)))
	if err != nil {
		log.Fatal(err)
	}

	sr.htmlNode = doc
}

//分析searchRequest裡面的htmlNode，將各欄位對應到courseDetail中，並儲存在course裡的courseXi中
func (course *course) find(n *html.Node) {
	var detail courseDetail
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "row" {
				count := 0
				for child := n.FirstChild; child != nil; child = child.NextSibling {
					for _, a1 := range child.Attr {
						if a1.Key == "class" && a1.Val == "td1" {
							detail.gradeNClass = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td2" {
							detail.creditNChooseSelect = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td3" {
							detail.courseID = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td4" {
							detail.courseName = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td5" {
							detail.teacher = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td7" {
							detail.courseInfo = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td8" {
							detail.remark = child.FirstChild.Data
							break
						}
						if a1.Key == "class" && a1.Val == "td9" {
							for _, a2 := range child.LastChild.Attr {
								if a2.Key == "href" {
									detail.courseURL = a2.Val
									break
								}
							}

						}
					}
					count++
				}
				course.courseXi = append(course.courseXi, detail)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		course.find(c)
	}
}

//輸出所有courseXi中的資料
func (course *course) print() {
	for _, value := range course.courseXi {
		fmt.Println(value.courseID, value.courseName, value.courseURL)
	}
}

//輸出到Excel檔案
func (course *course) exportToExcel(sheetName string) {
	if sheetName == "數位課綱" {
		sheetName = academicYear + semester + "數位課綱"
	}

	xlsx := excelize.NewFile()
	xlsx.SetSheetName("Sheet1", sheetName)
	xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{"年-班", "學分數/必選別", "科目序號", "科目名稱", "授課教師", "上課時間/地點", "備註", "數位課綱URL"})
	for index, element := range course.courseXi {
		strIndex := strconv.Itoa(index + 2)
		strIndex = "A" + strIndex
		urlIndex := "H" + strIndex
		if element.courseID != "" {
			xlsx.SetCellHyperLink(sheetName, urlIndex, element.courseURL, "External")
			xlsx.SetSheetRow(sheetName, strIndex, &[]interface{}{element.gradeNClass, element.creditNChooseSelect, element.courseID, element.courseName, element.teacher, element.courseInfo, element.remark, element.courseURL})
		}
	}

	// Save xlsx file by the given path.
	err := xlsx.SaveAs(sheetName + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
}
