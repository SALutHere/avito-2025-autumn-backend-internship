package main

import "github.com/SALutHere/avito-2025-autumn-backend-internship/internal/app"

const configPath = "internal/config/config.yaml"

func main() {
	app.Run(configPath)
}
