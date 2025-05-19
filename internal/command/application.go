package command

import "os"

type Application struct {
	Version       string
	GitCommitHash string
	WorkDir       string
}

func (a *Application) Init() {
	if a.Version == "" {
		a.Version = "0.0.0"
	}

	if a.GitCommitHash == "" {
		a.GitCommitHash = "<unknown>"
	}

	if a.WorkDir == "" {
		a.WorkDir, _ = os.Getwd()
	}
}
