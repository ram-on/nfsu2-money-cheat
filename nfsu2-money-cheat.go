package main

//
// Cheat that edits the money stored within the Need For Speed Underground 2 save file.
//

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// NFSU2_MONEY_FILE_OFFSET:  offset to where the money value is stored within the save file
const NFSU2_MONEY_FILE_OFFSET int64 = 0xA16A

//
// PrintUsage:  Print usage.
//
func PrintUsage() {
	programName := os.Args[0]
	fmt.Printf("Usage:  %s <command> <save_file>\n", programName)
	fmt.Println()
	fmt.Println("Command:")
	fmt.Println("    info         Display how much money is stored within <save_file>")
	fmt.Println("    edit         Edit the money within <save_file>.  The original <save_file> will be backed-up.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("    Edit the money available to 25000 within the given save file:")
	fmt.Printf("        %s edit -m 25000 -f 'C:\\Users\\MyUserName\\bin\\Need For Speed Underground 2\\save\\NFS Underground 2\\MySaveFile'\n", programName)
	fmt.Println()
	fmt.Println("    Display the money within the given save file:")
	fmt.Printf("        %s info -f 'C:\\Users\\MyUserName\\bin\\Need For Speed Underground 2\\save\\NFS Underground 2\\MySaveFile'\n", programName)
}

//
// printCurrentMoney:  Print the currently stored money within the provided save file.
//
func printCurrentMoney(saveFile string) {
	// open the save file
	file, err := os.OpenFile(saveFile, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open save file '%s'\n%s\n", saveFile, err)
		os.Exit(1)
	}
	defer file.Close()

	// moneyArr will hold the signed 32-bit integer (which represents the money store within the file)
	moneyArr := make([]byte, 4)

	// read the money stored in the save file
	_, err = file.ReadAt(moneyArr, NFSU2_MONEY_FILE_OFFSET)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read the money currently stored within the save file '%s'\n%s\n", saveFile, err)
		os.Exit(1)
	}

	// convert the money to int32 (as money in NFSU2 is stored as int32 - little endian)
	moneyStored := int32(binary.LittleEndian.Uint32(moneyArr))
	fmt.Printf("Money stored within save file = %d\n", moneyStored)
}

//
// backupSaveFile:  Backups up the provided save file.
//
func backupSaveFile(origSaveFile string) {
	dateTime := time.Now().Format(time.RFC3339)
	backupFilePath := origSaveFile + "." + strings.ReplaceAll(dateTime, ":", ".") + ".backup"

	// check if the save file is actually a file
	origFileStat, err := os.Stat(origSaveFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open save file '%s'\n%s\n", origSaveFile, err)
		os.Exit(1)
	}

	if !origFileStat.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "The provided file is not a regular file '%s'\n", origSaveFile)
		os.Exit(1)
	}

	// open the original (provided) save file
	orig, err := os.Open(origSaveFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while opening save file '%s'\n%s\n", origSaveFile, err)
		os.Exit(1)
	}
	defer orig.Close()

	// create the backup save file
	backup, err := os.Create(backupFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while creating backup save file '%s'\n%s\n", origSaveFile, err)
		os.Exit(1)
	}
	defer backup.Close()

	// copy the original into the backup
	if _, err := io.Copy(backup, orig); err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while backuping the save file '%s'\n%s\n", origSaveFile, err)
		os.Exit(1)
	}

	fmt.Println("Save file backed up at:  ", backupFilePath)
}

//
// editMoneySaveFile:  Edit the money being stored within the provided save file.
//
func editMoneySaveFile(money uint32, saveFile string) {
	// open the save file
	file, err := os.OpenFile(saveFile, os.O_RDWR, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open save file '%s'\n%s\n", saveFile, err)
		os.Exit(1)
	}
	defer file.Close()

	// convert the money into an array of a 32 bit integer (little endian)
	moneyArr := make([]byte, 4)
	binary.LittleEndian.PutUint32(moneyArr, money)

	// write the new money value within the save file
	_, err = file.WriteAt(moneyArr, NFSU2_MONEY_FILE_OFFSET)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write money into save file '%s'\n%s\n", saveFile, err)
		os.Exit(1)
	}

	fmt.Printf("Edit save file with money = %d.\n", money)
}

func main() {
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)
	infoSaveFilePathVal := infoCmd.String("f", "", "NFSU2 save file path")

	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	moneyVal := editCmd.Int("m", 0, "Money that will be saved in the provided save file.  Minumum amount is 0.  Maximum amount is 2147483647.")
	editSaveFilePathVal := editCmd.String("f", "", "NFSU2 save file path")

	if len(os.Args) == 1 {
		PrintUsage()
	} else if os.Args[1] == "info" {
		infoCmd.Parse(os.Args[2:])

		if *infoSaveFilePathVal == "" {
			fmt.Fprintf(os.Stderr, "Save file not provided\n")
			infoCmd.Usage()
			os.Exit(1)
		}

		printCurrentMoney(*infoSaveFilePathVal)
	} else if os.Args[1] == "edit" {
		editCmd.Parse(os.Args[2:])

		if *moneyVal < 0 {
			fmt.Fprintf(os.Stderr, "Money cannot be < 0\n")
			os.Exit(1)
		}

		if *editSaveFilePathVal == "" {
			fmt.Fprintf(os.Stderr, "Save file not provided\n")
			editCmd.Usage()
			os.Exit(1)
		}

		backupSaveFile(*editSaveFilePathVal)
		editMoneySaveFile(uint32(*moneyVal), *editSaveFilePathVal)
	}
}
