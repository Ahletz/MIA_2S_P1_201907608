package Tools

import (
	"MIA_2S_P1_201907608/Mount"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func mkdisk(commandArray []string) {
	salida_comando += "[MSJ] COMANDO MKDISK...\\n"

	// Variables para los valores de los parametros
	val_size := 0
	val_fit := ""
	val_unit := ""
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_fit := false
	band_unit := false
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				salida_comando += "[ERR] PARAMETRO SIZE INGRESADO YA\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			// ERROR de conversion
			if err != nil {
				band_error = true
				salida_comando += "[ERR] CONVERSION ENTERO\\n"
				break
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				band_error = true
				salida_comando += "[ERR] TAMAÑO NEGATIVO\\n"
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				salida_comando += "[ERR] PARAMETRO FIT INGRESADO YA \\n"
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				salida_comando += "[ERR] VALOR FIT NO VALIDO \\n"
				band_error = true
				break
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				salida_comando += "[ERR] PARAMETRO UNIT INGRESADO YA\\n"
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				salida_comando += "[ERR] PARAMETRO UNIT NO VALIDO\\n"
				band_error = true
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				salida_comando += "[ERR] PARAMETRO PATH INGRESADO YA\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			salida_comando += "[ERR] PARAMETRO INVALIDO\\n"
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		if band_path {
			// Verifico que el parametro "Size" (Obligatorio) este ingresado
			if band_size {
				total_size := 1024
				master_boot_record := MBR{}

				// Disco -> Archivo Binario
				crear_disco(val_path)

				// Fecha
				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")

				// Copio valor al Struct
				copy(master_boot_record.fecha_creacion[:], str_fecha)

				// Numero aleatorio
				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min

				// Copio valor al Struct del numero entero
				master_boot_record.dsk_asignature = num_random //al ser un numero entero es necesario copiarse en la estructura

				// Verifico si existe el parametro "Fit" (Opcional)
				if band_fit {
					// Copio valor al Struct
					copy(master_boot_record.fit[:], val_fit)
				} else {
					// Si no especifica -> "Primer Ajuste"
					copy(master_boot_record.fit[:], "F")
				}

				// Verifico si existe el parametro "Unit" (Opcional)
				if band_unit {
					// Megabytes
					if val_unit == "M" {
						master_boot_record.tamaño = val_size * 1024 * 1024 //definimos tamaño en numeros enteros
						total_size = val_size * 1024
					} else {
						// Kilobytes
						master_boot_record.tamaño = val_size * 1024 //definimos tamaño en numeros enteros
						total_size = val_size
					}
				} else {
					// Si no especifica -> Megabytes
					master_boot_record.tamaño = val_size * 1024 * 1024 //tamaño por defecto
					total_size = val_size * 1024
				}

				// Inicializar Parcticiones
				for i := 0; i < 4; i++ {
					copy(master_boot_record.partition[i].status[:], "0")
					copy(master_boot_record.partition[i].type_[:], "0")
					copy(master_boot_record.partition[i].fit[:], "0")
					master_boot_record.partition[i].start = -1 //indica inicio de particion
					master_boot_record.partition[i].size = 0   //tamaño de la particion
					copy(master_boot_record.partition[i].name[:], "")
					master_boot_record.partition[i].correlative = 0 //inicializacion de correlativo
					copy(master_boot_record.partition[i].id[:], "")
				}

				// Convierto de entero a string
				str_total_size := strconv.Itoa(total_size)

				// Comando para definir el tamaño (Kilobytes) y llenarlo de ceros
				cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+val_path+"\" bs=1024 count="+str_total_size)
				cmd.Dir = "/"
				_, err := cmd.Output()

				// ERROR
				if err != nil {
					salida_comando += "[ERR] EJECUTAR COMANDO\\n"
				}

				// Se escriben los datos en disco

				// Apertura del archivo
				f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

				// ERROR
				if err != nil {
					salida_comando += "[ERR] ABRIR ARCHIVO\\n"
				} else {
					// Conversion de struct a bytes
					mbr_byte := struct_a_bytes(master_boot_record)

					// Escribo el mbr desde el inicio del archivos
					f.Seek(0, io.SeekStart)
					f.Write(mbr_byte)
					f.Close()

					salida_comando += "[SUCCES] DISCO CREADO\\n"
				}
			}
		}
	}

	salida_comando += "[MSJ] FINALIZADO MKDISK\\n"
}

/* RMDISK 1.0 */
func rmdisk(commandArray []string) {
	salida_comando += "[MESJ] COMANDO RMDISK\\n"

	// Variables para los valores de los parametros
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				salida_comando += "[ERR] PARAMETRO PATH INGRESAO YA\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			salida_comando += "[ERR] PARAMETRO NO VALIDO\\n"
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		if band_path {
			// Si existe el archivo binario
			_, e := os.Stat(val_path)

			if e != nil {
				// Si no existe
				if os.IsNotExist(e) {
					salida_comando += "[ERR] NO EXISTE EL DISCO PARA ELIMINAR\\n"
					band_path = false
				}
			} else {
				// Elimina el archivo
				cmd := exec.Command("/bin/sh", "-c", "rm \""+val_path+"\"")
				cmd.Dir = "/"
				_, err := cmd.Output()

				// ERROR
				if err != nil {
					salida_comando += "[ERR] EJECUTAR COMANDO\\n"
				} else {
					salida_comando += "[SUCCES] ELIMINADO\\n"
				}

				band_path = false
			}
		} else {
			salida_comando += "[ERR] PATH INGRESADO YA\\n"
		}
	}

	salida_comando += "[MSJ] RMDISK FINALIZADO\\n"

}

func fdisk(commandArray []string) {
	salida_comando += "[MENSAJE] El comando FDISK aqui inicia\\n"

	// Variables para los valores de los parametros
	val_size := 0
	val_unit := ""
	val_path := ""
	val_type := ""
	val_fit := ""
	val_name := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_unit := false
	band_path := false
	band_type := false
	band_fit := false
	band_name := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				salida_comando += "[ERROR] El parametro -size ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			// ERROR de conversion
			if err != nil {
				salida_comando += "[ERROR] Al convertir a entero\\n"
				band_error = true
				break
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				band_error = true
				salida_comando += "[ERROR] El parametro -size es negativo...\\n"
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				salida_comando += "[ERROR] El parametro -unit ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				salida_comando += "[ERROR] El Valor del parametro -unit no es valido...\\n"
				band_error = true
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				salida_comando += "[ERROR] El parametro -path ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OPCIONAL -> TYPE */
		case strings.Contains(data, "type="):
			if band_type {
				salida_comando += "[ERROR] El parametro -type ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)

			if val_type == "p" || val_type == "e" || val_type == "l" {
				// Activo la bandera del parametro
				band_type = true
			} else {
				// Parametro no valido
				salida_comando += "[ERROR] El Valor del parametro -type no es valido...\\n"
				band_error = true
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				salida_comando += "[ERROR] El parametro -fit ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				salida_comando += "[ERROR] El Valor del parametro -fit no es valido...\\n"
				band_error = true
				break
			}
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				salida_comando += "[ERROR] El parametro -name ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			salida_comando += "[ERROR] Parametro no valido...\\n"
		}
	}

	// Verifico si no hay errores
	if !band_error {
		if band_size {
			if band_path {
				if band_name {
					if band_type {
						if val_type == "p" {
							// Primaria
							crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
						} else if val_type == "e" {
							// Extendida
							crear_particion_extendia(val_path, val_name, val_size, val_fit, val_unit)
						} else {
							// Logica
							crear_particion_logica(val_path, val_name, val_size, val_fit, val_unit)
						}
					} else {
						// Si no lo indica se tomara como Primaria
						crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
					}
				} else {
					salida_comando += "[ERROR] El parametro -name no fue ingresado\\n"
				}
			} else {
				salida_comando += "[ERROR] El parametro -path no fue ingresado\\n"
			}
		} else {
			salida_comando += "[ERROR] El parametro -size no fue ingresado\\n"
		}
	}

	salida_comando += "[MENSAJE] El comando FDISK aqui finaliza\\n"
}

/* MOUNT 1.0 */
func mount(commandArray []string) {
	salida_comando += "[MENSAJE] El comando MOUNT aqui inicia\\n"

	// Variables para los valores de los parametros
	val_path := ""
	val_name := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_name := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				salida_comando += "[ERROR] El parametro -path ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				salida_comando += "[ERROR] El parametro -name ya fue ingresado...\\n"
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			salida_comando += "[ERROR] Parametro no valido...\\n"
		}
	}

	// Si no hay reptidos
	if !band_error {
		//Parametro obligatorio
		if band_path {
			if band_name {
				index_p := buscar_particion_p_e(val_path, val_name)
				// Si existe
				if index_p != -1 {
					// Apertura del archivo
					f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

					if err == nil {
						mbr_empty := MBR{}

						// Calculo del tamaño de struct en bytes
						mbr2 := struct_a_bytes(mbr_empty)
						sstruct := len(mbr2)

						// Lectrura del archivo binario desde el inicio
						lectura := make([]byte, sstruct)
						f.Seek(0, io.SeekStart)
						f.Read(lectura)

						// Conversion de bytes a struct
						master_boot_record := bytes_a_struct_mbr(lectura)

						// Colocamos la particion ocupada
						copy(master_boot_record.partition[index_p].status[:], "2")

						// Conversion de struct a bytes
						mbr_byte := struct_a_bytes(master_boot_record)

						// Se posiciona al inicio del archivo para guardar la informacion del disco
						f.Seek(0, io.SeekStart)
						f.Write(mbr_byte)
						f.Close()

						if Mount.Buscar_particion(val_path, val_name, lista_montajes) {
							salida_comando += "[ERROR] La particion ya esta montada...\\n"
						} else {
							num := Mount.Buscar_numero(val_path, lista_montajes)
							letra := Mount.Buscar_letra(val_path, lista_montajes)
							id := "30" + strconv.Itoa(num) + letra

							var n *Mount.Nodo = Mount.New_nodo(id, val_path, val_name, letra, num)
							Mount.Insertar(n, lista_montajes)
							salida_comando += "[SUCCES] Particion montada con exito!\\n"
							salida_comando += Mount.Imprimir_contenido(lista_montajes)
						}
					} else {
						salida_comando += "[ERROR] No se encuentra el disco...\\n"
					}
				} else {
					//Posiblemente logica
					index_p := buscar_particion_l(val_path, val_name)
					if index_p != -1 {
						// Apertura del archivo
						f, err := os.OpenFile(val_path, os.O_RDWR, 0660)

						if err == nil {
							ebr_empty := EBR{}

							// Calculo del tamaño de struct en bytes
							ebr2 := struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)

							// Lectrura del archivo binario desde el inicio
							lectura := make([]byte, sstruct)
							f.Seek(int64(index_p), io.SeekStart)
							f.Read(lectura)

							// Conversion de bytes a struct
							extended_boot_record := bytes_a_struct_ebr(lectura)

							// Conversion de struct a bytes
							mbr_byte := struct_a_bytes(extended_boot_record)

							// Se posiciona al inicio del archivo para guardar la informacion del disco
							f.Seek(int64(index_p), io.SeekStart)
							f.Write(mbr_byte)
							f.Close()

							if Mount.Buscar_particion(val_path, val_name, lista_montajes) {
								salida_comando += "[ERROR] La particion ya esta montada...\\n"
							} else {
								num := Mount.Buscar_numero(val_path, lista_montajes)
								letra := Mount.Buscar_letra(val_path, lista_montajes)
								id := "30" + strconv.Itoa(num) + letra

								var n *Mount.Nodo = Mount.New_nodo(id, val_path, val_name, letra, num)
								Mount.Insertar(n, lista_montajes)
								salida_comando += "[SUCCES] Particion montada con exito!\\n"
								salida_comando += Mount.Imprimir_contenido(lista_montajes)
							}
						} else {
							salida_comando += "[ERROR] No se encuentra el disco...\\n"
						}

					} else {
						salida_comando += "[ERROR] No se encuentra la particion a montar...\\n"
					}
				}
			} else {
				salida_comando += "[ERROR] Parametro -name no definido...\\n"
			}
		} else {
			salida_comando += "[ERROR] Parametro -path no definido...\\n"
		}
	}

	salida_comando += "[MENSAJE] El comando MOUNT aqui finaliza\\n"
}
