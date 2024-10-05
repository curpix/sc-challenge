package folder_test

import (
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
