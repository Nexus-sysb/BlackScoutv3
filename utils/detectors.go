package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func DetectTechnologies(target string) []string {
	client := &http.Client{Timeout: 6 * time.Second}
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("User-Agent", RandomUserAgent())
	resp, err := client.Do(req)
	if err != nil {
		return []string{"Desconhecido"}
	}
	defer resp.Body.Close()

	var techs []string
	server := strings.ToLower(resp.Header.Get("Server"))
	xpb := strings.ToLower(resp.Header.Get("X-Powered-By"))
	if strings.Contains(server, "nginx") {
		techs = append(techs, "Nginx")
	}
	if strings.Contains(server, "apache") {
		techs = append(techs, "Apache")
	}
	if strings.Contains(xpb, "php") {
		techs = append(techs, "PHP")
	}
	if len(techs) == 0 {
		techs = append(techs, server, xpb)
	}
	fmt.Println("[+] Tecnologias detectadas:", techs)
	return techs
}

func BruteForceSubdomains(domain string) []string {
	subs := GetSubdomains(domain)
	if len(subs) > 0 {
		fmt.Printf("[+] Subdomínios via crt.sh: %d\n", len(subs))
		return subs
	}
	return []string{}
}

func ScanSensitiveFiles(target string) []string {
	found := FindSensitiveData(target)
	if len(found) > 0 {
		fmt.Printf("[+] Dados sensíveis encontrados: %d itens\n", len(found))
	}
	return found
}
