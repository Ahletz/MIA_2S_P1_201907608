package Tools

import (
	"MIA_2S_P1_201907608/Mount"
	"fmt"
	"strings"
)

// variables globales
var lista_montajes *Mount.Lista = Mount.New_lista()
var salida_comando string = ""
var graphDot string = ""

// separar comandos de consola
func Separar_cmd(cmd string) {
	linea := strings.Split(cmd, "\n") //separa por cada línea

	for i := 0; i < len(linea); i++ {
		if linea[i] != "" {
			separa_comando(linea[i])
			salida_comando += "\\n"
		}
	}
}

// Separa los diferentes comando con sus parametros
func separa_comando(comando string) {

	var commandArray []string
	// Elimina los saltos de linea y retornos de carro
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)

	// Banderas para verficar comentarios
	band_comentario := false

	if strings.Contains(comando, "#") {
		// Comentario
		band_comentario = true
		salida_comando += comando + "\\n"
	} else {
		// Comando con Parametros
		commandArray = strings.Split(comando, " -")
	}

	// Ejecuta el comando leido si no es un comentario
	if !band_comentario {
		ejecutar_comando(commandArray)
	}
}

func ejecutar_comando(commandArray []string) {
	// Convierte el comando a minusculas
	data := strings.ToLower(commandArray[0])

	// Identifica el comando
	if data == "mkdisk" {
		/* MKDISK */
		mkdisk(commandArray)
		fmt.Println("MKDISK")
	} else if data == "rmdisk" {
		/* RMDISK */
		rmdisk(commandArray)
		fmt.Println("RMDISK")
	} else if data == "fdisk" {
		/* FDISK */
		fdisk(commandArray)
		fmt.Println("FDISK")
	} else if data == "mount" {
		/* MOUNT */
		mount(commandArray)
		fmt.Println("MOUNT")
	} else if data == "rep" {
		/* REP */
		//rep(commandArray)
		fmt.Println("REP")
	} else {
		/* ERROR */
		salida_comando += "ERROR DE COMANDO \\n"
	}
}
