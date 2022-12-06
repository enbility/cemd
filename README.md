# cemd

[![Build Status](https://github.com/enbility/cemd/actions/workflows/default.yml/badge.svg?branch=dev)](https://github.com/enbility/cemd/actions/workflows/default.yml/badge.svg?branch=dev)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/enbility/cemd)
[![Go report](https://goreportcard.com/badge/github.com/enbility/cemd)](https://goreportcard.com/report/github.com/enbility/cemd)

The goal is to provide an EEBUS CEM implementation

## Introduction

This library provides a foundation to implement energy management solutions using the [eebus-go](https://github.com/enbility/eebus-go) library. It is designed to be included either directly into go projects, or it will be able to run as a daemon for other systems interact with (to be implemented).

These EEBUS use cases are already supported:

- E-Mobility:

  - EVSE Commissioning and Configuration V1.0.1
  - EV Commissioning and Configuration V1.0.1
  - EV Charging Electricity Measurement V1.0.1
  - EV State Of Charge V1.0.0 RC1
  - Optimization of Self Consumption During EV Charging V1.0.1b
  - Overload Protection by EV Charging Current Curtailment V1.0.1b

These use cases are currently planned to be supported in the future:

- E-Mobility:

  - Coordinated EV Charging V1.0.1
  - EV Charging Summary V1.0.1

More use cases and scenarios will hopefully follow in the future as well.

## Usage

```sh
Usage: go run cmd/main.go <serverport> <remoteski> <certfile> <keyfile>
```

Example certificate and key files are located in the keys folder

### Explanation

The remoteski is from the eebus service to connect to.
If no certfile or keyfile are provided, they are generated and printed in the console so they can be saved in a file and later used again. The local SKI is also printed.
