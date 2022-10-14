package diccionario

import (
	TDALista "diccionario/lista"
)

type dict[K comparable, V any] struct {
	tablaValores []TDALista.Lista[*elementoTabla[K, V]]
	elementos    int
}

type elementoTabla[K comparable, V any] struct {
	clave K
	valor V
}

type iteradorDict[K comparable, V any] struct {
	diccionario         *dict[K, V]
	iteradorListaActual TDALista.IteradorLista[elementoTabla[K, V]]
	posicionListaActual int
}

// PRIMITIVAS DICCIONARIO
// #############################################################################################################

func (diccionario *dict[K, V]) Guardar(clave K, dato V) {

}

func (diccionario *dict[K, V]) Pertenece(clave K) bool {
	return false
}

func (diccionario *dict[K, V]) Obtener(clave K) V {
	return diccionario.tablaValores[0].VerPrimero().valor
}

func (diccionario *dict[K, V]) Borrar(clave K) V {
	return diccionario.tablaValores[0].VerPrimero().valor
}

func (diccionario *dict[K, V]) Cantidad() int {
	return 0
}

// ITERADOR INTERNO
// #############################################################################################################

func (diccionario *dict[K, V]) Iterar(Sigue func(K, V) bool) {
	condCorte := false
	for i := 0; i < len(diccionario.tablaValores) && !condCorte; i++ {
		if !diccionario.tablaValores[i].EstaVacia() && !condCorte {
			for iter := diccionario.tablaValores[i].Iterador(); iter.HaySiguiente(); iter.Siguiente() {
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
func encontrarPrimeraLista[K comparable, V any](diccionario *dict[K, V]) TDALista.Lista[*elementoTabla[K, V]] {
	for i := range diccionario.tablaValores {
		if !diccionario.tablaValores[i].EstaVacia() {
			return diccionario.tablaValores[i]
		}
	}
	return diccionario.tablaValores[0]
}

func (diccionario *dict[K, V]) Iterador() IterDiccionario[K, V] {
	return &iteradorDict[K, V]{diccionario: diccionario, iteradorListaActual: encontrarPrimeraLista(diccionario).Iterador()}
}

func (iterador *iteradorDict[K, V]) HaySiguiente() bool {
	return iterador.iteradorListaActual.VerActual() != nil
}

func (iterador *iteradorDict[K, V]) VerActual() (K, V) {
	if iterador.iteradorListaActual.VerActual() == nil {
		panic("El iterador termino de iterar")
	}
	return iterador.iteradorListaActual.VerActual().clave, iterador.iteradorListaActual.VerActual().valor
}

func (iterador *iteradorDict[K, V]) Siguiente() K {
	claveActual, _ := iterador.VerActual()
	if !iterador.iteradorListaActual.HaySiguiente() {
		if iterador.posicionListaActual >= len(iterador.diccionario.tablaValores) {
			panic("El iterador termino de iterar")
		}
		iterador.posicionListaActual++
	}
	return claveActual
}
