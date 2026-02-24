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

    clients, _, err := internetgateway1.NewWANCommonInterfaceConfig1Clients()
    if err != nil {
        fmt.Println("Error while searching for Fritz!Box:", err)
        return
    }

    if len(clients) > 0 {
        client := clients[0]
        fmt.Println("Found a Fritz Box! Starting monitoring...")

        lastBytesSent, _ := client.GetTotalBytesSent()
        lastBytesReceived, _ := client.GetTotalBytesReceived()

        for {
            time.Sleep(time.Second * 5)
            currentBytesSent, _ := client.GetTotalBytesSent()
            currentBytesReceived, _ := client.GetTotalBytesReceived()
            
            if currentBytesSent >= lastBytesSent && currentBytesReceived >= lastBytesReceived {
                diffSent := currentBytesSent - lastBytesSent
                diffReceived := currentBytesReceived - lastBytesReceived

                mbitSent := (float64(diffSent) * 8) / (1000000 * 5)
                mbitReceived := (float64(diffReceived) * 8) / (1000000 * 5)

                fmt.Printf("Download: %.2f Mbit/s | Upload: %.2f Mbit/s\n", mbitReceived, mbitSent)

                p := influxdb2.NewPoint("stat", 
                    map[string]string{"source": "fritzbox"},
                    map[string]interface{}{"upload": mbitSent, "download": mbitReceived},
                    time.Now())

                if err := writeAPI.WritePoint(context.Background(), p); err != nil {
                    fmt.Printf("InfluxDB Error: %v\n", err)
                }
            } else {
                fmt.Println("Counter reset detected (Router rebooted or overflow). Skipping calculation...")
            }
            
            lastBytesSent = currentBytesSent
            lastBytesReceived = currentBytesReceived
        }
    } else {
        fmt.Println("No Fritz!Box found.")
    }
}