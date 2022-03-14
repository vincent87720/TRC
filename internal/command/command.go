package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	difference "github.com/vincent87720/TRC/internal/Difference"
	rapidprint "github.com/vincent87720/TRC/internal/RapidPrint"
	scorealert "github.com/vincent87720/TRC/internal/ScoreAlert"
	syllabusvideo "github.com/vincent87720/TRC/internal/SyllabusVideo"
	teacherdata "github.com/vincent87720/TRC/internal/TeacherData"
	"github.com/vincent87720/TRC/internal/file"
	gui "github.com/vincent87720/TRC/internal/gui"
	"github.com/vincent87720/TRC/internal/logging"
)

type command struct {
	commandString string
	commandAction string
	flagSet       *flag.FlagSet
}

type CommandSet struct {
	lyr1commandSet map[string][]command //第一層指令集
	lyr2commandSet map[string][]command //第二層指令集
}

var (
	INITPATH string
)

//constructer
func (cmdSet *CommandSet) CommandInit(f *Flags, path string) {
	cmdSet.lyr1commandSet = make(map[string][]command)
	cmdSet.lyr2commandSet = make(map[string][]command)

	INITPATH = path

	cmdSet.setLyr1Command("start", "啟動服務(啟動圖形化介面)")
	cmdSet.setLyr2Command("start", "gui", "啟動圖形化介面", f.startGuiFlagSet)
	cmdSet.setLyr1Command("download", "下載檔案(下載教師資料、數位課綱連結)")
	cmdSet.setLyr2Command("download", "video", "下載數位課綱影片連結", f.downloadVideoFlagSet)
	cmdSet.setLyr2Command("download", "object.Teacher", "下載教師資料", nil)
	cmdSet.setLyr1Command("split", "分割檔案(分割預警總表)")
	cmdSet.setLyr2Command("split", "scoreAlert", "分割預警總表", f.splitScoreAlertFlagSet)
	cmdSet.setLyr1Command("merge", "合併檔案(合併數位課綱資料、開課及製版數登記表)")
	cmdSet.setLyr2Command("merge", "video", "依教師合併數位課綱影片問題", f.mergeVideoFlagSet)
	cmdSet.setLyr2Command("merge", "object.Course", "依教師合併開課總表，製作為製版數登記表", f.mergeCourseFlagSet)
	cmdSet.setLyr1Command("calculate", "計算檔案(計算成績差分)")
	cmdSet.setLyr2Command("calculate", "difference", "計算成績差分", f.calcDifferentFlagSet)
}

//setLyr1Command 設定第一層的指令，commandString為指令字串，commandAction為指令說明
func (cmdSet *CommandSet) setLyr1Command(commandString string, commandAction string) {
	cmd := command{
		commandString: commandString,
		commandAction: commandAction,
	}
	cmdSet.lyr1commandSet[commandString] = append(cmdSet.lyr1commandSet[commandString], cmd)
}

//setLyr2Command 設定第二層的指令，layer1command為第一層指令，也就是此指令為哪一個父指令的子指令，commandString為指令字串，commandAction為指令說明
func (cmdSet *CommandSet) setLyr2Command(layer1command string, commandString string, commandAction string, flagSet *flag.FlagSet) {
	cmd := command{
		commandString: commandString,
		commandAction: commandAction,
		flagSet:       flagSet,
	}
	cmdSet.lyr2commandSet[layer1command] = append(cmdSet.lyr2commandSet[layer1command], cmd)
}

//layer1CommandUsage 顯示第一層指令的使用說明
func (cmdSet *CommandSet) layer1CommandUsage() {
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
func (cmdSet *CommandSet) layer2CommandUsage(layer1command string) {
	cmds, exist := cmdSet.lyr2commandSet[layer1command]
	if exist {
		fmt.Printf("Usage: trc %s <command> [<args>]", layer1command)
		fmt.Println()
		fmt.Println("Commands:")
		for _, cmd := range cmds {
			fmt.Printf("  %s [<args>]\t%s\n", cmd.commandString, cmd.commandAction)
			if cmd.flagSet != nil {
				for i := 0; i < cmd.flagSet.NFlag(); i++ {
					fmt.Println(cmd.flagSet.Arg(i))
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
func AnalyseCommand(cmdSet *CommandSet, f *Flags) {
	if len(os.Args) <= 1 || f.help {
		cmdSet.layer1CommandUsage()
	} else {
		switch os.Args[1] {
		case "start":
			if len(os.Args) <= 2 || f.help {
				cmdSet.layer2CommandUsage(os.Args[1])
			} else {
				switch os.Args[2] {
				case "gui":
					err := f.startGuiFlagSet.Parse(os.Args[3:])
					if err != nil {
						logging.Error.Printf("%+v\n", err)
					}

					if f.help == true {
						f.startGuiFlagSet.Usage()
						return
					}
					gui.StartWindow(INITPATH)
				default:
					cmdSet.layer2CommandUsage(os.Args[1])
				}
			}
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

					if f.youtubeAPIKey == "" {
						err = errors.WithStack(fmt.Errorf("youtubeAPIKey missing"))
						logging.Error.Printf("%+v\n", err)
						fmt.Println("\rError: 沒有指定youtubeAPIKey，使用-key指定參數")
						fmt.Println("\r> Downloading have failed")
					} else {
						fmt.Println(">> GetSyllabusVideoLink")
						if f.appendVideoInfo {
							inputFileInfo := file.FileInfo{
								FilePath:  f.downloadSVInputFile[0],
								FileName:  f.downloadSVInputFile[1],
								SheetName: f.downloadSVInputFile[2],
							}
							outputFileInfo := file.FileInfo{
								FilePath:  f.downloadSVOutputFile[0],
								FileName:  f.downloadSVOutputFile[1],
								SheetName: f.downloadSVOutputFile[2],
							}
							err = syllabusvideo.AppendSyllabusVideo_Command(inputFileInfo, outputFileInfo, f.academicYear, f.semester, f.youtubeAPIKey)
							if err != nil {
								logging.Error.Printf("%+v\n", err)
								fmt.Println("\r> Downloading have failed")
							}
						} else {
							outputFileInfo := file.FileInfo{
								FilePath:  f.downloadSVOutputFile[0],
								FileName:  f.downloadSVOutputFile[1],
								SheetName: f.downloadSVOutputFile[2],
							}
							err = syllabusvideo.GetSyllabusVideo_Command(outputFileInfo, f.youtubeAPIKey, f.academicYear, f.semester)
							if err != nil {
								logging.Error.Printf("%+v\n", err)
								fmt.Println("\r> Downloading have failed")
							}
						}
					}

				case "teacher":
					fmt.Println(">> GetTeacherData")
					err := teacherdata.GetTeacher_Command()
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

					scoreAlertFileInfo := file.FileInfo{
						FilePath:  f.scoreAlertFile[0],
						FileName:  f.scoreAlertFile[1],
						SheetName: f.scoreAlertFile[2],
					}
					templateFileInfo := file.FileInfo{
						FilePath:  f.exportTemplateFile[0],
						FileName:  f.exportTemplateFile[1],
						SheetName: f.exportTemplateFile[2],
					}
					teacherFileInfo := file.FileInfo{
						FilePath:  f.teacherInfoFile[0],
						FileName:  f.teacherInfoFile[1],
						SheetName: f.teacherInfoFile[2],
					}
					outputFileInfo := file.FileInfo{
						FilePath:  f.exportTemplateFile[0],
						FileName:  "",
						SheetName: "工作表",
					}

					err = scorealert.SplitScoreAlert_Command(scoreAlertFileInfo, templateFileInfo, teacherFileInfo, outputFileInfo)
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
					inputFileInfo := file.FileInfo{
						FilePath:  f.svInputFile[0],
						FileName:  f.svInputFile[1],
						SheetName: f.svInputFile[2],
					}
					outputFileInfo := file.FileInfo{
						FilePath:  f.svOutputFile[0],
						FileName:  f.svOutputFile[1],
						SheetName: f.svOutputFile[2],
					}
					teacherFileInfo := file.FileInfo{
						FilePath:  f.svTeacherFile[0],
						FileName:  f.svTeacherFile[1],
						SheetName: f.svTeacherFile[2],
					}
					err = syllabusvideo.MergeSyllabusVideoData_Command(inputFileInfo, outputFileInfo, teacherFileInfo, f.tfile)
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
					inputFileInfo := file.FileInfo{
						FilePath:  f.cdInputFile[0],
						FileName:  f.cdInputFile[1],
						SheetName: f.cdInputFile[2],
					}
					outputFileInfo := file.FileInfo{
						FilePath:  f.cdOutputFile[0],
						FileName:  f.cdOutputFile[1],
						SheetName: f.cdOutputFile[2],
					}
					err = rapidprint.MergeRapidPrintData_Command(inputFileInfo, outputFileInfo)
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
						inputFileInfo := file.FileInfo{
							FilePath:  f.inputFilePath,
							FileName:  "",
							SheetName: f.inputSheetName,
						}
						outputFileInfo := file.FileInfo{
							FilePath:  f.outputFilePath,
							FileName:  "",
							SheetName: f.outputSheetName,
						}
						err := difference.DifferenceCalculate_AllFile_Command(inputFileInfo, outputFileInfo)
						if err != nil {
							logging.Error.Printf("%+v\n", err)
						}
					} else {
						inputFileInfo := file.FileInfo{
							FilePath:  f.inputFilePath,
							FileName:  f.inputFileName,
							SheetName: f.inputSheetName,
						}
						outputFileInfo := file.FileInfo{
							FilePath:  f.outputFilePath,
							FileName:  "[DIFFERENCE]" + f.inputFileName,
							SheetName: f.outputSheetName,
						}
						err := difference.DifferenceCalculate_Command(inputFileInfo, outputFileInfo)
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
