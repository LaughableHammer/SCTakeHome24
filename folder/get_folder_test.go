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

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()
	defaultOrgId := uuid.FromStringOrNil(folder.DefaultOrgID)
	altOrgID := uuid.Must(uuid.NewV4())

	var testData = []folder.Folder{
		{
			Name:  "alpha",
			Paths: "alpha",
			OrgId: defaultOrgId,
		},
		{
			Name:  "bravo",
			Paths: "alpha.bravo",
			OrgId: defaultOrgId,
		},
		{
			Name:  "charlie",
			Paths: "alpha.bravo.charlie",
			OrgId: defaultOrgId,
		},
		{
			Name:  "delta",
			Paths: "alpha.delta",
			OrgId: defaultOrgId,
		},
		{
			Name:  "delta",
			Paths: "alpha.delta.bravo.delta",
			OrgId: defaultOrgId,
		},
		{
			Name:  "delta",
			Paths: "alpha.echo.delta",
			OrgId: defaultOrgId,
		},
		{
			Name:  "echo",
			Paths: "echo",
			OrgId: defaultOrgId,
		},
		{
			Name:  "foxtrot",
			Paths: "foxtrot",
			OrgId: altOrgID,
		},
	}

	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
		wantErr bool
	}{
		// Expect output: When 1 folder matches criteria
		{
			name:    "echo",
			orgID:   defaultOrgId,
			folders: testData,
			want: []folder.Folder{
				{
					Name:  "delta",
					Paths: "alpha.echo.delta",
					OrgId: defaultOrgId,
				},
			},
			wantErr: false,
		},
		// Expect output: When multiple folders match criteria
		{
			name:    "alpha",
			orgID:   defaultOrgId,
			folders: testData,
			want: []folder.Folder{
				{
					Name:  "bravo",
					Paths: "alpha.bravo",
					OrgId: defaultOrgId,
				},
				{
					Name:  "charlie",
					Paths: "alpha.bravo.charlie",
					OrgId: defaultOrgId,
				},
				{
					Name:  "delta",
					Paths: "alpha.delta",
					OrgId: defaultOrgId,
				},
				{
					Name:  "delta",
					Paths: "alpha.delta.bravo.delta",
					OrgId: defaultOrgId,
				},
				{
					Name:  "delta",
					Paths: "alpha.echo.delta",
					OrgId: defaultOrgId,
				},
			},
			wantErr: false,
		},
		// Expect output: When name is in upper case
		{
			name:    "EcHo",
			orgID:   defaultOrgId,
			folders: testData,
			want: []folder.Folder{
				{
					Name:  "delta",
					Paths: "alpha.echo.delta",
					OrgId: defaultOrgId,
				},
			},
			wantErr: false,
		},
		// Expect output: When the folder name repeats in the path at the leaf
		{
			name:    "delta",
			orgID:   defaultOrgId,
			folders: testData,
			want: []folder.Folder{
				{
					Name:  "delta",
					Paths: "alpha.delta.bravo.delta",
					OrgId: defaultOrgId,
				},
			},
			wantErr: false,
		},
		// Expect Empty: When leaf folder is requested
		{
			name:    "foxtrot",
			orgID:   altOrgID,
			folders: testData,
			want:    []folder.Folder{},
			wantErr: false,
		},
		// Expect Error: When invalid folder is requested
		{
			name:    "invalid_folder",
			orgID:   defaultOrgId,
			folders: testData,
			want:    []folder.Folder{},
			wantErr: true,
		},
		// Expect Error: When folder doesn't exist in the specified organisation
		{
			name:    "foxtrot",
			orgID:   defaultOrgId,
			folders: testData,
			want:    []folder.Folder{},
			wantErr: true,
		},
		// Expect Error: When invalid OrgId is provided
		{
			name:    "alpha",
			orgID:   uuid.Must(uuid.NewV4()),
			folders: testData,
			want:    []folder.Folder{},
			wantErr: true,
		},
		// Expect Error: When no folders exist
		{
			name:    "alpha",
			orgID:   defaultOrgId,
			folders: []folder.Folder{},
			want:    []folder.Folder{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			res, err := f.GetAllChildFolders(tt.orgID, tt.name)

			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got nil\n\n")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one\n\n")
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, res, "Unexpected result returned\n\n")
			}
		})
	}
}
