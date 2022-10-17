package diccionario

import (
	"fmt"
)

const (
	CAPACIDAD_INICIAL  = 127
	CAPACIDAD_MAXIMA   = 2157587
	ANTERIOR_PRIMO     = -1
	PROX_PRIMO         = 1
	MAX_FC             = 0.99
	MIN_FC             = 0.1
	FACTOR_REDIMENSION = 2
	TABLA_VACIA        = 0
	NO_EN_TABLA        = 0
	///////////////////////////////
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

// ###### Hashear clave ------------------------------------------------------------------------------------------

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

// ######## FUNCION 1 - JENKINS
/* Hash: jenkins one-at-a-time-hash
https://en.wikipedia.org/wiki/Jenkins_hash_function
*/

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

// ######### FUNCION 1 - ADLER32
/* ###### Hash: adler32
https://pkg.go.dev/hash/adler32
*/

//func funcionHash1(clave []byte, largo int) int {
//	posicion := hash2.Checksum(clave) % uint32(largo)
//	return int(posicion)
//}

//######## FUNCION 3 - CRC64

//func funcionHash1(clave []byte, largo int) int {
//	posicion := crc64.Checksum(clave, crc64.MakeTable(crc64.ISO)) % uint64(largo)
//	return int(posicion)
//}

/*######## FUNCION 5 - FNV
Hash: Fowler–Noll–Vo (FNV)
https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
Implementación de https://golangprojectstructure.com/
*/

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

////####### FUNCION SDBM
//
//func funcionHash2(clave []byte, largo int) int {
//	posicion := sdbmHash(clave) % uint64(largo)
//	return int(posicion)
//}
//
//func sdbmHash(data []byte) uint64 {
//	var hash uint64
//
//	for _, b := range data {
//		hash = uint64(b) + (hash << 6) + (hash << 16) - hash
//	}
//
//	return hash
//}

// ######## FUNCION 2 - DJB2

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

//############ REDIMENSION -----------------------------------------------------------------------------------------

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

//TODO cambiar numeros magicos
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

//TODO borrar codigo repetido
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

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (dict *dictImplementacion[K, V]) Guardar(claveAEvaluar K, dato V) {
	if sobrecarga(dict.elementos+1, len(dict.tabla)) {
		capacidad := dict.nuevaCapacidad(dict.primo, PROX_PRIMO)
		dict.redimensionar(capacidad)
	}

	dict.guardarEnTabla(dict.tabla, claveAEvaluar, dato)

}

// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (dict dictImplementacion[K, V]) Pertenece(clave K) bool {
	hash, _ := dict.buscar(dict.tabla, clave)
	return hash != NO_EN_TABLA
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (dict dictImplementacion[K, V]) Obtener(clave K) V {
	hash, indice := dict.buscar(dict.tabla, clave)
	if hash == NO_EN_TABLA {
		panic("La clave no pertenece al diccionario")
	}
	return dict.tabla[indice].valor
}

// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
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

// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (dict dictImplementacion[K, V]) Cantidad() int {
	return dict.elementos
}

// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo
func (dict dictImplementacion[K, V]) Iterar(visitar func(K, V) bool) {
	for i := 0; i < len(dict.tabla); i++ {
		if dict.tabla[i] != nil {
			if !visitar(dict.tabla[i].clave, dict.tabla[i].valor) {
				break
			}
		}
	}
}

// Iterador devuelve un IterDiccionario para este Diccionario
func (dict *dictImplementacion[K, V]) Iterador() IterDiccionario[K, V] {
	for i := range dict.tabla {
		if dict.tabla[i] != nil {
			return &iteradorDict[K, V]{diccionario: dict, posicion: i}
		}
	}
	return &iteradorDict[K, V]{diccionario: dict, posicion: len(dict.tabla)}
}

// HaySiguiente devuelve si hay más datos para ver. Esto es, si en el lugar donde se encuentra parado
// el iterador hay un elemento.
func (iter *iteradorDict[K, V]) HaySiguiente() bool {
	return iter.posicion < len(iter.diccionario.tabla)
}

// VerActual devuelve la clave y el dato del elemento actual en el que se encuentra posicionado el iterador.
// Si no HaySiguiente, debe entrar en pánico con el mensaje 'El iterador termino de iterar'
func (iter *iteradorDict[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iter.diccionario.tabla[iter.posicion].clave, iter.diccionario.tabla[iter.posicion].valor
}

// Siguiente si HaySiguiente, devuelve la clave actual (equivalente a VerActual, pero únicamente la clave), y
// además avanza al siguiente elemento en el diccionario. Si no HaySiguiente, entonces debe entrar en pánico con
// mensaje 'El iterador termino de iterar'
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
