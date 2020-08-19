package main

import (
	"testing"

	"github.com/Luxurioust/excelize"
	"github.com/stretchr/testify/assert"
)

func TestSetFile(t *testing.T) {
	assert := assert.New(t)
	testcase := []struct {
		filePath    string     //檔案路徑
		fileName    string     //檔案名稱
		sheetName   string     //工作表名稱
		firstRow    []string   //第一行
		dataRows    [][]string //資料
		newDataRows [][]string //新資料
		xlsx        *excelize.File
	}{
		{
			filePath:  "./",
			fileName:  "test.xlsx",
			sheetName: "sheet1",
		},
		{
			filePath:  "./",
			fileName:  "中文檔名測試.xlsx",
			sheetName: "sheet1",
		},
		{
			filePath:  "./",
			fileName:  "space test.xlsx",
			sheetName: "sheet1",
		},
	}
	for _, cases := range testcase {
		var f file
		f.setFile(cases.filePath, cases.fileName, cases.sheetName)
		assert.Equal(cases.filePath, f.filePath)
		assert.Equal(cases.fileName, f.fileName)
		assert.Equal(cases.sheetName, f.sheetName)
	}
}

func TestFillSliceLength(t *testing.T) {
	assert := assert.New(t)
	testcase := []struct {
		dataRows [][]string
		length   int
	}{
		{
			dataRows: [][]string{{"a", "b", "c", "d", "e"}, {"a", "b", "c", ""}, {"a", "b", "c", "d", "e"}},
			length:   5,
		}, {
			dataRows: [][]string{{"a", "b"}, {"a", "b", "c", "d", ""}, {"a", "b", "c", "d", "e"}},
			length:   5,
		}, {
			dataRows: [][]string{{"a", "b"}, {"a", "b"}, {"a", "b", "c"}},
			length:   5,
		},
	}
	for _, cases := range testcase {
		var f file
		f.dataRows = cases.dataRows
		f.fillSliceLength(cases.length)
		for _, ary := range cases.dataRows {
			assert.Equal(cases.length, len(ary))
		}
	}
}
