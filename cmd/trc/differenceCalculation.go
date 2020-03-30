package main

import (
	"fmt"
	"strconv"

	"github.com/Luxurioust/excelize"
)

func (dcf *dcFile) DifferenceCalculate(outputFile file) (err error) {
	err = dcf.readRawData()
	if err != nil {
		return err
	}
	err = dcf.setJudge()
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

	dcf.gbsd = make(map[student][]facing)

	startLoop := 4
	var s student
	if len(dcf.dataRows) <= 0 {
		return fmt.Errorf("dataRows has no data")
	}
	for index, value := range dcf.dataRows {

		//遇到第四行以前，或委員A的面向一未評分，或"綜合評語/備註："行跳過
		if index < startLoop || value[5] == "" || value[2] == "綜合評語/備註：" {
			continue
		}

		//在每位考生的第零列
		if value[0] != "" {
			s.studentID = value[0]   //指定考生編號
			s.studentName = value[1] //指定學生姓名

			//提取各面向平均欄位
			f, err := getScore(value, 4, "avg")
			if err != nil {
				return err
			}
			dcf.gbsd[s] = append(dcf.gbsd[s], f)
		}

		//提取單一面向不同評審的成績
		f, err := getScore(value, 0, value[2])
		if err != nil {
			return err
		}
		dcf.gbsd[s] = append(dcf.gbsd[s], f)
	}
	return nil
}

func (dcf *dcFile) setJudge() (err error) {
	if len(dcf.dataRows) <= 0 {
		return fmt.Errorf("dataRows has no data")
	}
	for i := 5; i < len(dcf.dataRows[1])-1; i += 5 {
		dcf.judge = append(dcf.judge, dcf.dataRows[1][i])
	}
	dcf.judgeNum = (len(dcf.dataRows[1]) - 6) / 5 //計算評委數量
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
	dcf.newDataRows = append(dcf.newDataRows, tempXi)

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
			tempInfoXi = append(tempDiffXi, key.studentID, key.studentName, fac.facingName)

			//歷遍差分Slice
			for _, diff := range fac.difference {
				diff, err := strconv.ParseFloat(fmt.Sprintf("%.2f", diff), 64) //取小數點後兩位
				if err != nil {
					return err
				}
				convStr := strconv.FormatFloat(diff, 'f', -1, 64)
				if diff > 9 || diff < -9 {
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

			dcf.newDataRows = append(dcf.newDataRows, tempInfoXi)
		}
	}
	return nil
}

//exportDataToExcel 將檔案匯出成xlsx檔案
func (dcf *dcFile) exportDataToExcel(outputFile file) (err error) {
	xlsx := excelize.NewFile()
	// sheetName := "工作表"
	xlsx.SetSheetName("Sheet1", outputFile.sheetName)
	err = xlsx.SetColWidth(outputFile.sheetName, "E", "E", 12)
	if err != nil {
		return err
	}

	//將f.newDataRows資料加入到xlsx內
	for index, value := range dcf.newDataRows {
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
