package SyllabusVideo

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
	teacherdata "github.com/vincent87720/TRC/internal/TeacherData"
	"github.com/vincent87720/TRC/internal/cmdline"
	"github.com/vincent87720/TRC/internal/file"
	"github.com/vincent87720/TRC/internal/logging"
	"github.com/vincent87720/TRC/internal/object"
)

type syllabusVideo struct {
	object.Course
	object.Teacher
	timeNPlace      string
	videoID         string
	videoTitle      string
	videoURL        string //數位課綱連結
	videoDuration   string //數位課綱影片長度(string)
	videoSeconds    int    //數位課綱影片長度(second)
	problemOfCourse string //影片問題
}

//svFile 數位課綱檔案
type svFile struct {
	file.WorksheetFile
	cdpCol       int                                //開課單位編號欄位(departmentID)
	cidCol       int                                //課程編號欄位(courseID)
	csnCol       int                                //課程名稱欄位(courseName)
	pocCol       int                                //影片問題欄位(problemOfCourse)
	cthCol       int                                //任課老師欄位(courseTeacher)
	gbtd         map[object.Teacher][]syllabusVideo //group by object.Teacher data
	mergedXi     [][]string                         //合併後的陣列
	maxCourseNum int                                //所有老師中，擁有最多科目的數量
}

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
// 			logging.Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}

//mergeSyllabusVideo 合併數位課綱紀錄表內可合併的內容
func MergeSyllabusVideoData_Command(inputFileInfo file.FileInfo, outputFileInfo file.FileInfo, teacherFileInfo file.FileInfo, tfile bool) (err error) {
	var inputFile svFile
	var outputFile file.WorksheetFile

	inputFile.SetFile(inputFileInfo.FilePath, inputFileInfo.FileName, inputFileInfo.SheetName)
	outputFile.SetFile(outputFileInfo.FilePath, outputFileInfo.FileName, outputFileInfo.SheetName)

	quit := make(chan int)
	defer close(quit)

	go cmdline.Spinner("Syllabus video file is loading...", 80*time.Millisecond, quit)
	err = inputFile.ReadRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading file have completed")

	go cmdline.Spinner("Data is splitting...", 80*time.Millisecond, quit)
	err = inputFile.groupByTeacher()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Splitting have completed")

	if tfile {
		//若使用教師名單檔案匯入所屬單位
		var teacherFile teacherdata.ThFile
		teacherFile.SetFile(teacherFileInfo.FilePath, teacherFileInfo.FileName, teacherFileInfo.SheetName)
		go cmdline.Spinner("Teacher file is loading...", 80*time.Millisecond, quit)
		err = teacherFile.ReadRawData()
		if err != nil {
			quit <- 1
			return err
		}
		err = teacherFile.GroupByTeacher()
		if err != nil {
			quit <- 1
			return err
		}
		quit <- 1
		fmt.Println("\r> Loading teacher file have completed")

		go cmdline.Spinner("Teacher info is matching...", 80*time.Millisecond, quit)
		err = inputFile.matchTeacherInfoFile(teacherFile)
		if err != nil {
			quit <- 1
			return err
		}
		quit <- 1
		fmt.Println("\r> Matching teacher info have completed")
	} else {
		//使用inputFile檔案內的所屬單位合併
		go cmdline.Spinner("Teacher info is matching...", 80*time.Millisecond, quit)
		err = inputFile.matchTeacherInfo()
		if err != nil {
			quit <- 1
			return err
		}
		quit <- 1
		fmt.Println("\r> Matching teacher info have completed")
	}

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
	fmt.Println("\r> Exporting have completed")
	return nil
}

func MergeSyllabusVideoData(progChan chan int, inputFile file.WorksheetFile, outputFile file.WorksheetFile) {

	svf := svFile{
		WorksheetFile: inputFile,
	}

	err := svf.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svf.groupByTeacher()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 2
	err = svf.matchTeacherInfo()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 3
	err = svf.transportToSlice()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 4
	err = svf.exportDataToExcel(outputFile)
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 5

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
// 			logging.Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func MergeSyllabusVideoDataByList(progChan chan int, inputFile file.WorksheetFile, outputFile file.WorksheetFile, teacherFile file.WorksheetFile) {
	svf := svFile{
		WorksheetFile: inputFile,
	}

	thf := teacherdata.ThFile{
		WorksheetFile: teacherFile,
	}

	err := svf.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svf.groupByTeacher()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = thf.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = thf.GroupByTeacher()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svf.matchTeacherInfoFile(thf)
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svf.transportToSlice()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svf.exportDataToExcel(outputFile)
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

//groupByTeacher 依照教師名稱將數位課綱資料分群
func (svf *svFile) groupByTeacher() (err error) {
	err = svf.FindCol("教師姓名", &svf.cthCol)
	if err != nil {
		return err
	}
	err = svf.FindCol("所屬單位", &svf.cdpCol)
	if err != nil {
		return err
	}
	err = svf.FindCol("科目序號", &svf.cidCol)
	if err != nil {
		return err
	}
	err = svf.FindCol("科目名稱", &svf.csnCol)
	if err != nil {
		return err
	}
	err = svf.FindCol("影片問題", &svf.pocCol)
	if err != nil {
		return err
	}

	svf.gbtd = make(map[object.Teacher][]syllabusVideo)
	if len(svf.DataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("DataRows has no data"))
		return err
	}
	for index, value := range svf.DataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		//尋找影片問題欄位有資料者
		if len(value) >= len(svf.FirstRow) && value[svf.pocCol] != "" {
			t := object.Teacher{
				TeacherName: value[svf.cthCol],
			}
			c := syllabusVideo{
				Course: object.Course{
					CourseID:   value[svf.cidCol],
					CourseName: value[svf.csnCol],
					Department: object.Department{
						DepartmentName: value[svf.cdpCol],
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
		t := object.Teacher{
			TeacherName: key.TeacherName,
			Department: object.Department{
				DepartmentName: value[0].Course.Department.DepartmentName,
			},
		}
		delete(svf.gbtd, key) //必須先刪除再加入，否則有可能誤刪
		svf.gbtd[t] = value
	}
	return nil
}

//matchTeacherInfoFile 使用額外輸入的teacherFile檔案進行教師比對合併
func (svf *svFile) matchTeacherInfoFile(teacherFile teacherdata.ThFile) (err error) {
	for key, value := range svf.gbtd {
		if len(teacherFile.TeacherMap[key.TeacherName]) > 0 {

			t := object.Teacher{
				TeacherName: key.TeacherName,
				Department: object.Department{
					DepartmentName: teacherFile.TeacherMap[key.TeacherName][0].Department.DepartmentName,
				},
			}
			delete(svf.gbtd, key) //必須先刪除再加入，否則有可能誤刪
			svf.gbtd[t] = value
		}
	}
	return nil
}

//transportToSlice 將map[object.Teacher][]syllabusVideo的資料轉換為二維陣列
func (svf *svFile) transportToSlice() (err error) {
	if len(svf.DataRows) <= 0 {
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
			tempXi = append(tempXi, key.TeacherName, key.Department.DepartmentName)

			//每列只能放9個，多的給下一圈執行
			for j := 0; j < 9; j++ {
				//放到最後一個為止
				if i*9+j >= len(value) {
					break
				}
				tempXi = append(tempXi, value[i*9+j].CourseName, value[i*9+j].CourseID, value[i*9+j].problemOfCourse)
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
func (svf *svFile) exportDataToExcel(outputFile file.WorksheetFile) (err error) {
	xlsx := excelize.NewFile()
	// SheetName := "工作表"
	xlsx.SetSheetName("Sheet1", outputFile.SheetName)

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
			err := xlsx.SetColStyle(outputFile.SheetName, mark, fillColorE9E7D6)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 2:
			err := xlsx.SetColStyle(outputFile.SheetName, mark, fillColorE0E4D6)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		case 0:
			err := xlsx.SetColStyle(outputFile.SheetName, mark, fillColorDADCD2)
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
	err = xlsx.SetSheetRow(outputFile.SheetName, "A1", &title)
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
		keys = append(keys, k.TeacherName)
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
