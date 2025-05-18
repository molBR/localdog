package statsd

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"strings"

	"openlocaldog/storage"
)

func StartUDPListener(port int) {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to start UDP listener: %v", err)
	}
	defer conn.Close()

	log.Printf("Listening for DogStatsD metrics on UDP %s", addr)

	buf := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading UDP: %v", err)
			continue
		}
		go parseMetricLine(string(buf[:n]))
	}
}
func parseMetricLine(line string) {
	scanner := bufio.NewScanner(strings.NewReader(line))
	for scanner.Scan() {
		text := scanner.Text()

		// Divide apenas na primeira ocorrência de ':'
		colonIndex := strings.Index(text, ":")
		if colonIndex == -1 {
			continue
		}

		name := text[:colonIndex]
		valuePart := text[colonIndex+1:]

		valFields := strings.Split(valuePart, "|")
		if len(valFields) < 2 {
			continue
		}

		value := valFields[0]
		typeMetric := valFields[1]

		tags := []string{}
		if len(valFields) > 2 && strings.HasPrefix(valFields[2], "#") {
			tags = strings.Split(strings.TrimPrefix(valFields[2], "#"), ",")
		}

		storage.AddMetric(name, value, typeMetric, tags)
	}
}
