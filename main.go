package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func main() {
	var target string
	var filename string

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter target domain or IP: ")
	target, _ = reader.ReadString('\n')
	target = strings.TrimSpace(target)

	fmt.Print("Enter path to .txt file containing target domains or IPs (leave empty if not using file): ")
	filename, _ = reader.ReadString('\n')
	filename = strings.TrimSpace(filename)

	if target == "" && filename == "" {
		fmt.Println("Error: Target or file is mandatory")
		return
	}

	var targets []string
	if target != "" {
		targets = append(targets, target)
	}

	if filename != "" {
		fileTargets, err := loadTargetsFromFile(filename)
		if err != nil {
			fmt.Printf("Error: Failed to load targets from file: %s\n", err)
			return
		}
		targets = append(targets, fileTargets...)
	}

	for _, t := range targets {
		ips, err := net.LookupIP(t)
		if err != nil {
			fmt.Printf("Error: Failed to perform DNS lookup for %s: %s\n", t, err)
			continue
		}

		domains := make([]string, 0)
		for _, ip := range ips {
			ptrs, err := net.LookupAddr(ip.String())
			if err != nil {
				fmt.Printf("Error: Failed to perform reverse DNS lookup for IP %s: %s\n", ip, err)
				continue
			}
			for _, ptr := range ptrs {
				domain := fmt.Sprintf("%s (%s)", ptr, ip.String())
				domains = append(domains, domain)
			}
		}

		if len(domains) == 0 {
			fmt.Printf("Found 0 domain for %s\n", t)
			continue
		}

		fmt.Printf("Found %d domain(s) for %s:\n", len(domains), t)
		printDomains(domains)

		var saveToFile string
		fmt.Print("Do you want to save the results to a file? (y/n): ")
		saveToFile, _ = reader.ReadString('\n')
		saveToFile = strings.ToLower(strings.TrimSpace(saveToFile))

		if saveToFile == "y" || saveToFile == "yes" {
			var fileName string
			fmt.Print("Enter the file name (include .txt extension): ")
			fileName, _ = reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)

			err := saveResultsToFile(fileName, domains)
			if err != nil {
				fmt.Printf("Error: Failed to save results to file: %s\n", err)
				continue
			}

			fmt.Printf("Results saved to %s\n", fileName)
		}
	}
}

func printDomains(domains []string) {
	color.Cyan("██████╗ ███████╗██╗   ██╗███████╗██████╗ ███████╗███████╗    ██╗██████╗ ")
	color.Cyan("██╔══██╗██╔════╝██║   ██║██╔════╝██╔══██╗██╔════╝██╔════╝    ██║██╔══██╗")
	color.Cyan("██████╔╝█████╗  ██║   ██║█████╗  ██████╔╝███████╗█████╗█████╗██║██████╔╝")
	color.Cyan("██╔══██╗██╔══╝  ╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██╔══╝╚════╝██║██╔═══╝ ")
	color.Cyan("██║  ██║███████╗ ╚████╔╝ ███████╗██║  ██║███████║███████╗    ██║██║   ")  
	color.Cyan("╚═╝  ╚═╝╚══════╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝    ╚═╝╚═╝ ")  
	color.Cyan("###################### Sicarios Reverse IP Scanner #########################")


	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "IP Address"})

	for _, domain := range domains {
		parts := strings.Split(domain, " ")
		ip := parts[len(parts)-1]
		domain := strings.Join(parts[:len(parts)-1], " ")
		table.Append([]string{domain, ip})
	}

	table.Render()
}

func saveResultsToFile(fileName string, domains []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, domain := range domains {
		_, err := file.WriteString(domain + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func loadTargetsFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	targets := make([]string, 0)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}
