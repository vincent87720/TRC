package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type rapidPrintXi []rapidPrint

func (rp rapidPrintXi) Len() int           { return len(rp) }
func (rp rapidPrintXi) Less(i, j int) bool { return rp[i].courseName < rp[j].courseName }
func (rp rapidPrintXi) Swap(i, j int)      { rp[i], rp[j] = rp[j], rp[i] }
func (rp rapidPrintXi) Sort()              { sort.Sort(rp) }

// MergeRapidPrintData 剖析並合併開課資料
// Goroutine interface for GUI
// For example:
//
//	var inputFile file
// 	var outputFile file
//
// 	inputFile.setFile("Your file path", "Your file name", "Your sheet name")
// 	outputFile.setFile("Your file path", "Your file name", "Your sheet name")
//
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
// 	go MergeRapidPrintData(errChan, exitChan, inputFile, outputFile)
//
// Loop:
//	for {
// 		select {
// 		case err := <-errChan:
// 			fmt.Println(err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func MergeRapidPrintData(errChan chan error, exitChan chan string, inputFile file, outputFile file) {

	rpf := rpFile{
		file: inputFile,
	}

	err := rpf.readRawData()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = rpf.fillSliceLength(15)
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = rpf.findColumn()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = rpf.groupByTeacher()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = rpf.mergeData()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = rpf.transportToSlice()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = rpf.exportDataToExcel(outputFile)
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}

	exitChan <- "exit"
	return
}

//findColumn 尋找欄位
func (rpf *rpFile) findColumn() (err error) {
	err = rpf.findCol("教師姓名", &rpf.trnCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("專兼任別", &rpf.tstCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("開課學制", &rpf.sysCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("開課系所", &rpf.csdCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("科目序號", &rpf.cidCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("系-組-年-班", &rpf.ifoCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("科目名稱", &rpf.csnCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("選修別", &rpf.ccsCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("學分", &rpf.cdtCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("時數", &rpf.ctmCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("星期-時間-教室", &rpf.wtrCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("選課人數", &rpf.nopCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("合班註記", &rpf.annCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("合班序號", &rpf.aidCol)
	if err != nil {
		return err
	}
	err = rpf.findCol("備註", &rpf.rmkCol)
	if err != nil {
		return err
	}
	return nil
}

//groupByTeacher 依照教師名稱將開課資料分群
func (rpf *rpFile) groupByTeacher() (err error) {

	r1, err := regexp.Compile(`\(.*?\s`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r2, err := regexp.Compile(`\(.*?\-`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r3, err := regexp.Compile(`[A-Z].*`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	rpf.gbtd = make(map[teacher][]rapidPrint)
	if len(rpf.dataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("dataRows has no data"))
		return err
	}
	for index, value := range rpf.dataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		//使用正規表達式拆解"星期-時間-教室"欄位
		thisTime := make([]string, 0)
		thisClassRoom := make([]string, 0)

		substrs := r1.FindAllString(value[rpf.wtrCol], -1)
		for _, loc := range substrs {
			loc = strings.Replace(loc, " ", "", -1)
			splitWtrLoc := r2.FindStringIndex(loc)
			thisTime = append(thisTime, loc[splitWtrLoc[0]:splitWtrLoc[1]-1])
			thisClassRoom = append(thisClassRoom, loc[splitWtrLoc[1]:])
		}

		//使用正規表達式拆解"科目名稱"欄位，若課程名稱包含副標題(ex.G3、J1...)則將其分開，並將課程名稱放入thisCourseName，副標題放入thisCourseSubName
		var thisCourseName string
		var thisCourseSubName string
		splitCsnLoc := r3.FindStringIndex(value[rpf.csnCol])
		if len(splitCsnLoc) > 0 {
			thisCourseName = value[rpf.csnCol][:splitCsnLoc[0]]
			thisCourseSubName = value[rpf.csnCol][splitCsnLoc[0]:]
		} else {
			thisCourseName = value[rpf.csnCol]
		}

		if value[rpf.trnCol] != "" {
			t := teacher{
				teacherName:  value[rpf.trnCol],
				teacherState: value[rpf.tstCol],
			}
			c := rapidPrint{
				course: course{
					system: value[rpf.sysCol],
					department: department{
						departmentName: value[rpf.csdCol],
						departmentID:   value[rpf.ifoCol][0:4],
					},
					courseID:      value[rpf.cidCol],
					courseName:    thisCourseName,
					courseSubName: thisCourseSubName,
					group:         value[rpf.ifoCol][5:8],
					grade:         value[rpf.ifoCol][9:10],
					class:         value[rpf.ifoCol][11:12],
					chooseSelect:  value[rpf.ccsCol],
					credit:        value[rpf.cdtCol],
					interval:      value[rpf.ctmCol],
					time:          thisTime,
					classRoom:     thisClassRoom,
					numOfPeople:   value[rpf.nopCol],
					annex:         value[rpf.annCol],
					annexID:       value[rpf.aidCol],
					remark:        value[rpf.rmkCol],
				},
			}
			rpf.gbtd[t] = append(rpf.gbtd[t], c)
		}
	}
	return nil
}

//mergeData 依照規則合併重複的科目名稱
func (rpf *rpFile) mergeData() (err error) {
	//對一位老師的所有科目進行處理
	mergedRpData := make(map[teacher][]rapidPrint)
	if len(rpf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for key, value := range rpf.gbtd {

		//排序
		var rpp rapidPrintXi
		rpp = value
		rpp.Sort()

		var nextLoop int
		nextLoop = 0

		//比對目前科目(rpp[i])與下一個科目(rpp[index])是否相同
		for index, _ := range rpp {
			if index < nextLoop {
				continue
			}

			tempSystemMap := make(map[string]int)     //暫存相同課程名稱的學制
			tempRemarkMap := make(map[string]int)     //暫存相同課程名稱的備註
			tempSubNameMap := make(map[string]int)    //暫存相同課程名稱的副標題
			tempDepartmentMap := make(map[string]int) //暫存相同課程名稱的開課系所
			tempCourseIDXi := make([]string, 0)       //暫存相同課程名稱的課程編號
			tempCourseTimeXi := make([]string, 0)     //暫存相同課程名稱的上課時間地點
			tempNumOfPeopleXi := make([]string, 0)    //暫存相同課程名稱的修課人數

			sameCourseNumCount := 0 //紀錄目前科目以外和目前科目可以合併的科目數量

			//將目前科目的資訊放入暫存變數中
			tempSystemMap[rpp[index].system] = 1
			tempRemarkMap[rpp[index].remark] = 1
			tempSubNameMap[rpp[index].courseSubName] = 1
			tempDepartmentMap[rpp[index].course.department.departmentName] = 1
			tempNumOfPeopleXi = append(tempNumOfPeopleXi, rpp[index].numOfPeople)
			for i := range rpp[index].time {
				tempCourseTimeXi = append(tempCourseTimeXi, rpp[index].time[i]+" "+rpp[index].classRoom[i])
			}

			//此為下一個科目(rpp[index])
			for nextIndex := index + 1; nextIndex < len(rpp); nextIndex++ {

				//若符合名稱相同及學制限制，則進行合併
				if rpp[index].courseName == rpp[nextIndex].courseName &&
					(((rpp[index].system == "大日" || rpp[index].system == "進" || rpp[index].system == "四技") && (rpp[nextIndex].system == "大日" || rpp[nextIndex].system == "進" || rpp[nextIndex].system == "四技")) ||
						((rpp[index].system == "博" || rpp[index].system == "日碩" || rpp[index].system == "碩") && (rpp[nextIndex].system == "博" || rpp[nextIndex].system == "日碩" || rpp[nextIndex].system == "碩")) ||
						((rpp[index].system == "日碩" || rpp[index].system == "碩在職") && (rpp[nextIndex].system == "日碩" || rpp[nextIndex].system == "碩在職"))) {

					//設定學制的map為1
					tempSystemMap[rpp[nextIndex].system] = 1

					//設定副標題的map為1，若已存在則累加
					v1, found := tempSubNameMap[rpp[nextIndex].courseSubName]
					if found {
						tempSubNameMap[rpp[nextIndex].courseSubName] = v1 + 1
					} else {
						tempSubNameMap[rpp[nextIndex].courseSubName] = 1
					}

					//設定備註的map為1
					tempRemarkMap[rpp[nextIndex].remark] = 1

					//設定學系的map為1
					tempDepartmentMap[rpp[nextIndex].course.department.departmentName] = 1

					//暫存人數欄位到tempNumOfPeopleXi
					tempNumOfPeopleXi = append(tempNumOfPeopleXi, rpp[nextIndex].numOfPeople)

					//暫存科目序號到tempCourseIDXi
					tempCourseIDXi = append(tempCourseIDXi, rpp[nextIndex].courseID)

					//暫存時間及教室欄位到tempCourseTimeXi
					for i := range rpp[nextIndex].time {
						tempCourseTimeXi = append(tempCourseTimeXi, rpp[nextIndex].time[i]+" "+rpp[nextIndex].classRoom[i])
					}

					sameCourseNumCount++

				} else {
					break
				}
			}

			//串接學制map到tempSystemMap
			if len(tempSystemMap) > 0 {
				var newSystem string
				var count bool
				count = false
				for key := range tempSystemMap {
					if count == false {
						newSystem = newSystem + key
						count = true
					} else {
						newSystem = newSystem + "," + key
					}

				}
				rpp[index].system = newSystem
			}

			//串接系所
			if len(tempDepartmentMap) > 0 {
				var newDep string
				var count bool
				count = false
				for key := range tempDepartmentMap {
					if count == false {
						newDep = newDep + key
						count = true
					} else {
						newDep = newDep + "," + key
					}

				}
				rpp[index].course.department.departmentName = newDep
			}

			//在科目名稱欄位加上科目名稱數量及串接副標題
			if len(tempSubNameMap) > 0 {
				var newSubName string
				count := false
				for key, value := range tempSubNameMap {
					if count == false {
						newSubName = newSubName + key + "*" + strconv.Itoa(value)
						count = true
					} else {
						newSubName = newSubName + "," + key + "*" + strconv.Itoa(value)
					}
				}
				rpp[index].courseName = rpp[index].courseName + newSubName
			}

			//串接備註
			if len(tempRemarkMap) > 0 {
				var newRemark string
				count := false
				for key, _ := range tempRemarkMap {
					if key != "" {
						if count == false {
							newRemark = newRemark + key
							count = true
						} else {
							newRemark = newRemark + "、" + key
						}
					}
				}
				rpp[index].remark = newRemark
			}

			//串接科目序號
			if len(tempCourseIDXi) > 0 {
				var newCourseID string
				if len(rpp[index].remark) == 0 {
					newCourseID = rpp[index].remark + "合"

				} else {
					newCourseID = rpp[index].remark + "、合"
				}

				for index, idXi := range tempCourseIDXi {
					if index == 0 {
						newCourseID = newCourseID + idXi
					} else {
						newCourseID = newCourseID + "," + idXi
					}
				}
				rpp[index].remark = newCourseID
			}

			//串接時間
			if len(tempCourseTimeXi) > 0 {
				var newTime string
				for index, timeXi := range tempCourseTimeXi {
					if index == 0 {
						newTime = newTime + timeXi
					} else {
						newTime = newTime + "," + timeXi
					}
				}
				rpp[index].timeNClassRoom = newTime
			}

			//串接人數
			if len(tempNumOfPeopleXi) > 0 {
				var newNOP string
				for index, nopXi := range tempNumOfPeopleXi {
					if index == 0 {
						newNOP = newNOP + nopXi
					} else {
						newNOP = newNOP + "+" + nopXi
					}
				}
				rpp[index].numOfPeople = newNOP
			}

			//設定下次要執行的圈數
			nextLoop = index + sameCourseNumCount + 1

			//將tempMerged加入到mergedRpData
			mergedRpData[key] = append(mergedRpData[key], rpp[index])
		}
	}

	//將合併過後的檔案加入
	rpf.gbtd = mergedRpData
	return nil
}

func (rpf *rpFile) transportToSlice() (err error) {
	//設定第一列
	rpf.newDataRows = append(rpf.newDataRows, rpf.firstRow)
	if len(rpf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for key, rpXiValue := range rpf.gbtd {
		for _, rpValue := range rpXiValue {
			tempXi := make([]string, 0)
			tempXi = append(tempXi, key.teacherName, key.teacherState, rpValue.system, rpValue.course.departmentName, rpValue.courseID, rpValue.course.departmentID+"-"+rpValue.group+"-"+rpValue.grade+"-"+rpValue.class, rpValue.courseName, rpValue.chooseSelect, rpValue.credit, rpValue.interval, rpValue.timeNClassRoom, rpValue.numOfPeople, rpValue.annex, rpValue.annexID, rpValue.remark)
			rpf.newDataRows = append(rpf.newDataRows, tempXi)
		}
	}
	return nil
}
