package Difference

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
	"github.com/vincent87720/TRC/internal/cmdline"
	file "github.com/vincent87720/TRC/internal/file"
	"github.com/vincent87720/TRC/internal/logging"
	"github.com/vincent87720/TRC/internal/object"
)

type facing struct {
	facingName string    //面向名稱
	score      []float64 //各評委分數
	difference []float64 //差分
}

//dcFile 差分計算檔案
type dcFile struct {
	file.WorksheetFile
	judgeNum int                         //評審數量
	judge    []string                    //評審
	gbsd     map[object.Student][]facing //group by object.Student data
}

//calcDifferenceAllFile 將目錄內所有尺規評量成績評分表進行差分計算
func DifferenceCalculate_AllFile_Command(inputFileInfo file.FileInfo, outputFileInfo file.FileInfo) (err error) {

	// in := make(chan int, 10)
	// quit := make(chan int)
	loopCount := 1
	xlsxFileXi := make([]string, 0)
	allFiles, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}

	for _, fi := range allFiles {
		match, _ := regexp.MatchString(`\.xlsx`, fi.Name())
		if match {
			xlsxFileXi = append(xlsxFileXi, fi.Name())
		}
	}

	// go percentViewer(len(xlsxFileXi), 80*time.Millisecond, in, quit)

	for _, xlsxName := range xlsxFileXi {
		// in <- loopCount
		var inputFile dcFile
		var outputFile file.WorksheetFile
		inputFile.SetFile(inputFileInfo.FilePath, xlsxName, inputFileInfo.SheetName)
		outputFile.SetFile(outputFileInfo.FilePath, "[DIFFERENCE]"+xlsxName, outputFileInfo.SheetName)
		err = inputFile.differenceCalculate(outputFile)
		if err != nil {
			fmt.Println("\r>[" + strconv.Itoa(loopCount) + "][Fail]\t\t" + xlsxName)
			logging.Error.Printf("%+v\n", err)
		} else {
			fmt.Println("\r>[" + strconv.Itoa(loopCount) + "][Success]\t" + xlsxName)
		}
		loopCount++
	}

	// quit <- 1
	// close(quit)
	// close(in)
	return nil
}

//calcDifferenceAllFile 將目錄內特定尺規評量成績評分表進行差分計算
func DifferenceCalculate_Command(inputFileInfo file.FileInfo, outputFileInfo file.FileInfo) (err error) {
	var inputFile dcFile
	var outputFile file.WorksheetFile

	inputFile.SetFile(inputFileInfo.FilePath, inputFileInfo.FileName, inputFileInfo.SheetName)
	outputFile.SetFile(outputFileInfo.FilePath, outputFileInfo.FileName, outputFileInfo.SheetName)
	// outputFile.SetFile(f.outputFilePath, f.outputFileName, f.outputSheetName)

	quit := make(chan int)
	defer close(quit)

	go cmdline.Spinner("Data file is loading...", 80*time.Millisecond, quit)
	err = inputFile.ReadRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading have completed")

	go cmdline.Spinner("Data is grouping by student...", 80*time.Millisecond, quit)
	err = inputFile.setJudge()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.fetchLongestLength()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.groupByStudent()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Grouping have completed")

	go cmdline.Spinner("Calculating difference...", 80*time.Millisecond, quit)
	err = inputFile.calculateDifference()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Calculating have completed")

	go cmdline.Spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = inputFile.transportToSlice()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.ExportDataToExcel(outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Export completed")
	return nil
}

func DifferenceCalculate(progChan chan int, inputFile file.WorksheetFile, outputFile file.WorksheetFile) {

	dcf := dcFile{
		WorksheetFile: inputFile,
	}

	err := dcf.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dcf.setJudge()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dcf.fetchLongestLength()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dcf.groupByStudent()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dcf.calculateDifference()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dcf.transportToSlice()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dcf.exportDataToExcel(outputFile)
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

func (dcf *dcFile) differenceCalculate(outputFile file.WorksheetFile) (err error) {
	err = dcf.ReadRawData()
	if err != nil {
		return err
	}
	err = dcf.setJudge()
	if err != nil {
		return err
	}
	err = dcf.fetchLongestLength()
	if err != nil {
		return err
	}
	err = dcf.groupByStudent()
	if err != nil {
		return err
	}
	err = dcf.calculateDifference()
	if err != nil {
		return err
	}
	err = dcf.transportToSlice()
	if err != nil {
		return err
	}
	err = dcf.exportDataToExcel(outputFile)
	if err != nil {
		return err
	}
	return nil
}

func (dcf *dcFile) groupByStudent() (err error) {

	//提取各評審成績，將其放入面向中，value傳入一行成績，shift代表要取得從第5格往右移多少格的資料，facingName可指定面向名稱
	getScore := func(value []string, shift int, facingName string) (f facing, err error) {

		scoreXi := make([]float64, 0)

		for i := 5; i < len(value)-1; i += 5 {
			convFloat, err := strconv.ParseFloat(value[i+shift], 32)
			if err != nil {
				return f, err
			}
			scoreXi = append(scoreXi, convFloat)
		}
		f = facing{
			facingName: facingName,
			score:      scoreXi,
		}
		return f, nil
	}

	dcf.gbsd = make(map[object.Student][]facing)

	startLoop := 4
	var s object.Student
	if len(dcf.DataRows) <= 0 {
		return fmt.Errorf("DataRows has no data")
	}
	for index, value := range dcf.DataRows {

		//遇到第四行以前，或委員A的面向一未評分，或"綜合評語/備註："行跳過
		if index < startLoop || len(value) < 6 || value[5] == "" || value[2] == "綜合評語/備註：" {
			continue
		}

		//在每位考生的第零列
		if value[0] != "" {
			s.StudentID = value[0]   //指定考生編號
			s.StudentName = value[1] //指定學生姓名

			//提取各面向平均欄位
			f, err := getScore(value, 4, "avg")
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
			dcf.gbsd[s] = append(dcf.gbsd[s], f)
		}

		//提取單一面向不同評審的成績
		f, err := getScore(value, 0, value[2])
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		dcf.gbsd[s] = append(dcf.gbsd[s], f)
	}
	return nil
}

func (dcf *dcFile) setJudge() (err error) {
	if len(dcf.DataRows) <= 0 {
		return fmt.Errorf("DataRows has no data")
	}
	for i := 5; i < len(dcf.DataRows[1])-1; i += 5 {
		dcf.judge = append(dcf.judge, dcf.DataRows[1][i])
	}
	dcf.judgeNum = (len(dcf.DataRows[1]) - 6) / 5 //計算評委數量
	return nil
}

func (dcf *dcFile) calculateDifference() (err error) {
	if len(dcf.gbsd) <= 0 {
		return fmt.Errorf("gbsd has no data")
	}
	for key, facXi := range dcf.gbsd {
		for index, fac := range facXi {
			for i, _ := range fac.score {
				for j := i + 1; j < len(fac.score); j++ {
					dcf.gbsd[key][index].difference = append(dcf.gbsd[key][index].difference, fac.score[i]-fac.score[j])
				}
			}
		}
	}
	return nil
}

func (dcf *dcFile) transportToSlice() (err error) {
	tempXi := make([]string, 0)

	//加入標題列
	tempXi = append(tempXi, "編號", "考生", "面向", "成績為零", "差分超過±10")
	if len(dcf.judge) <= 0 {
		return fmt.Errorf("var 'judge' has no data")
	}
	for i := 0; i < len(dcf.judge); i++ {
		for j := i + 1; j < len(dcf.judge); j++ {
			tempXi = append(tempXi, dcf.judge[i]+"-"+dcf.judge[j])
		}
	}
	dcf.NewDataRows = append(dcf.NewDataRows, tempXi)

	//加入差分結果
	if len(dcf.gbsd) <= 0 {
		return fmt.Errorf("gbsd has no data")
	}
	for key, facXi := range dcf.gbsd {
		for _, fac := range facXi {
			facingHasZero := false //紀錄該面向是否有成績為零的項目
			outOfRange := false    //紀錄該面向是否有差分超過±10的項目
			tempInfoXi := make([]string, 0)
			tempDiffXi := make([]string, 0)
			tempInfoXi = append(tempDiffXi, key.StudentID, key.StudentName, fac.facingName)

			//歷遍差分Slice
			for _, diff := range fac.difference {
				diff, err := strconv.ParseFloat(fmt.Sprintf("%.2f", diff), 64) //取小數點後兩位
				if err != nil {
					return err
				}
				convStr := strconv.FormatFloat(diff, 'f', -1, 64)
				if diff >= 10 || diff <= -10 {
					convStr = "[" + convStr + "]"
					outOfRange = true //該面向中有超過差分的成績
				}
				tempDiffXi = append(tempDiffXi, convStr)
			}

			//歷遍成績Slice
			if len(fac.score) <= 0 {
				return fmt.Errorf("score slice has no data")
			}
			for _, scr := range fac.score {
				if scr == 0 {
					facingHasZero = true //該面向中有成績為零的項目
				}
			}

			//若紀錄為true則標記Zero或OutOfRange，若無則填空
			if facingHasZero == true {
				tempInfoXi = append(tempInfoXi, "Zero")
			} else {
				tempInfoXi = append(tempInfoXi, "")
			}
			if outOfRange == true {
				tempInfoXi = append(tempInfoXi, "OutOfRange")
			} else {
				tempInfoXi = append(tempInfoXi, "")
			}

			//合併tempInfoXi和tempDiffXi
			for _, value := range tempDiffXi {
				tempInfoXi = append(tempInfoXi, value)
			}

			dcf.NewDataRows = append(dcf.NewDataRows, tempInfoXi)
		}
	}
	return nil
}

//exportDataToExcel 將檔案匯出成xlsx檔案
func (dcf *dcFile) exportDataToExcel(outputFile file.WorksheetFile) (err error) {
	xlsx := excelize.NewFile()
	// SheetName := "工作表"
	xlsx.SetSheetName("Sheet1", outputFile.SheetName)
	err = xlsx.SetColWidth(outputFile.SheetName, "E", "E", 12)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	//將f.newDataRows資料加入到xlsx內
	for index, value := range dcf.NewDataRows {
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

//fetchLongestLength 取得所有字串陣列中最長且不是空白的長度
func (dcf *dcFile) fetchLongestLength() (err error) {
	lastPos := 0

	if len(dcf.DataRows) <= 0 {
		return fmt.Errorf("DataRows has no data")
	}

	//尋找所有字串陣列中最長且不是空白的長度
	for _, value := range dcf.DataRows {
		for i := lastPos; i < len(value); i++ {
			if value[i] != "" && lastPos < i {
				lastPos = i
			}
		}
	}

	longestLen := lastPos + 1

	//將所有字串陣列長度變成與longestPos+1相同
	for index, value := range dcf.DataRows {
		if len(value) > longestLen {
			dcf.DataRows[index] = dcf.DataRows[index][0:longestLen]
		} else if len(value) < longestLen {
			//若比longestLen小則將長度補足
			for len(dcf.DataRows[index]) < longestLen {
				dcf.DataRows[index] = append(dcf.DataRows[index], "")
			}
		}
	}
	return nil
}
