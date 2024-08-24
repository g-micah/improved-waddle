package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/xuri/excelize/v2"
)

func main() {
	address := ""
	username := ""
	password := ""
	path := ""
	// Start
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("---------------------")
	fmt.Println("Start of Micah's FTP App to get files and sizes")
	fmt.Println("---------------------")
	// Get Server Addr
	if address == "" {
		fmt.Println("Enter full address, including port, of ftp server. Ex: ftp.example.org:21")
		fmt.Print("-> ")
		scanner.Scan()
		address = scanner.Text()
		fmt.Println("Connecting to address via ftp: " + address)
	}
	// Attempt connection
	c, err := ftp.Dial(address, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		errorTime(err)
	}
	// Get Server Username
	if username == "" {
		fmt.Println("Enter username:")
		fmt.Print("-> ")
		scanner.Scan()
		username = scanner.Text()
		username = strings.Replace(username, "\n", "", -1)
	}
	// Get Server Password
	if password == "" {
		fmt.Println("Enter password:")
		fmt.Print("-> ")
		scanner.Scan()
		password = scanner.Text()
		fmt.Println("Logging in with user [" + username + "] and password [" + password + "]")
	}
	// Login
	err = c.Login(username, password)
	if err != nil {
		errorTime(err)
	}
	retry := true
	var files []*ftp.Entry
	// Loop for retry
	for retry {
		// Get full folder path
		if path == "" {
			fmt.Println("Enter full filepath to run program in. ex: /folder1/folder2/folder3 OR ~/folder2/folder")
			fmt.Print("-> ")
			scanner.Scan()
			path = scanner.Text()
		}
		// Get all files in specific folder
		err = c.ChangeDir(path)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Try again with a new directory...")
			retry = true
			path = ""
			continue
		}
		files, err = c.List(path)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Try again with a new directory...")
			retry = true
			path = ""
			continue
		}
		// Print to screen/text file / excel file
		fmt.Println("Reading files found in " + path)
		fmt.Print("Files found: ")
		fmt.Println(len(files))
		retry = len(files) == 0
		if retry {
			fmt.Println("No files found in directory... Try again with a new directory...")
			path = ""
		}
	}
	// Sort array
	sort.Slice(files, func(i int, j int) bool {
		return int(files[i].Size) > int(files[j].Size)
	})
	outputString := "File Name | File Size (KB)\n"
	for _, element := range files {
		if element.Type.String() == "folder" {
			element.Size = uint64(1)
		}
		outputString += element.Name + " | " + fmt.Sprintf("%.2f", convertByteToKb(int64(element.Size))) + "\n"
	}
	// Create .txt
	err = os.WriteFile("output.txt", []byte(outputString), 0644)
	if err != nil {
		errorTime(err)
	}
	fmt.Println("Generated file output.txt")
	// Create Excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "Filename")
	f.SetCellValue("Sheet1", "B1", "File Size (KB)")
	for index, element := range files {
		column1, _ := excelize.CoordinatesToCellName(1, index+2)
		column2, _ := excelize.CoordinatesToCellName(2, index+2)
		f.SetCellValue("Sheet1", column1, element.Name)
		f.SetCellValue("Sheet1", column2, convertByteToKb(int64(element.Size)))
	}
	// Save spreadsheet by the given path.
	err = f.SaveAs("output.xlsx")
	if err != nil {
		errorTime(err)
	}
	fmt.Println("Generated file output.xlsx")
	if err := c.Quit(); err != nil {
		errorTime(err)
	}
	fmt.Println("---------------------")
	fmt.Print(" == ENTER to END ==")
	scanner.Scan()
	scanner.Text()
}

func errorTime(err error) {
	fmt.Println("---------------------")
	fmt.Println("Error")
	fmt.Println(err.Error())
	fmt.Println("---------------------")
	fmt.Print(" == ENTER to END ==")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scanner.Text()
	log.Fatal(err)
}

func convertByteToKb(arg int64) float64 {
	return math.Round(float64(arg)/10.24) / 100
}
