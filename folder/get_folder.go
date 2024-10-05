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
	return findFoldersByOrgId(&f.folders, orgID)
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	var sameNamedFolders []Folder = findFoldersByName(&f.folders, name)

	if len(sameNamedFolders) == 0 {
		return []Folder{}, errors.New("folder does not exist")
	}

	var orgSameNamedFolders []Folder = findFoldersByOrgId(&sameNamedFolders, orgID)

	if len(orgSameNamedFolders) == 0 {
		return []Folder{}, errors.New("folder does not exist in the specified organization")
	}

	res := []Folder{}
	addedFilePaths := map[string]bool{}
	for _, parent := range orgSameNamedFolders {
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

func isChildFolder(parent *Folder, child *Folder) bool {
	return (child.OrgId == parent.OrgId &&
		child.Paths != parent.Paths && strings.HasPrefix(child.Paths, parent.Paths))
}
