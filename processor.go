package wits

import "net/http"

// ServerPostProcessor Server的扩展
type ServerPostProcessor func(*http.Server)

func (e *Engine) RegistryServerPostProcessor(processor ServerPostProcessor) {
	e.serverPostProcessor = append(e.serverPostProcessor, processor)
}
func (e *Engine) serverHandler() {
	for _, processor := range e.serverPostProcessor {
		processor(e.srv)
	}
}
