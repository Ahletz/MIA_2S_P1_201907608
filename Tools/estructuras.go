package Tools

type MBR struct {
	tamaño         int
	fecha_creacion [45]byte
	dsk_asignature int
	fit            [1]byte
	partition      [4]Partition //arreglo con conteniod de 4 particiones
}

type Partition struct {
	status      [50]byte
	type_       [1]byte
	fit         [1]byte
	start       int
	size        int
	name        [16]byte
	correlative int
	id          [4]byte
}

type EBR struct {
	mount [1]byte
	fit   [1]byte
	start int
	s     int
	next  int
	name  [16]byte
}
