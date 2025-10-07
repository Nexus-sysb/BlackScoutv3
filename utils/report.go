package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"blackscout/config"
)

type Report struct {
	URL           string   `json:"url"`
	Subdomains    []string `json:"subdomains"`
	AdminPanels   []string `json:"admin_panels"`
	SQLInjection  bool     `json:"sql_injection"`
	XSSVulnerable bool     `json:"xss_vulnerable"`
	SensitiveData []string `json:"sensitive_data"`
}

// Função para salvar o relatório em formato JSON
func SaveReport(report Report) error {
	timestamp := time.Now().Format("20060102_150405")
	fileName := report.URL + "_" + timestamp + ".json"
	filePath := filepath.Join(config.OutputFolder, fileName)

	// Garante que o diretório de saída existe
	if err := os.MkdirAll(config.OutputFolder, os.ModePerm); err != nil {
		return err
	}

	// Cria e escreve o arquivo JSON
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
