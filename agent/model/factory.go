package model

import "github.com/sashabaranov/go-openai"

type AIClient interface {
	Send(input string) (string, error)
}

type AI struct {
	client *openai.Client
	model  string
}
