package folder

import (
	"errors"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	return findFoldersByOrgId(&f.folders, orgID)
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	var sameNamedFolders []Folder = findFoldersByName(&f.folders, name)

	if len(sameNamedFolders) == 0 {
		return []Folder{}, errors.New("folder does not exist")
	}

	var parentFolders []Folder = findFoldersByOrgId(&sameNamedFolders, orgID)

	if len(parentFolders) == 0 {
		return []Folder{}, errors.New("folder does not exist in the specified organization")
	}

	res := []Folder{}
	addedFilePaths := map[string]bool{}
	for _, parent := range parentFolders {
		for _, f := range f.folders {
			_, alreadyAdded := addedFilePaths[f.Paths]
			if !alreadyAdded && isChildFolder(&parent, &f) {
				res = append(res, f)
				addedFilePaths[f.Paths] = true
			}
		}
	}

	return res, nil
}
