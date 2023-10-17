package formatter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ms-henglu/pal/types"
	"github.com/ms-henglu/pal/utils"
)

var _ Formatter = MarkdownFormatter{}

type MarkdownFormatter struct {
}

func (m MarkdownFormatter) Format(r types.RequestTrace) string {
	content := markdownTemplate
	content = strings.ReplaceAll(content, "{Time}", r.TimeStamp.Format("15:04:05"))
	content = strings.ReplaceAll(content, "{Method}", r.Method)
	content = strings.ReplaceAll(content, "{Host}", r.Host)
	urlStr := r.Url
	parsedUrl, err := url.Parse(r.Url)
	if err == nil {
		urlStr = parsedUrl.Path
		if value := parsedUrl.Query()["api-version"]; len(value) > 0 {
			urlStr += "?api-version=" + value[0]
		}
	}
	content = strings.ReplaceAll(content, "{Url}", urlStr)
	content = strings.ReplaceAll(content, "{StatusCode}", fmt.Sprintf("%d", r.StatusCode))
	content = strings.ReplaceAll(content, "{StatusMessage}", http.StatusText(r.StatusCode))
	content = strings.ReplaceAll(content, "{RequestHeaders}", m.formatHeaders(r.Request.Headers))
	content = strings.ReplaceAll(content, "{RequestBody}", utils.JsonPretty(r.Request.Body))
	if r.Response == nil {
		content = strings.ReplaceAll(content, "{ResponseHeaders}", "")
		content = strings.ReplaceAll(content, "{ResponseBody}", "")
		return content
	}
	content = strings.ReplaceAll(content, "{ResponseHeaders}", m.formatHeaders(r.Response.Headers))
	content = strings.ReplaceAll(content, "{ResponseBody}", utils.JsonPretty(r.Response.Body))
	return content
}

func (m MarkdownFormatter) formatHeaders(headers map[string]string) string {
	var content string
	for k, v := range headers {
		content += fmt.Sprintf("| %s | %s |\n", k, v)
	}
	return content
}

const markdownTemplate = `
##### <!--- {Method} {Url} {StatusCode} -->
<details>
  <summary>
    {Time} {Method} {Host} {Url} {StatusCode}
  </summary>
<blockquote>
<details>
  <summary><strong>Request</strong></summary>

| Header | Value |
| ----- | ----- |
{RequestHeaders}

Request Body:
` + "```" + `json
{RequestBody}
` + "```" + `

</details>
<details>
  <summary><strong>Response</strong></summary>

  **Response Status: {StatusCode} {StatusMessage}**

| Header | Value |
| ----- | ----- |
{ResponseHeaders}

Response Body:
` + "```" + `json
{ResponseBody}
` + "```" + `

</details>
</blockquote>
</details>

-----

`
