package Tools

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func crear_disco(ruta string) {
	aux, err := filepath.Abs(ruta)

	// ERROR
	if err != nil {
		salida_comando += "[ERROR] Al abrir el archivo\\n"
	}

	// Crea el directiorio de forma recursiva
	cmd1 := exec.Command("/bin/sh", "-c", "echo 253097 | sudo -S mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err = cmd1.Output()

	// ERROR
	if err != nil {
		salida_comando += "[ERROR] Al ejecutar el comando\\n"
	}

	// Da los permisos al directorio
	cmd2 := exec.Command("/bin/sh", "-c", "echo 253097 | sudo -S chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err = cmd2.Output()

	// ERROR
	if err != nil {
		salida_comando += "[ERROR] Error al ejecutar el comando\\n"
	}

	// Verifica si existe la ruta para el archivo
	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			salida_comando += "[FAILURE] No se pudo crear el disco...\\n"
		}
	}
}

// Codifica de Struct a []Bytes
func struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	// ERROR
	if err != nil && err != io.EOF {
		salida_comando += "[ERROR] Al codificar de struct a bytes \n"
	}

	return buf.Bytes()
}

// Decodifica de [] Bytes a Struct
func bytes_a_struct_mbr(s []byte) MBR {
	p := MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		salida_comando += "[ERROR] Al decodificar a MBR\\n"
	}

	return p
}

// Decodifica de [] Bytes a Struct
func bytes_a_struct_ebr(s []byte) EBR {
	p := EBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		salida_comando += "[ERROR] AL decodificar a EBR\n"
	}

	return p
}

// Verifica si el nombre de la particion esta disponible
func existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := MBR{}
	ebr_empty := EBR{}

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		// Procedo a leer el archivo

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
			s_part_name := ""
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_name = string(master_boot_record.partition[i].name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				// Verifico si ya existe una particion con ese nombre
				if s_part_name == nombre {
					f.Close()
					return true
				}

				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.partition[i].type_[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				// Verifico si de tipo extendida
				if s_part_type == "e" {
					extendida = i
				}
			}

			// Si es extendida
			if extendida != -1 {
				// Obtengo el inicio de la particion
				s_part_start := master_boot_record.partition[extendida].start

				// Obtengo el espacio de la partcion
				s_part_size := master_boot_record.partition[extendida].size

				// Calculo del tamaño de struct en bytes
				ebr2 := struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				// Lectrura de conjunto de bytes en archivo binario
				lectura := make([]byte, sstruct)
				// Lee a partir del inicio de la particion
				n_leidos, _ := f.Read(lectura)

				// Posicion actual en el archivo
				f.Seek(int64(s_part_start), io.SeekStart)

				// Posicion actual en el archivo
				pos_actual, _ := f.Seek(0, io.SeekCurrent)

				// Lectrura de conjunto de bytes desde el inicio de la particion
				for n_leidos != 0 && (pos_actual < int64(s_part_size+s_part_start)) {
					// Lectrura de conjunto de bytes en archivo binario
					lectura := make([]byte, sstruct)
					// Lee a partir del inicio de la particion
					n_leidos, _ = f.Read(lectura)

					// Posicion actual en el archivo
					pos_actual, _ = f.Seek(0, io.SeekCurrent)

					// Conversion de bytes a struct
					extended_boot_record := bytes_a_struct_ebr(lectura)

					if extended_boot_record.s == 0 {
						break
					} else {
						// Antes de comparar limpio la cadena
						s_part_name = string(extended_boot_record.name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						// Verifico si ya existe una particion con ese nombre
						if s_part_name == nombre {
							f.Close()
							return true
						}

						// Obtengo el espacio utilizado
						s_part_next := extended_boot_record.next

						// Si ya termino
						if s_part_next != -1 {
							f.Close()
							return false
						}
					}
				}
			}
		} else {
			salida_comando += "[ERROR] el disco se encuentra vacio...\\n"
		}
	}

	f.Close()
	return false
}

// Busca particiones Primarias o Extendidas
func buscar_particion_p_e(direccion string, nombre string) int {
	// Apertura del archivo
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

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

		s_part_status := ""
		s_part_name := ""

		// Recorro las 4 particiones
		for i := 0; i < 4; i++ {
			// Antes de comparar limpio la cadena
			s_part_status = string(master_boot_record.partition[i].status[:])
			s_part_status = strings.Trim(s_part_status, "\x00")

			if s_part_status != "1" {
				// Antes de comparar limpio la cadena
				s_part_name = string(master_boot_record.partition[i].name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				if s_part_name == nombre {
					return i
				}
			}

		}
	}

	return -1
}

// Busca particiones Logicas
func buscar_particion_l(direccion string, nombre string) int {
	// Apertura del archivo
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		extendida := -1
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

		s_part_type := ""

		// Recorro las 4 particiones
		for i := 0; i < 4; i++ {
			// Antes de comparar limpio la cadena
			s_part_type = string(master_boot_record.partition[i].type_[:])
			s_part_type = strings.Trim(s_part_type, "\x00")

			if s_part_type != "e" {
				extendida = i
				break
			}
		}

		// Si hay extendida
		if extendida != -1 {
			ebr_empty := EBR{}

			ebr2 := struct_a_bytes(ebr_empty)
			sstruct := len(ebr2)

			// Lectrura del archivo binario desde el inicio
			lectura := make([]byte, sstruct)

			s_part_start := master_boot_record.partition[extendida].start
			f.Seek(int64(s_part_start), io.SeekStart)

			n_leidos, _ := f.Read(lectura)
			pos_actual, _ := f.Seek(0, io.SeekCurrent)

			s_part_size := master_boot_record.partition[extendida].start

			for (n_leidos != 0) && (pos_actual < int64(s_part_start+s_part_size)) {
				n_leidos, _ = f.Read(lectura)
				pos_actual, _ = f.Seek(0, io.SeekCurrent)

				// Conversion de bytes a struct
				extended_boot_record := bytes_a_struct_ebr(lectura)

				s_part_name_ext := string(extended_boot_record.name[:])
				s_part_name_ext = strings.Trim(s_part_name_ext, "\x00")

				ebr_empty_byte := struct_a_bytes(ebr_empty)

				if s_part_name_ext == nombre {
					return int(pos_actual) - len(ebr_empty_byte)
				}
			}
		}
		f.Close()
	}

	return -1
}
