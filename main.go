import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"fmt"
)

type meme struct {
	Title string `json:"Title"`
	Link  string `json:"link"`
}
type document struct {
	Body  string
	Title string
}

func RenderApi(documentChan chan string, res *http.Response) {

	// decode the api
	body, _ := ioutil.ReadAll(res.Body)
	var api []meme
	json.Unmarshal(body, &api)
	// parse the api
	var document string
	for _, i := range api {
		document += `
    <div>
      <h1>` +
			i.Title +
			` </h1> <img src="` + i.Link + `" width=500>
    </div>`

	}

	documentChan <- (document)
}
func RenderPage(w http.ResponseWriter) {

	url := "https://api-bruh-2.elpanajose.repl.co/memes"

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)
	documentChan := make(chan string)
	go RenderApi(documentChan, res)
	//open the template
	t, err := template.ParseFiles("page.html")
	if err != nil {
		w.Write([]byte("some error"))
		return
	}
	// render the template
	t.Execute(w, document{
		Title: "<h1></h1>",
		Body:  <-documentChan,
	})

}

func main() {
  fmt.Println("server on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RenderPage(w)
	})
	log.Println(http.ListenAndServe(":8080", nil))
}
