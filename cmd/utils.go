package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

func RenderIISConfigTemplate(templatePath, outputPath, port string) error {
	var tmpl *template.Template
	var err error

	// Try to use embedded files first
	if embeddedData != nil {
		iisConfigFS := embeddedData.GetIISConfigFS()
		// Extract filename from path (e.g., "internal/config/iis/iis_default.config" -> "iis_default.config")
		filename := filepath.Base(templatePath)
		templateContent, err := fs.ReadFile(iisConfigFS, filename)
		if err == nil {
			tmpl, err = template.New(filename).Parse(string(templateContent))
			if err != nil {
				return err
			}
		}
	}

	// Fallback to filesystem if embedded file not found or not available
	if tmpl == nil {
		tmpl, err = template.ParseFiles(templatePath)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, map[string]string{"Port": port})
}
