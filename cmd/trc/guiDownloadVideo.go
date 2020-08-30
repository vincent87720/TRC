package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runDownloadVideoDialog(owner walk.Form, fi *DownloadVideoInfo, iconFilePath string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var db *walk.DataBinder

	var year *walk.LineEdit
	var semester *walk.LineEdit
	var youtubeAPIKey *walk.LineEdit
	var inputFilePath *walk.LineEdit
	var inputSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo:   &dlg,
		Title:      "下載數位課綱資料",
		Icon:       iconFilePath,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:       Font{Family: "Microsoft JhengHei", PointSize: 9},
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "downloadSyllabusVideo",
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
						DataMember: "Append",
						Buttons: []RadioButton{
							RadioButton{
								Name:  "new",
								Text:  "下載並建立新檔案",
								Value: false,
							},
							RadioButton{
								Name:  "append",
								Text:  "下載並合併檔案",
								Value: true,
							},
						},
					},
					HSpacer{},
				},
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "學年",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &year,
						MinSize:  Size{Width: 250, Height: 0},
						Text:     Bind("Year"),
					},

					Label{
						Text: "學期",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &semester,
						MinSize:  Size{Width: 250, Height: 0},
						Text:     Bind("Semester"),
					},

					Label{
						Text: "YoutubeAPI Key",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &youtubeAPIKey,
						MinSize:  Size{Width: 250, Height: 0},
						Text:     Bind("Key"),
					},
				},
			},
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text:    "數位課綱",
						Font:    labelFont,
						Visible: Bind("append.Checked"),
					},
					LineEdit{
						AssignTo: &inputFilePath,
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
						Text:     Bind("InputPath"),
						Visible:  Bind("append.Checked"),
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, inputFilePath, inputSheetSelector)
						},
						Visible: Bind("append.Checked"),
					},
					ComboBox{
						AssignTo:      &inputSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("InputSheet"),
						Visible:       Bind("append.Checked"),
					},
				},
				Visible: Bind("append.Checked"),
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
								Error.Printf("%+v\n", err)
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
