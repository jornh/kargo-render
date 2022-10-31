package metadata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akuityio/bookkeeper/internal/file"
	"github.com/stretchr/testify/require"
)

func TestLoadTargetBranchMetadata(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func() string
		assertions func(*TargetBranchMetadata, error)
	}{
		{
			name: "metadata does not exist",
			setup: func() string {
				repoDir, err := os.MkdirTemp("", "")
				require.NoError(t, err)
				return repoDir
			},
			assertions: func(md *TargetBranchMetadata, err error) {
				require.NoError(t, err)
				require.Nil(t, md)
			},
		},
		{
			name: "invalid YAML",
			setup: func() string {
				repoDir, err := os.MkdirTemp("", "")
				require.NoError(t, err)
				bkDir := filepath.Join(repoDir, ".bookkeeper")
				err = os.Mkdir(bkDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(
					filepath.Join(bkDir, "metadata.yaml"),
					[]byte("bogus"),
					0600,
				)
				require.NoError(t, err)
				return repoDir
			},
			assertions: func(_ *TargetBranchMetadata, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "error unmarshaling branch metadata")
			},
		},
		{
			name: "valid YAML",
			setup: func() string {
				repoDir, err := os.MkdirTemp("", "")
				require.NoError(t, err)
				bkDir := filepath.Join(repoDir, ".bookkeeper")
				err = os.Mkdir(bkDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(
					filepath.Join(bkDir, "metadata.yaml"),
					[]byte(""), // An empty file should actually be valid
					0600,
				)
				require.NoError(t, err)
				return repoDir
			},
			assertions: func(_ *TargetBranchMetadata, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			md, err := LoadTargetBranchMetadata(testCase.setup())
			testCase.assertions(md, err)
		})
	}
}

func TestWriteTargetBranchMetadata(t *testing.T) {
	repoDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	err = WriteTargetBranchMetadata(
		TargetBranchMetadata{
			SourceCommit: "1234567",
		},
		repoDir,
	)
	require.NoError(t, err)
	exists, err :=
		file.Exists(filepath.Join(repoDir, ".bookkeeper", "metadata.yaml"))
	require.NoError(t, err)
	require.True(t, exists)
}
