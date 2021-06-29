package server

import (
	"go.uber.org/zap"
	"html/template"
)

type Application struct {
	Logger        *zap.SugaredLogger
	TemplateCache map[string]*template.Template
}
