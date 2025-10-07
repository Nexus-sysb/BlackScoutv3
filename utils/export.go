package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"blackscout/config"
)

func ExportResults(results []string, prefix string) error {
	if _, err := os.Stat(config.OutputFolder); os.IsNotExist(err) {
		os.MkdirAll(config.OutputFolder, 0755)
	}
	ts := time.Now().Format("20060102_150405")
	txt := filepath.Join(config.OutputFolder, fmt.Sprintf("%s_%s.txt", prefix, ts))
	jsonf := filepath.Join(config.OutputFolder, fmt.Sprintf("%s_%s.json", prefix, ts))

	// Save txt
	f1, err := os.Create(txt)
	if err == nil {
		defer f1.Close()
		for _, r := range results {
			_, _ = f1.WriteString(r + "\n")
		}
		fmt.Printf("\n[✓] Exportado TXT: %s\n", txt)
	} else {
		fmt.Printf("\n[!] Falha ao criar %s: %v\n", txt, err)
	}

	// Save json
	f2, err := os.Create(jsonf)
	if err == nil {
		defer f2.Close()
		enc := json.NewEncoder(f2)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
		fmt.Printf("[✓] Exportado JSON: %s\n", jsonf)
	} else {
		fmt.Printf("[!] Falha ao criar %s: %v\n", jsonf, err)
	}

	return nil
}
