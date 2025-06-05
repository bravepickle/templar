package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplication(t *testing.T) {
	must := require.New(t)
	app := Application{
		Version:       "1.0",
		GitCommitHash: "123",
	}

	app.Init()
	must.Equal(app.Version, "1.0")
	must.Equal(app.GitCommitHash, "123")

	app = Application{}
	app.Init()
	must.Equal(app.Version, "dev")
	must.Equal(app.GitCommitHash, "<unknown>")
}
