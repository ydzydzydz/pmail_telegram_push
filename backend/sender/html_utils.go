package sender

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/microcosm-cc/bluemonday"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

var (
	languageClassRegex = regexp.MustCompile(`^language-[\w-]+$`)
	leadingWhitespace  = regexp.MustCompile(`[\t\x20]*\n`)
	trailingWhitespace = regexp.MustCompile(`\n[\t\x20]*`)
	multipleNewlines   = regexp.MustCompile(`\n{2,}`)
	multipleSpaces     = regexp.MustCompile(`[\t\x20]{2,}`)
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
	p.AllowAttrs("class").Matching(languageClassRegex).OnElements("code")
	// <blockquote>Block quotation started\nBlock quotation continued\nThe last line of the block quotation</blockquote>
	p.AllowElements("blockquote")
	// 移除标签时添加空格，解决 a 标签粘在一起
	p.AddSpaceWhenStrippingTag(true)
	return p.Sanitize(content)
}

// removeExtraSpace 移除多余空格
// 删除标签时会替换为空格，多个连续空格影响显示效果
func removeExtraSpace(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = leadingWhitespace.ReplaceAllString(content, "\n")
	content = trailingWhitespace.ReplaceAllString(content, "\n")
	content = multipleNewlines.ReplaceAllString(content, "\n")
	content = multipleSpaces.ReplaceAllString(content, " ")
	return content
}

// getSubjectText 获取主题文本
func getSubjectText(email *parsemail.Email) string {
	if len(email.Subject) <= 0 {
		return ""
	}
	return fmt.Sprintf("🔖 主题：<b>%s</b>\n", email.Subject)
}

// getFromText 获取发件人文本
func getFromText(email *parsemail.Email) string {
	if len(email.From.EmailAddress) <= 0 {
		return ""
	}
	return fmt.Sprintf("📤 发件：&#60;%s&#62;\n", email.From.EmailAddress)
}

// getToText 获取收件人文本
func getToText(email *parsemail.Email) string {
	if len(email.To) <= 0 {
		return ""
	}
	text := "📥 收件："
	for _, to := range email.To {
		text += fmt.Sprintf("&#60;%s&#62; ", to.EmailAddress)
	}
	text += "\n"
	return text
}

// getCcText 获取抄送人文本
func getCcText(email *parsemail.Email) string {
	if len(email.Cc) <= 0 {
		return ""
	}
	text := "📋 抄送："
	for _, cc := range email.Cc {
		text += fmt.Sprintf("&#60;%s&#62; ", cc.EmailAddress)
	}
	text += "\n"
	return text
}

// getBccText 获取密送人文本
func getBccText(email *parsemail.Email) string {
	if len(email.Bcc) <= 0 {
		return ""
	}
	text := "🕵️ 密送："
	for _, bcc := range email.Bcc {
		text += fmt.Sprintf("&#60;%s&#62; ", bcc.EmailAddress)
	}
	text += "\n"
	return text
}

// getAttachmentsText 获取附件文本
func getAttachmentsText(email *parsemail.Email) string {
	if len(email.Attachments) <= 0 {
		return ""
	}
	return fmt.Sprintf("📎 附件：%d 个\n", len(email.Attachments))
}

// getContentText 获取邮件内容文本
func getContentText(email *parsemail.Email, setting *model.PluginTelegramPushSettingModel) string {
	if !setting.ShowContent {
		return ""
	}
	if len(email.Text) > 0 {
		return string(email.Text)
	}
	if len(email.HTML) > 0 {
		return removeHTMLTag(string(email.HTML))
	}
	return ""
}

// getSpoilerText 获取spoiler文本
func getSpoilerText(text string, setting *model.PluginTelegramPushSettingModel) string {
	if !setting.SpoilerContent {
		return text
	}
	return fmt.Sprintf("<tg-spoiler>%s</tg-spoiler>", text)
}
