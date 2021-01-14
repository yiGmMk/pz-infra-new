package excelUtil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gyf841010/pz-infra-new/commonUtil"

	"github.com/tealeg/xlsx"
)

type ExcelContent struct {
	Header  []string
	Content [][]interface{}
}

func GenerateExcel(content *ExcelContent, f string) {
	dir, _ := filepath.Split(f)
	os.Mkdir(dir, 0777)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet1")

	row = sheet.AddRow()
	for _, header := range content.Header {
		cell = row.AddCell()
		cell.Value = header
	}

	for _, rowInContent := range content.Content {
		row = sheet.AddRow()
		for _, v := range rowInContent {
			cell = row.AddCell()
			switch v.(type) {
			case int:
				cell.SetInt(v.(int))
			case float64:
				cell.SetFloat(v.(float64))
			case string:
				value := v.(string)
				if strings.HasSuffix(value, "%") {
					floatValue := commonUtil.FloatValue(value[:len(value)-1])
					floatValue = floatValue / 100
					cell.SetFloatWithFormat(floatValue, "0.00%")
				} else {
					cell.SetString(v.(string))
				}
			}
		}
	}

	err = file.Save(f)
	if err != nil {
		fmt.Printf(err.Error())
	}

}
