package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// feel free to change how the unit test is structured
func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	// Sample test data
	testData := folder.GetSampleData()
	defaultOrgId := uuid.FromStringOrNil(folder.DefaultOrgID)

	// Check case where input is valid
	f := folder.NewDriver(testData)
	res := f.GetFoldersByOrgID(defaultOrgId)

	// Verify correct output for realistic data set
	for _, folder := range res {
		assert.Equal(t, folder.OrgId, defaultOrgId,
			"Folders with invalid incorrect OrgId were returned")
	}

	errorTests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		{
			name:    "Expect empty: Invalid OrgId Provided",
			orgID:   uuid.Must(uuid.NewV4()),
			folders: testData,
			want:    []folder.Folder{},
		},
		{
			name:    "Expect empty: Empty folders input",
			orgID:   uuid.Must(uuid.NewV4()),
			folders: []folder.Folder{},
			want:    []folder.Folder{},
		},
		// Path is irrelevant so can be ommitted
		{
			name:  "Expect output: Single Folder given",
			orgID: defaultOrgId,
			folders: []folder.Folder{
				{OrgId: defaultOrgId, Name: "Folder 1"},
			},
			want: []folder.Folder{
				{OrgId: defaultOrgId, Name: "Folder 1"},
			},
		},
		{
			name:  "Expect output: Multiple folders given",
			orgID: defaultOrgId,
			folders: []folder.Folder{
				{OrgId: defaultOrgId, Name: "Folder 1"},
				{OrgId: uuid.Must(uuid.NewV4()), Name: "Folder 2"},
				{OrgId: defaultOrgId, Name: "Folder 3"},
			},
			want: []folder.Folder{
				{OrgId: defaultOrgId, Name: "Folder 1"},
				{OrgId: defaultOrgId, Name: "Folder 3"},
			},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			res := f.GetFoldersByOrgID(tt.orgID)
			assert.Equal(t, res, tt.want, tt.name) // Checks if output == expected
		})
	}
}

func Benchmark_folder_GetFoldersByOrgID(b *testing.B) {
	testData := folder.GetSampleData()
	defaultOrgId := uuid.FromStringOrNil(folder.DefaultOrgID)
	f := folder.NewDriver(testData)

	for i := 0; i < b.N; i++ {
		f.GetFoldersByOrgID(defaultOrgId)
	}
}

// TODO: write tests for GetAllChildFolders
//	Have coverage testing as well
//  Have benchmark testing
