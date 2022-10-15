package diccionario

import (
	"fmt"
	hash2 "hash/adler32"
	"hash/maphash"
)

const (
	CAPACIDAD_INICIAL  = 3
	MAX_FC             = 0.8
	MIN_FC             = 0.1
	FACTOR_REDIMENSION = 2
	TABLA_VACIA        = 0
	NO_EN_TABLA        = 0
	LIM_REHASH         = 100
)

type dictImplementacion[K comparable, V any] struct {
	tabla     []*elementoTabla[K, V]
	elementos int
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
	//ver de pasar la funcion por parametro?
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

func funcionHash2(clave []byte, largo int) int {
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

// ######### FUNCION 2 - ADLER32
/* ###### Hash: adler32
https://pkg.go.dev/hash/adler32
*/

func funcionHash1(clave []byte, largo int) int {
	posicion := hash2.Checksum(clave) % uint32(largo)
	return int(posicion)
}

// ######## FUNCION 3 - ???

func funcionHash3(clave []byte, largo int) int {
	var h maphash.Hash
	valorHasheado, _ := h.Write(clave)
	posicion := valorHasheado % largo
	return posicion
}

//############ REDIMENSION -----------------------------------------------------------------------------------------

func esPrimo(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func proximoPrimo(n int) int {
	if esPrimo(n) {
		return n
	}
	return proximoPrimo(n + 1)
}

func anteriorPrimo(n int) int {
	if esPrimo(n) {
		return n
	}
	return anteriorPrimo(n - 1)
}

func pocaCarga(elementos int, largoTabla int) bool {
	return float32(elementos)/float32(largoTabla) < MIN_FC && largoTabla/FACTOR_REDIMENSION > CAPACIDAD_INICIAL
}

func (dict *dictImplementacion[K, V]) redimensionar(nuevaCapacidad int) {
	nuevaTabla := crearTabla[K, V](nuevaCapacidad)

	for iter := dict.Iterador(); iter.HaySiguiente(); {
		clave, valor := iter.VerActual()
		dict.guardarEnTabla(nuevaTabla, clave, valor)
		iter.Siguiente()
	}

	dict.tabla = nuevaTabla
}

//func sobrecarga (elementos int, largoTabla int) bool {
//	return float32(elementos) / float32(largoTabla) > MAX_FC
//}

//TODO cambiar numeros magicos
func (dict *dictImplementacion[K, V]) buscar(tabla []*elementoTabla[K, V], clave K) (int, int) {
	claveEnByte := convertirABytes(clave)

	posicionFuncion1 := posicionEnTabla(1, claveEnByte, len(tabla))
	if dict.elementos == TABLA_VACIA {
		return NO_EN_TABLA, posicionFuncion1
	}
	if tabla[posicionFuncion1] != nil && tabla[posicionFuncion1].clave == clave {
		return 1, posicionFuncion1
	}

	posicionFuncion2 := posicionEnTabla(2, claveEnByte, len(tabla))
	if tabla[posicionFuncion2] != nil && tabla[posicionFuncion2].clave == clave {
		return 2, posicionFuncion2
	}

	posicionFuncion3 := posicionEnTabla(3, claveEnByte, len(tabla))
	if tabla[posicionFuncion3] != nil && tabla[posicionFuncion3].clave == clave {
		return 3, posicionFuncion3
	}

	return NO_EN_TABLA, posicionFuncion1
}

//TODO borrar codigo repetido
func (dict *dictImplementacion[K, V]) guardarEnOcupado(tabla []*elementoTabla[K, V], elemento *elementoTabla[K, V], claveOriginal K, cnt int) bool {
	if elemento.clave == claveOriginal || cnt > LIM_REHASH {
		return false
	}
	cnt++
	claveEnByte := convertirABytes(elemento.clave)

	if elemento.opcion == 3 {
		indice := posicionEnTabla(1, claveEnByte, len(dict.tabla))
		if dict.tabla[indice] != nil {
			elementoAMover := tabla[indice]
			return dict.guardarEnOcupado(tabla, elementoAMover, claveOriginal, cnt)
		}
		tabla[indice] = &elementoTabla[K, V]{clave: elemento.clave, valor: elemento.valor, opcion: 1}
	}

	if elemento.opcion == 1 {
		indice := posicionEnTabla(2, claveEnByte, len(dict.tabla))
		if dict.tabla[indice] != nil {
			elementoAMover := tabla[indice]
			return dict.guardarEnOcupado(tabla, elementoAMover, claveOriginal, cnt)
		}
		tabla[indice] = &elementoTabla[K, V]{clave: elemento.clave, valor: elemento.valor, opcion: 2}
	}

	if elemento.opcion == 2 {
		indice := posicionEnTabla(3, claveEnByte, len(dict.tabla))
		if dict.tabla[indice] != nil {
			elementoAMover := tabla[indice]
			return dict.guardarEnOcupado(tabla, elementoAMover, claveOriginal, cnt)
		}
		tabla[indice] = &elementoTabla[K, V]{clave: elemento.clave, valor: elemento.valor, opcion: 3}
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
			nuevaCapacidad := proximoPrimo(len(tabla) * FACTOR_REDIMENSION)
			dict.redimensionar(nuevaCapacidad)
			dict.Guardar(claveAEvaluar, dato)
		}
	}

}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (dict *dictImplementacion[K, V]) Guardar(claveAEvaluar K, dato V) {
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
		nuevaCapacidad := anteriorPrimo(len(dict.tabla) / FACTOR_REDIMENSION)
		dict.redimensionar(nuevaCapacidad)
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
	return &iteradorDict[K, V]{diccionario: dict, posicion: len(dict.tabla) - 1}
}

// HaySiguiente devuelve si hay más datos para ver. Esto es, si en el lugar donde se encuentra parado
// el iterador hay un elemento.
func (iter *iteradorDict[K, V]) HaySiguiente() bool {
	return iter.diccionario.tabla[iter.posicion] != nil && iter.posicion < len(iter.diccionario.tabla)
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

	for i := iter.posicion + 1; i < len(iter.diccionario.tabla); i++ {
		iter.posicion++
		if iter.diccionario.tabla[iter.posicion] != nil {
			break
		}
		//iter.posicion++
	}

	return iter.diccionario.tabla[posActual].clave
}
