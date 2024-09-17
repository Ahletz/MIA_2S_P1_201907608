package Tools

import (
	"io"
	"os"
	"strconv"
	"strings"
)

// Crea la Particion Extendida
func crear_particion_extendia(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := MBR{}

	// Verifico si tiene Ajuste
	if fit != "" {
		aux_fit = fit
	} else {
		// Por default es Peor ajuste
		aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		salida_comando += "[ERROR] Al abrir el archivo\\n"
	} else {
		// Procede a leer el archivo
		band_particion := false
		band_extendida := false
		num_particion := 0

		// Calculo del tamaño de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		lectura := make([]byte, sstruct)
		f.Seek(0, io.SeekStart)
		f.Read(lectura)

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// Si el disco esta creado
		if master_boot_record.tamaño > 0 {
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.partition[i].type_[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					band_extendida = true
					break
				}
			}

			// Si no es extendida
			if !band_extendida {
				s_part_start := 0
				s_part_status := ""
				s_part_size := 0

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Antes de comparar limpio la cadena
					s_part_start = master_boot_record.partition[i].start

					s_part_status = string(master_boot_record.partition[i].status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					s_part_size = master_boot_record.partition[i].size

					// Verifica si existe una particion disponible
					if s_part_start == -1 || (s_part_status == "1" && s_part_size >= size_bytes) {
						band_particion = true
						num_particion = i
						break
					}
				}

				// Si hay una particion
				if band_particion {
					espacio_usado := 0

					// Recorro las 4 particiones
					for i := 0; i < 4; i++ {
						s_part_status = string(master_boot_record.partition[i].status[:])
						s_part_status = strings.Trim(s_part_status, "\x00")

						if s_part_status != "1" {
							// Obtengo el espacio utilizado
							s_part_size = master_boot_record.partition[i].size

							// Le sumo el valor al espacio
							espacio_usado += s_part_size
						}
					}

					/* Tamaño del disco */

					// Obtengo el tamaño del disco
					s_tamaño_disco := master_boot_record.tamaño

					espacio_disponible := s_tamaño_disco - espacio_usado

					salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(espacio_disponible) + " Bytes\\n"
					salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

					// Verifico que haya espacio suficiente
					if espacio_disponible >= size_bytes {
						// Verifico si no existe una particion con ese nombre
						if !existe_particion(direccion, nombre) {
							// Antes de comparar limpio la cadena
							s_dsk_fit := string(master_boot_record.fit[:])
							s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

							/*  Primer Ajuste  */
							if s_dsk_fit == "f" {
								copy(master_boot_record.partition[num_particion].type_[:], "e")
								copy(master_boot_record.partition[num_particion].fit[:], aux_fit)

								// Si esta iniciando
								if num_particion == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									master_boot_record.partition[num_particion].start = len(mbr_empty_byte)
								} else {
									// Obtengo el inicio de la particion anterior
									s_part_start_ant := master_boot_record.partition[num_particion-1].start

									// Obtengo el tamaño de la particion anterior
									s_part_size_ant := master_boot_record.partition[num_particion-1].size

									master_boot_record.partition[num_particion].start = s_part_start_ant + s_part_size_ant
								}

								master_boot_record.partition[num_particion].size = size_bytes
								copy(master_boot_record.partition[num_particion].status[:], "0")
								copy(master_boot_record.partition[num_particion].name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := struct_a_bytes(master_boot_record)

								// Escribe en la posicion inicial del archivo
								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								// Obtengo el tamaño de la particion
								s_part_start = master_boot_record.partition[num_particion].start

								// Se posiciona en el inicio de la particion
								f.Seek(int64(s_part_start), io.SeekStart)

								extended_boot_record := EBR{}
								copy(extended_boot_record.fit[:], aux_fit)
								extended_boot_record.start = s_part_start
								extended_boot_record.s = 0
								extended_boot_record.next = -1
								copy(extended_boot_record.name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo llena de unos
								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}

								salida_comando += "[SUCCES] La Particion extendida fue creada con exito!\\n"
							} else if s_dsk_fit == "b" {
								/*  Mejor Ajuste  */
								best_index := num_particion

								// Variables para conversiones
								s_part_start_act := 0
								s_part_status_act := ""
								s_part_size_act := 0
								s_part_start_best := 0
								s_part_start_best_ant := 0
								s_part_size_best := 0
								s_part_size_best_ant := 0

								for i := 0; i < 4; i++ {
									// Obtengo el inicio de la particion actual
									s_part_start_act = master_boot_record.partition[i].start

									// Obtengo el size de la particion actual
									s_part_status_act = string(master_boot_record.partition[i].status[:])
									// Le quito los caracteres null
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									// Obtengo la posicion de la particion actual
									s_part_size_act = master_boot_record.partition[i].size
									if s_part_start_act == -1 || (s_part_status_act == "1" && s_part_size_act >= size_bytes) {
										if i != num_particion {
											// Obtengo el tamaño de la particion del mejor indice
											s_part_size_best = master_boot_record.partition[best_index].size

											// Obtengo la posicion de la particion actual
											s_part_size_act = master_boot_record.partition[i].size

											if s_part_size_best > s_part_size_act {
												best_index = i
												break
											}
										}
									}
								}

								// Extendida
								copy(master_boot_record.partition[best_index].type_[:], "e")
								copy(master_boot_record.partition[best_index].fit[:], aux_fit)

								// Si esta iniciando
								if best_index == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									master_boot_record.partition[best_index].start = len(mbr_empty_byte)
								} else {
									// Obtengo el inicio de la particion actual
									s_part_start_best_ant = master_boot_record.partition[best_index-1].start

									// Obtengo el inicio de la particion actual
									s_part_size_best_ant = master_boot_record.partition[best_index-1].size

									master_boot_record.partition[best_index].start = s_part_start_best_ant + s_part_size_best_ant
								}

								master_boot_record.partition[best_index].size = size_bytes
								copy(master_boot_record.partition[best_index].status[:], "0")
								copy(master_boot_record.partition[best_index].name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := struct_a_bytes(master_boot_record)

								// Se escribe al inicio del archivo
								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								// Obtengo el inicio de la particion best
								s_part_start_best = master_boot_record.partition[best_index].start

								// Se posiciona en el inicio de la particion
								f.Seek(int64(s_part_start_best), io.SeekStart)

								extended_boot_record := EBR{}
								copy(extended_boot_record.fit[:], aux_fit)
								extended_boot_record.start = s_part_start_best
								extended_boot_record.s = 0
								extended_boot_record.next = -1
								copy(extended_boot_record.name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo llena de unos
								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}

								salida_comando += "[SUCCES] La Particion extendida fue creada con exito!\\n"
							} else {
								/*  Peor ajuste  */
								worst_index := num_particion

								// Variables para conversiones
								s_part_start_act := 0
								s_part_status_act := ""
								s_part_size_act := 0
								s_part_start_worst := 0
								s_part_start_worst_ant := 0
								s_part_size_worst := 0
								s_part_size_worst_ant := 0

								for i := 0; i < 4; i++ {
									// Obtengo el inicio de la particion actual
									s_part_start_act = master_boot_record.partition[i].start

									// Obtengo el size de la particion actual
									s_part_status_act = string(master_boot_record.partition[i].status[:])
									// Le quito los caracteres null
									s_part_status_act = strings.Trim(s_part_status_act, "\x00")

									// Obtengo la posicion de la particion actual
									s_part_size_act = master_boot_record.partition[i].size

									if s_part_start_act == -1 || (s_part_status_act == "1" && s_part_size_act >= size_bytes) {
										if i != num_particion {
											// Obtengo el tamaño de la particion del mejor indice
											s_part_size_worst = master_boot_record.partition[worst_index].size

											// Obtengo la posicion de la particion actual
											s_part_size_act = master_boot_record.partition[i].size
											if s_part_size_worst < s_part_size_act {
												worst_index = i
												break
											}
										}
									}
								}

								// Particiones Extendidas
								copy(master_boot_record.partition[worst_index].type_[:], "e")
								copy(master_boot_record.partition[worst_index].fit[:], aux_fit)

								// Se esta iniciando
								if worst_index == 0 {
									// Guardo el inicio de la particion y dejo un espacio de separacion
									mbr_empty_byte := struct_a_bytes(mbr_empty)
									master_boot_record.partition[worst_index].start = len(mbr_empty_byte)
								} else {
									// Obtengo el inicio de la particion actual
									s_part_start_worst_ant = master_boot_record.partition[worst_index-1].start

									// Obtengo el inicio de la particion actual
									s_part_size_worst_ant = master_boot_record.partition[worst_index-1].size

									master_boot_record.partition[worst_index].start = s_part_start_worst_ant + s_part_size_worst_ant
								}

								master_boot_record.partition[worst_index].size = size_bytes
								copy(master_boot_record.partition[worst_index].status[:], "0")
								copy(master_boot_record.partition[worst_index].name[:], nombre)

								// Se guarda de nuevo el MBR

								// Conversion de struct a bytes
								mbr_byte := struct_a_bytes(master_boot_record)

								// Se escribe desde el inicio del archivo
								f.Seek(0, io.SeekStart)
								f.Write(mbr_byte)

								// Obtengo el inicio de la particion best
								s_part_start_worst = master_boot_record.partition[worst_index].start

								// Se posiciona en el inicio de la particion
								f.Seek(int64(s_part_start_worst), io.SeekStart)

								extended_boot_record := EBR{}
								copy(extended_boot_record.fit[:], aux_fit)
								extended_boot_record.start = s_part_start_worst
								extended_boot_record.s = 0
								extended_boot_record.next = -1
								copy(extended_boot_record.name[:], "")
								ebr_byte := struct_a_bytes(extended_boot_record)
								f.Write(ebr_byte)

								// Lo llena de unos
								for i := 0; i < (size_bytes - len(ebr_byte)); i++ {
									f.Write([]byte{1})
								}

								salida_comando += "[SUCCES] La Particion extendida fue creada con exito!\\n"
							}
						} else {
							salida_comando += "[ERROR] Ya existe una particion creada con ese nombre...\\n"
						}
					} else {
						salida_comando += "[ERROR] La particion que desea crear excede el espacio disponible...\\n"
					}
				} else {
					salida_comando += "[ERROR] La suma de particiones primarias y extendidas no debe exceder de 4 particiones...\\n"
					salida_comando += "[MENSAJE] Se recomienda eliminar alguna particion para poder crear otra particion primaria o extendida\\n"
				}
			} else {
				salida_comando += "[ERROR] Solo puede haber una particion extendida por disco...\\n"
			}
		} else {
			salida_comando += "[ERROR] el disco se encuentra vacio...\\n"
		}
		f.Close()
	}
}
