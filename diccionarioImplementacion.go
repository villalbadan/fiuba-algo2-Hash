package diccionario

import (
	TDALista "diccionario/lista"
	"fmt"
)

const (
	CAPACIDAD_INICIAL  = 101
	MAX_FC             = 0.8
	MIN_FC             = 0.2
	FACTOR_REDIMENSION = 2
)

type dictImplementacion[K comparable, V any] struct {
	tablaValores []TDALista.Lista[*elementoTabla[K, V]]
	elementos    int
}

type elementoTabla[K comparable, V any] struct {
	clave K
	valor V
}

type iteradorDict[K comparable, V any] struct {
	diccionario         *dictImplementacion[K, V]
	iteradorListaActual TDALista.IteradorLista[*elementoTabla[K, V]]
	posicionListaActual int
}

func crearTabla[K comparable, V any](capacidad int) []TDALista.Lista[*elementoTabla[K, V]] {
	tabla := make([]TDALista.Lista[*elementoTabla[K, V]], capacidad)
	for i := range tabla {
		tabla[i] = TDALista.CrearListaEnlazada[*elementoTabla[K, V]]()
	}
	return tabla
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	dict := new(dictImplementacion[K, V])
	dict.tablaValores = crearTabla[K, V](CAPACIDAD_INICIAL)
	return dict
}

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func posicionEnTabla[K comparable](clave K, largo int) int {
	return jenkins(convertirABytes(clave), largo)
}

func jenkins(clave []byte, largo int) int {
	var hash byte
	for _, b := range clave {
		hash += b
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return int(hash) % largo
}

//func sdbmHash(data []byte) uint64 {
//	var hash uint64
//
//	for _, b := range data {
//		hash = uint64(b) + (hash << 6) + (hash << 16) - hash
//	}
//	return hash
//}

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

func (dict dictImplementacion[K, V]) buscar(lista TDALista.Lista[*elementoTabla[K, V]], clave K) TDALista.IteradorLista[*elementoTabla[K, V]] {
	//Este es un iterador de las listas del dict
	for iter := lista.Iterador(); iter.HaySiguiente(); {
		if iter.VerActual().clave == clave {
			return iter
		}
		iter.HaySiguiente()
	}
	return nil
}

func (dict *dictImplementacion[K, V]) guardarEnTabla(tabla []TDALista.Lista[*elementoTabla[K, V]], indice int, clave K, dato V) {
	//Este iterador es de LISTAS, NO de diccionario
	for iter := tabla[indice].Iterador(); iter.HaySiguiente(); {
		if iter.VerActual().clave == clave {
			iter.VerActual().valor = dato
			return
		}
		iter.Siguiente()
	}

	//lista vacia o clave no esta
	tabla[indice].InsertarUltimo(&elementoTabla[K, V]{clave, dato})
	dict.elementos++
}

func (dict *dictImplementacion[K, V]) redimensionar(nuevaCapacidad int) {
	nuevaTabla := crearTabla[K, V](nuevaCapacidad)

	for iter := dict.Iterador(); iter.HaySiguiente(); {
		clave, valor := iter.VerActual()
		index := posicionEnTabla(clave, nuevaCapacidad)
		dict.guardarEnTabla(nuevaTabla, index, clave, valor)
		iter.Siguiente()
	}

	dict.tablaValores = nuevaTabla
}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (dict *dictImplementacion[K, V]) Guardar(clave K, dato V) {

	if (float32(len(dict.tablaValores)) / float32(dict.Cantidad())) > MAX_FC {
		nuevaCapacidad := proximoPrimo(len(dict.tablaValores) * FACTOR_REDIMENSION)
		dict.redimensionar(nuevaCapacidad)
	}

	index := posicionEnTabla(clave, len(dict.tablaValores))
	dict.guardarEnTabla(dict.tablaValores, index, clave, dato)
}

// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (dict dictImplementacion[K, V]) Pertenece(clave K) bool {
	index := posicionEnTabla(clave, len(dict.tablaValores))
	if !dict.tablaValores[index].EstaVacia() {
		if dict.buscar(dict.tablaValores[index], clave) != nil {
			return true
		}
	}
	return false
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (dict dictImplementacion[K, V]) Obtener(clave K) V {
	index := posicionEnTabla(clave, len(dict.tablaValores))
	if !dict.tablaValores[index].EstaVacia() {
		dato := dict.buscar(dict.tablaValores[index], clave)
		if dato != nil {
			return dato.VerActual().valor
		}
	}
	panic("La clave no pertenece al diccionario")
}

// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (dict *dictImplementacion[K, V]) Borrar(clave K) V {
	index := posicionEnTabla(clave, len(dict.tablaValores))
	if !dict.tablaValores[index].EstaVacia() {
		dato := dict.buscar(dict.tablaValores[index], clave)

		if dato != nil {
			borrado := dato.Borrar()
			dict.elementos--

			if (float32(len(dict.tablaValores)) / float32(dict.Cantidad())) < MIN_FC {
				nuevaCapacidad := proximoPrimo(len(dict.tablaValores) / FACTOR_REDIMENSION)
				dict.redimensionar(nuevaCapacidad)
			}
			return borrado.valor
		}
	}
	panic("La clave no pertenece al diccionario")
}

// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (dict dictImplementacion[K, V]) Cantidad() int {
	return dict.elementos
}

// ITERADOR INTERNO
// #############################################################################################################

func (dict *dictImplementacion[K, V]) Iterar(Sigue func(K, V) bool) {
	condCorte := false
	for i := 0; i < len(dict.tablaValores) && !condCorte; i++ {
		if !dict.tablaValores[i].EstaVacia() && !condCorte {
			for iter := dict.tablaValores[i].Iterador(); iter.HaySiguiente(); iter.Siguiente() {
				if !Sigue(iter.VerActual().clave, iter.VerActual().valor) {
					condCorte = true
					break
				}
			}
		}
	}
}

// ITERADOR EXTERNO
// #############################################################################################################

func (dict *dictImplementacion[K, V]) encontrarPrimeraLista() TDALista.Lista[*elementoTabla[K, V]] {
	for i := range dict.tablaValores {
		if !dict.tablaValores[i].EstaVacia() {
			return dict.tablaValores[i]
		}
	}
	return dict.tablaValores[0]
}

func (dict *dictImplementacion[K, V]) Iterador() IterDiccionario[K, V] {
	return &iteradorDict[K, V]{
		diccionario:         dict,
		iteradorListaActual: dict.encontrarPrimeraLista().Iterador(),
	}
}

func (iter *iteradorDict[K, V]) HaySiguiente() bool {
	return iter.iteradorListaActual.HaySiguiente()
}

func (iter *iteradorDict[K, V]) VerActual() (K, V) {
	if !iter.iteradorListaActual.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iter.iteradorListaActual.VerActual().clave, iter.iteradorListaActual.VerActual().valor
}

func (iter *iteradorDict[K, V]) Siguiente() K {
	claveActual, _ := iter.VerActual()
	if !iter.iteradorListaActual.HaySiguiente() {
		if iter.posicionListaActual >= len(iter.diccionario.tablaValores) {
			panic("El iterador termino de iterar")
		}
		iter.posicionListaActual++
	}
	return claveActual
}
