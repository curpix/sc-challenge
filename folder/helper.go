package folder

import "github.com/gofrs/uuid"

func filterFolders(folders *[]Folder, predicate func(Folder) bool) []Folder {
	result := []Folder{}
	for _, f := range *folders {
		if predicate(f) {
			result = append(result, f)
		}
	}
	return result
}

func findFoldersByName(folders *[]Folder, name string) []Folder {
	return filterFolders(folders, func(f Folder) bool {
		return f.Name == name
	})
}

func findFoldersByOrgId(folders *[]Folder, orgId uuid.UUID) []Folder {
	return filterFolders(folders, func(f Folder) bool {
		return f.OrgId == orgId
	})
}
