package folder

import "github.com/gofrs/uuid"

// IDriver provides utility functions for manipulating folder structures
type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []Folder

	// GetAllChildFolders returns all child folders of a specific folder within same organisation.
	// Duplicate files with same names but different paths are possible,
	// where all possible children are returned without duplicates
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// MoveFolder moves a folder to a new destination.
	// Assumption: Names are unique, as cannot distinguish between different paths.
	// E.g. We have a, c.a and d. moveFolder("d", "a") is ambiguous, as we don't know
	// which a to move to.
	MoveFolder(name string, dst string) ([]Folder, error)
}

// A driver which stores folders
type driver struct {
	folders []Folder
}

// Creates a driver to execute utility functions
func NewDriver(folders []Folder) IDriver {
	return &driver{
		folders: folders,
	}
}
