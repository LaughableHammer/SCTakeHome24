package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

	// Load folder data
	res := folder.GetAllFolders()
	folderDriver := folder.NewDriver(res)

	// CLI
	fmt.Println("\n---------------------")
	fmt.Println("CLI File Explorer")
	fmt.Println("---------------------")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. List all folders")
		fmt.Println("2. List folders by OrgID")
		fmt.Println("3. Move a folder")
		fmt.Println("4. Exit")
		fmt.Print("\nEnter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			// List all folders
			folder.PrettyPrint(res)

		case "2":
			// List folders by OrgID
			fmt.Print("Enter OrgID (or press Enter for default): ")
			orgIDInput, _ := reader.ReadString('\n')
			orgIDInput = strings.TrimSpace(orgIDInput)

			if orgIDInput != "" {
				orgID = uuid.FromStringOrNil(orgIDInput)
			}
			orgFolders := folderDriver.GetFoldersByOrgID(orgID)
			fmt.Printf("\nFolders for orgID: %s\n", orgID)
			folder.PrettyPrint(orgFolders)

		case "3":
			// Move a folder
			fmt.Print("Enter the name of the folder to move: ")
			sourceFolder, _ := reader.ReadString('\n')
			sourceFolder = strings.TrimSpace(sourceFolder)

			fmt.Print("Enter the destination folder name: ")
			destinationFolder, _ := reader.ReadString('\n')
			destinationFolder = strings.TrimSpace(destinationFolder)

			// Attempt to move the folder
			movedFolders, err := folderDriver.MoveFolder(sourceFolder, destinationFolder)
			if err != nil {
				fmt.Printf("Error moving folder: %v\n", err)
			} else {
				fmt.Printf("Folder %s moved to %s successfully!\n", sourceFolder, destinationFolder)
				folder.PrettyPrint(movedFolders)
			}

		case "4":
			// Exit the program
			fmt.Println("Exiting the program.")
			os.Exit(0)

		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}
