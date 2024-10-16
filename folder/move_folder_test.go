package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()

	defaultOrgID := uuid.FromStringOrNil(folder.DefaultOrgID)
	altOrgID := uuid.Must(uuid.NewV4())

	testData := []folder.Folder{
		{
			Name:  "alpha",
			Paths: "alpha",
			OrgId: defaultOrgID,
		},
		{
			Name:  "bravo",
			Paths: "alpha.bravo",
			OrgId: defaultOrgID,
		},
		{
			Name:  "charlie",
			Paths: "alpha.bravo.charlie",
			OrgId: defaultOrgID,
		},
		{
			Name:  "delta",
			Paths: "alpha.delta",
			OrgId: defaultOrgID,
		},
		{
			Name:  "echo",
			Paths: "alpha.delta.echo",
			OrgId: defaultOrgID,
		},
		{
			Name:  "foxtrot",
			Paths: "foxtrot",
			OrgId: altOrgID,
		},
		{
			Name:  "golf",
			Paths: "golf",
			OrgId: defaultOrgID,
		},
	}

	tests := []struct {
		name        string
		destination string
		want        []folder.Folder
		wantErr     bool
	}{
		// Expect output: When moving bravo to delta
		{
			name:        "bravo",
			destination: "delta",
			want: []folder.Folder{
				{
					Name:  "alpha",
					Paths: "alpha",
					OrgId: defaultOrgID,
				},
				{
					Name:  "bravo",
					Paths: "alpha.delta.bravo",
					OrgId: defaultOrgID,
				},
				{
					Name:  "charlie",
					Paths: "alpha.delta.bravo.charlie",
					OrgId: defaultOrgID,
				},
				{
					Name:  "delta",
					Paths: "alpha.delta",
					OrgId: defaultOrgID,
				},
				{
					Name:  "echo",
					Paths: "alpha.delta.echo",
					OrgId: defaultOrgID,
				},
				{
					Name:  "foxtrot",
					Paths: "foxtrot",
					OrgId: altOrgID,
				},
				{
					Name:  "golf",
					Paths: "golf",
					OrgId: defaultOrgID,
				},
			},
			wantErr: false,
		},
		// Expect Error: When moving a folder into a different OrgId
		{
			name:        "bravo",
			destination: "foxtrot",
			want:        []folder.Folder{},
			wantErr:     true,
		},
		// Expect Error: When moving a folder into itself
		{
			name:        "bravo",
			destination: "bravo",
			want:        []folder.Folder{},
			wantErr:     true,
		},
		// Expect Error: When moving a folder to its child
		{
			name:        "bravo",
			destination: "charlie",
			want:        []folder.Folder{},
			wantErr:     true,
		},
		// Expect Error: When source folder doesn't exist
		{
			name:        "invalid_folder",
			destination: "delta",
			want:        []folder.Folder{},
			wantErr:     true,
		},
		// Expect Error: When moving to non-existent destination folder
		{
			name:        "bravo",
			destination: "invalid_folder",
			want:        []folder.Folder{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(testData)
			res, err := f.MoveFolder(tt.name, tt.destination)

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
