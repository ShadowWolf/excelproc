package main

import (
	"encoding/csv"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
)

import "github.com/gin-gonic/gin"

func downloadFile(url, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}(file)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(b io.ReadCloser) {
		err := b.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		var details struct {
			SourceFile                  string `json:"sourceFile" binding:"required"`
			DestinationFile             string `json:"destinationFile" binding:"required"`
			SourceConnectionString      string `json:"sourceConnectionString" binding:"required"`
			DestinationConnectionString string `json:"destinationConnectionString"`
		}

		if c.Bind(&details) == nil {
			client, err := azblob.NewClientFromConnectionString(details.SourceConnectionString, nil)
		}
	})

	file, err := excelize.OpenFile("C:\\Users\\shado\\Downloads\\testoutput.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := file.GetRows("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}

	outFile, err := os.Create("C:\\Users\\shado\\Downloads\\testoutput-sheet2.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	csvWriter := csv.NewWriter(outFile)

	for _, row := range rows {
		err = csvWriter.Write(row)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
