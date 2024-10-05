package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_GetAllFolders(t *testing.T) {
	t.Parallel()

	var result []folder.Folder = folder.GetAllFolders()
	assert.NotNil(t, result, "Folders are returned")
}

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	var validOrgId = uuid.FromStringOrNil("c59cc5c1-9b81-4d00-95e3-22c6efdaf134")

	tests := [...]struct {
		testName string
		orgID    uuid.UUID
		folders  []folder.Folder
		want     []folder.Folder
	}{
		{
			testName: "No folders gives empty slice",
			orgID:    validOrgId,
			folders:  []folder.Folder{},
			want:     []folder.Folder{},
		},
		{
			testName: "orgID not present gives empty slice",
			orgID:    validOrgId,
			folders: []folder.Folder{
				{
					Name:  "fold",
					OrgId: uuid.FromStringOrNil("5e35ff8f-21dd-4ed5-b861-ba93dbcdadc3"),
					Paths: "fold",
				},
			},
			want: []folder.Folder{},
		},
		{
			testName: "orgID present in root folder",
			orgID:    validOrgId,
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
			},
			want: []folder.Folder{{Name: "fold", OrgId: validOrgId, Paths: "fold"}},
		},
		{
			testName: "orgID present in parent and child folder - return both",
			orgID:    validOrgId,
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "child-fold", OrgId: validOrgId, Paths: "fold.child-fold"},
			},
			want: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "child-fold", OrgId: validOrgId, Paths: "fold.child-fold"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			folders := f.GetFoldersByOrgID(tt.orgID)
			assert.Equal(t, tt.want, folders, tt.testName)
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	var validOrgId = uuid.FromStringOrNil("c59cc5c1-9b81-4d00-95e3-22c6efdaf134")

	tests := [...]struct {
		testName string
		orgID    uuid.UUID
		name     string
		folders  []folder.Folder
		want     []folder.Folder
		err      error
	}{
		{
			testName: "Folder with no children, empty slice returned",
			orgID:    validOrgId,
			name:     "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
			},
			want: []folder.Folder{},
		},
		{
			testName: "Folder with multiple children, all returned",
			orgID:    validOrgId,
			name:     "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c2"},
			},
			want: []folder.Folder{
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c2"},
			},
		},
		{
			testName: "Folder with nested children, all returned",
			orgID:    validOrgId,
			name:     "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c1.c2"},
			},
			want: []folder.Folder{
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c1.c2"},
			},
		},
		{
			testName: "Folder with parent and child, only child returned",
			orgID:    validOrgId,
			name:     "c1",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c1.c2"},
			},
			want: []folder.Folder{
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c1.c2"},
			},
		},
		{
			testName: "Child of same file name, but different organisation - not returned",
			orgID:    validOrgId,
			name:     "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "fold", OrgId: uuid.FromStringOrNil("5e35ff8f-21dd-4ed5-b861-ba93dbcdadc3"), Paths: "fold"},
				{Name: "c1", OrgId: uuid.FromStringOrNil("5e35ff8f-21dd-4ed5-b861-ba93dbcdadc3"), Paths: "fold.c1"},
			},
			want: []folder.Folder{},
		},
		{
			testName: "Error: Folder does not exist",
			orgID:    validOrgId,
			name:     "missing",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
			},
			err: errors.New(""),
		},
		{
			testName: "Error: Folder belongs to another organisation",
			orgID:    validOrgId,
			name:     "fold",
			folders: []folder.Folder{
				{
					Name:  "fold",
					OrgId: uuid.FromStringOrNil("5e35ff8f-21dd-4ed5-b861-ba93dbcdadc3"),
					Paths: "fold",
				},
			},
			err: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			folders, err := f.GetAllChildFolders(tt.orgID, tt.name)

			if tt.err == nil {
				assert.NoError(t, err, tt.testName)
				assert.Equal(t, tt.want, folders, tt.testName)
			} else {
				assert.Error(t, err, tt.testName)
			}
		})
	}
}
