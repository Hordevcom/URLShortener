// Package main запускает multichecker со множеством анализаторов.
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"

	"github.com/Hordevcom/URLShortener/cmd/staticlint/noexit"
	"github.com/gostaticanalysis/nilerr"
)

func main() {
	var analyzers []*analysis.Analyzer

	// SA анализаторы
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	// Публичные
	analyzers = append(analyzers, nilerr.Analyzer)

	// Кастомный
	analyzers = append(analyzers, noexit.Analyzer)

	multichecker.Main(analyzers...)
}
