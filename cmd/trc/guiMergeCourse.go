package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func RunMergeCourseDialog(owner walk.Form, fi *NormalFileInfo, iconFilePath string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var db *walk.DataBinder

	var inputFilePath *walk.LineEdit

	var inputSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo:   &dlg,
		Title:      "MergeCourseData",
		Icon:       iconFilePath,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:       Font{Family: "Microsoft JhengHei", PointSize: 9},
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "mergeCourseInfo",
			DataSource:     fi,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{300, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text: "開課總表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &inputFilePath,
						Text:     Bind("InputPath"),
						MinSize:  Size{250, 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							OnOpenFileButtonClicked(owner, inputFilePath, inputSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &inputSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("InputSheet"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
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
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
