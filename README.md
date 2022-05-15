# eebus-go-cem

The goal is to provide a basic EEBUS CEM implementation

## Introduction

This repository contains:

- Uses the eebus-go EEBUS stack at <https://github.com/DerAndereAndi/eebus-go>
- Initially working to support on EVSE and EV related use cases, meaning to work with an EVSE that supports EEBUS
- ... work in progress

## Usage

```sh
Usage: go run cmd/main.go <serverport> <remoteski> <certfile> <keyfile>
```

Example certificate and key files are located in the keys folder

### Explanation

The remoteski is from the eebus service to connect to.
If no certfile or keyfile are provided, they are generated and printed in the console so they can be saved in a file and later used again. The local SKI is also printed.
