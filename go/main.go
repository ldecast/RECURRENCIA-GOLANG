package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
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

type ResultFile struct {
	Results []Result
	Cola    int
}

var cantidad_monos, tamano_cola, n_r int
var url_inicial, nombre_archivo string
var queuesize int
var origen bool = true
var anterior string
var resultFile ResultFile

func worker(jobs <-chan Task, results chan<- Task, id int) {
	for j := range jobs {
		Url := j.Url
		Nr := j.Referencias

		conteo_palabras := 0
		conteo_url := 0
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
			conteo_url++
			if conteo < Nr {
				aux = e.Request.AbsoluteURL(e.Attr("href"))
				results <- Task{aux, Nr - 1}
				conteo = conteo + 1
			}
		})
		c.Visit(Url)
		if origen {
			anterior = "0"
			origen = false
		}
		after := j
		// AQUI SE RECUPERAN TODOS LOS DATOS
		result := Result{anterior, conteo_palabras, conteo_url, newSha(after.Url), after.Url, strconv.Itoa(id)}
		anterior = newSha(j.Url)
		conteo_url = 0
		queuesize = len(jobs)
		escribirArchivo(result)
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
		result_file.Close()
	}
	result_file.Close()
}

func escribirArchivo(contenido Result) {
	resultFile.Results = append(resultFile.Results, contenido)
	file, _ := json.MarshalIndent(resultFile, "", " ")
	err1 := ioutil.WriteFile(nombre_archivo+".json", file, 0644)
	if err1 != nil {
		fmt.Println(err1)
	}
}

func main() {
	init_values()
	jobs := make(chan Task, tamano_cola)
	results := make(chan Task, 1000000)
	for i := 0; i < cantidad_monos; i++ {
		go worker(jobs, results, i)
	}

	jobs <- Task{url_inicial, n_r}

	for r := range results {
		jobs <- r
	}
}
