package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type person struct {
	name           string
	id             int
	savedImagePath string
}

var people []*person

func createImageHandler(writer http.ResponseWriter, request *http.Request) {

	params := mux.Vars(request)
	for _, p := range people {
		if p.name == params["person"] {
			outfile, err := os.Create(p.name + ".jpg")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			m := image.NewRGBA(image.Rect(0, 0, 240, 240))
			blue := color.RGBA{0, 0, 255, 255}
			draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
			var opt jpeg.Options

			opt.Quality = 80

			err = jpeg.Encode(outfile, m, &opt) // put quality to 80%
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Generated image to %s.jpg \n", p.name)
			p.savedImagePath = p.name + ".jpg"
			return
		}
	}

}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("Hello API"))
}

func viewImageHandler(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for _, p := range people {
		if p.name == params["person"] {

			imgfile, err := os.Open(p.savedImagePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			img, _ := jpeg.Decode(imgfile)

			writeImage(writer, &img)
			return

		}
	}
}


func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func accountHandler(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for _, p := range people {
		if p.name == params["person"] {
			_, err := writer.Write([]byte("Hello " + p.name + " your account number is " + strconv.Itoa(p.id) ))
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	_, err := writer.Write([]byte("Sorry I couldnt find this name: " + params["person"]))
	if err != nil {
		log.Fatal(err)
	}
	return
}


func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	_, _ = tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	_, _ = fmt.Fprintf(w, "Successfully Uploaded File\n")
}



func main() {

	jane := &person{
		name:           "jane",
		id:             0,
	}

	nick := &person{
		name: "nick",
		id:   1,
	}



	people = append(people, jane)
	people = append(people, nick)

	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/account/{person}/upload", uploadFile)

	router.HandleFunc("/account/{person}", accountHandler)
	router.HandleFunc("/account/{person}/create", createImageHandler)
	router.HandleFunc("/account/{person}/view", viewImageHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func check(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
