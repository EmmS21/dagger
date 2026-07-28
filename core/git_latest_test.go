package core

import (
	"testing"

	"github.com/dagger/dagger/util/gitutil"
	"github.com/stretchr/testify/require"
)

func TestSelectLatestGitRef(t *testing.T) {
	t.Parallel()

	const (
		headCommit       = "0000000000000000000000000000000000000001"
		stableCommit     = "0000000000000000000000000000000000000002"
		prereleaseCommit = "0000000000000000000000000000000000000003"
	)
	remote := &gitutil.Remote{
		Refs: []*gitutil.Ref{
			{Name: "HEAD", SHA: headCommit},
			{Name: "refs/heads/main", SHA: headCommit},
			{Name: "refs/tags/v1.9.0", SHA: "0000000000000000000000000000000000000004"},
			{Name: "refs/tags/2.0.0", SHA: stableCommit},
			{Name: "refs/tags/v3.0.0-rc.1", SHA: prereleaseCommit},
			{Name: "refs/tags/not-a-release", SHA: "0000000000000000000000000000000000000005"},
		},
		Symrefs: map[string]string{"HEAD": "refs/heads/main"},
	}

	t.Run("stable releases by default", func(t *testing.T) {
		ref, err := SelectLatestGitRef(remote, false)
		require.NoError(t, err)
		require.Equal(t, "refs/tags/2.0.0", ref.Name)
		require.Equal(t, stableCommit, ref.SHA)
	})

	t.Run("subreleases when requested", func(t *testing.T) {
		ref, err := SelectLatestGitRef(remote, true)
		require.NoError(t, err)
		require.Equal(t, "refs/tags/v3.0.0-rc.1", ref.Name)
		require.Equal(t, prereleaseCommit, ref.SHA)
	})
}

func TestSelectLatestGitRefFallsBackToHead(t *testing.T) {
	t.Parallel()

	const headCommit = "0000000000000000000000000000000000000001"
	remote := &gitutil.Remote{
		Refs: []*gitutil.Ref{
			{Name: "HEAD", SHA: headCommit},
			{Name: "refs/heads/trunk", SHA: headCommit},
			{Name: "refs/tags/nightly", SHA: "0000000000000000000000000000000000000002"},
		},
		Symrefs: map[string]string{"HEAD": "refs/heads/trunk"},
	}

	ref, err := SelectLatestGitRef(remote, false)
	require.NoError(t, err)
	require.Equal(t, "refs/heads/trunk", ref.Name)
	require.Equal(t, headCommit, ref.SHA)
}

func TestGitRefPinRoundTrip(t *testing.T) {
	t.Parallel()

	ref := &gitutil.Ref{
		Name: "refs/tags/v1.2.3",
		SHA:  "0123456789abcdef0123456789abcdef01234567",
	}
	pin, err := EncodeGitRefPin(ref)
	require.NoError(t, err)
	require.Equal(t, "refs/tags/v1.2.3@0123456789abcdef0123456789abcdef01234567", pin)

	decoded, err := DecodeGitRefPin(pin)
	require.NoError(t, err)
	require.Equal(t, ref, decoded)
}
