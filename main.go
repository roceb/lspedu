package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"roceb/lspedu/analysis"
	"roceb/lspedu/lsp"
	"roceb/lspedu/rpc"
)

func main() {
	logger := getLogger("/tmp/lsp.txt")
	logger.Println("Logger is running")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	writer := os.Stdout
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error: %s\n", err)
			continue
		}
		handleMessage(logger, writer, state, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, content []byte) {
	logger.Printf("Received msg with method: %s\n", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Could not unmarshal content from request Err: %s", err)
		}
		logger.Printf("Connected to: %s\n%s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error: textDocument/didOpen: %s", err)
			return
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.Uri)
		diagnostics := state.OpenDocument(request.Params.TextDocument.Uri, request.Params.TextDocument.Text)
		writeResponse(writer, lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.Uri,
				Diagnostics: diagnostics,
			},
		})
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error: textDocument/didChange: %s", err)
			return
		}
		logger.Printf("Changes: %s", request.Params.TextDocument.Uri)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.Uri, change.Text)
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.Uri,
					Diagnostics: diagnostics,
				},
			})
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error: textDocument/hover: %s", err)
			return
		}
		response := state.Hover(request.ID, request.Params.TextDocument.Uri, request.Params.Position)
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error: textDocument/definition: %s", err)
			return
		}
		response := state.Definition(request.ID, request.Params.TextDocument.Uri, request.Params.Position)
		writeResponse(writer, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error: textDocument/codeAction: %s", err)
			return
		}
		response := state.TextDocumentCodeAction(logger, request.ID, request.Params.TextDocument.Uri)
		writeResponse(writer, response)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error: textDocument/codeAction: %s", err)
			return
		}
		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.Uri)
		writeResponse(writer, response)

	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)

	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic("Bad logging file")
	}
	return log.New(logfile, "[LSP]", log.Ldate|log.Ltime|log.Lshortfile)
}
