package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runDownloadTeacherDialog(owner walk.Form, iconFilePath string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var masterFilePath *walk.LineEdit
	var teacherFilePath *walk.LineEdit
	var templateFilePath *walk.LineEdit

	var masterSheetSelector *walk.ComboBox
	var teacherSheetSelector *walk.ComboBox
	var templateSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo:      &dlg,
		Title:         "DownloadTeacherList",
		Icon:          iconFilePath,
		Background:    SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:          Font{Family: "Microsoft JhengHei", PointSize: 9},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{Width: 300, Height: 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text: "預警總表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &masterFilePath,
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, masterFilePath, masterSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &masterSheetSelector,
						Editable:      false,
						BindingMember: "key",
						DisplayMember: "Name",
					},

					Label{
						Text: "教師名單",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &teacherFilePath,
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, teacherFilePath, teacherSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &teacherSheetSelector,
						Editable:      false,
						BindingMember: "key",
						DisplayMember: "Name",
					},

					Label{
						Text: "空白分表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &templateFilePath,
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, templateFilePath, templateSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &templateSheetSelector,
						Editable:      false,
						BindingMember: "key",
						DisplayMember: "Name",
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
							fmt.Println(masterFilePath.Text())
							fmt.Println(teacherFilePath.Text())
							fmt.Println(templateFilePath.Text())

							fmt.Println(masterSheetSelector.Text())
							fmt.Println(teacherSheetSelector.Text())
							fmt.Println(templateSheetSelector.Text())

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
