package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
	"github.com/vincent87720/TRC/internal/logging"
)

type FileInfo struct {
	FilePath  string //檔案路徑
	FileName  string //檔案名稱
	SheetName string //工作表名稱
}

//File 檔案
type WorksheetFile struct {
	FilePath    string     //檔案路徑
	FileName    string     //檔案名稱
	SheetName   string     //工作表名稱
	FirstRow    []string   //第一行
	DataRows    [][]string //資料
	NewDataRows [][]string //新資料
	Xlsx        *excelize.File
}

func GetInitialPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logging.Error.Fatalf("%+v\n", err)
	}
	INITPATH := filepath.ToSlash(dir)
	return INITPATH
}

//SetFile 設定檔案資訊
func (f *WorksheetFile) SetFile(FilePath string, FileName string, SheetName string) {
	f.FilePath = FilePath
	f.FileName = FileName
	f.SheetName = SheetName
}

//FillSliceLength 補足slice到指定大小
func (f *WorksheetFile) FillSliceLength(length int) (err error) {
	if len(f.DataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("DataRows has no data"))
		return err
	}
	for index, _ := range f.DataRows {
		for len(f.DataRows[index]) < length {
			f.DataRows[index] = append(f.DataRows[index], "")
		}
	}
	return nil
}

//ReadRawData 讀入檔案資訊
func (f *WorksheetFile) ReadRawData() (err error) {
	xlsx, err := excelize.OpenFile(f.FilePath + f.FileName)
	if err != nil {
		fmt.Println("\rError: 無法開啟\"" + f.FileName + "\"檔案，請確認檔案名稱是否正確")
		err = errors.WithStack(err)
		return err
	}

	f.Xlsx = xlsx
	// SheetName := "工作表"
	f.DataRows, err = xlsx.GetRows(f.SheetName)
	if err != nil {
		fmt.Println("\rError: 無法讀取\"" + f.FileName + "\"檔案內的工作表")
		err = errors.WithStack(err)
		return err
	}

	f.FirstRow = f.DataRows[0]
	return nil
}

//ExportDataToExcel 將檔案匯出成xlsx檔案
func (f *WorksheetFile) ExportDataToExcel(outputFile WorksheetFile) (err error) {
	xlsx := excelize.NewFile()
	// SheetName := "工作表"
	xlsx.SetSheetName("Sheet1", outputFile.SheetName)

	//將f.newDataRows資料加入到xlsx內
	if len(f.NewDataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("NewDataRows has no data"))
		return err
	}
	for index, value := range f.NewDataRows {
		row := strconv.Itoa(index + 1)
		position := "A" + row

		err = xlsx.SetSheetRow(outputFile.SheetName, position, &value)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
	}

	//使用路徑及檔名匯出檔案
	err = xlsx.SaveAs(outputFile.FilePath + outputFile.FileName)
	if err != nil {
		fmt.Println("\rError: 無法將檔案\"" + outputFile.FileName + "\"儲存在\"" + outputFile.FilePath + "\"目錄內")
		err = errors.WithStack(err)
		return err
	}
	return nil
}

//FindCol 尋找檔案內第一列與columnText相符合的儲存格
func (f *WorksheetFile) FindCol(columnText string, result *int) (err error) {
	*result = -1 //初始值為-1，若沒找到相對應的字串便會顯示-1
	//尋找"教師姓名"欄位
	if len(f.DataRows[0]) <= 0 {
		err = errors.WithStack(fmt.Errorf("DataRows has no data"))
		return err
	}
	for index, value := range f.DataRows[0] {
		if value == columnText {
			*result = index
			break
		}
	}
	if *result == -1 {
		fmt.Printf("\rError: \"%s\" column not found\n", columnText)
		err = errors.WithStack(fmt.Errorf("\"%s\" column not found", columnText))
		return err
	}
	return nil
}

//FindCol 尋找檔案內第一列與columnText相符合的所有儲存格
func (f *WorksheetFile) FindAllCol(columnText string, result *[]int) (err error) {
	//*result = -1 //初始值為-1，若沒找到相對應的字串便會顯示-1
	//尋找"教師姓名"欄位
	if len(f.DataRows[0]) <= 0 {
		err = errors.WithStack(fmt.Errorf("DataRows has no data"))
		return err
	}
	for index, value := range f.DataRows[0] {
		if value == columnText {
			*result = append(*result, index)
		}
	}
	if len(*result) == 0 {
		fmt.Printf("\rError: \"%s\" column not found\n", columnText)
		err = errors.WithStack(fmt.Errorf("\"%s\" column not found", columnText))
		return err
	}
	return nil
}

//PrintRawData 輸出檔案資訊
func (f *WorksheetFile) PrintRawData() {
	for _, value := range f.DataRows {
		fmt.Println(len(value))
	}
}

//PrintNewRawData 輸出檔案資訊
func (f *WorksheetFile) PrintNewRawData() {
	for _, value := range f.NewDataRows {
		fmt.Println(value)
	}
}
