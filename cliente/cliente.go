package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Usa: go run cliente <rut> <correo> <ruta>")
		os.Exit(1)
	}

	rut := os.Args[1]
	correo := os.Args[2]
	ruta := os.Args[3]

	fmt.Println("RUT:", rut)
	fmt.Println("Correo:", correo)
	fmt.Println("Ruta:", ruta)
}
