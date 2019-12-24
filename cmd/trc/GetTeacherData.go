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

	"github.com/Luxurioust/excelize"
)

type searchTeacherRequest struct {
	dept       string
	deptItem   string
	searchItem string
	sitemap    []byte
}

type teacher struct {
	teacherXi []string
}

func (req *searchTeacherRequest) getTeacherData() {
	v := url.Values{}
	//post form data
	v.Add("dept", req.dept)
	v.Add("dept_item", req.deptItem)
	v.Add("search_item", req.searchItem)
	v.Add("te_name", "")
	v.Add("addScholarship1", "")
	v.Add("addScholarship2", "")
	v.Add("ext", "")

	//送出請求並將返回結果放入res
	res, err := http.Post("http://people.dyu.edu.tw/index.php", "application/x-www-form-urlencoded", bytes.NewBufferString(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	req.sitemap = sitemap
}

func (t *teacher) parseTeacherData(deptItem string, sitemap []byte) {
	//搜尋sitemap內字串，將括號內的教師資料放入tempstring
	status := 0 //0=none,1=ready,2=startreading
	var tempstring []byte
	for index, element := range sitemap {
		if status == 0 {
			if string(sitemap[index:index+48]) == "teacher_array[teacher_array.length] = new Array(" {
				status = 1
				continue
			}
		} else if status == 1 {
			if string(element) == "(" {
				status = 2
				tempstring = nil
				continue
			}
		} else if status == 2 {
			if string(element) == ")" {
				status = 0

				// var cusDeptID string
				// switch deptItem {
				// case "1":
				// 	cusDeptID = "01"
				// case "2":
				// 	cusDeptID = "02"
				// case "3":
				// 	cusDeptID = "03"
				// case "4":
				// 	cusDeptID = "04"
				// case "5":
				// 	cusDeptID = "05"
				// case "6":
				// 	cusDeptID = "06"
				// case "9":
				// 	cusDeptID = "09"
				// }
				t.teacherXi = append(t.teacherXi, "\""+deptItem+"\","+string(tempstring))
				// fmt.Println("\r" + string(tempstring))
				continue
			}
			tempstring = append(tempstring, element)
		}
	}
}

func (t *teacher) ExportToExcel(filename string) {
	sheetName := "工作表"
	xlsx := excelize.NewFile()
	xlsx.SetSheetName("Sheet1", sheetName)
	xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{"學院編號", "教師編號", "教師姓名", "教師系所", "所屬單位編號", "任職狀態", "職稱", "最後更新日期", "teno", "分機"})
	for index, element := range t.teacherXi {
		s := strings.Replace(element, "\"", "", -1)
		sXi := strings.Split(s, ",")
		strIndex := strconv.Itoa(index + 1)
		strIndex = "A" + strIndex
		if sXi[0] != "" && sXi != nil {
			xlsx.SetSheetRow(sheetName, strIndex, &[]interface{}{sXi[0], sXi[1], sXi[2], sXi[3], sXi[4], sXi[5], sXi[6], sXi[7], sXi[8], sXi[9]})
		}
	}

	// Save xlsx file by the given path.
	err := xlsx.SaveAs(filename + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func (t *teacher) print() {
	for _, value := range t.teacherXi {
		fmt.Println(value)
	}
}
