package lsp

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem
}
type TextDocumentItem struct {
	Uri        string `json:"uri,omitempty"`
	LanguageID string `json:"languageid,omitempty"`
	Version    int    `json:"version,omitempty"`
	Text       string `json:"text,omitempty"`
}

type DidOpenTextDocumentNotification struct {
	Notification
	Params DidOpenTextDocumentParams
}
type TextDocumentDidChangeNotification struct {
	Notification
	Params DidChangeTextDocumentParams
}
type DidChangeTextDocumentParams struct {
	TextDocument   VersionTextDocumentIdentifier    `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}
type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}
type Location struct {
	Uri   string `json:"uri"`
	Range Range  `json:"range"`
}
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}
type VersionTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}
type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}
type HoverParams struct {
	TextDocumentPositionParams
}
type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}
type HoverResult struct {
	Contents string `json:"contents"`
}
type DefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}
type DefinitionParams struct {
	TextDocumentPositionParams
}
type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}
type CodeActionRequest struct {
	Request
	Params CodeActionParams `json:"params"`
}
type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}
type CodeActionContext struct {
	// TODO: Added more context here
}
type CodeAction struct {
	Title   string         `json:"title"`
	Edit    *WorkspaceEdit `json:"edit,omitempty"`
	Command *Command       `json:"command,omitempty"`
}
type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}
type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

type TextDocumentCodeActionResponse struct {
	Response
	Result []CodeAction `json:"result,omitempty"`
}
type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}
type CompletionParams struct {
	TextDocumentPositionParams
}
type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}
type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
}
