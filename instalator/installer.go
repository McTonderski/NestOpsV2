package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	releaseURL = "https://github.com/youruser/homelab/releases/latest/download/homelab.zip"
	installDir = "./homelab" // Można zmienić na /opt/homelab lub %ProgramFiles%
)

func main() {
	fmt.Println("🔍 Sprawdzanie Docker...")

	if !isDockerInstalled() {
		fmt.Println("🚧 Docker nie znaleziony. Instaluję...")
		if err := installDocker(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Błąd instalacji Dockera: %v\n", err)
			return
		}
	} else {
		fmt.Println("✅ Docker jest już zainstalowany.")
	}

	fmt.Println("📦 Pobieranie aplikacji...")
	if err := downloadAndExtract(releaseURL, installDir); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Błąd podczas pobierania/rozpakowywania: %v\n", err)
		return
	}

	fmt.Println("🚀 Uruchamianie aplikacji...")
	if err := runApp(installDir); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Nie udało się uruchomić aplikacji: %v\n", err)
	}
}

func isDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func installDocker() error {
	switch runtime.GOOS {
	case "linux":
		script := `#!/bin/bash
			set -e
			curl -fsSL https://get.docker.com -o get-docker.sh
			sh get-docker.sh
			rm get-docker.sh
		`
		tmp := "/tmp/install_docker.sh"
		os.WriteFile(tmp, []byte(script), 0755)
		cmd := exec.Command("bash", tmp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	case "darwin":
		fmt.Println("🧑‍💻 Proszę zainstalować Docker Desktop ręcznie z https://www.docker.com/products/docker-desktop")
		return nil
	case "windows":
		fmt.Println("🧑‍💻 Proszę zainstalować Docker Desktop ręcznie z https://www.docker.com/products/docker-desktop")
		return nil
	default:
		return fmt.Errorf("Nieobsługiwany system: %s", runtime.GOOS)
	}
}
func downloadAndExtract(url, target string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmp := "app.zip"
	out, _ := os.Create(tmp)
	io.Copy(out, resp.Body)
	out.Close()

	return unzip(tmp, target)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		outFile, _ := os.Create(fpath)
		inFile, _ := f.Open()
		io.Copy(outFile, inFile)
		outFile.Close()
		inFile.Close()

		// Jeśli to binarka, nadaj prawa wykonania
		if f.Name == "homelab" || f.Name == "homelab.exe" {
			os.Chmod(fpath, 0755)
		}
	}
	return nil
}

func runApp(dir string) error {
	binary := filepath.Join(dir, "homelab")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	cmd := exec.Command(binary)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
