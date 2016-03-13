package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
)

var (
	addr         string
	videoDir     string
	heartbeatDir string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "hostname:port")
	flag.StringVar(&videoDir, "video", "video", "video dir")
	flag.StringVar(&heartbeatDir, "heartbeat", "heartbeat", "heartbeat dir")
	flag.Parse()
}

func basic2(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				log.Println(r)
				http.Error(w, "Error", http.StatusInternalServerError)
			}
		}()

		log.Println(r.URL, r.Method)
		handler.ServeHTTP(w, r)
	})
}

func basic(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return basic2(http.HandlerFunc(handler))
}

func uploaheartbeat(w http.ResponseWriter, r *http.Request) {

	tag := r.FormValue("tag")

	f, _, err := r.FormFile("file")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	start, err := strconv.Atoi(data[0][0])
	if err != nil {
		panic(err)
	}

	originalDir := filepath.Join(heartbeatDir)
	os.MkdirAll(originalDir, os.ModePerm)

	fn := fmt.Sprintf("%s.csv", tag)
	wf, err := os.Create(path.Join(originalDir, fn))
	defer wf.Close()
	log.Println("Upload", fn)
	wcsv := csv.NewWriter(wf)
	for _, d := range data {
		n, err := strconv.Atoi(d[0])
		if err != nil {
			panic(err)
		}

		wcsv.Write([]string{
			fmt.Sprintf("%d", n-start),
			d[1],
		})
	}
	wcsv.Flush()
	//io.Copy(wf, f)

	// デバッグ用
	if r.FormValue("debug") != "" {
		http.Redirect(w, r, "/debug.html", http.StatusFound)
	}
}

func uploadvideo(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", r)
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	}()

	tag := r.FormValue("tag")

	f, _, err := r.FormFile("file")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	originalDir := filepath.Join(videoDir, "original", tag)
	os.MkdirAll(originalDir, os.ModePerm)

	list, err := ioutil.ReadDir(originalDir)
	if err != nil {
		panic(err)
	}

	fn := fmt.Sprintf("%s_%d.mp4", tag, len(list)+1)
	wf, err := os.Create(path.Join(originalDir, fn))
	defer wf.Close()
	log.Println("Upload", fn)
	io.Copy(wf, f)

	catfn := filepath.Join(videoDir, tag+".mp4")
	if len(list) <= 0 {
		log.Println("First vide upload")
		if err := exec.Command("cp", filepath.Join(originalDir, fn), catfn).Run(); err != nil {
			panic(err)
		}
	} else {
		log.Println("video concat")
		// MP4Box.exe -add "in1.mp4" -cat "in2.mp4" -new "out.mp4"
		if err := exec.Command("MP4Box", "-add", catfn, "-cat", filepath.Join(originalDir, fn), "-new", catfn+"_").Run(); err != nil {
			panic(err)
		}
		os.Rename(catfn+"_", catfn)
	}

	log.Println("convert video to images")
	imgDir := filepath.Join("static", "image", tag)
	os.MkdirAll(imgDir, os.ModePerm)
	// ffmpeg -i 元動画.avi -ss 144 -t 148 -r 24 -f image2 %06d.jpg
	if err := exec.Command("ffmpeg", "-i", catfn, "-r", "30", "-f", "image2", filepath.Join(imgDir, "%06d.jpg")).Run(); err != nil {
		panic(err)
	}

	// デバッグ用
	if r.FormValue("debug") != "" {
		http.Redirect(w, r, "/debug.html", http.StatusFound)
	}
}

func imagecount(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", r)
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	}()

	tag := r.FormValue("tag")
	imgDir := filepath.Join("static", "image", tag)
	list, err := ioutil.ReadDir(imgDir)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%d", len(list))
}

func main() {

	http.Handle("/api/uploadvideo", basic(uploadvideo))
	http.Handle("/api/uploadheartbeat", basic(uploaheartbeat))
	http.Handle("/api/imagecount", basic(imagecount))
	http.Handle("/", basic2(http.FileServer(http.Dir("static"))))

	log.Println(addr)
	http.ListenAndServe(addr, nil)
}
