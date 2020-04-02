package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Luxurioust/excelize"
	"github.com/pkg/errors"
)

//MergeSyllabusVideoData 合併數位課綱資料
func (svf *svFile) MergeSyllabusVideoData(outputFile file) (err error) {
	err = svf.readRawData()
	if err != nil {
		return err
	}
	err = svf.groupByTeacher()
	if err != nil {
		return err
	}
	err = svf.transportToSlice()
	if err != nil {
		return err
	}
	err = svf.exportDataToExcel(outputFile)
	if err != nil {
		return err
	}
	return nil
}

//groupByTeacher 依照教師名稱將數位課綱資料分群
func (svf *svFile) groupByTeacher() (err error) {
	err = svf.findCol("教師姓名", &svf.cthCol)
	if err != nil {
		return err
	}
	err = svf.findCol("開課系所", &svf.depCol)
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
					department: department{
						departmentName: value[svf.depCol],
					},
					courseID:   value[svf.cidCol],
					courseName: value[svf.csnCol],
				},
				problemOfCourse: value[svf.pocCol],
			}
			svf.gbtd[t] = append(svf.gbtd[t], c)
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
		tempXi := make([]string, 0)
		tempXi = append(tempXi, key.teacherName)

		coutCourseNum := 0 //計算每位老師的科目數量

		for _, array := range value {
			tempXi = append(tempXi, array.departmentName, array.courseID, array.courseName, array.problemOfCourse)
			coutCourseNum++
		}

		//若目前老師的科目數量大於最大數量，將其設為最大數量
		if coutCourseNum > svf.maxCourseNum {
			svf.maxCourseNum = coutCourseNum
		}

		svf.mergedXi = append(svf.mergedXi, tempXi)
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
	fillColorEFECD7, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#EFECD7"],"pattern":1}}`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
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

	for i := 66; i < 66+svf.maxCourseNum*4; i++ {
		if i < 91 {
			mark = string(i)
		} else if i < 117 {
			mark = "A" + string(i-26)
		} else if i < 143 {
			mark = "B" + string(i-26*2)
		} else if i < 169 {
			mark = "C" + string(i-26*3)
		}

		switch i % 4 {
		case 2:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorEFECD7)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 3:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorE9E7D6)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 0:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorE0E4D6)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 1:
			err := xlsx.SetColStyle(outputFile.sheetName, mark, fillColorDADCD2)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
	}

	//設定第一列
	title := make([]string, 0)
	title = append(title, "授課教師")
	for i := 0; i < svf.maxCourseNum; i++ {
		title = append(title, "開課系所", "科目序號", "科目名稱", "影片問題")
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
