package Tools

import (
	"io"
	"os"
	"strconv"
	"strings"
)

// Crea la Particion Logica
func crear_particion_logica(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := MBR{}
	ebr_empty := EBR{}

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
		salida_comando += "[ERROR] No existe el disco duro con ese nombre...\\n"
	} else {
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
			num_extendida := -1

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.partition[i].type_[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					num_extendida = i
					break
				}
			}

			if !existe_particion(direccion, nombre) {
				if num_extendida != -1 {
					s_part_start := master_boot_record.partition[num_extendida].start

					cont := s_part_start

					// Se posiciona en el inicio de la particion
					f.Seek(int64(cont), io.SeekStart)

					// Calculo del tamaño de struct en bytes
					ebr2 := struct_a_bytes(ebr_empty)
					sstruct := len(ebr2)

					// Lectrura del archivo binario desde el inicio
					lectura := make([]byte, sstruct)
					f.Read(lectura)

					// Conversion de bytes a struct
					extended_boot_record := bytes_a_struct_ebr(lectura)

					// Obtencion de datos
					s_part_size_ext := extended_boot_record.s

					if s_part_size_ext == 0 {
						// Obtencion de datos
						s_part_size := master_boot_record.partition[num_extendida].size

						salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(s_part_size) + " Bytes\\n"
						salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

						// Si excede el tamaño de la extendida
						if s_part_size < size_bytes {
							salida_comando += "[ERROR] La particion logica a crear excede el espacio disponible de la particion extendida...\\n"
						} else {
							copy(extended_boot_record.fit[:], aux_fit)

							// Posicion actual en el archivo
							pos_actual, _ := f.Seek(0, io.SeekCurrent)
							ebr_empty_byte := struct_a_bytes(ebr_empty)

							extended_boot_record.start = int(pos_actual) - len(ebr_empty_byte)
							extended_boot_record.s = size_bytes
							extended_boot_record.next = -1
							copy(extended_boot_record.name[:], nombre)

							// Obtencion de datos
							s_part_start := master_boot_record.partition[num_extendida].start

							// Se posiciona en el inicio de la particion
							ebr_byte := struct_a_bytes(extended_boot_record)
							f.Seek(int64(s_part_start), io.SeekStart)
							f.Write(ebr_byte)

							salida_comando += "[SUCCES] La Particion logica fue creada con exito!\\n"
						}
					} else {
						// Obtencion de datos
						s_part_size := master_boot_record.partition[num_extendida].size

						// Obtencion de datos
						s_part_start := master_boot_record.partition[num_extendida].start
						salida_comando += "[ESPACIO DISPONIBLE] " + strconv.Itoa(s_part_size+s_part_start) + " Bytes\\n"
						salida_comando += "[ESPACIO NECESARIO] " + strconv.Itoa(size_bytes) + " Bytes\\n"

						// Obtencion de datos
						s_part_next := extended_boot_record.next

						pos_actual, _ := f.Seek(0, io.SeekCurrent)

						for (s_part_next != -1) && (int(pos_actual) < (s_part_size + s_part_start)) {
							// Se posiciona en el inicio de la particion
							f.Seek(int64(s_part_next), io.SeekStart)

							// Calculo del tamaño de struct en bytes
							ebr2 := struct_a_bytes(ebr_empty)
							sstruct := len(ebr2)

							// Lectrura del archivo binario desde el inicio
							lectura := make([]byte, sstruct)
							f.Read(lectura)

							// Posicion actual
							pos_actual, _ = f.Seek(0, io.SeekCurrent)

							// Conversion de bytes a struct
							extended_boot_record = bytes_a_struct_ebr(lectura)

							if extended_boot_record.next == 0 {
								break
							}

							// Obtencion de datos
							s_part_next = extended_boot_record.next
						}

						// Obtencion de datos
						s_part_start_ext := extended_boot_record.start
						// Obtencion de datos
						s_part_size_ext := extended_boot_record.s
						// Obtencion de datos
						s_part_size_mbr := master_boot_record.partition[num_extendida].size
						// Obtencion de datos
						s_part_start_mbr := master_boot_record.partition[num_extendida].start

						espacio_necesario := s_part_start_ext + s_part_size_ext + size_bytes

						if espacio_necesario <= (s_part_size_mbr + s_part_start_mbr) {
							extended_boot_record.next = s_part_start_ext + s_part_size_ext

							// Escribo el nedxto del ultimo ebr
							pos_actual, _ = f.Seek(0, io.SeekCurrent)
							ebr_byte := struct_a_bytes(extended_boot_record)
							// Escribo el next del ultimo EBR
							f.Seek(int64(int(pos_actual)-len(ebr_byte)), io.SeekStart)
							f.Write(ebr_byte)

							// Escribo el nuevo EBR
							f.Seek(int64(s_part_start_ext+s_part_size_ext), io.SeekStart)
							copy(extended_boot_record.fit[:], aux_fit)
							// Posicion actual del archivo
							pos_actual, _ = f.Seek(0, io.SeekCurrent)
							extended_boot_record.start = int(pos_actual)
							extended_boot_record.s = size_bytes
							extended_boot_record.next = -1
							copy(extended_boot_record.name[:], nombre)
							ebr_byte = struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							salida_comando += "[SUCCES] La Particion logica fue creada con exito!\\n"
						} else {
							salida_comando += "[ERROR] La particion logica a crear excede el espacio disponible de la particion extendida...\\n"
						}
					}
				} else {
					salida_comando += "[ERROR] No se puede crear una particion logica si no hay una extendida...\\n"
				}
			} else {
				salida_comando += "[ERROR] Ya existe una particion con ese nombre...\\n"
			}
		} else {
			salida_comando += "[ERROR] el disco se encuentra vacio...\\n"
		}
		f.Close()
	}
}
