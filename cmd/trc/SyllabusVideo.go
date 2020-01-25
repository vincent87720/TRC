package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
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
	year                string //學年度
	semester            string //學期
	system              string //學制
	college             string //開課學院
	department          string //開課系所
	gradeNClass         string //年-班
	creditNChooseSelect string //學分-必選別
	courseID            string //課程序號
	courseName          string //課程名稱
	teacher             string //授課教師
	courseInfo          string //上課時間/地點
	remark              string //備註
	courseURL           string //數位課綱連結
	videoProblem        string //數位課綱影片問題
}

type merge struct {
	departmentColumnNum   int //開課系所欄位編號
	teacherColumnNum      int //授課教師欄位編號
	courseIDColumnNum     int //課程序號欄位編號
	courseNameColumnNum   int //課程名稱欄位編號
	videoProblemColumnNum int //影片問題欄位編號
	teacherMap            map[string][]mergeDetail
	syllabusVideoRows     [][]string
	courseXi              []courseDetail
}

type mergeDetail struct {
	teacher      string //授課教師
	department   string //開課系所
	courseID     string //課程序號
	courseName   string //課程名稱
	videoProblem string //數位課綱影片問題
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
	style, err := xlsx.NewStyle(`{"font":{"color":"#1265BE","underline":"single"}}`)
	if err != nil {
		fmt.Println(err)
	}
	xlsx.SetColStyle(sheetName, "H", style)
	for index, element := range course.courseXi {
		strIndex := strconv.Itoa(index + 2)
		strIndex = "A" + strIndex
		urlIndex := "H" + strIndex
		if element.courseID != "" {
			xlsx.SetCellHyperLink(sheetName, urlIndex, element.courseURL, "External")
			err = xlsx.SetCellStyle(sheetName, urlIndex, urlIndex, style)
			if err != nil {
				fmt.Println(err)
			}
			xlsx.SetSheetRow(sheetName, strIndex, &[]interface{}{element.gradeNClass, element.creditNChooseSelect, element.courseID, element.courseName, element.teacher, element.courseInfo, element.remark, element.courseURL})
		}
	}

	// Save xlsx file by the given path.
	err = xlsx.SaveAs(sheetName + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

//讀入數位課綱資料
func (m *merge) loadSyllabusVideoList() error {
	videoListXlsx, err := excelize.OpenFile(inputFilePath)
	if err != nil {
		fmt.Println("\rERROR:", err)
		os.Exit(2)
	}
	videoListXlsxSheetName := "工作表"
	m.syllabusVideoRows, err = videoListXlsx.GetRows(videoListXlsxSheetName)
	if err != nil {
		fmt.Println("\rERROR:找不到\"數位課綱\"檔案內的\"工作表\"")
		os.Exit(2)
	}
	// fmt.Println(syllabusVideoRows[0])
	for index, value := range m.syllabusVideoRows[0] {
		switch value {
		case "開課系所":
			if m.departmentColumnNum == 0 {
				m.departmentColumnNum = index
			}
		case "教師姓名":
			if m.teacherColumnNum == 0 {
				m.teacherColumnNum = index
			}
		case "科目序號":
			if m.courseIDColumnNum == 0 {
				m.courseIDColumnNum = index
			}
		case "科目名稱":
			if m.courseNameColumnNum == 0 {
				m.courseNameColumnNum = index
			}
		case "影片問題":
			if m.videoProblemColumnNum == 0 {
				m.videoProblemColumnNum = index
			}
		}
	}
	// fmt.Println(m.departmentColumnNum, m.teacherColumnNum, m.courseIDColumnNum, m.courseNameColumnNum, m.videoProblemColumnNum)
	return nil
}

//依教師合併科目
func (m *merge) mergeSyllabusVideoList() error {
	m.teacherMap = make(map[string][]mergeDetail)
	for _, value := range m.syllabusVideoRows {

		//尋找影片問題欄位有資料者
		if len(value) >= m.videoProblemColumnNum && value[m.videoProblemColumnNum] != "" {
			md := mergeDetail{
				teacher:      value[m.teacherColumnNum],
				department:   value[m.departmentColumnNum],
				courseID:     value[m.courseIDColumnNum],
				courseName:   value[m.courseNameColumnNum],
				videoProblem: value[m.videoProblemColumnNum],
			}
			if value[m.teacherColumnNum] == "教師姓名" {
				continue //讀取到標題列略過該迴圈
			}
			m.teacherMap[value[m.teacherColumnNum]] = append(m.teacherMap[value[m.teacherColumnNum]], md)
		}
	}

	return nil
}

//輸出合併後的Excel檔案
func (m *merge) exportMergedCourseDataToExcel(fileName string) error {
	xlsx := excelize.NewFile()
	sheetName := "工作表"
	xlsx.SetSheetName("Sheet1", sheetName)

	var column int
	var mark string //當超過Z時會變成AA，超過AZ會變成BA，此變數標記目前標記為何

	//更改工作表網底
	// fillColor1E3048, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#1E3048"],"pattern":1},"font":{"color":"#FFFFFF"}}`)
	fillColorEBF0F3, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#EBF0F3"],"pattern":1}}`)
	fillColorCCDBE2, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#CCDBE2"],"pattern":1}}`)
	fillColor658FA7, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#658FA7"],"pattern":1}}`)

	for i := 67; i < 148; {
		if i < 91 {
			mark = ""
		} else if i < 117 {
			mark = "A"
		} else if i < 143 {
			mark = "B"
		} else if i < 169 {
			mark = "C"
		}
		xlsx.SetColStyle(sheetName, mark+string(i), fillColorEBF0F3)
		xlsx.SetColStyle(sheetName, mark+string(i+1), fillColorCCDBE2)
		xlsx.SetColStyle(sheetName, mark+string(i+2), fillColor658FA7)
		i = i + 3
	}

	xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{"授課教師", "開課系所", "科目序號1", "科目名稱1", "影片問題1", "科目序號2", "科目名稱2", "影片問題2", "科目序號3", "科目名稱3", "影片問題3", "科目序號4", "科目名稱4", "影片問題4", "科目序號5", "科目名稱5", "影片問題5", "科目序號6", "科目名稱6", "影片問題6", "科目序號7", "科目名稱7", "影片問題7", "科目序號8", "科目名稱8", "影片問題8", "科目序號9", "科目名稱9", "影片問題9", "科目序號10", "科目名稱10", "影片問題10"})

	keys := make([]string, 0, len(m.teacherMap))
	for k := range m.teacherMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	row := 2 //Row 2
	for _, key := range keys {
		column = 65 //Column A
		mark = ""
		for cindex, celement := range m.teacherMap[key] {
			if cindex == 0 {
				xlsx.SetCellValue(sheetName, string(column)+strconv.Itoa(row), celement.teacher)
				column++
				xlsx.SetCellValue(sheetName, string(column)+strconv.Itoa(row), celement.department)
				column++
			}
			if column < 90 {
				mark = ""
			} else if column < 116 {
				mark = "A"
			} else if column < 142 {
				mark = "B"
			}
			xlsx.SetCellValue(sheetName, mark+string(column)+strconv.Itoa(row), celement.courseID)
			column++
			xlsx.SetCellValue(sheetName, mark+string(column)+strconv.Itoa(row), celement.courseName)
			column++
			xlsx.SetCellValue(sheetName, mark+string(column)+strconv.Itoa(row), celement.videoProblem)
			column++
		}
		row++
	}

	// Save xlsx file by the given path.
	err := xlsx.SaveAs(fileName + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
