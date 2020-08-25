package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"regexp"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

type NormalFileInfo struct {
	Path  string
	Sheet string
}
type SplitFileInfo struct {
	MasterPath    string
	MasterSheet   string
	TeacherPath   string
	TeacherSheet  string
	TemplatePath  string
	TemplateSheet string
}

func RunMainWindow(){

	font := Font{Family: "Microsoft JhengHei", PointSize: 12}

	var mw *walk.MainWindow

	r, err := regexp.Compile(`(.*\\)([^\\]*.xlsx)`)
	if err != nil {
		log.Print(err)
	}

	if _, err := (MainWindow{
		AssignTo:   &mw,
		Title:      "TRC",
		Icon:       "./assets/guiImage/blockchain-blueblue.png",
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Size:       Size{1000, 200},
		MinSize:    Size{1000, 200},
		Layout:     HBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					//download
					PushButton{
						Text:    "Download Teacher",
						Image:   "./assets/guiImage/Those_Icons-download-32.png",
						Font:    font,
						MinSize: Size{200, 50},
						OnClicked: func() {
							if cmd, err := RunDownloadTeacherDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {

							}
						},
					},

					//download
					PushButton{
						Text:    "Download Video   ",
						Image:   "./assets/guiImage/Those_Icons-download-32.png",
						Font:    font,
						MinSize: Size{200, 50},
						OnClicked: func() {
							if cmd, err := RunDownloadVideoDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {

							}
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					//split
					PushButton{
						Text:    "Split ScoreAlert",
						Image:   "./assets/guiImage/Those_Icons-split-32.png",
						Font:    font,
						MinSize: Size{200, 50},
						OnClicked: func() {
							fi := new(SplitFileInfo)
							if cmd, err := RunSplitScoreAlertDialog(mw, fi); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								var masterFile saFile
								var templateFile file
								var teacherFile thFile
								var outputFile file
								masterPathXi := r.FindStringSubmatch(fi.MasterPath)
								teacherPathXi := r.FindStringSubmatch(fi.TeacherPath)
								templatePathXi := r.FindStringSubmatch(fi.TemplatePath)
								masterFile.setFile(masterPathXi[1], masterPathXi[2], fi.MasterSheet)
								templateFile.setFile(templatePathXi[1], templatePathXi[2], fi.TemplateSheet)
								teacherFile.setFile(teacherPathXi[1], teacherPathXi[2], fi.TeacherSheet)
								go SplitScoreAlertData(masterFile, templateFile, teacherFile, outputFile)
							}
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					//calculate
					PushButton{
						Text:    "Calculate Difference",
						Image:   "./assets/guiImage/Pixel_Perfect-calculate-32.png",
						Font:    font,
						MinSize: Size{200, 50},
						OnClicked: func() {
							if cmd, err := RunCalculateDifferenceDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {

							}
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					//merge
					PushButton{
						Text:    "Merge Course",
						Image:   "./assets/guiImage/Those_Icons-merge-32.png",
						Font:    font,
						MinSize: Size{200, 50},
						OnClicked: func() {
							if cmd, err := RunMergeCourseDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {

							}
						},
					},

					//merge
					PushButton{
						Text:    "Merge Video  ",
						Image:   "./assets/guiImage/Those_Icons-merge-32.png",
						Font:    font,
						MinSize: Size{200, 50},
						OnClicked: func() {
							if cmd, err := RunMergeVideoDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {

							}
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func ExportAssets() {

	fileXi, err := AssetDir("")
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range fileXi {
		xi, err := Asset(value)
		if err != nil {
			fmt.Println(err)
		}
		// convert []byte to image for saving to file
		img, _, _ := image.Decode(bytes.NewReader(xi))

		//save the imgByte to file
		out, err := os.Create("./assets/" + value)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = png.Encode(out, img)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func OnOpenFileButtonClicked(owner walk.Form, filePath *walk.LineEdit, selector *walk.ComboBox) {
	dlg := new(walk.FileDialog)
	dlg.Title = "Open File"
	dlg.Filter = "Excel檔案 (*.xlsx)|*.xlsx|所有檔案 (*.*)|*.*"

	if ok, err := dlg.ShowOpen(owner); err != nil {
		fmt.Fprintln(os.Stderr, "Error:Open master file")
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "Cancel file selection")
		return
	}

	filePath.SetText(dlg.FilePath)

	keys := []*DropDownItem{}

	xlsx, err := excelize.OpenFile(dlg.FilePath)
	if err != nil {
		fmt.Println(err)
	}
	xlsxSht := xlsx.GetSheetMap()
	for idx, val := range xlsxSht {
		keys = append(keys, &DropDownItem{Key: idx, Name: val})
	}
	selector.SetModel(keys)
}
