package main

//go:generate go run main.go

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func main() {
	buildContract("erc165", "0.4.20", 200)
	buildContract("erc20", "0.8.0", 200)
	buildContract("erc721", "0.4.24", 200)
	buildContract("erc1155", "0.5.9", 200)
}

func buildContract(pkg, solcVer string, optimize int) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	contractDir := path.Dir(cwd)

	outputDir, err := os.MkdirTemp("", "contract-building-*")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(outputDir)

	// build solidity code

	if solcVer != "" {
		solcVer = fmt.Sprintf(":%s", solcVer)
	}

	args := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/contract.sol:ro", path.Join(path.Join(contractDir, pkg), "contract.sol")),
		"-v", fmt.Sprintf("%s:/output", outputDir),
		fmt.Sprintf("ethereum/solc%s", solcVer),
		"--combined-json", "abi,bin",
		"-o", "/output",
	}

	if optimize > 0 {
		args = append(args, "--optimize")
		args = append(args, "--optimize-runs", strconv.Itoa(optimize))
	}

	args = append(args, "/contract.sol")

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("building contract source code:\n$ docker %s\n\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// generate go code

	args = []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/combined.json:ro", path.Join(outputDir, "combined.json")),
		"-v", fmt.Sprintf("%s:/output", path.Join(outputDir, "go")),
		"ethereum/client-go:alltools-latest",
		"abigen",
		"--combined-json", "/combined.json",
		"--pkg", pkg,
		"--out", "/output/contract.go",
	}

	cmd = exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("generating contract go code:\n$ docker %s\n\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// make the go file readable

	args = []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%s:/output", path.Join(outputDir, "go")),
		"busybox",
		"chmod", "o+r", "/output/contract.go",
	}

	cmd = exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("setting permissions:\n$ docker %s\n\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	src, err := os.Open(path.Join(outputDir, "go", "contract.go"))
	if err != nil {
		panic(err)
	}

	defer src.Close()

	dst, err := os.OpenFile(path.Join(contractDir, pkg, "contract.go"), os.O_CREATE|os.O_TRUNC|os.O_SYNC|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		panic(err)
	}

	fmt.Printf("contract/%s/contract.go: file has been generated successfully.\n", pkg)
}
