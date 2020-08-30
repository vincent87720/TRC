package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runCalculateDifferenceDialog(owner walk.Form, fi *CalculateDifferenceInfo, iconFilePath string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var db *walk.DataBinder

	var inputDirPath *walk.LineEdit

	var inputFilePath *walk.LineEdit

	var inputSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo:   &dlg,
		Title:      "計算成績差分",
		Icon:       iconFilePath,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:       Font{Family: "Microsoft JhengHei", PointSize: 9},
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "calculateDifferenceInfo",
			DataSource:     fi,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{Width: 300, Height: 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					RadioButtonGroup{
						DataMember: "CalcAll",
						Buttons: []RadioButton{
							RadioButton{
								Name:  "single",
								Text:  "單一檔案",
								Value: false,
							},
							RadioButton{
								Name:  "multiple",
								Text:  "多重檔案",
								Value: true,
							},
						},
					},
					HSpacer{},
				},
			},

			Composite{
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{
						Text:    "文件路徑",
						Font:    labelFont,
						Visible: Bind("multiple.Checked"),
					},
					LineEdit{
						AssignTo: &inputDirPath,
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
						Text:     Bind("Folder"),
						Visible:  Bind("multiple.Checked"),
					},
					PushButton{
						Text: "選擇資料夾",
						OnClicked: func() {
							onBrowseFolderButtonClicked(owner, inputDirPath)
						},
						Visible: Bind("multiple.Checked"),
					},
				},
				Visible: Bind("multiple.Checked"),
			},
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text:    "評量尺規成績表",
						Font:    labelFont,
						Visible: Bind("single.Checked"),
					},
					LineEdit{
						AssignTo: &inputFilePath,
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
						Text:     Bind("InputPath"),
						Visible:  Bind("single.Checked"),
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, inputFilePath, inputSheetSelector)
						},
						Visible: Bind("single.Checked"),
					},
					ComboBox{
						AssignTo:      &inputSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("InputSheet"),
						Visible:       Bind("single.Checked"),
					},
				},
				Visible: Bind("single.Checked"),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "確定",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
