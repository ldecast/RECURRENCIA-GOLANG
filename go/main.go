package main

import (
	"fmt"
	"log"
	"os"
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
			fmt.Println(e)
		})

		c.OnHTML("div#mw-content-text p a", func(e *colly.HTMLElement) {
			//fmt.Println(e.Attr("href"))
			if conteo < Nr {
				fmt.Println(e.Request.AbsoluteURL(e.Attr("href")))
				aux = e.Request.AbsoluteURL(e.Attr("href"))
				results <- Task{aux, Nr - 1}
				conteo = conteo + 1
			}
		})

		c.Visit(Url)
		time.Sleep(time.Duration(conteo_palabras/500) * time.Second)
	}
}

func init_values() {
	fmt.Println("Práctica #2 | Sistemas Operativos 2 | G14")
	fmt.Println("")

	fmt.Print("Cantidad de monos buscadores: ")
	fmt.Scan(&cantidad_monos)

	fmt.Print("Tamaño de la cola de espera: ")
	fmt.Scan(&tamano_cola)

	fmt.Print("Nr: ")
	fmt.Scan(&n_r)

	fmt.Print("URL inicial: ")
	fmt.Scan(&url_inicial)

	fmt.Print("Nombre del archivo para el resultado de la búsqueda: ")
	fmt.Scan(&nombre_archivo)

	result_file, e := os.Create(nombre_archivo + ".json")
	if e != nil {
		log.Fatal(e)
	}
	result_file.Close()
}

func main() {
	init_values()
	jobs := make(chan Task, 100)
	results := make(chan Task, 100)

	go worker(jobs, results)

	jobs <- Task{url_inicial, n_r}

	for r := range results {
		fmt.Println(<-results)
		fmt.Println("Visitando resultados")
		jobs <- r
	}
}
