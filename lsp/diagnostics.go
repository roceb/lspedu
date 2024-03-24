package lsp

type PublishDiagnosticsNotification struct {
	Notification
	Params PublishDiagnosticsParams `json:"params"`
}

type PublishDiagnosticsParams struct {
	URI         string        `json:"uri"`
	Diagnostics []Diagnostics `json:"diagnostics"`
}

type Diagnostics struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}
