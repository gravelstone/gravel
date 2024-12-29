# Gravel - Telegram Package for Go

A Go package for interacting with the Telegram Bot API, providing a client wrapper, data models, and utility functions to simplify bot development.

This package is designed to be used as part of a larger project or as a standalone module to help integrate Telegram bots into Go applications.

## Overview

This package includes three main components:

1. **Client**: A wrapper around the Telegram Bot API to interact with Telegram servers.
2. **Types**: Data Types used for parsing and handling Telegram updates.
3. **Utils**: Utility functions for commonly needed operations like message formatting.
4. **Gravel**: Telegram Bot Client and configuration of Gravel.

## Installation

You can install the package by using the following command:

```bash
go get github.com/gravelstone/gravel
```

## Documentation

1. Client
   The Gravel client provides methods to interact with the Telegram API. Initialize it using:

```go
client := gravel.NewGravel(token, debug)
```

Parameters: \
token: The bot token provided by BotFather. \
debug: Enable or disable debug logs (true/false).\
2. Fetching Updates
Receive messages and commands from users using:

```go
updates, err := client.GetUpdates()
```

Returns:
A slice of Update objects containing messages, commands, or callback queries.

```go
package main

import (
"fmt"
"log"

    "github.com/gravelstone/gravel"

)

func main() {

    t := "TELEGRAM_TOKEN"

    client := gravel.NewGravel(t, false)

    for {
    	updates, err := client.GetUpdates()
    	if err != nil {
    		log.Printf("Error fetching updates: %v", err)
    		continue
    	}

    	for _, update := range updates {
    		if update.Message != nil {
    			if update.Message.IsCommand() {
    				err = client.SendMessage(update.Message.Chat.ID, "Hello! You sent: "+update.Message.Text)
    			}

    			if update.Message.Text == "ping" {
    				usr, _ := client.GetUserInfo(update.Message.Chat.ID)
    				client.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Pong! Hello, %v!", usr.FirstName))
    			}

    			if err != nil {
    				log.Printf("Error sending message: %v", err)
    			}
    		}
    	}
    }

}

```
