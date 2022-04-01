package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Mono struct {
	Id string
}

type Task struct {
	Url         string
	Referencias int
}

type Result struct {
	Origen          string `json:"origen"`
	Conteo_Palabras int    `json:"conteo_palabras"`
	Conteo_Enlaces  int    `json:"conteo_enlaces"`
	Sha             string `json:"sha"`
	Url             string `json:"url"`
	Mono            string `json:"mono"`
}

var cantidad_monos, tamano_cola, n_r int
var url_inicial, nombre_archivo string
var queue []string

func worker(jobs <-chan Task, results chan<- Task, id int) {
	for j := range jobs {
		Url := j.Url
		Nr := j.Referencias

		conteo_palabras := 0
		conteo := 0
		aux := ""

		c := colly.NewCollector()
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		c.OnHTML("div#mw-content-text p", func(e *colly.HTMLElement) {
			conteo_palabras += len(strings.Split(e.Text, " "))
		})

		c.OnHTML("div#mw-content-text p a", func(e *colly.HTMLElement) {
			if conteo < Nr {
				aux = e.Request.AbsoluteURL(e.Attr("href"))
				results <- Task{aux, Nr - 1}
				conteo = conteo + 1
			}
		})
		c.Visit(Url)
		fmt.Println("mono: 				" + strconv.Itoa(id))
		conteo_string := strconv.Itoa(conteo_palabras)
		after := <-jobs
		tofile := "{\n\"origen\":\"" + newSha(Url) + "\",\n\"conteo_palabras\":" + conteo_string + ",\n\"conteo_enlaces\":" + strconv.Itoa(conteo) + ",\n\"sha\":\"" + newSha(after.Url) + "\",\n\"url\":\"" + (after.Url) + "\",\n\"id_mono\":\"monin" + strconv.Itoa(id) + "\"\n},"
		if len(queue) > 0 {
			queue = queue[:len(queue)-1]
		}
		escribirArchivo(tofile)
		fmt.Println(queue)
		time.Sleep(time.Duration(500/500) * time.Second)
	}
}

func newSha(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func init_values() {
	fmt.Println("Práctica #2 | Sistemas Operativos 2 | G14")
	fmt.Println("")

	fmt.Print("Cantidad de monos buscadores: ")
	// fmt.Scan(&cantidad_monos)
	cantidad_monos = 3

	fmt.Print("Tamaño de la cola de espera: ")
	// fmt.Scan(&tamano_cola)
	tamano_cola = 3

	fmt.Print("Nr: ")
	// fmt.Scan(&n_r)
	n_r = 3

	fmt.Print("URL inicial: ")
	// fmt.Scan(&url_inicial)
	url_inicial = "https://es.wikipedia.org/wiki/Panavia_Tornado"

	fmt.Print("Nombre del archivo para el resultado de la búsqueda: ")
	// fmt.Scan(&nombre_archivo)
	nombre_archivo = "res"

	result_file, e := os.Create(nombre_archivo + ".json")
	if e != nil {
		log.Fatal(e)
		result_file.Close()
	}
	result_file.Close()
}

func escribirArchivo(contenido string) {
	file, err := os.OpenFile(nombre_archivo+".json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Print("el archivo no se ha creado")
		return
	}
	anterior, err := ioutil.ReadFile(nombre_archivo + ".json")
	lenstr := strconv.Itoa(len(queue))
	if len(anterior) > 14 {
		desfase := 13 + len(lenstr)
		anterior = anterior[:len(anterior)-desfase]
	}
	_, err1 := file.WriteString(string(anterior) + contenido + "\n{\"cola\": " + lenstr + "}\n]\n")
	if err1 != nil {
		fmt.Println(err1)
	}
	defer file.Close()
}

func main() {
	init_values()
	jobs := make(chan Task, 100)
	results := make(chan Task, 100)
	escribirArchivo("[")
	for i := 0; i < cantidad_monos; i++ {
		go worker(jobs, results, i)
	}

	jobs <- Task{url_inicial, n_r}

	for r := range results {
		// before := <-results
		// escribirArchivo("origen: " + newSha(before.Url))
		// fmt.Println("Visitando resultados")
		// fmt.Println("id ", r)
		// fmt.Println(<-results)
		fmt.Println(r)
		jobs <- r
		jobs <- r
		queue = append([]string{r.Url}, queue...)
		fmt.Println(queue)
	}
}

//Este comentario es de Benaventi
//va mira, lo que falta es la cola, que esa creo que la podemos hacer facil, eso solo tenemos que hacer uC
// creo que son las url segun faltan de hacer, o sea entra un mono y tiene que buscar una url, pero hay ciertos momentos en los que la cola tiene varias url
// a las cuales tiene que visitar, entonces eso es todo, si queres seguime y probamos
// en ese caso no podriamos usar un queue qye ya tenga golang vdd?
// si podemos, yo te iba a decir que usaramos un slice, porque ahi podemos usar append, creo que push y pop, porque algo asi necesitamos
// creo que mas que nada es un push y pop
// lo que no se es donde podriamos ir mostrando la lista, si mostrarla cada vez que agregemos un elemento, o irla escribiendo en un archivo, eso si no se

//ps creo que si se podria en un archivo y que se vaya leyendo._
// ahora con el queue lo que te pregunto es que si tiene que ser algun tipo de archvio especial o puede ser un string
// o si golang tiene algun tipo para URL?
// yo digo que string igual lo vamos a escribir despues

// solo voy a probar que onda con las colas, porque creo que podemos usar apend y extract, ahorita reviso dame chance
// vaa pero no seria en un queue sino en otra estructura o si?
// yo digo que si, o cual pensabas usar vos? no es una lista como tal, es un slice, que es diferente, pero en este caso nos sirve
// no se si hay tipo queue en go, alguno ya implementado, porque si no tendriamos que crearlo nosotros

//no he usado un slice pero si crees que nos sirve mas que un queue puedo ver como implementarlo
