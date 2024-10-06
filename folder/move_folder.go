package folder

import (
	"errors"
	"strings"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	if name == dst {
		return []Folder{}, errors.New("cannot move a folder to itself")
	}

	var srcFolders []Folder = findFoldersByName(&f.folders, name)
	var dstFolders []Folder = findFoldersByName(&f.folders, dst)

	if len(srcFolders) == 0 {
		return []Folder{}, errors.New("source folder does not exist")
	}

	if len(dstFolders) == 0 {
		return []Folder{}, errors.New("destination folder does not exist")
	}

	// Refer to Assumption
	var srcFolder Folder = srcFolders[0]
	var dstFolder Folder = dstFolders[0]

	if srcFolder.OrgId != dstFolder.OrgId {
		return []Folder{}, errors.New("cannot move a folder to a different organization")
	}

	if isChildFolder(&srcFolder, &dstFolder) {
		return []Folder{}, errors.New("cannot move a folder to a child of itself")
	}

	resultAfterMove := []Folder{}
	for _, folder := range f.folders {
		if isMovedFolder(&folder, &srcFolder) {
			movedFolder := getNewMovedFolder(folder, srcFolder, dstFolder)
			resultAfterMove = append(resultAfterMove, movedFolder)
		} else {
			resultAfterMove = append(resultAfterMove, folder)
		}
	}

	return resultAfterMove, nil
}

func isMovedFolder(folder *Folder, srcFolder *Folder) bool {
	return srcFolder.Name == folder.Name || // Refer to assumption
		isChildFolder(srcFolder, folder)
}

func getNewMovedFolder(folder Folder, srcFolder Folder, dstFolder Folder) Folder {
	var relativePathFromSrc, _ = strings.CutPrefix(folder.Paths, srcFolder.Paths)
	var newPath string = dstFolder.Paths + "." + srcFolder.Name
	if relativePathFromSrc != "" {
		newPath += relativePathFromSrc
	}

	return Folder{
		Name:  folder.Name,
		OrgId: folder.OrgId,
		Paths: newPath,
	}
}
