package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

func threadPrompts(threadUrl string) ([]string, error) {
	err := s.Open()
	if err != nil {
		return nil, err
	}

	// Get the thread ID from the URL
	i := strings.LastIndex(threadUrl, "/")
	if len(threadUrl) < i+1 {
		return nil, errors.New("invalid thread URL")
	}
	threadId := threadUrl[i+1:]

	r := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ", "\\", "")
	msgsToUse := make([]string, 0)
	var beforeId string

	// Get all messages from the thread
	for {
		msgs, err := s.ChannelMessages(threadId, 100, beforeId, "", "")
		if err != nil {
			return nil, err
		}

		for _, m := range msgs {
			if m.Content == "" || m.Author.Bot {
				continue
			}

			content := r.Replace(m.ContentWithMentionsReplaced())
			msgStr := fmt.Sprintf("%s: %s", m.Author.Username, content)
			msgsToUse = append(msgsToUse, msgStr)
		}

		// Update the beforeId for the next iteration
		beforeId = msgs[len(msgs)-1].ID
		if len(msgs) < 100 {
			break
		}
	}

	reverseSlice(msgsToUse)

	// Build into prompts
	prompts := make([]string, 0)
	promptPrefix := "Here are the messages to summarize:\n\n"

	var b bytes.Buffer
	b.WriteString(promptPrefix)

	tks, _ := tk.Encode(promptPrefix, true)
	tokenCounter := len(tks)

	for _, m := range msgsToUse {
		tks, _ := tk.Encode(m, false)

		if tokenCounter+len(tks) > CONTEXT_LENGTH-1_024 {
			prompts = append(prompts, b.String())
			b.Reset()
			b.WriteString(promptPrefix)
			tokenCounter = len(tks)
		}

		b.WriteString(m + "\n")
		tokenCounter += len(tks)
	}

	prompts = append(prompts, b.String())

	return prompts, nil
}

// Reverse the order of elements in a slice.
func reverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
