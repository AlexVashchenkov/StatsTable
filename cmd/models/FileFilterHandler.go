package models

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

//go:embed extensions/language_extensions.json
var embeddedLanguageExtensions []byte

type LanguageConfig struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

type FileFilterHandler struct {
	next Handler
}

func (h *FileFilterHandler) Handle(context *Context) error {
	if len(context.Extensions) > 0 {
		context.Files = filterFilesByExtensions(context.Files, context.Extensions)
	}

	if len(context.Languages) > 0 {
		context.Files, _ = filterFilesByLanguages(context.Files, context.Languages)
	}

	if len(context.Exclude) > 0 {
		context.Files = filterFilesByGlobPatterns(context.Files, context.Exclude)
	}

	if len(context.RestrictTo) > 0 {
		context.Files = filterFilesByRestrictions(context.Files, context.RestrictTo)
	}

	return h.next.Handle(context)
}

func (h *FileFilterHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func filterFilesByExtensions(files []string, extensions []string) []string {
	var filtered []string

	for _, extension := range extensions {
		for _, file := range files {
			if strings.Contains(file, extension) {
				filtered = append(filtered, file)
			}
		}
	}

	return filtered
}

func filterFilesByRestrictions(files []string, restrictions []string) []string {
	var filtered []string
	for _, file := range files {
		restrictFile := true
		for _, pattern := range restrictions {
			if match, _ := filepath.Match(pattern, file); match {
				restrictFile = false
				break
			}
		}
		if !restrictFile {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func filterFilesByGlobPatterns(files []string, exclude []string) []string {
	var filtered []string
	for _, file := range files {
		excludeFile := false
		for _, pattern := range exclude {
			if match, _ := filepath.Match(pattern, file); match {
				excludeFile = true
				break
			}
		}
		if !excludeFile {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func filterFilesByLanguages(files []string, languages []string) ([]string, error) {
	languageMap, err := loadLanguageExtensions()
	if err != nil {
		return nil, fmt.Errorf("error loading language extensions: %v", err)
	}

	if len(languages) == 0 {
		return files, nil
	}

	extensionSet := make(map[string]struct{})
	for _, language := range languages {
		extensions, exists := languageMap[strings.ToLower(language)]
		if !exists {
			continue
		}
		for _, ext := range extensions {
			extensionSet[ext] = struct{}{}
		}
	}

	var filtered []string
	for _, file := range files {
		if _, ok := extensionSet[filepath.Ext(file)]; ok {
			filtered = append(filtered, file)
		}
	}

	return filtered, nil
}

func loadLanguageExtensions() (map[string][]string, error) {
	var configs []LanguageConfig
	if err := json.Unmarshal(embeddedLanguageExtensions, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode embedded language extensions: %w", err)
	}

	languageMap := make(map[string][]string)
	for _, config := range configs {
		languageMap[strings.ToLower(config.Name)] = config.Extensions
	}

	return languageMap, nil
}
