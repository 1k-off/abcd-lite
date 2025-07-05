package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"text/template"
)

func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random token")
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}

func RenderIISConfigTemplate(templatePath, outputPath, port string) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, map[string]string{"Port": port})
}
