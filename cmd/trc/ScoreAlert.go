package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Luxurioust/excelize"
)

type scoreAlert struct {
	teacherMap     map[string][]string
	scoreAlertRows [][]string
}

//load teacher data list
func (sa *scoreAlert) loadTeacherInfo() (err error) {
	sa.teacherMap = make(map[string][]string)
	teacherXlsx, err := excelize.OpenFile(teacherInfoFilePath)
	if err != nil {
		return err
	}
	teacherXlsxSheetName := "工作表"
	teacherRows, err := teacherXlsx.GetRows(teacherXlsxSheetName)
	if err != nil {
		fmt.Println("\rERROR:找不到\"教師名單\"檔案內的\"工作表\"")
		os.Exit(2)
	}
	for _, value := range teacherRows {
		if value[2] != "" {
			sa.teacherMap[value[2]] = value
		}
	}
	return nil
}

//load score alert master list
func (sa *scoreAlert) loadScoreAlertList() (err error) {
	scoreAlert, err := excelize.OpenFile(scoreAlertFilePath)
	scoreAlertSheetName := "工作表"
	if err != nil {
		return err
	}
	sa.scoreAlertRows, err = scoreAlert.GetRows(scoreAlertSheetName)
	if err != nil {
		fmt.Println("\rERROR:找不到\"預警總表\"檔案內的\"工作表\"")
		os.Exit(2)
	}
	return nil
}

func (sa *scoreAlert) splitScoreAlertData() {
	lastTeacher := "nil"
	var xiOfrowXi [][]string
	for rowsIndex, rowValue := range sa.scoreAlertRows {
		if rowsIndex == 0 {
			continue
		}
		if rowsIndex == 1 {
			xiOfrowXi = append(xiOfrowXi, rowValue)
			lastTeacher = rowValue[3]
			continue
		}
		if rowValue[3] != lastTeacher || rowsIndex == len(sa.scoreAlertRows)-1 {
			//Find teacher data
			var filename string
			if len(sa.teacherMap[lastTeacher]) == 0 || sa.teacherMap[lastTeacher] == nil {
				filename = lastTeacher + ".xlsx"
			} else {
				filename = sa.teacherMap[lastTeacher][0] + "_" + sa.teacherMap[lastTeacher][3] + "_" + sa.teacherMap[lastTeacher][2] + ".xlsx"
			}

			sa.exportToExcel(filename, xiOfrowXi)

			xiOfrowXi = nil
		}
		xiOfrowXi = append(xiOfrowXi, rowValue)
		lastTeacher = rowValue[3]
	}
}

func (sa *scoreAlert) exportToExcel(sheetName string, xiOfrowXi [][]string) (err error) {
	xlsxOutputFile, err := excelize.OpenFile(exportTemplateFilePath)
	if err != nil {
		return err
	}

	xlsxOutputFileExists := false
	for _, name := range xlsxOutputFile.GetSheetMap() {
		if name == "工作表" {
			xlsxOutputFileExists = true
		}
	}
	if !xlsxOutputFileExists {
		fmt.Println("\rERROR:找不到\"空白預警分表\"檔案內的\"工作表\"")
		os.Exit(2)
	}

	xlsxOutputFileSheetName := "工作表"
	for index, value := range xiOfrowXi {
		strIndex := strconv.Itoa(index + 2)
		strIndex = "A" + strIndex
		if value[0] != "" && value != nil {
			xlsxOutputFile.SetSheetRow(xlsxOutputFileSheetName, strIndex, &[]interface{}{value[0], value[1], value[2], value[3], value[4], value[5], value[6]})
		}
	}

	xlsxOutputFile.SaveAs(sheetName)
	return nil
}
