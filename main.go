package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SimpleMsg allows the app to extract just the bytes sent from the device.
type SimpleMsg struct {
	UplinkMessage struct {
		DecodedPayload struct {
			Bytes []byte `json:"bytes"`
		} `json:"decoded_payload"`
	} `json:"uplink_message"`
}

func main() {
	var broker string
	var username string
	var password string
	var topic string
	var help bool
	var raw bool

	flag.StringVar(&broker, "b", "127.0.0.1:1883", "Broker URL")
	flag.StringVar(&username, "u", "", "Username")
	flag.StringVar(&password, "p", "", "Password")
	flag.StringVar(&topic, "t", "#", "Topic to subscribe to")
	flag.BoolVar(&raw, "raw", false, "Dump raw message")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)

	messages := make(chan []byte)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		messages <- msg.Payload()
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to broker: %v\n", token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe(topic, byte(0), nil); token.Wait() && token.Error() != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe: %v\n", token.Error())
		os.Exit(1)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Subscribed to %s ...\n", topic)
	for {
		select {
		case <-sigs:
			client.Disconnect(250)
			fmt.Println("Disconnected")
			os.Exit(0)
		case data := <-messages:
			switch {
			case raw:
				msg, err := parseRawMessage(data)
				if err != nil {
					continue
				}
				fmt.Println("Message:")
				fmt.Println(msg)

			default:
				msg, err := parseMessage(data)
				if err != nil {
					continue
				}
				fmt.Println("Message:")
				fmt.Printf("%s\n", string(msg.UplinkMessage.DecodedPayload.Bytes))
			}
		}
	}
}

func parseMessage(data []byte) (SimpleMsg, error) {
	var dest SimpleMsg
	if err := json.Unmarshal(data, &dest); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse message: %v\n", err)
		return dest, err
	}
	return dest, nil
}

func parseRawMessage(data []byte) (string, error) {
	var dest bytes.Buffer
	if err := json.Indent(&dest, data, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse message: %v\n", err)
		return "", err
	}
	return dest.String(), nil
}
