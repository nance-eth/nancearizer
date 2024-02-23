package main

import (
	"log"
	"text/template"
)

var (
	proposalSystemPrompt string
	proposalUserTmpl     *template.Template // User prompt template for summarizing a proposal.
)

func init() {
	var err error

	generalTips := `Don't gloss over important details. Speak in specific, topic-relevant terminology. Focus on action items, deliverables, and major discussion points. Order your summary by importance, with the most important information first.`

	// syntaxTip := `Mention user IDs with the "<@ID>" syntax. For example, the user ID 145386154785505280 would be formatted as <@145386154785505280>. Mention channels using the "<#ID>" syntax. For example, the channel ID 775859454780244031 would be formatted as <#775859454780244031>.`

	proposalSystemPrompt = `Briefly summarize the provided governance proposal in a single paragraph with markdown formatting. ` + generalTips

	proposalUserStr := "Here is the proposal to summarize.\n\n{{ .Body }}"
	if proposalUserTmpl, err = template.New("proposalUser").Parse(proposalUserStr); err != nil {
		log.Fatalf("Error parsing proposal template: %v\n", err)
	}
}
