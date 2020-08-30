package main

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

//collegeJSON 從學校取得的單位及所屬單位關聯(UnmarshalJSON用)
type collegeJSON struct {
	Unitid  string
	Contain []departmentJSON
}

//departmentJSON 一個collegeJSON包含多個departmentJSON(UnmarshalJSON用)
type departmentJSON struct {
	Unitid string
}

//unitTitleJSON 單位中英文名稱及簡稱等詳細資訊(UnmarshalJSON用)
type unitTitleJSON struct {
	Title_tw       string
	Title_short_tw string
	Title_en       string
	Title_short_en string
}

//teacherJSON 教師資訊(UnmarshalJSON用)
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

//unitJSON 一個teacherJSON可包含多個unitJSON(UnmarshalJSON用)
type unitJSON struct {
	Unitid    string
	Unittitle string
}

// GetTeacher 下載教師資料
// Goroutine interface for GUI
// For example:
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
// 	var outputFile file
//
// 	outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), "教師名單.xlsx", "工作表")
//
// 	go GetTeacher(errChan, exitChan, outputFile)
// Loop:
// 	for {
// 		select {
// 		case err := <-errChan:
// 			Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func GetTeacher(progChan chan int, outputFile file) {
	var tdreq getTDRequest
	var udreq getUDRequest
	var cudreq getCUDRequest
	var inputFile downloadTeacherFile

	udreq.setURL("https://lg.dyu.edu.tw/get_unit_title.php")
	err := udreq.sendPostRequest()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = udreq.parseData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	cudreq.setURL("http://lg.dyu.edu.tw/search_unit.php")
	err = cudreq.sendPostRequest()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = cudreq.parseData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	tdreq.setURL("https://lg.dyu.edu.tw/search_teacher.php")
	err = tdreq.sendPostRequest()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = tdreq.parseData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = inputFile.transportToSlice(&tdreq, &udreq, &cudreq)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = inputFile.exportDataToExcel(outputFile)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

//parseData 解析取得的教師資訊
func (tdreq *getTDRequest) parseData() (err error) {
	tdreq.teacherInfo = make([]teacherJSON, 0)
	err = json.Unmarshal(tdreq.responseData, &tdreq.teacherInfo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

//parseData 解析取得的單位詳細資訊(單位中英文名稱及簡寫)
func (udreq *getUDRequest) parseData() (err error) {
	udreq.unitTitle = make(map[string]unitTitleJSON, 0)
	err = json.Unmarshal(udreq.responseData, &udreq.unitTitle)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

//parseData 解析取得的單位所屬關係
func (cudreq *getCUDRequest) parseData() (err error) {
	cudreq.unitInfo = make([]collegeJSON, 0)
	err = json.Unmarshal(cudreq.responseData, &cudreq.unitInfo)
	if err != nil {
		err = errors.WithStack(err)
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
	err = errors.WithStack(fmt.Errorf("College not found"))
	return "", err
}

//transportToSlice 將教師資料放入dtFile的newDataRows中，以便使用exportDataToExcel方法輸出
func (dtFile *downloadTeacherFile) transportToSlice(tdreq *getTDRequest, unitData *getUDRequest, containUnitData *getCUDRequest) (err error) {
	dtFile.newDataRows = make([][]string, 0)
	dtFile.newDataRows = append(dtFile.newDataRows, []string{"學院編號", "教師編號", "教師姓名", "所屬單位編號", "所屬單位名稱", "任職狀態", "職稱", "最後更新日期", "分機", "空間代號", "Mail"})

	if len(tdreq.teacherInfo) <= 0 {
		err = errors.WithStack(fmt.Errorf("teacherInfo hasno data"))
		return err
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
