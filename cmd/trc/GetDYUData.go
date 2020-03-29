package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/html"
)

type request struct {
	url          string
	values       url.Values
	responseData []byte
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
	smye     string //academicYear
	smty     string //semester
	day      string
	lesson   string
	htmlNode *html.Node
}

//collegeJSON 從學校取得的單位及所屬單位關聯
type collegeJSON struct {
	Unitid  string
	Contain []departmentJSON
}

//departmentJSON 一個collegeJSON包含多個departmentJSON
type departmentJSON struct {
	Unitid string
}

//unitTitleJSON 單位中英文名稱及簡稱等詳細資訊
type unitTitleJSON struct {
	Title_tw       string
	Title_short_tw string
	Title_en       string
	Title_short_en string
}

//teacherJSON 教師資訊
type teacherJSON struct {
	Teno      string     //教師編號
	Name      string     //教師姓名
	Unit      []unitJSON //所屬單位
	Worktime  string     //任職狀態(專兼任)
	Title     string     //職稱
	Update    string     //更新日期
	Ext       string     //分機
	Room      string     //空間代號
	Mail      string     //Mail
	Specialty string     //專長
	Image     string     //圖片
}

//unitJSON 一個teacherJSON可包含多個unitJSON
type unitJSON struct {
	Unitid    string
	Unittitle string
}

func (req *request) newRequest() {
	req.values = url.Values{}
	req.responseData = make([]byte, 0)
}

//setURL 設定請求網址
func (req *request) setURL(url string) {
	req.url = url
}

//setURLValues 設定發送(數位課綱影片連結)請求的參數
func (svlreq *getSVLRequest) setURLValues(academicYear string, semester string, day string, lesson string) (err error) {
	yearInt, err := strconv.Atoi(academicYear)
	if err != nil || yearInt < 0 {
		fmt.Printf("無法解析\"%s\"，請輸入合法的年分", academicYear)
		return fmt.Errorf("Incorrect year value")
	} else {
		svlreq.values.Add("smye", strconv.Itoa(yearInt))
	}

	semesterInt, err := strconv.Atoi(semester)
	if err != nil || semesterInt < 0 {
		fmt.Printf("無法解析\"%s\"，請輸入合法的學期", semester)
		return fmt.Errorf("Incorrect year value")
	} else {
		svlreq.values.Add("smty", strconv.Itoa(semesterInt))
	}

	svlreq.values.Add("str_time", day+"sec"+lesson)
	return nil
}

//sendRequest 發送請求
func (req *request) sendRequest() (err error) {

	//送出請求並將返回結果放入res
	res, err := http.Post(req.url, "application/x-www-form-urlencoded", bytes.NewBufferString(req.values.Encode()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//取得res的body並放入sitemap
	req.responseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}

//parseData 解析取得的教師資訊
func (tdreq *getTDRequest) parseData() (err error) {
	tdreq.teacherInfo = make([]teacherJSON, 0)
	err = json.Unmarshal(tdreq.responseData, &tdreq.teacherInfo)
	if err != nil {
		return err
	}
	return nil
}

//parseData 解析取得的單位詳細資訊(單位中英文名稱及簡寫)
func (udreq *getUDRequest) parseData() (err error) {
	udreq.unitTitle = make(map[string]unitTitleJSON, 0)
	err = json.Unmarshal(udreq.responseData, &udreq.unitTitle)
	if err != nil {
		return err
	}
	return nil
}

//parseData 解析取得的單位所屬關係
func (cudreq *getCUDRequest) parseData() (err error) {
	cudreq.unitInfo = make([]collegeJSON, 0)
	err = json.Unmarshal(cudreq.responseData, &cudreq.unitInfo)
	if err != nil {
		return err
	}
	return nil
}

//findCollege 查詢所屬單位(學院)，department傳入單位(系所)代號，college返回所屬單位(學院)代號
func (cudreq *getCUDRequest) findCollege(department string) (college string, err error) {
	for _, coll := range cudreq.unitInfo {
		for _, dep := range coll.Contain {
			if coll.Unitid == department || dep.Unitid == department {
				college = coll.Unitid
				return college, nil
			}
		}
	}
	return "", fmt.Errorf("College not found")
}

//transportToSlice 將教師資料放入dtFile的newDataRows中，以便使用exportDataToExcel方法輸出
func (dtFile *downloadTeacherFile) transportToSlice(tdreq *getTDRequest, unitData *getUDRequest, containUnitData *getCUDRequest) (err error) {
	dtFile.newDataRows = make([][]string, 0)
	dtFile.newDataRows = append(dtFile.newDataRows, []string{"學院編號", "教師編號", "教師姓名", "所屬單位編號", "所屬單位名稱", "任職狀態", "職稱", "最後更新日期", "分機", "空間代號", "Mail"})

	if len(tdreq.teacherInfo) <= 0 {
		return fmt.Errorf("teacherInfo hasno data")
	}
	for _, value := range tdreq.teacherInfo {
		//無所屬單位
		if len(value.Unit) == 0 {
			dtFile.newDataRows = append(dtFile.newDataRows, []string{"0", value.Teno, value.Name, "", "", value.Worktime, value.Title, value.Update, value.Ext, value.Room, value.Mail})
			continue
		}
		for _, unit := range value.Unit {
			tempXi := make([]string, 0)
			collID, err := containUnitData.findCollege(unit.Unitid)
			if err != nil {
				//College not found
				tempXi = append(tempXi, "0")
			} else {
				switch collID {
				case "2000":
					tempXi = append(tempXi, "1") //工學院院部
				case "2003":
					tempXi = append(tempXi, "2") //管理學院院部
				case "2004":
					tempXi = append(tempXi, "3") //設計暨藝術學院院部
				case "2005":
					tempXi = append(tempXi, "4") //外語學院院部
				case "2006":
					tempXi = append(tempXi, "5") //生物科技暨資源學院院部
				case "2007":
					tempXi = append(tempXi, "6") //觀光餐旅學院院部
				case "2008":
					tempXi = append(tempXi, "7") //護理暨健康學院院部
				case "9999":
					tempXi = append(tempXi, "8") //其他
				}
			}
			tempXi = append(tempXi, value.Teno, value.Name, unit.Unitid, unitData.unitTitle[unit.Unitid].Title_tw, value.Worktime, value.Title, value.Update, value.Ext, value.Room, value.Mail)
			dtFile.newDataRows = append(dtFile.newDataRows, tempXi)
		}
	}
	return nil
}
