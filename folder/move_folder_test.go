package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()

	var validOrgId = uuid.FromStringOrNil("c59cc5c1-9b81-4d00-95e3-22c6efdaf134")

	tests := [...]struct {
		testName string
		src      string
		dst      string
		folders  []folder.Folder
		want     []folder.Folder
		err      error
	}{
		{
			testName: "Move folder with no children",
			src:      "c1",
			dst:      "b",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "a", OrgId: validOrgId, Paths: "a"},
				{Name: "b", OrgId: validOrgId, Paths: "a.b"},
			},
			want: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "a.b.c1"},
				{Name: "a", OrgId: validOrgId, Paths: "a"},
				{Name: "b", OrgId: validOrgId, Paths: "a.b"},
			},
		},
		{
			testName: "Move root folder with multiple children",
			src:      "fold",
			dst:      "a",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c2"},
				{Name: "a", OrgId: validOrgId, Paths: "a"},
				{Name: "b", OrgId: validOrgId, Paths: "a.b"},
			},
			want: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "a.fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "a.fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "a.fold.c2"},
				{Name: "a", OrgId: validOrgId, Paths: "a"},
				{Name: "b", OrgId: validOrgId, Paths: "a.b"},
			},
		},
		{
			testName: "Move folder with nested children",
			src:      "c",
			dst:      "a",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "b", OrgId: validOrgId, Paths: "fold.b"},
				{Name: "c", OrgId: validOrgId, Paths: "fold.c"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c.c2"},
				{Name: "c1-nest", OrgId: validOrgId, Paths: "fold.c.c1.c1-nest"},
				{Name: "c2-nest", OrgId: validOrgId, Paths: "fold.c.c2.c2-nest"},
				{Name: "a", OrgId: validOrgId, Paths: "a"},
			},
			want: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "b", OrgId: validOrgId, Paths: "fold.b"},
				{Name: "c", OrgId: validOrgId, Paths: "a.c"},
				{Name: "c1", OrgId: validOrgId, Paths: "a.c.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "a.c.c2"},
				{Name: "c1-nest", OrgId: validOrgId, Paths: "a.c.c1.c1-nest"},
				{Name: "c2-nest", OrgId: validOrgId, Paths: "a.c.c2.c2-nest"},
				{Name: "a", OrgId: validOrgId, Paths: "a"},
			},
		},
		{
			testName: "Move child folder to its parent",
			src:      "c1",
			dst:      "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "b", OrgId: validOrgId, Paths: "fold.b"},
				{Name: "c", OrgId: validOrgId, Paths: "fold.c"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c.c2"},
				{Name: "c1-nest", OrgId: validOrgId, Paths: "fold.c.c1.c1-nest"},
			},
			want: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "b", OrgId: validOrgId, Paths: "fold.b"},
				{Name: "c", OrgId: validOrgId, Paths: "fold.c"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c.c2"},
				{Name: "c1-nest", OrgId: validOrgId, Paths: "fold.c1.c1-nest"},
			},
		},
		{
			testName: "Error: Cannot move a folder to a child of itself",
			src:      "c1",
			dst:      "c3",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "c1", OrgId: validOrgId, Paths: "fold.c1"},
				{Name: "c2", OrgId: validOrgId, Paths: "fold.c1.c2"},
				{Name: "c3", OrgId: validOrgId, Paths: "fold.c1.c2.c3"},
			},
			err: errors.New(""),
		},
		{
			testName: "Error: Cannot move a folder to itself",
			src:      "fold",
			dst:      "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
			},
			err: errors.New(""),
		},
		{
			testName: "Error: Cannot move a folder to a different organization",
			src:      "fold",
			dst:      "a",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
				{Name: "a", OrgId: uuid.FromStringOrNil("5e35ff8f-21dd-4ed5-b861-ba93dbcdadc3"), Paths: "a"},
			},
			err: errors.New(""),
		},
		{
			testName: "Error: Source folder does not exist",
			src:      "missing",
			dst:      "fold",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
			},
			err: errors.New(""),
		},
		{
			testName: "Error: Destination folder does not exist",
			src:      "fold",
			dst:      "missing",
			folders: []folder.Folder{
				{Name: "fold", OrgId: validOrgId, Paths: "fold"},
			},
			err: errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			folders, err := f.MoveFolder(tt.src, tt.dst)

			if tt.err == nil {
				assert.NoError(t, err, tt.testName)
				assert.Equal(t, tt.want, folders, tt.testName)
			} else {
				assert.Error(t, err, tt.testName)
			}
		})
	}
}
