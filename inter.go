package main

import (
	"bufio"
	"fmt"
	"strings"
)

// funcion para analizar comando entrantes
func Comando(comandos string) {

	// Crear un nuevo escáner para leer el string línea por línea
	scanner := bufio.NewScanner(strings.NewReader(comandos))

	fmt.Println(comandos)

	// Leer cada línea
	for scanner.Scan() {
		linea := scanner.Text()
		arreglo := strings.Split(linea, " ")
		fmt.Println("COMMANDO: ", arreglo[0])
	}

	// Comprobar si hay algún error durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer la cadena:", err)
	}
}

// Función que se exporta y puede ser llamada desde otro archivo
func Saludar(nombre string) {
	fmt.Println("Hola,", nombre)
}
