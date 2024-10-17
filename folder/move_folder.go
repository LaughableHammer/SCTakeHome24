package folder

import (
	"errors"
	"strings"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	var srcFolder *Folder
	var dstFolder *Folder

	for i := range f.folders {
		if f.folders[i].Name == name {
			srcFolder = &f.folders[i]
		}
		if f.folders[i].Name == dst {
			dstFolder = &f.folders[i]
		}
	}

	if srcFolder == nil {
		return []Folder{}, errors.New("source folder does not exist")
	}

	if dstFolder == nil {
		return []Folder{}, errors.New("destination folder does not exist")
	}

	if srcFolder.Name == dstFolder.Name {
		return []Folder{}, errors.New("cannot move a folder to itself")
	}

	if strings.HasPrefix(dstFolder.Paths, srcFolder.Paths+".") {
		return []Folder{}, errors.New("cannot move a folder to a child of itself")
	}

	if srcFolder.OrgId != dstFolder.OrgId {
		return []Folder{}, errors.New("cannot move a folder to a different organization")
	}

	// Update the paths for the source folder and its subfolders
	oldPath := srcFolder.Paths
	newPath := dstFolder.Paths + "." + srcFolder.Name

	for i := range f.folders {
		if strings.HasPrefix(f.folders[i].Paths, oldPath) {
			f.folders[i].Paths = strings.Replace(f.folders[i].Paths, oldPath, newPath, 1)
		}
	}

	return f.folders, nil
}
