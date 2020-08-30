package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func (f *flags) initFlag() (err error) {
	helpFlag := flag.Lookup("help")
	if helpFlag == nil {
		flag.BoolVar(&f.help, "help", false, "Usage")
	}

	f.startGuiFlagSet = flag.NewFlagSet("start gui", flag.ExitOnError)
	err = f.setStartGuiFlag()
	if err != nil {
		return nil
	}
	f.downloadVideoFlagSet = flag.NewFlagSet("download video", flag.ExitOnError)
	err = f.setDownloadVideoFlag()
	if err != nil {
		return nil
	}
	f.splitScoreAlertFlagSet = flag.NewFlagSet("split scoreAlert", flag.ExitOnError)
	err = f.setSplitScoreAlertFlag()
	if err != nil {
		return nil
	}
	f.mergeVideoFlagSet = flag.NewFlagSet("merge video", flag.ExitOnError)
	err = f.setMergeVideoFlag()
	if err != nil {
		return nil
	}
	f.mergeCourseFlagSet = flag.NewFlagSet("merge course", flag.ExitOnError)
	err = f.setMergeCourseFlag()
	if err != nil {
		return nil
	}
	f.calcDifferentFlagSet = flag.NewFlagSet("calculate different", flag.ExitOnError)
	err = f.setCalcDifferentFlag()
	if err != nil {
		return nil
	}
	return nil
}

func (f *flags) setStartGuiFlag() (err error) {

	f.startGuiFlagSet.BoolVar(&f.help, "help", false, "Usage")
	return nil
}

func (f *flags) setDownloadVideoFlag() (err error) {
	//set default video parameter
	defaultAcademicYear := ""
	defaultSemester := ""
	if time.Now().Month() <= 7 {
		defaultAcademicYear = strconv.Itoa(time.Now().Year() - 1911 - 1)
		defaultSemester = "2"
	} else {
		defaultAcademicYear = strconv.Itoa(time.Now().Year() - 1911)
		defaultSemester = "1"
	}

	path, err := os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	f.downloadSVInputFile = make([]string, 3)
	f.downloadSVOutputFile = make([]string, 3)
	f.downloadVideoFlagSet.StringVar(&f.academicYear, "year", defaultAcademicYear, "設定學年度，預設為當前學年度")
	f.downloadVideoFlagSet.StringVar(&f.semester, "semester", defaultSemester, "設定學期，預設為當前學期")
	f.downloadVideoFlagSet.StringVar(&f.youtubeAPIKey, "key", "", "設定YoutubeAPIKey")
	f.downloadVideoFlagSet.StringVar(&f.downloadSVInputFile[0], "inPath", path+"\\", "輸入檔案路徑")
	f.downloadVideoFlagSet.StringVar(&f.downloadSVInputFile[1], "inName", "數位課綱.xlsx", "輸入檔案名稱")
	f.downloadVideoFlagSet.StringVar(&f.downloadSVInputFile[2], "inSheet", "工作表", "輸入工作表名稱")
	f.downloadVideoFlagSet.StringVar(&f.downloadSVOutputFile[0], "outPath", path+"\\", "輸出檔案路徑")
	f.downloadVideoFlagSet.StringVar(&f.downloadSVOutputFile[1], "outName", "數位課綱.xlsx", "輸出檔案名稱")
	f.downloadVideoFlagSet.StringVar(&f.downloadSVOutputFile[2], "outSheet", "工作表", "輸出工作表名稱")
	f.downloadVideoFlagSet.BoolVar(&f.appendVideoInfo, "append", false, "在原有檔案內增加影片資訊")
	f.downloadVideoFlagSet.BoolVar(&f.help, "help", false, "Usage")
	return nil
}

func (f *flags) setSplitScoreAlertFlag() (err error) {
	path, err := os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	f.scoreAlertFile = make([]string, 3)
	f.teacherInfoFile = make([]string, 3)
	f.exportTemplateFile = make([]string, 3)
	f.splitScoreAlertFlagSet.StringVar(&f.scoreAlertFile[0], "masterPath", path+"\\", "設定預警總表檔案路徑")
	f.splitScoreAlertFlagSet.StringVar(&f.scoreAlertFile[1], "masterName", "預警總表.xlsx", "設定預警總表檔案名稱")
	f.splitScoreAlertFlagSet.StringVar(&f.scoreAlertFile[2], "masterSheet", "工作表", "設定預警總表工作表名稱")
	f.splitScoreAlertFlagSet.StringVar(&f.teacherInfoFile[0], "teacherPath", path+"\\", "設定教師名單檔案路徑")
	f.splitScoreAlertFlagSet.StringVar(&f.teacherInfoFile[1], "teacherName", "教師名單.xlsx", "設定教師名單檔案名稱")
	f.splitScoreAlertFlagSet.StringVar(&f.teacherInfoFile[2], "teacherSheet", "工作表", "設定教師名單工作表名稱")
	f.splitScoreAlertFlagSet.StringVar(&f.exportTemplateFile[0], "templatePath", path+"\\", "設定空白分表檔案名稱")
	f.splitScoreAlertFlagSet.StringVar(&f.exportTemplateFile[1], "templateName", "空白.xlsx", "設定空白分表檔案名稱")
	f.splitScoreAlertFlagSet.StringVar(&f.exportTemplateFile[2], "templateSheet", "工作表", "設定空白分表檔案名稱")
	f.splitScoreAlertFlagSet.BoolVar(&f.help, "help", false, "Usage")
	return nil
}

func (f *flags) setMergeVideoFlag() (err error) {
	path, err := os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	f.svInputFile = make([]string, 3)
	f.svOutputFile = make([]string, 3)
	f.svTeacherFile = make([]string, 3)
	f.mergeVideoFlagSet.StringVar(&f.svInputFile[0], "inPath", path+"\\", "輸入檔案路徑")
	f.mergeVideoFlagSet.StringVar(&f.svInputFile[1], "inName", "數位課綱.xlsx", "輸入檔案名稱")
	f.mergeVideoFlagSet.StringVar(&f.svInputFile[2], "inSheet", "工作表", "輸入工作表名稱")
	f.mergeVideoFlagSet.StringVar(&f.svOutputFile[0], "outPath", path+"\\", "輸出檔案路徑")
	f.mergeVideoFlagSet.StringVar(&f.svOutputFile[1], "outName", "[MERGENCE]數位課綱.xlsx", "輸出檔案名稱")
	f.mergeVideoFlagSet.StringVar(&f.svOutputFile[2], "outSheet", "工作表", "輸出工作表名稱")
	f.mergeVideoFlagSet.StringVar(&f.svTeacherFile[0], "tfPath", path+"\\", "設定教師名單檔案路徑")
	f.mergeVideoFlagSet.StringVar(&f.svTeacherFile[1], "tfName", "教師名單.xlsx", "設定教師名單檔案名稱")
	f.mergeVideoFlagSet.StringVar(&f.svTeacherFile[2], "tfSheet", "工作表", "設定教師名單工作表名稱")
	f.mergeVideoFlagSet.BoolVar(&f.tfile, "tfile", false, "使用教師名單合併")
	f.mergeVideoFlagSet.BoolVar(&f.help, "help", false, "Usage")
	return nil
}

func (f *flags) setMergeCourseFlag() (err error) {
	path, err := os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	f.cdInputFile = make([]string, 3)
	f.cdOutputFile = make([]string, 3)
	f.mergeCourseFlagSet.StringVar(&f.cdInputFile[0], "inPath", path+"\\", "輸入檔案路徑")
	f.mergeCourseFlagSet.StringVar(&f.cdInputFile[1], "inName", "開課總表.xlsx", "輸入檔案名稱")
	f.mergeCourseFlagSet.StringVar(&f.cdInputFile[2], "inSheet", "工作表", "輸入工作表名稱")
	f.mergeCourseFlagSet.StringVar(&f.cdOutputFile[0], "outPath", path+"\\", "輸出檔案路徑")
	f.mergeCourseFlagSet.StringVar(&f.cdOutputFile[1], "outName", "[MERGENCE]開課總表.xlsx", "輸出檔案名稱")
	f.mergeCourseFlagSet.StringVar(&f.cdOutputFile[2], "outSheet", "工作表", "輸出工作表名稱")
	f.mergeCourseFlagSet.BoolVar(&f.help, "help", false, "Usage")
	return nil
}

func (f *flags) setCalcDifferentFlag() (err error) {
	path, err := os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	f.calcDifferentFlagSet.StringVar(&f.inputFilePath, "inPath", path+"\\", "輸入檔案路徑")
	f.calcDifferentFlagSet.StringVar(&f.inputFileName, "inName", "評分表.xlsx", "輸入檔案名稱")
	f.calcDifferentFlagSet.StringVar(&f.inputSheetName, "inSheet", "學系彙整版", "輸入工作表名稱")
	f.calcDifferentFlagSet.StringVar(&f.outputFilePath, "outPath", path+"\\", "輸出檔案路徑")
	f.calcDifferentFlagSet.StringVar(&f.outputFileName, "outName", "成績差分.xlsx", "輸出檔案名稱")
	f.calcDifferentFlagSet.StringVar(&f.outputSheetName, "outSheet", "工作表", "輸出工作表名稱")
	f.calcDifferentFlagSet.BoolVar(&f.readAllFilesInDir, "A", false, "讀取目錄內所有檔案")
	f.calcDifferentFlagSet.BoolVar(&f.help, "help", false, "Usage")
	return nil
}
