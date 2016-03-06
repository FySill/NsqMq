package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
  "github.com/bitly/go-nsq"
  "flag"
  "github.com/rednut/nsq-cli/lib/nsqdflags"
)

var (
    msg = flag.String("msg", "", "message to publish")
)

func init() {
    flag.StringVar(msg, "m", *msg, "message to publish")
}



func SetCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:  "name",
		Value: "tu",
	}
	http.SetCookie(w, &cookie)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile("views/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("home.html file error %v", err), 500)
	}
	username := r.PostFormValue("username")
	pwd := r.PostFormValue("pwd")


	flag.Parse()
	nsqd.CheckFlags()
	nsqTopic   := nsqd.GetTopic()
  nsqAddress := nsqd.GetNsqdAddress();

	if "" == *msg {
      log.Fatal("ERROR: missing required 'msg' parameter");
  }



  fmt.Printf("PRODUCER: nsqd=%s, topic=%s, msg=%s\n", nsqAddress, nsqTopic, *msg)


  config := nsq.NewConfig()
  w, _ := nsq.NewProducer(nsqAddress, config)


  err := w.Publish(nsqTopic, []byte(*msg))
  if err != nil {
      log.Fatal("Could not connect ", err)
  }

  w.Stop()

  fmt.Printf("DONE\n\n")
	fmt.Fprintf(w, string(webpage))
}

func main() {




	port := 8090
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()
  mux.Handle("/", http.HandlerFunc(HomeHandler))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Print("Listening on port: " + portstring + "...")
	err := http.ListenAndServe(":"+portstring, mux)
	if err != nil {
		log.Print("ListenAndServe error: ", err)
	}
}
