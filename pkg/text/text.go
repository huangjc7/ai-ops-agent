package text

import "strings"

// ExtractAllResults 提取<result>中的json
func ExtractAllResults(reply string) []string {
	var results []string
	startTag := "<result>"
	endTag := "</result>"

	for {
		start := strings.Index(reply, startTag)
		end := strings.Index(reply, endTag)

		if start == -1 || end == -1 || start >= end {
			break
		}

		content := reply[start+len(startTag) : end]
		results = append(results, content)

		// 移除已提取部分，继续找下一段
		reply = reply[end+len(endTag):]
	}

	return results
}
