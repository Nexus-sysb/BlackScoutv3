package core

import (
	"blackscout/utils"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TargetURL string
	Threads   int
	DelayMs   int
	Mode      int
}

func GetConfig() Config {
	utils.ClearScreen()
	utils.Banner()

	fmt.Println(`
[1] Scan de Endpoints
[2] Scan de APIs
[3] Reconhecimento Completo
[4] Análise de JavaScript
[5] Detectar Tecnologias
[6] Verificação de Subdomínios
[7] Verificação de Arquivos Sensíveis
[0] Sair
`)
	var config Config

	for {
		val := utils.ReadInput("Escolha uma opção: ")
		mode, err := strconv.Atoi(val)
		if err == nil && mode >= 0 && mode <= 7 {
			config.Mode = mode
			break
		}
		fmt.Println("Opção inválida.")
	}

	if config.Mode == 0 {
		fmt.Println("Saindo...")
		os.Exit(0)
	}

	config.TargetURL = utils.ReadInput("URL alvo (https://site.com): ")

	for {
		t := utils.ReadInput("Threads (ex: 10): ")
		if v, err := strconv.Atoi(t); err == nil && v > 0 {
			config.Threads = v
			break
		}
	}

	for {
		d := utils.ReadInput("Delay entre requisições (ms): ")
		if v, err := strconv.Atoi(d); err == nil && v >= 0 {
			config.DelayMs = v
			break
		}
	}

	return config
}

func Run(initial Config) {
	config := initial
	for {
		start := time.Now()
		go utils.ShowLiveProgress(&utils.TotalRequests, start)

		crawler, err := NewCrawler(config.TargetURL, config.Threads, config.DelayMs)
		if err != nil {
			fmt.Println("[!] Erro ao iniciar crawler:", err)
		} else {
			switch config.Mode {
			case 1:
				results := crawler.Start(false, false)
				_ = utils.ExportResults(results, "endpoints")
			case 2:
				results := crawler.Start(true, false)
				_ = utils.ExportResults(results, "apis")
			case 3:
				results := crawler.Start(true, true)
				_ = utils.ExportResults(results, "fullscan")
			case 4:
				results := crawler.Start(false, true)
				_ = utils.ExportResults(results, "javascript")
			case 5:
				_ = utils.DetectTechnologies(config.TargetURL)
			case 6:
				_ = utils.BruteForceSubdomains(config.TargetURL)
			case 7:
				_ = utils.ScanSensitiveFiles(config.TargetURL)
			}
		}

		fmt.Println("\n[✓] Scan finalizado. Retornando ao menu...")
		time.Sleep(2 * time.Second)
		config = GetConfig()
	}
}
