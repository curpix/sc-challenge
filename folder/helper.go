package folder

import (
	"strings"

	"github.com/gofrs/uuid"
)

// filterFolders filters folders with a predicate
func filterFolders(folders *[]Folder, predicate func(Folder) bool) []Folder {
	result := []Folder{}
	for _, f := range *folders {
		if predicate(f) {
			result = append(result, f)
		}
	}
	return result
}

// findFoldersByName filters via folder name
func findFoldersByName(folders *[]Folder, name string) []Folder {
	return filterFolders(folders, func(f Folder) bool {
		return f.Name == name
	})
}

// findFoldersByOrgId filters via folder orgId
func findFoldersByOrgId(folders *[]Folder, orgId uuid.UUID) []Folder {
	return filterFolders(folders, func(f Folder) bool {
		return f.OrgId == orgId
	})
}

// isChildFolder checks if the first folder is a parent, and the second is a child
func isChildFolder(parent *Folder, child *Folder) bool {
	return (child.OrgId == parent.OrgId &&
		child.Paths != parent.Paths && strings.HasPrefix(child.Paths, parent.Paths))
}
