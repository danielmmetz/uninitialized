package uninitialized

import (
	"github.com/danielmmetz/uninitialized/internal/uninitialized"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("uninitialized", New)
}

func New(settings any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

type Plugin struct{}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		uninitialized.NewAnalyzer(),
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
