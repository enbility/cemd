# cemd

[![Build Status](https://github.com/enbility/cemd/actions/workflows/default.yml/badge.svg?branch=dev)](https://github.com/enbility/cemd/actions/workflows/default.yml/badge.svg?branch=dev)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/enbility/cemd)
[![Coverage Status](https://coveralls.io/repos/github/enbility/cemd/badge.svg?branch=dev)](https://coveralls.io/github/enbility/cemd?branch=dev)
[![Go report](https://goreportcard.com/badge/github.com/enbility/cemd)](https://goreportcard.com/report/github.com/enbility/cemd)
[![CodeFactor](https://www.codefactor.io/repository/github/enbility/cemd/badge)](https://www.codefactor.io/repository/github/enbility/cemd)

The goal is to provide an EEBUS CEM implementation

## Introduction

This library provides a foundation to implement energy management solutions using the [eebus-go](https://github.com/enbility/eebus-go) library. It is designed to be included either directly into go projects, or it will be able to run as a daemon for other systems interact with (to be implemented).

## Packages

- `api`: API interface definitions
- `cem`: Central CEM implementation which needs to be used by a HEMS implementation
- `cmd`: Example project
- `uccevc`: Use Case Coordinated EV Charging V1.0.1
- `ucevcc`: Use Case EV Commissioning and Configuration V1.0.1
- `ucevcem`: Use Case EV Charging Electricity Measurement V1.0.1
- `ucevsecc`: Use Case EVSE Commissioning and Configuration V1.0.1
- `ucevsoc`: Use Case EV State Of Charge V1.0.0 RC1
- `uclpc`: Use Case Limitation of Power Consumption V1.0.0 as a Energy Guard
- `uclpcserver`: Use Case Limitation of Power Consumption V1.0.0 as a Controllable System
- `ucmgcp`: Use Case Monitoring of Grid Connection Point V1.0.0
- `ucmpc`: Use Case Monitoring of Power Consumption V1.0.0 as a Monitoring Appliance
- `ucopev`: Use Case Overload Protection by EV Charging Current Curtailment V1.0.1b
- `ucoscev`: Use Case Optimization of Self Consumption During EV Charging V1.0.1b
- `ucvabd`: Use Case Visualization of Aggregated Battery Data V1.0.0 RC1 as a Visualization Appliance
- `ucvapd`: Use Case Visualization of Aggregated Photovoltaic Data V1.0.0 RC1 as a Visualization Appliance
- `util`: various internal helpers

## Usage

Run the following command to see all the options:

```sh
Usage: go run cmd/main.go
```

Example certificate and key files are located in the keys folder. If no certificate and key are provided in the options, new ones will be generated in the current folder.

### Explanation

The remoteski is from the eebus service to connect to.
If no certfile or keyfile are provided, they are generated and printed in the console so they can be saved in a file and later used again. The local SKI is also printed.
