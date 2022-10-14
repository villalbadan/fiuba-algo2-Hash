package diccionario

import (
	TDALista "Hash/lista"
	"fmt"
)

const (
	CAPACIDAD_INICIAL = 101
	FACTOR_DE_CARGA   = 0.8
	FACTOR_REDIMENSION = 2
)

type diccionarioImplementacion[K comparable, V any] struct {
	tablaValores []TDALista.Lista[elementoTabla[K, V]]
	elementos    int
}

type elementoTabla[K comparable, V any] struct {
	clave K
	valor V
}

func crearTabla[K comparable, V any](capacidad int) []TDALista.Lista[elementoTabla[K, V]]{
	return make([]TDALista.Lista[elementoTabla[K, V]], capacidad)
	//for i := range dict.tablaValores {
	//lista := TDALista.CrearListaEnlazada[elementoTabla[K, V]]()
	//} >>> no se si es necesario

}

func CrearHash[K comparable, V any]() Diccionario[K, V]{
	dict := new(diccionarioImplementacion[K, V])
	dict.tablaValores = crearTabla(CAPACIDAD_INICIAL)
	return dict
}

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func posicionEnTabla[K comparable](clave K, largo int) int {
	return funcionHash(convertirABytes(clave)) / largo
}

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

//me falta una vuelta de tuerca para que esto funcione reutilizando el codigo de Guardar
//ahora mismo esta haciendo cualquier cosa
func (dict *diccionarioImplementacion[K, V]) redimensionar() {
	nuevaCapacidad := proximoPrimo(len(dict.tablaValores) * FACTOR_REDIMENSION)
	tablaActual := dict.tablaValores
	nuevaTabla := crearTabla(nuevaCapacidad)

	for iter := dict.Iterador(); iter.HaySiguiente(); {
		index := posicionEnTabla(iter.VerActual().clave, nuevaCapacidad)
		dict.guardarEnTabla(nuevaTabla, index, iter.VerActual().clave, iter.VerActual().valor)
		iter.Siguiente()
	}

	dict.tablaValores = nuevaTabla
}

func (dict diccionarioImplementacion[K, V]) buscar(index int, clave K) IterDiccionario[K,V] {
	//Este iterador es de DICT
	for iter := dict.tablaValores[index].Iterador(); iter.HaySiguiente(); {
		if iter.VerActual().clave == clave {
			return iter
		}
		iter.HaySiguiente()
	}
	return nil
}

func (dict *diccionarioImplementacion[K, V]) guardarEnTabla(tabla []TDALista.Lista[elementoTabla[K, V]], indice int, clave K, dato V) {
	//!! Este iterador es de lista, NO DICT.
	for iter := tabla[indice].Iterador(); iter.HaySiguiente(); {
		if iter.VerActual().clave == clave {
			iter.VerActual().valor = dato
			return
		}
		iter.Siguiente()
	}

	//lista vacia o clave no esta
	tabla[indice].InsertarUltimo(elementoTabla[K, V]{clave, dato})
}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (dict *diccionarioImplementacion[K, V]) Guardar(clave K, dato V) {

	if (float32(len(dict.tablaValores)) / float32(dict.Cantidad())) > FACTOR_DE_CARGA {
		redimensionar()
	}

	index := posicionEnTabla(clave, len(dict.tablaValores))
	dict.guardarEnTabla(dict.tablaValores, index, clave, dato)
	//if dict.tablaValores[index] == nil {
	//	lista := TDALista.CrearListaEnlazada[elementoTabla[K, V]]()
	//	lista.InsertarUltimo(elementoTabla[K, V]{clave, dato})
	//	dict.tablaValores[index] = lista
	//} else {}


}

// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (dict diccionarioImplementacion[K, V]) Pertenece(clave K) bool {
	index := posicionEnTabla(clave, len(dict.tablaValores))
	if !dict.tablaValores[index].EstaVacia() {
		if  dict.buscar(index, clave) != nil {
				return true
			}
	}
	return false
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (dict diccionarioImplementacion[K, V]) Obtener(clave K) V {
	index := posicionEnTabla(clave, len(dict.tablaValores))
	if !dict.tablaValores[index].EstaVacia() {
		dato := dict.buscar(index, clave)
		if dato != nil {
			return dato.VerActual().valor
		}
	}
	panic("La clave no pertenece al diccionario")
}

// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (dict *diccionarioImplementacion[K, V]) Borrar(clave K) V {
	index := posicionEnTabla(clave, len(dict.tablaValores))
	if !dict.tablaValores[index].EstaVacia() {
		dato := dict.buscar(index, clave)
		if dato != nil {
				return dato.Borrar().valor
		}
	}
	panic("La clave no pertenece al diccionario")
}

// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (dict diccionarioImplementacion[K, V]) Cantidad() int {
	return dict.elementos
}


// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo
Iterar(func(clave K, dato V) bool)

// Iterador devuelve un IterDiccionario para este Diccionario
Iterador() IterDiccionario[K, V]