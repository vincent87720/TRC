package ScoreAlert

import (
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
	teacherdata "github.com/vincent87720/TRC/internal/TeacherData"
	cmdline "github.com/vincent87720/TRC/internal/cmdline"
	"github.com/vincent87720/TRC/internal/file"
	"github.com/vincent87720/TRC/internal/logging"
	"github.com/vincent87720/TRC/internal/object"
)

type scoreAllert struct {
	object.Course
	object.Student
	allertReason string //預警原因
	tutorMethod  string //輔導方式
}

//saFile 成績預警檔案
type saFile struct {
	file.WorksheetFile
	csdCol int                              //開課單位(courseDepartment)
	cidCol int                              //科目序號(courseID)
	csnCol int                              //課程名稱(courseName)
	cthCol int                              //任課老師(courseTeacher)
	sidCol int                              //學生學號(studentID)
	stnCol int                              //學生姓名(studentName)
	alrCol int                              //預警原由(alertReason)
	gbtd   map[object.Teacher][]scoreAllert //group by object.Teacher data
}

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
// 			logging.Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func SplitScoreAlert(progChan chan int, inputFile file.WorksheetFile, templateFile file.WorksheetFile, teacherFile file.WorksheetFile, outputFile file.WorksheetFile) {
	saf := saFile{
		WorksheetFile: inputFile,
	}
	thf := teacherdata.ThFile{
		WorksheetFile: teacherFile,
	}

	err := saf.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = saf.groupByTeacher()
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
	err = templateFile.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = saf.exportDataToExcel(templateFile, thf, outputFile)
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

//splitScoreAlert 將成績預警總表切分為各個小分表
func SplitScoreAlert_Command(masterFileInfo file.FileInfo, templateFileInfo file.FileInfo, teacherFileInfo file.FileInfo, outputFileInfo file.FileInfo) (err error) {
	var masterFile saFile
	var templateFile file.WorksheetFile
	var teacherFile teacherdata.ThFile
	var outputFile file.WorksheetFile

	masterFile.SetFile(masterFileInfo.FilePath, masterFile.FileName, masterFile.SheetName)
	templateFile.SetFile(templateFile.FilePath, teacherFile.FileName, teacherFile.SheetName)
	teacherFile.SetFile(teacherFile.FilePath, teacherFile.FileName, teacherFile.SheetName)
	outputFile.SetFile(outputFile.FilePath, outputFile.FileName, outputFile.SheetName)

	quit := make(chan int)
	defer close(quit)

	go cmdline.Spinner("Master file is loading...", 80*time.Millisecond, quit)
	err = masterFile.ReadRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading master file have completed")

	go cmdline.Spinner("Teacher file is loading...", 80*time.Millisecond, quit)
	err = teacherFile.ReadRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading teacher file have completed")

	go cmdline.Spinner("Template file is loading...", 80*time.Millisecond, quit)
	err = templateFile.ReadRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading template file have completed")

	go cmdline.Spinner("Data is splitting...", 80*time.Millisecond, quit)
	err = masterFile.groupByTeacher()
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
	fmt.Println("\r> Splitting have completed")

	go cmdline.Spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = masterFile.exportDataToExcel(templateFile, teacherFile, outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Exporting have completed")
	return nil

}

//groupByTeacher 依照教師名稱將預警總表資料分群
func (saf *saFile) groupByTeacher() (err error) {
	err = saf.FindCol("開課學系", &saf.csdCol)
	if err != nil {
		return err
	}
	saf.FindCol("科目序號", &saf.cidCol)
	saf.FindCol("預警科目", &saf.csnCol)
	saf.FindCol("授課教師", &saf.cthCol)
	saf.FindCol("學號", &saf.sidCol)
	saf.FindCol("學生姓名", &saf.stnCol)
	saf.FindCol("預警原由", &saf.alrCol)

	saf.gbtd = make(map[object.Teacher][]scoreAllert)
	if len(saf.DataRows[0]) <= 0 {
		err = errors.WithStack(fmt.Errorf("dataRows has no data"))
		return err
	}
	for index, value := range saf.DataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		if value[saf.cthCol] != "" {
			t := object.Teacher{
				TeacherName: value[saf.cthCol],
			}
			sa := scoreAllert{
				Course: object.Course{
					Department: object.Department{
						DepartmentName: value[saf.csdCol],
					},
					CourseID:   value[saf.cidCol],
					CourseName: value[saf.csnCol],
				},
				Student: object.Student{
					StudentID:   value[saf.sidCol],
					StudentName: value[saf.stnCol],
				},
				allertReason: value[saf.alrCol],
			}
			saf.gbtd[t] = append(saf.gbtd[t], sa)
		}
	}
	return nil
}

//exportDataToExcel 匯出預警分表
func (saf *saFile) exportDataToExcel(templateFile file.WorksheetFile, teacherFile teacherdata.ThFile, outputFile file.WorksheetFile) (err error) {

	if len(saf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for key, value := range saf.gbtd {
		xlsx, err := excelize.OpenFile(templateFile.FilePath + templateFile.FileName)
		if err != nil {
			err = errors.WithStack(err)
			return err
		}
		if len(teacherFile.TeacherMap[key.TeacherName]) > 0 {
			//教師存在於名單中，設定檔名為"學院編號(int)_系所名稱(string)_教師姓名(string).xlsx"
			clgID := teacherFile.TeacherMap[key.TeacherName][0].CollegeID
			depName := teacherFile.TeacherMap[key.TeacherName][0].DepartmentName
			thrName := key.TeacherName
			outputFile.FileName = clgID + "_" + depName + "_" + thrName + ".xlsx"
		} else {
			//教師不存在於名單中，設定檔名為"教師姓名(string).xlsx"
			outputFile.FileName = key.TeacherName + ".xlsx"
		}

		for index, val := range value {
			row := strconv.Itoa(index + 2)
			position := "A" + row
			err = xlsx.SetSheetRow(templateFile.SheetName, position, &[]interface{}{val.Course.Department.DepartmentName, val.CourseID, val.CourseName, key.TeacherName, val.StudentID, val.StudentName, val.allertReason})
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
	}
	return nil
}
