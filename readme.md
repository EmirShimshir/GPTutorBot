# Telegram GPTutorBot

Status of Last Deployment:
[![Build Status](https://github.com/EmirShimshir/GPTutorBot/workflows/GPTutorBot-Build-and-Push/badge.svg?branch=main)](https://github.com/EmirShimshir/GPTutorBot/actions)

[@GPTutorBot](https://t.me/GPTutorBot?start=github) is a telegram bot that 
will help you solve any academic task using **ChatGPT** and **Tesseract 5**

# Contents
1. [Run](#Run)
2. [Implementation](#Implementation)

# Run

To run application you need:
- .env file
- SSL certificates

## Env configuration

To run API locally you should create your own **.env** file

Example `.env`:

```env
TELEGRAM_TOKEN=<your-token>
CHAT_GPT_API_KEY=<your-key>
UMONEY_WALLET=<your-wallet>
KEY_PAYMENT=<your-key>
IP=<your-host-ip>
```

## SSL certificates

To create your SSL certificates:

```
make generate-certs
```

## Local run

```
make run
```

With Docker

```
make up
```

With Telegram: [@GPTutorBot](https://t.me/GPTutorBot?start=github)

# Implementation

- Golang
- Telegram Bot API
- OpenAI API
- Tesseract 5
- Using Umoney WebHooks for payment
- Using Golang file database as a main data storage
- Creating my own subscription system
- Creating my own UTM tags collection system 
- Creating my own referral system
- Creating admin panel with database CRUD operations
- Clean architecture design
- Env based application configuration
- Run with docker-compose
- Full automated CI/CD process

## Project structure

```
.
├── .github          // GitHub Actions CI CD pipeline
├── .bin             // app binary files
├── certs            // SSL certificates
├── logs             // directory to store logs
├── config           // directory to store bot config
├── images           // directory to store bot images
├── db               // directory to store local db data
├── cmd              // entry point
│   └── bot          // main application package
├── internal
│    ├── config      // config loading utils
│    ├── domain      // all business entities and dto's
│    ├── openai      // ChatGPT API module
│    ├── repository  // database repository layer
│    ├── server      // payment WebHook server
│    ├── service     // business logic services layer
│    ├── telegram    // Telegram Bot API module
│    └── tesseract   // Tesseract 5 module
├── ci-cd            // CI/CD docker-compose.yaml file
├── deploy           // Docker building utils
└──
```
