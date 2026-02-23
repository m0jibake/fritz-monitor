package main

import (
	"os"
	"context"
	"fmt"
	"time"
	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/influxdata/influxdb-client-go/v2" 
)

func main() {

	influx_token := os.Getenv("INFLUX_TOKEN")

		if influx_token == "" {
			fmt.Println("Error: No INFLUX_TOKEN env variable has been set.")
			return
		}


	clientDB := influxdb2.NewClient("http://192.168.178.49:8086", influx_token)
	writeAPI := clientDB.WriteAPIBlocking("justme", "mbit-monitoring")

	// 1. Suche nach dem Service "WANCommonInterfaceConfig1"
	clients, errors, err := internetgateway1.NewWANCommonInterfaceConfig1Clients()
	
	if err != nil {
		fmt.Println("Error while reading data:", err)
    	return
	}

	if len(clients) > 0 {
		client := clients[0] // Wir nehmen die erste FRITZ!Box, die wir finden
		fmt.Println("Found a Fritz Box!")

		lastBytesSent, _ := client.GetTotalBytesReceived()
		lastBytesReceived, _ := client.GetTotalBytesSent()

		time.Sleep(time.Second * 1)
		for {
			time.Sleep(time.Second * 5)
			currentBytesSent, _ := client.GetTotalBytesSent()
			currentBytesReceived, _ := client.GetTotalBytesReceived()
			
			diffSent := currentBytesSent - lastBytesSent
			diffReceived := currentBytesReceived - lastBytesReceived

			mbitSent := (float64(diffSent) * 8) / (1000000 * 5)
			mbitReceived := (float64(diffReceived) * 8) / (1000000 * 5)

			fmt.Printf("Current Download Speed: %.2f Mbit/s\n", mbitSent)
			fmt.Printf("Current Upload Speed: %.2f Mbit/s\n", mbitReceived)

			p := influxdb2.NewPoint("stat", 
			map[string]string{"source": "fritzbox"},
			  map[string]interface{}{"sent": mbitSent, "received": mbitReceived},
  			time.Now())

			err = writeAPI.WritePoint(context.Background(), p)
			if err != nil {
				fmt.Printf("InfluxDB Error: %v\n", err)
			}
			
			lastBytesSent = currentBytesSent
			lastBytesReceived = currentBytesReceived
		}

    
	} else {
		fmt.Println("Keine FRITZ!Box gefunden. Fehler:", errors)
	}
}