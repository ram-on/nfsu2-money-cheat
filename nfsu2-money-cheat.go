package main

//
// Cheat that edits the money stored within the Need For Speed Underground 2 save file.
//

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// NFSU2_MONEY_FILE_OFFSET:  offset to where the money value is stored within the save file
	NFSU2_MONEY_FILE_OFFSET int64 = 0xA16A

	// program version
	VERSION = "2.0"
)

// set to true if the program launched the TUI (text user interface); false otherwise
var using_tui bool = false

// PrintUsage:  Print usage.
func PrintUsage() {
	programName := os.Args[0]
	fmt.Printf("Usage:  %s <command> <save_file>\n", programName)
	fmt.Println()
	fmt.Println("Command:")
	fmt.Println("    info         Display how much money is stored within <save_file>")
	fmt.Println("    edit         Edit the money within <save_file>.  The original <save_file> will be backed-up.")
	fmt.Println("    help         Print this help.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("    Edit the money available to 25000 within the given save file:")
	fmt.Printf("        %s edit -m 25000 -f 'C:\\Users\\MyUserName\\bin\\Need For Speed Underground 2\\save\\NFS Underground 2\\MySaveFile'\n", programName)
	fmt.Println()
	fmt.Println("    Display the money within the given save file:")
	fmt.Printf("        %s info -f 'C:\\Users\\MyUserName\\bin\\Need For Speed Underground 2\\save\\NFS Underground 2\\MySaveFile'\n", programName)
}

// PrintCurrentMoney:  Print the currently stored money within the provided save file.
func PrintCurrentMoney(saveFile string) {
	// open the save file
	file, err := os.OpenFile(saveFile, os.O_RDONLY, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open save file '%s': %s\n", saveFile, err)
		Exit(1)
	}
	defer file.Close()

	// moneyArr will hold the signed 32-bit integer (which represents the money store within the file)
	moneyArr := make([]byte, 4)

	// read the money stored in the save file
	_, err = file.ReadAt(moneyArr, NFSU2_MONEY_FILE_OFFSET)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read the money currently stored within the save file '%s': %s\n", saveFile, err)
		Exit(1)
	}

	// convert the money to int32 (as money in NFSU2 is stored as int32 - little endian)
	moneyStored := int32(binary.LittleEndian.Uint32(moneyArr))
	fmt.Printf("Money stored within save file = %d\n", moneyStored)
}

// BackupSaveFile:  Backups up the provided save file.
func BackupSaveFile(origSaveFile string) {
	dateTime := time.Now().Format(time.RFC3339)
	backupFilePath := origSaveFile + "." + strings.ReplaceAll(dateTime, ":", ".") + ".backup"

	// check if the save file is actually a file
	origFileStat, err := os.Stat(origSaveFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open save file '%s': %s\n", origSaveFile, err)
		Exit(1)
	}

	if !origFileStat.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "The provided file is not a regular file '%s'\n", origSaveFile)
		Exit(1)
	}

	// open the original (provided) save file
	orig, err := os.Open(origSaveFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while opening save file '%s': %s\n", origSaveFile, err)
		Exit(1)
	}
	defer orig.Close()

	// create the backup save file
	backup, err := os.Create(backupFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while creating backup save file '%s': %s\n", origSaveFile, err)
		Exit(1)
	}
	defer backup.Close()

	// copy the original into the backup
	if _, err := io.Copy(backup, orig); err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while backuping the save file '%s': %s\n", origSaveFile, err)
		Exit(1)
	}

	fmt.Println("Save file backed up at:  ", backupFilePath)
}

// EditMoneySaveFile:  Edit the money being stored within the provided save file.
func EditMoneySaveFile(money uint32, saveFile string) {
	// open the save file
	file, err := os.OpenFile(saveFile, os.O_RDWR, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open save file '%s': %s\n", saveFile, err)
		Exit(1)
	}
	defer file.Close()

	// convert the money into an array of a 32 bit integer (little endian)
	moneyArr := make([]byte, 4)
	binary.LittleEndian.PutUint32(moneyArr, money)

	// write the new money value within the save file
	_, err = file.WriteAt(moneyArr, NFSU2_MONEY_FILE_OFFSET)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write money into save file '%s': %s\n", saveFile, err)
		Exit(1)
	}

	fmt.Printf("Placed %v money within your save.  Enjoy!\n", money)
}

func Exit(exitCode int) {
	if using_tui {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("~~ Press <enter> to exit ~~")
		reader.ReadString('\n')
	}

	os.Exit(1)
}

func LaunchTui() {
	using_tui = true

	fmt.Println("----------------------------------------------------")
	fmt.Println("     NEED FOR SPEED UNDERGROUND 2 - MONEY CHEAT     ")
	fmt.Println("----------------------------------------------------")
	fmt.Println()
	fmt.Println("Version:", VERSION)
	fmt.Println("URL:    ", "https://github.com/ram-on/nfsu2-money-cheat")
	fmt.Println()
	fmt.Println("Press <CTRL+C> to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	saveFilePath := ""

	// get the save file path
	for {
		fmt.Print("  > Input full path of the NFSU2 Save File:  ")

		scanner.Scan()
		saveFilePath = strings.TrimSpace(scanner.Text())

		if len(saveFilePath) > 0 {
			PrintCurrentMoney(saveFilePath)
			break
		}
	}

	fmt.Println()

	// get the new money value and save it within the Save File
	for {
		fmt.Print("  > Input the money to be save in the Save File:  ")
		scanner.Scan()
		moneyStr := strings.TrimSpace(scanner.Text())

		if len(moneyStr) > 0 {
			if money, err := strconv.ParseUint(moneyStr, 10, 32); err == nil {
				EditMoneySaveFile(uint32(money), saveFilePath)
				break
			} else {
				fmt.Println("Please input a valid number.")
			}
		}
	}
}

func main() {
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)
	infoSaveFilePathVal := infoCmd.String("f", "", "NFSU2 save file path")

	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	moneyVal := editCmd.Int("m", 0, "Money that will be saved in the provided save file.  Minumum amount is 0.  Maximum amount is 2147483647.")
	editSaveFilePathVal := editCmd.String("f", "", "NFSU2 save file path")

	if len(os.Args) == 1 {
		LaunchTui()
	} else if os.Args[1] == "info" {
		infoCmd.Parse(os.Args[2:])

		if *infoSaveFilePathVal == "" {
			fmt.Fprintf(os.Stderr, "Save file not provided\n")
			infoCmd.Usage()
			Exit(1)
		}

		PrintCurrentMoney(*infoSaveFilePathVal)
	} else if os.Args[1] == "edit" {
		editCmd.Parse(os.Args[2:])

		if *moneyVal < 0 {
			fmt.Fprintf(os.Stderr, "Money cannot be < 0\n")
			Exit(1)
		}

		if *editSaveFilePathVal == "" {
			fmt.Fprintf(os.Stderr, "Save file not provided\n")
			editCmd.Usage()
			Exit(1)
		}

		BackupSaveFile(*editSaveFilePathVal)
		EditMoneySaveFile(uint32(*moneyVal), *editSaveFilePathVal)
	} else {
		PrintUsage()
	}
}
