package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "pull":
		handlePull()
	case "push":
		handlePush()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Erreur: commande inconnue '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: gnit <command> [options]")
	fmt.Println()
	fmt.Println("Commandes disponibles:")
	fmt.Println("  pull <file>     Récupère un fichier depuis le repository")
	fmt.Println("  push <message>  Pousse les changements avec un message de commit")
	fmt.Println("  help            Affiche cette aide")
	fmt.Println()
	fmt.Println("Exemples:")
	fmt.Println("  gnit pull example.gno")
	fmt.Println("  gnit push \"Mon message de commit\"")
}

func handlePull() {
	if len(os.Args) < 3 {
		fmt.Println("Erreur: nom de fichier requis pour pull")
		fmt.Println("Usage: gnit pull <file>")
		os.Exit(1)
	}

	filename := os.Args[2]

	fmt.Printf("Pulling '%s'...\n", filename)

	cmd := exec.Command("gnokey", "query", "vm/qeval",
		"-data", fmt.Sprintf("gno.land/r/example.Pull(\"%s\")", filename),
		"-remote", "tcp://127.0.0.1:26657")

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}

	content, err := parseHexOutput(string(output))
	if err != nil {
		fmt.Printf("Erreur lors du parsing de la réponse: %v\n", err)
		os.Exit(1)
	}

	if len(content) == 0 {
		fmt.Printf("Fichier '%s' non trouvé ou vide\n", filename)
		os.Exit(1)
	}

	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Printf("Erreur lors de l'écriture du fichier: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Fichier '%s' récupéré avec succès (%d bytes)\n", filename, len(content))
}

func handlePush() {
	if len(os.Args) < 3 {
		fmt.Println("Erreur: message de commit requis pour push")
		fmt.Println("Usage: gnit push \"<message>\"")
		os.Exit(1)
	}

	message := strings.Join(os.Args[2:], " ")
	message = strings.Trim(message, "\"")

	fmt.Printf("Pushing avec le message: '%s'...\n", message)

	files := make(map[string][]byte)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Printf("Avertissement: impossible de lire %s: %v\n", path, readErr)
				return nil
			}
			files[path] = content
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Erreur lors de la lecture des fichiers: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("Aucun fichier .gno trouvé à pousser")
		os.Exit(1)
	}

	fmt.Printf("Fichiers à pousser: %d\n", len(files))
	for filename := range files {
		fmt.Printf("  - %s\n", filename)
	}

	var filesData strings.Builder
	for filename, content := range files {
		filesData.WriteString(filename)
		filesData.WriteString("|")
		filesData.WriteString(string(content))
		filesData.WriteString("\n")
	}

	cmd := exec.Command("gnokey", "maketx", "call",
		"-pkgpath", "gno.land/r/example",
		"-func", "Push",
		"-args", message,
		"-args", filesData.String(),
		"-gas-fee", "1000000ugnot",
		"-gas-wanted", "50000000",
		"-send", "",
		"-broadcast",
		"-chainid", "dev",
		"-remote", "tcp://127.0.0.1:26657",
		"test")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Push réussi!\n")
}

func parseHexOutput(output string) ([]byte, error) {
	lines := strings.Split(output, "\n")
	var dataLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			dataLine = line
			break
		}
	}

	if dataLine == "" {
		return nil, fmt.Errorf("ligne data: non trouvée")
	}

	if strings.Contains(dataLine, "(nil []uint8)") {
		return []byte{}, nil
	}

	re := regexp.MustCompile(`slice\[0x([0-9a-fA-F]+)\]`)
	matches := re.FindStringSubmatch(dataLine)

	if len(matches) < 2 {
		return nil, fmt.Errorf("format de sortie non reconnu: %s", dataLine)
	}

	data, err := hex.DecodeString(matches[1])
	if err != nil {
		return nil, fmt.Errorf("erreur de décodage hex: %v", err)
	}

	return data, nil
}