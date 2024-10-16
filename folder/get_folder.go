package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}

// Checks if a particular folder exists
func DoesFolderExist(folders []Folder, name string) bool {
	for _, folder := range folders {
		// Check if the folder.Paths contains the name
		if strings.Contains(folder.Paths, name) {
			return true
		}
	}

	return false
}

// Checks if a particular folder has child folders
func HasSubFolders(Path string, name string) bool {
	return strings.Contains(Path, name+".")
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	name = strings.ToLower(name)
	folders := f.GetFoldersByOrgID(orgID)
	res := []Folder{}

	if len(f.folders) == 0 {
		return res, errors.New("no folders exist")
	}

	if len(folders) == 0 {
		return res, errors.New("no such orgid exists")
	}

	if !DoesFolderExist(f.folders, name) {
		return res, errors.New("invalid folder requested")
	} else if !DoesFolderExist(folders, name) {
		return res, errors.New("folder doesn't exist in specified organisation")
	}

	for _, folder := range folders {
		if HasSubFolders(folder.Paths, name) {
			res = append(res, folder)
		}
	}

	return res, nil
}
