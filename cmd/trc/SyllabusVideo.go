package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
)

// MergeSyllabusVideoData 合併數位課綱資料
// Goroutine interface for GUI
// For example:
// 	var inputFile file
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
// 	go MergeSyllabusVideoData(errChan, exitChan, inputFile, outputFile)
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
func MergeSyllabusVideoData(errChan chan error, exitChan chan string, inputFile file, outputFile file) {

	svf := svFile{
		file: inputFile,
	}

	err := svf.readRawData()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.groupByTeacher()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.matchTeacherInfo()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.transportToSlice()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.exportDataToExcel(outputFile)
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}

	exitChan <- "exit"
	return
}

// MergeSyllabusVideoDataByList 使用教師名單合併數位課綱資料
// Goroutine interface for GUI
// For example:
// 	var inputFile file
// 	var outputFile file
//	var teacherFile file
//
// 	inputFile.setFile("Your file path", "Your file name", "Your sheet name")
// 	outputFile.setFile("Your file path", "Your file name", "Your sheet name")
//	teacherFile.setFile(teacherPathXi[1], teacherPathXi[2], fi.TeacherSheet)
//
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
//	go MergeSyllabusVideoDataByList(errChan, exitChan, inputFile, outputFile, teacherFile)
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
func MergeSyllabusVideoDataByList(errChan chan error, exitChan chan string, inputFile file, outputFile file, teacherFile file) {
	svf := svFile{
		file: inputFile,
	}

	thf := thFile{
		file: teacherFile,
	}

	err := svf.readRawData()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.groupByTeacher()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = thf.readRawData()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = thf.groupByTeacher()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.matchTeacherInfoFile(thf)
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.transportToSlice()
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}
	err = svf.exportDataToExcel(outputFile)
	if err != nil {
		errChan <- err
		exitChan <- "exit"
		return
	}

	exitChan <- "exit"
	return
}

//groupByTeacher 依照教師名稱將數位課綱資料分群
func (svf *svFile) groupByTeacher() (err error) {
	err = svf.findCol("教師姓名", &svf.cthCol)
	if err != nil {
		return err
	}
	err = svf.findCol("所屬單位", &svf.cdpCol)
	if err != nil {
		return err
	}
	err = svf.findCol("科目序號", &svf.cidCol)
	if err != nil {
		return err
	}
	err = svf.findCol("科目名稱", &svf.csnCol)
	if err != nil {
		return err
	}
	err = svf.findCol("影片問題", &svf.pocCol)
	if err != nil {
		return err
	}

	svf.gbtd = make(map[teacher][]syllabusVideo)
	if len(svf.dataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("dataRows has no data"))
		return err
	}
	for index, value := range svf.dataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		//尋找影片問題欄位有資料者
		if len(value) >= len(svf.firstRow) && value[svf.pocCol] != "" {
			t := teacher{
				teacherName: value[svf.cthCol],
			}
			c := syllabusVideo{
				course: course{
					courseID:   value[svf.cidCol],
					courseName: value[svf.csnCol],
					department: department{
						departmentName: value[svf.cdpCol],
					},
				},
				problemOfCourse: value[svf.pocCol],
			}
			svf.gbtd[t] = append(svf.gbtd[t], c)
		}
	}
	return nil
}

//matchTeacherInfo 使用原本svFile檔案內的所屬單位進行教師比對合併
func (svf *svFile) matchTeacherInfo() (err error) {
	for key, value := range svf.gbtd {
		t := teacher{
			teacherName: key.teacherName,
			department: department{
				departmentName: value[0].course.department.departmentName,
			},
		}
		delete(svf.gbtd, key) //必須先刪除再加入，否則有可能誤刪
		svf.gbtd[t] = value
	}
	return nil
}

//matchTeacherInfoFile 使用額外輸入的teacherFile檔案進行教師比對合併
func (svf *svFile) matchTeacherInfoFile(teacherFile thFile) (err error) {
	for key, value := range svf.gbtd {
		if len(teacherFile.teacherMap[key.teacherName]) > 0 {

			t := teacher{
				teacherName: key.teacherName,
				department: department{
					departmentName: teacherFile.teacherMap[key.teacherName][0].department.departmentName,
				},
			}
			delete(svf.gbtd, key) //必須先刪除再加入，否則有可能誤刪
			svf.gbtd[t] = value
		}
	}
	return nil
}

//transportToSlice 將map[teacher][]syllabusVideo的資料轉換為二維陣列
func (svf *svFile) transportToSlice() (err error) {
	if len(svf.dataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for key, value := range svf.gbtd {

		//計算若目前老師的科目數量會佔幾列
		rowNum := len(value) / 9
		if len(value)%9 != 0 {
			rowNum++
		}

		//將rowNum列加入到svf.mergedXi中
		for i := 0; i < rowNum; i++ {
			tempXi := make([]string, 0)
			tempXi = append(tempXi, key.teacherName, key.department.departmentName)

			//每列只能放9個，多的給下一圈執行
			for j := 0; j < 9; j++ {
				//放到最後一個為止
				if i*9+j >= len(value) {
					break
				}
				tempXi = append(tempXi, value[i*9+j].courseName, value[i*9+j].courseID, value[i*9+j].problemOfCourse)
			}

			//若目前老師的科目數量大於最大數量，將其設為最大數量
			if len(value) > svf.maxCourseNum {
				svf.maxCourseNum = len(value)
			}

			svf.mergedXi = append(svf.mergedXi, tempXi)
		}

	}
	return nil
}

//exportDataToExcel 將資料匯出至xlsx檔案
func (svf *svFile) exportDataToExcel(outputFile file) (err error) {
	xlsx := excelize.NewFile()
	// sheetName := "工作表"
	xlsx.SetSheetName("Sheet1", outputFile.sheetName)

	//更改工作表網底
	var mark string //當超過Z時會變成AA，超過AZ會變成BA，此變數標記目前標記為何
	// fillColorEFECD7, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#EFECD7"],"pattern":1}}`)
	// if err != nil {
	// 	err = errors.WithStack(err)
	// 	return err
	// }
	fillColorE9E7D6, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#E9E7D6"],"pattern":1}}`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	fillColorE0E4D6, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#E0E4D6"],"pattern":1}}`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	fillColorDADCD2, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#DADCD2"],"pattern":1}}`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	for i := 67; i < 67+svf.maxCourseNum*3; i++ {
		if i < 91 {
			mark = string(rune(i))
		} else if i < 117 {
			mark = "A" + string(rune(i-26))
		} else if i < 143 {
			mark = "B" + string(rune(i-26*2))
		} else if i < 169 {
			mark = "C" + string(rune(i-26*3))
		}

		switch i % 3 {
		case 1:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorE9E7D6)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 2:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorE0E4D6)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 0:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorDADCD2)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
	}

	//設定第一列
	title := make([]string, 0)
	title = append(title, "授課教師", "所屬單位")
	for i := 0; i < svf.maxCourseNum; i++ {
		title = append(title, "科目名稱", "科目序號", "影片問題")
	}
	err = xlsx.SetSheetRow(outputFile.sheetName, "A1", &title)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	//依照unicode排序gptd map
	keys := make([]string, 0, len(svf.gbtd))
	if len(svf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for k := range svf.gbtd {
		keys = append(keys, k.teacherName)
	}
	sort.Strings(keys)

	//將svf.mergedXi資料加入到xlsx內
	if len(svf.mergedXi) <= 0 {
		err = errors.WithStack(fmt.Errorf("mergedXi has no data"))
		return err
	}
	for index, value := range svf.mergedXi {
		row := strconv.Itoa(index + 2)
		position := "A" + row

		err = xlsx.SetSheetRow(outputFile.sheetName, position, &value)
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
	return nil
}
