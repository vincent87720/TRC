package main

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
)

// SplitScoreAlertData 分割預警總表
// 使用教師名單(teacherFile)內的資訊建立檔名，以空白表格(templateFile)作為模板並另存至outputFile設定的路徑
// Goroutine interface for GUI
// For example:
// 	var masterFile file
// 	var templateFile file
// 	var teacherFile file
// 	var outputFile file
//
// 	masterFile.setFile("Your file path", "Your file name", "Your sheet name")
// 	templateFile.setFile("Your file path", "Your file name", "Your sheet name")
// 	teacherFile.setFile("Your file path", "Your file name", "Your sheet name")
// 	outputFile.setFile("Your file path", "Your file name", "Your sheet name")
//
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
// 	go SplitScoreAlertData(errChan, exitChan, masterFile, templateFile, teacherFile, outputFile)
//
// Loop:
// 	for {
// 		select {
// 		case err := <-errChan:
// 			Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func SplitScoreAlertData(progChan chan int, inputFile file, templateFile file, teacherFile file, outputFile file) {
	saf := saFile{
		file: inputFile,
	}
	thf := thFile{
		file: teacherFile,
	}

	err := saf.readRawData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = saf.groupByTeacher()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = thf.readRawData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = thf.groupByTeacher()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = templateFile.readRawData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = saf.exportDataToExcel(templateFile, thf, outputFile)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

//groupByTeacher 依照教師名稱將預警總表資料分群
func (saf *saFile) groupByTeacher() (err error) {
	err = saf.findCol("開課學系", &saf.csdCol)
	if err != nil {
		return err
	}
	saf.findCol("科目序號", &saf.cidCol)
	saf.findCol("預警科目", &saf.csnCol)
	saf.findCol("授課教師", &saf.cthCol)
	saf.findCol("學號", &saf.sidCol)
	saf.findCol("學生姓名", &saf.stnCol)
	saf.findCol("預警原由", &saf.alrCol)

	saf.gbtd = make(map[teacher][]scoreAllert)
	if len(saf.dataRows[0]) <= 0 {
		err = errors.WithStack(fmt.Errorf("dataRows has no data"))
		return err
	}
	for index, value := range saf.dataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		if value[saf.cthCol] != "" {
			t := teacher{
				teacherName: value[saf.cthCol],
			}
			sa := scoreAllert{
				course: course{
					department: department{
						departmentName: value[saf.csdCol],
					},
					courseID:   value[saf.cidCol],
					courseName: value[saf.csnCol],
				},
				student: student{
					studentID:   value[saf.sidCol],
					studentName: value[saf.stnCol],
				},
				allertReason: value[saf.alrCol],
			}
			saf.gbtd[t] = append(saf.gbtd[t], sa)
		}
	}
	return nil
}

//groupByTeacher 依照教師名稱將教師資料分群
func (thr *thFile) groupByTeacher() (err error) {
	err = thr.findCol("學院編號", &thr.didCol)
	if err != nil {
		return err
	}
	thr.findCol("教師編號", &thr.tidCol)
	thr.findCol("教師姓名", &thr.trnCol)
	thr.findCol("所屬單位名稱", &thr.tdpCol)

	thr.teacherMap = make(map[string][]teacher)
	if len(thr.dataRows[0]) <= 0 {
		err = errors.WithStack(fmt.Errorf("teacher list dataRows has no data"))
		return err
	}
	for index, value := range thr.dataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		if value[thr.tidCol] != "" {
			t := teacher{
				department: department{
					college: college{
						collegeID: value[thr.didCol],
					},
					departmentName: value[thr.tdpCol],
				},
				teacherID:   value[thr.tidCol],
				teacherName: value[thr.trnCol],
			}
			thr.teacherMap[t.teacherName] = append(thr.teacherMap[t.teacherName], t)
		}
	}
	return nil
}

//exportDataToExcel 匯出預警分表
func (saf *saFile) exportDataToExcel(templateFile file, teacherFile thFile, outputFile file) (err error) {

	if len(saf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for key, value := range saf.gbtd {
		xlsx, err := excelize.OpenFile(templateFile.filePath + templateFile.fileName)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		if len(teacherFile.teacherMap[key.teacherName]) > 0 {
			//教師存在於名單中，設定檔名為"學院編號(int)_系所名稱(string)_教師姓名(string).xlsx"
			clgID := teacherFile.teacherMap[key.teacherName][0].collegeID
			depName := teacherFile.teacherMap[key.teacherName][0].departmentName
			thrName := key.teacherName
			outputFile.fileName = clgID + "_" + depName + "_" + thrName + ".xlsx"
		} else {
			//教師不存在於名單中，設定檔名為"教師姓名(string).xlsx"
			outputFile.fileName = key.teacherName + ".xlsx"
		}

		for index, val := range value {
			row := strconv.Itoa(index + 2)
			position := "A" + row
			err = xlsx.SetSheetRow(templateFile.sheetName, position, &[]interface{}{val.course.department.departmentName, val.courseID, val.courseName, key.teacherName, val.studentID, val.studentName, val.allertReason})
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}

		//使用路徑及檔名匯出檔案
		err = xlsx.SaveAs(outputFile.filePath + outputFile.fileName)
		if err != nil {
			fmt.Println("\rError: 無法將檔案\"" + outputFile.fileName + "\"儲存在\"" + outputFile.filePath + "\"目錄內")
			err = errors.WithStack(err)
			return err
		}
	}
	return nil
}
