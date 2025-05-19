package command

type Application struct {
	Version       string
	GitCommitHash string
}

func (a *Application) Init() {
	if a.Version == "" {
		a.Version = "dev"
	}

	if a.GitCommitHash == "" {
		a.GitCommitHash = "<unknown>"
	}
}
