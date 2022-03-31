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

func worker(jobs <-chan Task, results chan<- Task) {
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
			// fmt.Println(e)
			// fmt.Println("conteo de palabras", conteo_palabras)
		})

		c.OnHTML("div#mw-content-text p a", func(e *colly.HTMLElement) {
			// fmt.Println(e.Attr("href"), "href ")
			if conteo < Nr {
				aux = e.Request.AbsoluteURL(e.Attr("href"))
				// fmt.Println(aux)
				results <- Task{aux, Nr - 1}
				conteo = conteo + 1
			}
		})
		c.Visit(Url)
		conteo_string := strconv.Itoa(conteo_palabras)
		after := <-jobs
		escribirArchivo("{")
		escribirArchivo("\"origen\":\"" + newSha(Url) + "\",")
		escribirArchivo("\"conteo_palabras\":" + conteo_string + ",")
		escribirArchivo("\"conteo_enlaces\":" + strconv.Itoa(conteo) + ",")
		escribirArchivo("\"sha\":\"" + newSha(after.Url) + "\",")
		escribirArchivo("\"url\":\"" + (after.Url) + "\"")
		escribirArchivo("},")
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
	cantidad_monos = 1

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
	if len(anterior) > 2 {
		anterior = anterior[:len(anterior)-2]
	}
	_, err1 := file.WriteString(string(anterior) + contenido + "\n]\n")
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
	go worker(jobs, results)

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
	}
}
