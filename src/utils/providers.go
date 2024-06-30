package utils

import (
	"fmt"
	"main/config"
	"strings"
)

func BuildClassificationPrompt(return_prompt string) string {
	var serviceList []string
	if config.BaseLLM != nil {
		serviceList = append(serviceList, "* LLM - Use large language models with advanced reasoning capabilities with no online dependencies")
	}
	if config.BaseSearch != nil {
		serviceList = append(serviceList, "* Search - AI powered search engines for when the AI requires information from the web")
	}

	// Create the prompt using fmt
	prompt := fmt.Sprintf("You are a Classificiation AI, your job is to classify a given text to identify what service the text intends to invoke. The available services are \n%s\n%s", strings.Join(serviceList, "\n"), return_prompt)

	return prompt
}

func SequenceReturnPrompt(start_seq string, end_seq string) string {
	return fmt.Sprintf("Your responses are to be given prefixed by %s and suffixed by %s and only contain the name of a given service", start_seq, end_seq)
}

func FindSubstring(input, startSeq, endSeq string) (string, bool) {
	startIndex := strings.Index(input, startSeq)
	if startIndex == -1 {
		return "", false // startSeq not found
	}
	startIndex += len(startSeq) // Move to the end of startSeq

	endIndex := strings.Index(input[startIndex:], endSeq)
	if endIndex == -1 {
		return "", false // endSeq not found
	}

	return input[startIndex : startIndex+endIndex], true
}
