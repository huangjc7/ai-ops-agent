package prompt

type PromptTemplate struct {
	System string
	User   string
}

// InputClassResult 定义结构体接收
var InputClassResult struct {
	Type string `json:"type"`
}

type Suggestion struct {
	Desc string `json:"desc"` // 用途说明
	Cmd  string `json:"cmd"`  // shell 命令
}
type SuggestionList []Suggestion
