package diccionario

import (
	"fmt"
)

const (
	CAPACIDAD_INICIAL  = 127
	CAPACIDAD_MAXIMA   = 2157587
	ANTERIOR_PRIMO     = -1
	PROX_PRIMO         = 1
	MAX_FC             = 0.91
	MIN_FC             = 0.1
	FACTOR_REDIMENSION = 2
	TABLA_VACIA        = 0
	NO_EN_TABLA        = 0

	PRIMER_HASH  = 1
	SEGUNDO_HASH = 2
	ULTIMO_HASH  = 3
)

type dictImplementacion[K comparable, V any] struct {
	tabla     []*elementoTabla[K, V]
	elementos int
	primo     int
}

type elementoTabla[K comparable, V any] struct {
	clave  K
	valor  V
	opcion int
}

type iteradorDict[K comparable, V any] struct {
	diccionario *dictImplementacion[K, V]
	posicion    int
}

func crearTabla[K comparable, V any](capacidad int) []*elementoTabla[K, V] {
	return make([]*elementoTabla[K, V], capacidad)
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	dict := new(dictImplementacion[K, V])
	dict.tabla = crearTabla[K, V](CAPACIDAD_INICIAL)
	return dict
}

// // ###################################### HASHEAR CLAVE ####################################################

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func posicionEnTabla(opcion int, claveEnBytes []byte, largo int) int {
	switch opcion {
	case 1:
		return funcionHash1(claveEnBytes, largo)
	case 2:
		return funcionHash2(claveEnBytes, largo)
	case 3:
		return funcionHash3(claveEnBytes, largo)
	default:
		panic("HASH INVALIDO")
	}
}

func funcionHash1(clave []byte, largo int) int {
	posicion := djb2(clave) % uint32(largo)
	return int(posicion)
}

func djb2(data []byte) uint32 {
	hash := uint32(5381)

	for _, b := range data {
		hash += uint32(b) + hash + hash<<5
	}

	return hash
}

func funcionHash2(clave []byte, largo int) int {
	posicion := fvnHash(clave) % uint64(largo)
	return int(posicion)
}

const (
	uint64Offset uint64 = 0xcbf29ce484222325
	uint64Prime  uint64 = 0x00000100000001b3
)

func fvnHash(data []byte) (hash uint64) {
	hash = uint64Offset

	for _, b := range data {
		hash ^= uint64(b)
		hash *= uint64Prime
	}

	return
}

func funcionHash3(clave []byte, largo int) int {
	posicion := jenkins(clave) % uint64(largo)
	return int(posicion)
}

func jenkins(clave []byte) uint64 {
	var hash uint64
	for _, b := range clave {
		hash += uint64(b)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return hash
}

//// ######################################### REDIMENSION ###################################################

func (dict *dictImplementacion[K, V]) nuevaCapacidad(pos int, movimiento int) int {
	arrayPrimos := []int{CAPACIDAD_INICIAL, 257, 523, 1049, 2099, 4201, 8419, 16843, 33703, 67409, 134837, 269683, 539389, 1078787, CAPACIDAD_MAXIMA}

	if pos+movimiento > len(arrayPrimos) {
		panic("Por ahora no puede hacer un hash tan grande, lo sentimos.")
	}

	dict.primo = dict.primo + movimiento
	return arrayPrimos[pos+movimiento]
}

func pocaCarga(elementos int, largoTabla int) bool {
	return float32(elementos)/float32(largoTabla) < MIN_FC && largoTabla/FACTOR_REDIMENSION > CAPACIDAD_INICIAL
}

func (dict *dictImplementacion[K, V]) redimensionar(nuevaCapacidad int) {
	var cantidad int
	nuevaTabla := crearTabla[K, V](nuevaCapacidad)

	for iter := dict.Iterador(); iter.HaySiguiente(); {
		clave, valor := iter.VerActual()
		dict.guardarEnTabla(nuevaTabla, clave, valor)
		cantidad++
		iter.Siguiente()
	}
	dict.elementos = cantidad
	dict.tabla = nuevaTabla
}

func sobrecarga(elementos int, largoTabla int) bool {
	return float32(elementos)/float32(largoTabla) >= MAX_FC
}

// ###################################### BUSQUEDA Y GUARDADO #################################################

func (dict *dictImplementacion[K, V]) buscar(tabla []*elementoTabla[K, V], clave K) (int, int) {
	claveEnByte := convertirABytes(clave)

	for i := PRIMER_HASH; i <= ULTIMO_HASH; i++ {
		posicion := posicionEnTabla(i, claveEnByte, len(tabla))
		if dict.elementos == TABLA_VACIA {
			return NO_EN_TABLA, posicion
		}

		if tabla[posicion] != nil && tabla[posicion].clave == clave {
			return i, posicion
		}
	}

	return NO_EN_TABLA, posicionEnTabla(PRIMER_HASH, claveEnByte, len(tabla))
}

func (dict *dictImplementacion[K, V]) guardarEnOcupado(tabla []*elementoTabla[K, V], elemento *elementoTabla[K, V], claveOriginal K, cnt int) bool {
	if elemento.clave == claveOriginal {
		return false
	}
	cnt++
	claveEnByte := convertirABytes(elemento.clave)

	nuevaOpcion := elemento.opcion + 1
	if nuevaOpcion > ULTIMO_HASH {
		nuevaOpcion = PRIMER_HASH
	}
	indice := posicionEnTabla(nuevaOpcion, claveEnByte, len(tabla))
	elementoAMover := tabla[indice]
	(tabla[indice]) = &elementoTabla[K, V]{clave: elemento.clave, valor: elemento.valor, opcion: nuevaOpcion}

	if elementoAMover != nil {
		return dict.guardarEnOcupado(tabla, elementoAMover, claveOriginal, cnt)
	}

	return true
}

func (dict *dictImplementacion[K, V]) guardarEnTabla(tabla []*elementoTabla[K, V], claveAEvaluar K, dato V) {
	hash, indice := dict.buscar(tabla, claveAEvaluar)
	//CLAVE EXISTE: actualizamos
	if hash != NO_EN_TABLA {
		tabla[indice].valor = dato
		return
	}

	//CLAVE NO EXISTE:

	//Guardamos elemento original:
	elementoAMover := tabla[indice]
	tabla[indice] = &elementoTabla[K, V]{clave: claveAEvaluar, valor: dato, opcion: 1}
	dict.elementos++

	//Posición no vacia, comenzamos a mover
	if elementoAMover != nil {
		if !dict.guardarEnOcupado(tabla, elementoAMover, claveAEvaluar, 0) {
			capacidad := dict.nuevaCapacidad(dict.primo, PROX_PRIMO)
			dict.redimensionar(capacidad)
			dict.Guardar(claveAEvaluar, dato)
		}
	}

}

// ################################### PRIMITIVAS DICCIONARIO #################################################

func (dict *dictImplementacion[K, V]) Guardar(claveAEvaluar K, dato V) {
	if sobrecarga(dict.elementos+1, len(dict.tabla)) {
		capacidad := dict.nuevaCapacidad(dict.primo, PROX_PRIMO)
		dict.redimensionar(capacidad)
	}

	dict.guardarEnTabla(dict.tabla, claveAEvaluar, dato)

}

func (dict dictImplementacion[K, V]) Pertenece(clave K) bool {
	hash, _ := dict.buscar(dict.tabla, clave)
	return hash != NO_EN_TABLA
}

func (dict dictImplementacion[K, V]) Obtener(clave K) V {
	hash, indice := dict.buscar(dict.tabla, clave)
	if hash == NO_EN_TABLA {
		panic("La clave no pertenece al diccionario")
	}
	return dict.tabla[indice].valor
}

func (dict *dictImplementacion[K, V]) Borrar(clave K) V {
	hash, indice := dict.buscar(dict.tabla, clave)
	if hash == NO_EN_TABLA {
		panic("La clave no pertenece al diccionario")
	}

	borrado := dict.tabla[indice]
	dict.tabla[indice] = nil
	dict.elementos--

	if pocaCarga(dict.elementos, len(dict.tabla)) {
		capacidad := dict.nuevaCapacidad(dict.primo, ANTERIOR_PRIMO)
		dict.redimensionar(capacidad)
	}

	return borrado.valor
}

func (dict dictImplementacion[K, V]) Cantidad() int {
	return dict.elementos
}

func (dict dictImplementacion[K, V]) Iterar(visitar func(K, V) bool) {
	for i := 0; i < len(dict.tabla); i++ {
		if dict.tabla[i] != nil {
			if !visitar(dict.tabla[i].clave, dict.tabla[i].valor) {
				break
			}
		}
	}
}

// ################################### PRIMITIVAS ITERADOR ###################################################

func (dict *dictImplementacion[K, V]) Iterador() IterDiccionario[K, V] {
	for i := range dict.tabla {
		if dict.tabla[i] != nil {
			return &iteradorDict[K, V]{diccionario: dict, posicion: i}
		}
	}
	return &iteradorDict[K, V]{diccionario: dict, posicion: len(dict.tabla)}
}

func (iter *iteradorDict[K, V]) HaySiguiente() bool {
	return iter.posicion < len(iter.diccionario.tabla)
}

func (iter *iteradorDict[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iter.diccionario.tabla[iter.posicion].clave, iter.diccionario.tabla[iter.posicion].valor
}

func (iter *iteradorDict[K, V]) Siguiente() K {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	posActual := iter.posicion

	iter.posicion++
	if iter.posicion < len(iter.diccionario.tabla) {
		for i := iter.posicion; i < len(iter.diccionario.tabla); i++ {
			if iter.diccionario.tabla[iter.posicion] != nil {
				break
			}
			iter.posicion++
		}
	}

	return iter.diccionario.tabla[posActual].clave
}
