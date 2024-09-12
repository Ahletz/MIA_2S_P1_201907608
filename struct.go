package main

type MBR struct {
	tamaño         int
	fecha_creacion string
	dsk_asignature int
	fit            byte
	partition      [4]Partition //arreglo con conteniod de 4 particiones
}

type Partition struct {
	status      byte
	type_       byte
	fit         byte
	start       int
	s           int
	name        [16]byte
	correlative int
	id          [4]byte
}

type EBR struct {
	mount byte
	fit   byte
	start int
	s     int
	next  int
	name  [16]byte
}

//Estructuras para carpetas y archivos

type SuperBloque struct {
	s_filesystem_type   int
	s_inodes_count      int
	s_blocks_count      int
	s_free_blocks_count int
	s_free_inodes_count int
	s_mtimetimeÚltima   string
	s_umtimetimeÚltima  string
	s_mnt_count         int
	s_magic             int
	s_inode_s           int
	s_block_s           int
	s_firts_ino         int
	s_first_blo         int
	s_bm_inode_start    int
	s_bm_block_start    int
	s_inode_start       int
	s_block_start       int
}

type Inodos struct {
	i_uid   int
	i_gid   int
	i_s     int
	i_atime string
	i_ctime string
	i_mtime string
	i_block [15]int
	i_type  byte
	i_perm  [3]byte
}

type B_content struct {
	b_name  [12]byte
	b_inodo int
}

type B_carpeta struct {
	b_content [4]B_content //contendra 4 estructuras B_content
}

type B_archivos struct {
	b_content [64]B_content
}

type B_pointer struct {
	b_pointer [16]int
}
