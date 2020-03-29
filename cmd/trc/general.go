package main

import (
	"fmt"
	"strconv"

	"github.com/Luxurioust/excelize"
)

//setFile 設定檔案資訊
func (f *file) setFile(filePath string, fileName string, sheetName string) {
	f.filePath = filePath
	f.fileName = fileName
	f.sheetName = sheetName
}

//fillSliceLength 補足slice到指定大小
func (f *file) fillSliceLength(length int) (err error) {
	if len(f.dataRows) <= 0 {
		return fmt.Errorf("dataRows has no data")
	}
	for index := range f.dataRows {
		for len(f.dataRows[index]) < length {
			f.dataRows[index] = append(f.dataRows[index], "")
		}
	}
	return nil
}

//readRawData 讀入檔案資訊
func (f *file) readRawData() (err error) {
	xlsx, err := excelize.OpenFile(f.filePath + f.fileName)
	if err != nil {
		fmt.Println("\rError: 找不到\"" + f.fileName + "\"檔案，請確認檔案名稱是否正確")
		return err
	}

	f.xlsx = xlsx
	// sheetName := "工作表"
	f.dataRows, err = xlsx.GetRows(f.sheetName)
	if err != nil {
		fmt.Println("\rError: 找不到\"" + f.fileName + "\"檔案內的\"工作表\"")
		return err
	}

	f.firstRow = f.dataRows[0]
	return nil
}

//exportDataToExcel 將檔案匯出成xlsx檔案
func (f *file) exportDataToExcel(outputFile file) (err error) {
	xlsx := excelize.NewFile()
	// sheetName := "工作表"
	xlsx.SetSheetName("Sheet1", outputFile.sheetName)

	//將f.newDataRows資料加入到xlsx內
	if len(f.newDataRows) <= 0 {
		return fmt.Errorf("newDataRows has no data")
	}
	for index, value := range f.newDataRows {
		row := strconv.Itoa(index + 1)
		position := "A" + row

		err = xlsx.SetSheetRow(outputFile.sheetName, position, &value)
		if err != nil {
			return err
		}
	}

	//使用路徑及檔名匯出檔案
	err = xlsx.SaveAs(outputFile.filePath + outputFile.fileName)
	if err != nil {
		fmt.Println("\rError: 無法將檔案\"" + outputFile.fileName + "\"儲存在\"" + outputFile.filePath + "\"目錄內")
		return err
	}
	return nil
}

//findCol 尋找檔案內第一列與columnText相符合的儲存格
func (f *file) findCol(columnText string, result *int) (err error) {
	*result = -1 //初始值為-1，若沒找到相對應的字串便會顯示-1
	//尋找"教師姓名"欄位
	if len(f.dataRows[0]) <= 0 {
		return fmt.Errorf("dataRows has no data")
	}
	for index, value := range f.dataRows[0] {
		if value == columnText {
			*result = index
		}
	}
	if *result == -1 {
		fmt.Printf("\rError: \"%s\" column not found\n", columnText)
		return fmt.Errorf("\"%s\" column not found", columnText)
	}
	return nil
}

//printRawData 輸出檔案資訊
func (f *file) printRawData() {
	for _, value := range f.dataRows {
		fmt.Println(value)
	}
}

func (f *file) printNewRawData() {
	for _, value := range f.newDataRows {
		fmt.Println(value)
	}
}
