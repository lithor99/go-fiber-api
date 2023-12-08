package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/responses"
)

// upload single file
func UploadSingleFile(c *fiber.Ctx) error {
	// Get first file from form field "document":
	file, err := c.FormFile("file")

	// Check for errors:
	if err == nil {
		replacer := strings.NewReplacer("/", "", "\\", "", " ", "")
		newFileName := replacer.Replace(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./uploads/single_files/%s", newFileName))
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "file": fmt.Sprintf("/uploads/single_files/%s", newFileName)})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "file": file.Filename})
}

// upload multi file
func UploadMultiFile(c *fiber.Ctx) error {
	// Get first file from form field "document":
	files, err := c.MultipartForm()

	if err == nil {
		var filenames []string
		for _, fileHeaders := range files.File {
			for _, fileHeader := range fileHeaders {
				replacer := strings.NewReplacer("/", "", "\\", "", " ", "")
				newFileName := replacer.Replace(fileHeader.Filename)
				filenames = append(filenames, "upload/multi_files/"+newFileName)
				c.SaveFile(fileHeader, fmt.Sprintf("./uploads/multi_files/%s", newFileName))
			}
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "files": filenames})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error"})
}

type ExcelResults struct {
	Index  string `json:"index"`
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
	Value4 string `json:"value4"`
}

// upload excel file
func UploadExcelFile(c *fiber.Ctx) error {
	// Get first file from form field "document":
	file, err := c.FormFile("file")
	// Check for errors:
	if err == nil {
		gen := strconv.Itoa(int(time.Now().UnixMilli()))
		replacer := strings.NewReplacer("/", "", "\\", "", " ", "")
		newFileName := replacer.Replace(file.Filename)

		ext1 := filepath.Ext(newFileName)
		if ext1 != ".xlsx" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "File must be excel file"})
		}
		c.SaveFile(file, fmt.Sprintf("./uploads/excel_files/%s.%s", gen, newFileName))
		excel, err := excelize.OpenFile(fmt.Sprintf("./uploads/excel_files/%s.%s", gen, newFileName))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}
		var excel_results []ExcelResults
		excel_result := new(ExcelResults)
		i := 1
		for {
			var startASCIINum int = 65
			for x := 0; x < 5; x++ {
				cell := string(rune(startASCIINum+x)) + strconv.Itoa(i)
				value := excel.GetCellValue("Sheet1", cell)
				if x == 0 {
					excel_result.Index = value
				}
				if x == 1 {
					excel_result.Value1 = value
				}
				if x == 2 {
					excel_result.Value2 = value
				}
				if x == 3 {
					excel_result.Value3 = value
				}
				if x == 4 {
					excel_result.Value4 = value
				}
			}

			newExcelResult := ExcelResults{
				Index:  excel_result.Index,
				Value1: excel_result.Value1,
				Value2: excel_result.Value2,
				Value3: excel_result.Value3,
				Value4: excel_result.Value4,
			}
			if newExcelResult == (ExcelResults{}) {
				break
			}
			i++
			excel_results = append(excel_results, newExcelResult)
		}
		return c.Status(http.StatusOK).JSON(responses.ExcelData{Status: "success", Total: len(excel_results), Value: &fiber.Map{"data": excel_results}})

	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "file": file.Filename})
}
