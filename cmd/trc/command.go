package main

import (
	"TRC/pkg/logging"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/fatih/color"
)

//constructer
func (cmdSet *commandSet) commandSet() {
	cmdSet.lyr1commandSet = make(map[string][]command)
	cmdSet.lyr2commandSet = make(map[string][]command)
}

//setLyr1Command 設定第一層的指令，commandString為指令字串，commandAction為指令說明
func (cmdSet *commandSet) setLyr1Command(commandString string, commandAction string) {
	cmd := command{
		commandString: commandString,
		commandAction: commandAction,
	}
	cmdSet.lyr1commandSet[commandString] = append(cmdSet.lyr1commandSet[commandString], cmd)
}

//setLyr2Command 設定第二層的指令，layer1command為第一層指令，也就是此指令為哪一個父指令的子指令，commandString為指令字串，commandAction為指令說明
func (cmdSet *commandSet) setLyr2Command(layer1command string, commandString string, commandAction string, flagSet *flag.FlagSet) {
	cmd := command{
		commandString: commandString,
		commandAction: commandAction,
		flagSet:       flagSet,
	}
	cmdSet.lyr2commandSet[layer1command] = append(cmdSet.lyr2commandSet[layer1command], cmd)
}

//layer1CommandUsage 顯示第一層指令的使用說明
func (cmdSet *commandSet) layer1CommandUsage() {
	fmt.Println("Usage: trc <command> [<args>]")
	fmt.Println()
	fmt.Println("Commands:")
	for _, cmdXi := range cmdSet.lyr1commandSet {
		for _, cmd := range cmdXi {
			fmt.Printf("  %s [<args>]\t        %s\n", cmd.commandString, cmd.commandAction)
		}
	}
	fmt.Println()
	color.Set(color.FgHiYellow)
	fmt.Println("Warning: 參數不可包含空白字元")
	color.Set(color.FgWhite)
}

//layer1CommandUsage 顯示第二層指令的使用說明
func (cmdSet *commandSet) layer2CommandUsage(layer1command string) {
	cmds, exist := cmdSet.lyr2commandSet[layer1command]
	if exist {
		fmt.Printf("Usage: trc %s <command> [<args>]", layer1command)
		fmt.Println()
		fmt.Println("Commands:")
		for _, cmd := range cmds {
			fmt.Printf("  %s [<args>]\t%s\n", cmd.commandString, cmd.commandAction)
			if cmd.flagSet != nil {
				for key, value := range cmd.flagSet.GetFlags() {
					fmt.Println("        -" + key + " " + value.Usage)
				}
			}
		}
		fmt.Println()
		color.Set(color.FgHiYellow)
		fmt.Println("Warning: 參數不可包含空白字元")
		color.Set(color.FgWhite)
	}
}

//analyseCommand 分析命令並觸發相對應的方法
func analyseCommand(cmdSet *commandSet, f *flags) {
	if len(os.Args) <= 1 || f.help {
		cmdSet.layer1CommandUsage()
	} else {
		switch os.Args[1] {
		case "download":
			if len(os.Args) <= 2 || f.help {
				cmdSet.layer2CommandUsage(os.Args[1])
			} else {
				switch os.Args[2] {
				case "video":
					err := f.downloadVideoFlagSet.Parse(os.Args[3:])
					if err != nil {
						logging.Error.Printf("%+v\n", err)
					}

					if f.help == true {
						f.downloadVideoFlagSet.Usage()
						return
					}

					fmt.Println(">> GetSyllabusVideoLink")
					getSyllabusVideo()
				case "teacher":
					fmt.Println(">> GetTeacherData")
					err := getTeacher()
					if err != nil {
						logging.Error.Printf("%+v\n", err)
						fmt.Println("\r> Downloading have failed")
					}
				default:
					cmdSet.layer2CommandUsage(os.Args[1])
				}
			}

		case "split":

			if len(os.Args) <= 2 || f.help {
				cmdSet.layer2CommandUsage(os.Args[1])
			} else {
				switch os.Args[2] {
				case "scoreAlert":
					err := f.splitScoreAlertFlagSet.Parse(os.Args[3:])
					if err != nil {
						logging.Error.Printf("%+v\n", err)
					}
					if f.help {
						f.splitScoreAlertFlagSet.Usage()
						return
					}

					fmt.Println(">> SplitScoreAlertFile")
					err = splitScoreAlertFile(f)
					if err != nil {
						logging.Error.Printf("%+v\n", err)
						fmt.Println("\r> Splitting have failed")
					}

				default:
					cmdSet.layer2CommandUsage(os.Args[1])
				}
			}

		case "merge":

			if len(os.Args) <= 2 || f.help {
				cmdSet.layer2CommandUsage(os.Args[1])
			} else {
				switch os.Args[2] {
				case "video":
					err := f.mergeVideoFlagSet.Parse(os.Args[3:])
					if err != nil {
						logging.Error.Printf("%+v\n", err)
					}
					if f.help {
						f.mergeVideoFlagSet.Usage()
						return
					}

					fmt.Println(">> MergeSyllabusVideoLink")
					err = mergeSyllabusVideo(f)
					if err != nil {
						logging.Error.Printf("%+v\n", err)
						fmt.Println("\r> Merging have failed")
					}
				case "course":
					err := f.mergeCourseFlagSet.Parse(os.Args[3:])
					if err != nil {
						logging.Error.Printf("%+v\n", err)
					}
					if f.help {
						f.mergeCourseFlagSet.Usage()
						return
					}

					fmt.Println(">> MergeCourseData")
					err = mergeCourseData(f)
					if err != nil {
						logging.Error.Printf("%+v\n", err)
						fmt.Println("\r> Merging have failed")
					}
				default:
					cmdSet.layer2CommandUsage(os.Args[1])
				}
			}

		case "calculate":

			if len(os.Args) <= 2 || f.help {
				cmdSet.layer2CommandUsage(os.Args[1])
			} else {
				switch os.Args[2] {
				case "difference":
					err := f.calcDifferentFlagSet.Parse(os.Args[3:])

					if err != nil {
						logging.Error.Printf("%+v\n", err)
					}
					if f.help {
						f.calcDifferentFlagSet.Usage()
						return
					}

					fmt.Println(">> CalculateDifference")

					if f.readAllFilesInDir {
						err := calcDifferenceAllFile(f)
						if err != nil {
							logging.Error.Printf("%+v\n", err)
							fmt.Println("\r> Calculating have failed")
						}
					} else {
						err := calcDifference(f)
						if err != nil {
							logging.Error.Printf("%+v\n", err)
							fmt.Println("\r> Calculating have failed")
						}
					}

				default:
					cmdSet.layer2CommandUsage(os.Args[1])
				}
			}
		default:
			fmt.Println()
			color.Set(color.FgHiYellow)
			fmt.Printf("Warning: %q is not a valid command.\n", os.Args[1])
			color.Set(color.FgWhite)
			cmdSet.layer1CommandUsage()
		}
	}

}

//顯示loading圈圈，status傳入要顯示在圈圈後的文字，delay傳入更新時間間隔，ch傳入結束訊號
func spinner(status string, delay time.Duration, ch chan int) {
	for {
		select {
		case <-ch:
			fmt.Printf("\r                                              ")
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c %s", r, status)
				time.Sleep(delay)
			}
		}
	}
}

//顯示進度(幾分之幾)，denominator傳入總數量(分母)，delay傳入更新時間間隔，in傳入目前進度數量(分子)，quit傳入結束訊號
func percentViewer(denominator int, delay time.Duration, in chan int, quit chan int) {
	var numerator int
	for {

		select {
		case <-quit:
			fmt.Printf("\r                                              ")
			return
		case numerator = <-in:
			fmt.Printf("\r>%d/%d", numerator, denominator)
		default:
			fmt.Printf("\r>%d/%d", numerator, denominator)
		}
	}
}

//getSyllabusVideo 從學校網站抓取數位課綱影片連結
func getSyllabusVideo() (err error) {
	var svlreq getSVLRequest
	svlreq.newRequest()
	svlreq.setURL("http://syl.dyu.edu.tw/sl_cour_time.php?itimestamp=" + string(int32(time.Now().Unix())))
	err = svlreq.setURLValues("108", "1", "'1','2','3','4','5','6','7'", "'1','2','3','4','N','5','6','7','8','9','A','B','C','D','E'")
	if err != nil {
		return err
	}
	err = svlreq.sendRequest()
	if err != nil {
		return err
	}
	return nil
}

//getTeacher 從學校網站抓取教師資料
func getTeacher() (err error) {
	var tdreq getTDRequest
	var udreq getUDRequest
	var cudreq getCUDRequest
	var inputFile downloadTeacherFile
	var outputFile file
	outputFile.setFile("./", "教師名單.xlsx", "工作表")

	quit := make(chan int)
	defer close(quit)

	go spinner("Unit data is downloading...", 80*time.Millisecond, quit)
	udreq.setURL("https://lg.dyu.edu.tw/get_unit_title.php")
	err = udreq.sendRequest()
	if err != nil {
		quit <- 1
		return err
	}
	err = udreq.parseData()
	if err != nil {
		quit <- 1
		return err
	}
	cudreq.setURL("https://lg.dyu.edu.tw/APP/getFile.php?URL=http://lg.dyu.edu.tw/search_unit.php")
	err = cudreq.sendRequest()
	if err != nil {
		quit <- 1
		return err
	}
	err = cudreq.parseData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Downloading unit data have completed")

	go spinner("Teacher data is downloading...", 80*time.Millisecond, quit)
	tdreq.setURL("https://lg.dyu.edu.tw/search_teacher.php")
	err = tdreq.sendRequest()
	if err != nil {
		quit <- 1
		return err
	}
	err = tdreq.parseData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Downloading teacher data have completed")

	go spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = inputFile.transportToSlice(&tdreq, &udreq, &cudreq)
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.exportDataToExcel(outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Exporting have completed")

	return nil
}

//splitScoreAlertFile 將成績預警總表切分為各個小分表
func splitScoreAlertFile(f *flags) (err error) {
	var masterFile saFile
	var templateFile file
	var teacherFile thFile
	var outputFile file
	masterFile.setFile(f.scoreAlertFile[0], f.scoreAlertFile[1], f.scoreAlertFile[2])
	templateFile.setFile(f.exportTemplateFile[0], f.exportTemplateFile[1], f.exportTemplateFile[2])
	teacherFile.setFile(f.teacherInfoFile[0], f.teacherInfoFile[1], f.teacherInfoFile[2])
	outputFile.setFile(f.exportTemplateFile[0], "", "工作表")

	quit := make(chan int)
	defer close(quit)

	go spinner("Master file is loading...", 80*time.Millisecond, quit)
	err = masterFile.readRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading master file have completed")

	go spinner("Teacher file is loading...", 80*time.Millisecond, quit)
	err = teacherFile.readRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading teacher file have completed")

	go spinner("Template file is loading...", 80*time.Millisecond, quit)
	err = templateFile.readRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading template file have completed")

	go spinner("Data is splitting...", 80*time.Millisecond, quit)
	err = masterFile.groupByTeacher()
	if err != nil {
		quit <- 1
		return err
	}
	err = teacherFile.groupByTeacher()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Splitting have completed")

	go spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = masterFile.exportDataToExcel(templateFile, teacherFile, outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Exporting have completed")
	return nil

}

//mergeSyllabusVideo 合併數位課綱紀錄表內可合併的內容
func mergeSyllabusVideo(f *flags) (err error) {
	var inputFile svFile
	var outputFile file

	inputFile.setFile(f.svInputFile[0], f.svInputFile[1], f.svInputFile[2])
	outputFile.setFile(f.svOutputFile[0], f.svOutputFile[1], f.svOutputFile[2])

	quit := make(chan int)
	defer close(quit)

	go spinner("Syllabus video file is loading...", 80*time.Millisecond, quit)
	err = inputFile.readRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading file have completed")

	go spinner("Data is splitting...", 80*time.Millisecond, quit)
	err = inputFile.groupByTeacher()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Splitting have completed")

	go spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = inputFile.transportToSlice()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.exportDataToExcel(outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Exporting have completed")
	return nil
}

//mergeCourseData 合併開課總表內可合併的課程
func mergeCourseData(f *flags) (err error) {

	var inputFile rpFile
	var outputFile file
	inputFile.setFile(f.cdInputFile[0], f.cdInputFile[1], f.cdInputFile[2])
	outputFile.setFile(f.cdOutputFile[0], f.cdOutputFile[1], f.cdOutputFile[2])

	quit := make(chan int)
	defer close(quit)

	go spinner("Course data file is loading...", 80*time.Millisecond, quit)
	err = inputFile.readRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading file have completed")

	go spinner("Data is preprocessing...", 80*time.Millisecond, quit)
	err = inputFile.fillSliceLength(15)
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.findColumn()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Preprocessing have completed")

	go spinner("Data is grouping by teacher...", 80*time.Millisecond, quit)
	err = inputFile.groupByTeacher()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Grouping have completed")

	go spinner("Data is merging...", 80*time.Millisecond, quit)
	err = inputFile.mergeData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Merging have completed")

	go spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = inputFile.transportToSlice()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.exportDataToExcel(outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Exporting have completed")
	return nil
}

//calcDifferenceAllFile 將目錄內特定尺規評量成績評分表進行差分計算
func calcDifference(f *flags) (err error) {
	var inputFile dcFile
	var outputFile file

	inputFile.setFile(f.inputFilePath, f.inputFileName, f.inputSheetName)
	outputFile.setFile(f.outputFilePath, "[DIFFERENCE]"+f.inputFileName, f.outputSheetName)
	// outputFile.setFile(f.outputFilePath, f.outputFileName, f.outputSheetName)

	quit := make(chan int)
	defer close(quit)

	go spinner("Data file is loading...", 80*time.Millisecond, quit)
	err = inputFile.readRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading have completed")

	go spinner("Data is grouping by student...", 80*time.Millisecond, quit)
	err = inputFile.setJudge()
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

	go spinner("Calculating difference...", 80*time.Millisecond, quit)
	err = inputFile.calculateDifference()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Calculating have completed")

	go spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = inputFile.transportToSlice()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.exportDataToExcel(outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Export completed")
	return nil
}

//calcDifferenceAllFile 將目錄內所有尺規評量成績評分表進行差分計算
func calcDifferenceAllFile(f *flags) (err error) {

	in := make(chan int, 10)
	quit := make(chan int)
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

	go percentViewer(len(xlsxFileXi), 80*time.Millisecond, in, quit)

	for _, xlsxName := range xlsxFileXi {
		in <- loopCount
		var inputFile dcFile
		var outputFile file
		inputFile.setFile(f.inputFilePath, xlsxName, f.inputSheetName)
		outputFile.setFile(f.outputFilePath, "[DIFFERENCE]"+xlsxName, f.outputSheetName)
		err = inputFile.DifferenceCalculate(outputFile)
		if err != nil {
			return err
		}
		loopCount++
	}

	quit <- 1
	fmt.Println("\r> Exporting have completed")
	close(quit)
	close(in)
	return nil
}
