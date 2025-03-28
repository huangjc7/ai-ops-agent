package cmd

type commandFlags struct {
	apiKey    string
	model     string
	baseURL   string
	aiMaxStep int
}

var cmdFlags commandFlags
