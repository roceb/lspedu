package analysis

import (
	"fmt"
	"log"
	"strings"

	"roceb/lspedu/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostics {
	diagnostics := []lsp.Diagnostics{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostics{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "LSPEDU",
				Message:  "Not the best editor",
			})
		}
	}
	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostics {
	s.Documents[uri] = text
	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostics {
	s.Documents[uri] = text
	return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]
	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			Uri: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(log *log.Logger, id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Remove VS Code",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})
			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "V* C***",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor VS Code",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}
	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
	log.Printf("this is the response %+v", response)
	return response
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim user",
			Detail:        "My favorite editor",
			Documentation: "Use this to tell everyone why you love neovim",
		},
	}
	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
	return response
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
