package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var startIP string
	fmt.Print("Введите начальный IP адрес: ")
	fmt.Scanln(&startIP)

	ip := net.ParseIP(startIP)      
	mask := net.CIDRMask(24, 32)    
	activeIPs := make([]string, 0) 

	if ip == nil {
		fmt.Println("Некорректный IP адрес")
		return
	}

	
	ipNet := &net.IPNet{
		IP:   ip,
		Mask: mask,
	}

	
	for ip := ip.Mask(mask); ipNet.Contains(ip); incIP(ip) {
		if isIPReachable(ip) {
			activeIP := ip.String()
			activeIPs = append(activeIPs, activeIP)
			fmt.Println("[+] Host Live:", activeIP)
		}
	}

	// Создаем новый файл CSV
	file, err := os.Create("active_ips.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	
	writer := csv.NewWriter(file)
	defer writer.Flush()

	
	for _, ip := range activeIPs {
		err := writer.Write([]string{ip})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Активные IP адреса сохранены в файл active_ips.csv")
}


func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}


func isIPReachable(ip net.IP) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", ip.String()), 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}
