package command_test

import (
	"context"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
	"github.com/pkg/errors"

	// Prevents a circular import with "github.com/docker/cli/internal/test"
	. "github.com/docker/cli/cli/command"
	"github.com/docker/cli/internal/test"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type fakeClient struct {
	client.Client
	infoFunc func() (types.Info, error)
}

func (cli *fakeClient) Info(_ context.Context) (types.Info, error) {
	if cli.infoFunc != nil {
		return cli.infoFunc()
	}
	return types.Info{}, nil
}

func TestElectAuthServer(t *testing.T) {
	testCases := []struct {
		expectedAuthServer string
		expectedWarning    string
		infoFunc           func() (types.Info, error)
	}{
		{
			expectedAuthServer: "https://index.docker.io/v1/",
			expectedWarning:    "",
			infoFunc: func() (types.Info, error) {
				return types.Info{IndexServerAddress: "https://index.docker.io/v1/"}, nil
			},
		},
		{
			expectedAuthServer: "https://index.docker.io/v1/",
			expectedWarning:    "Empty registry endpoint from daemon",
			infoFunc: func() (types.Info, error) {
				return types.Info{IndexServerAddress: ""}, nil
			},
		},
		{
			expectedAuthServer: "https://foo.bar",
			expectedWarning:    "",
			infoFunc: func() (types.Info, error) {
				return types.Info{IndexServerAddress: "https://foo.bar"}, nil
			},
		},
		{
			expectedAuthServer: "https://index.docker.io/v1/",
			expectedWarning:    "failed to get default registry endpoint from daemon",
			infoFunc: func() (types.Info, error) {
				return types.Info{}, errors.Errorf("error getting info")
			},
		},
	}
	for _, tc := range testCases {
		cli := test.NewFakeCli(&fakeClient{infoFunc: tc.infoFunc})
		server := ElectAuthServer(context.Background(), cli)
		assert.Check(t, is.Equal(tc.expectedAuthServer, server))
		actual := cli.ErrBuffer().String()
		if tc.expectedWarning == "" {
			assert.Check(t, is.Len(actual, 0))
		} else {
			assert.Check(t, is.Contains(actual, tc.expectedWarning))
		}
	}
}
