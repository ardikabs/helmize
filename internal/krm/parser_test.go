package krm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/ardikabs/helmize/internal/errs"
	"github.com/ardikabs/helmize/internal/krm"
	"github.com/sters/yaml-diff/yamldiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKRM_Parse(t *testing.T) {

	testcases := []struct {
		name    string
		wantErr error
	}{
		{
			name: "default",
		},
		{
			name:    "invalid-input-kind",
			wantErr: errs.ErrInvalidObject,
		},
		{
			name:    "invalid-input-apiVersion",
			wantErr: errs.ErrInvalidObject,
		},
		{
			name:    "invalid-input-noname",
			wantErr: errs.ErrInvalidObject,
		},
		{
			name:    "invalid-input-nonamespace",
			wantErr: errs.ErrInvalidObject,
		},
		{
			name:    "invalid-input-norepo",
			wantErr: errs.ErrInvalidObject,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in, err := os.ReadFile(filepath.Join("testdata", tc.name, "input.yaml"))
			require.NoError(t, err)

			out, err := fn.Run(fn.ResourceListProcessorFunc(krm.Parse), in)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			} else {
				require.NoError(t, err)
			}

			rl, _ := fn.ParseResourceList(out)
			assert.NotEqual(t, len(rl.Items), 0)

			want, err := os.ReadFile(filepath.Join("testdata", tc.name, "want.yaml"))
			require.NoError(t, err)

			checkDiff(t, string(want), rl.Items.String())
		})
	}
}

func checkDiff(t *testing.T, expected, actual string) {
	var result string
	status := yamldiff.DiffStatusSame

	expectedYAML, err := yamldiff.Load(expected)
	require.NoError(t, err)

	actualYAML, err := yamldiff.Load(actual)
	require.NoError(t, err)

	for _, diff := range yamldiff.Do(expectedYAML, actualYAML, yamldiff.EmptyAsNull()) {
		result += diff.Dump()
		result += "\n"

		if diff.Status() != yamldiff.DiffStatusSame {
			status = diff.Status()
		}
	}

	if status != yamldiff.DiffStatusSame {
		t.Errorf("mismatch (-expected +actual):\n%s\n", result)
	}
}
