package main

import (
	"log"
	"text/template"
)

var (
	proposalTmpl *template.Template // Template for summarizing a proposal.
)

func init() {
	var err error

	/*r := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ", "\\", "")

	// template.FuncMap declaration for use throughout
	funcMap := template.FuncMap{
		"cleanText": func(m string) string {
			return r.Replace(m)
		},
	}*/

	generalTips := `Don't gloss over important details. Speak in specific, topic relevant terminology. Focus on action items, deliverables, and major discussion points. Ignore casual conversations and off-topic discussions. Include important dates, links, and quotes in your summary.`

	syntaxTip := `Mention user IDs with the "<@ID>" syntax. For example, the user ID 145386154785505280 would be formatted as <@145386154785505280>. Mention channels using the "<#ID>" syntax. For example, the channel ID 775859454780244031 would be formatted as <#775859454780244031>.`

	orderTip := `Order your summary by importance, with the most important information first.`

	// Set up proposalTmpl.
	proposalStr := `Summarize the provided governance proposal, "{{ .ProposalTitle }}" as clearly as possible with markdown formatting. ` + generalTips + " " + syntaxTip + " " + orderTip
	if proposalTmpl, err = template.New("proposalPrompt").Parse(proposalStr); err != nil {
		log.Fatalf("Error parsing proposal template: %v\n", err)
	}
}
