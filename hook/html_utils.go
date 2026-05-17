package hook

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// removeHTMLTag 移除 HTML 标签
// 保留 Telegram 支持的标签
// https://core.telegram.org/bots/api#sendmessage
func removeHTMLTag(content string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	// <b>bold</b>
	p.AllowElements("b")
	// <strong>bold</strong>
	p.AllowElements("strong")
	// <i>italic</i>
	p.AllowElements("i")
	// <i>italic</i>
	p.AllowElements("em")
	// <u>underline</u>
	p.AllowElements("u")
	// <ins>underline</ins>
	p.AllowElements("ins")
	// <s>strikethrough</s>
	p.AllowElements("s")
	// <strike>strikethrough</strike>
	p.AllowElements("strike")
	// <del>strikethrough</del>
	p.AllowElements("del")
	// <span class="tg-spoiler">spoiler</span>
	// p.AllowAttrs("class").Matching(regexp.MustCompile(`^tg-spoiler$`)).OnElements("span")
	// <a href="http://www.example.com/">inline URL</a>
	p.AllowAttrs("href").OnElements("a")
	p.RequireNoFollowOnLinks(false)
	p.AllowURLSchemes("http", "https", "tg")
	// <tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>
	// p.AllowAttrs("emoji-id").Matching(regexp.MustCompile(`^\d+$`)).OnElements("tg-emoji")
	// <code>inline fixed-width code</code>
	p.AllowElements("code")
	// <pre>pre-formatted fixed-width code block</pre>
	p.AllowElements("pre")
	// <pre><code class="language-python">pre-formatted fixed-width code block written in the Python programming language</code></pre>
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^language-[\w-]+$`)).OnElements("code")
	// <blockquote>Block quotation started\nBlock quotation continued\nThe last line of the block quotation</blockquote>
	p.AllowElements("blockquote")
	// 移除标签时添加空格，解决 a 标签粘在一起
	p.AddSpaceWhenStrippingTag(true)
	return p.Sanitize(content)
}

// removeExtraSpace 移除多余空格
// 删除标签时会替换为空格，多个连续空格影响显示效果
func removeExtraSpace(content string) string {
	// 规范换行符
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	// 去除换行两侧的空格/制表符
	// 左侧空白 + 换行
	content = regexp.MustCompile(`[\t\x20]*\n`).ReplaceAllString(content, "\n")
	// 换行 + 右侧空白
	content = regexp.MustCompile(`\n[\t\x20]*`).ReplaceAllString(content, "\n")

	// 多个连续换行替换为一个换行
	content = regexp.MustCompile(`\n{2,}`).ReplaceAllString(content, "\n")

	// 多个连续空格/制表符替换为一个空格（不影响换行）
	content = regexp.MustCompile(`[\t\x20]{2,}`).ReplaceAllString(content, " ")

	return content
}
