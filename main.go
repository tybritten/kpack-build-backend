package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	test_k8s()
	user, pass := get_logins()
	router := httprouter.New()
	router.GET("/image", BasicAuth(get_image_list, user, pass))
	router.GET("/image/:user/:code", BasicAuth(get_image, user, pass))
	router.POST("/image/:user/:code", BasicAuth(create_update_image, user, pass))
	router.DELETE("/image/:user/:code", BasicAuth(del_image, user, pass))
	log.Fatal(http.ListenAndServe(":8080", router))

}

func create_update_image(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("user") + "-" + ps.ByName("code")
	imagestatus, err := get_image_status(name)
	var req ImgReq
	reqbody, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		fmt.Fprintf(w, "Image Failed to create/update with error: "+err.Error())
		log.Println("Image Failed to create/update with error: ", err)
	}

	json.Unmarshal(reqbody, &req)
	req.Name = name
	if err != nil {
		tag := "us.gcr.io/pgtm-tbritten/" + req.Name
		if err := create_git_image(req.Name, req.Repo, req.Revision, tag); err != nil {
			fmt.Fprintf(w, "Image Failed to Create with error: "+err.Error())
			log.Println("Image Failed to Create with error: ", err)
		}
	} else {
		if err := update_git_image(req.Name, req.Repo, req.Revision); err != nil {
			fmt.Fprintf(w, "Image Failed to Update with error: "+err.Error())
			log.Println("Image Failed to Update with error: ", err)
		}
	}
	imagestatus, err = get_image_status(req.Name)
	if err != nil {
		fmt.Fprintf(w, "Failed to return image with error: "+err.Error())
		return
	}
	imgreturn, _ := json.MarshalIndent(imagestatus, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(imgreturn))
	log.Println("Image successfully created/updated", imagestatus.Metadata.GetName())
}

func get_image(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	image := ps.ByName("user") + "-" + ps.ByName("code")

	if len(image) < 1 {
		log.Println("Image Name is missing")
		fmt.Fprintf(w, "Image Name is missing")
		return
	}
	imagestatus, err := get_image_status(image)
	if err != nil {
		fmt.Fprintf(w, "Image Failed to Create with error: "+err.Error())
		return
	}
	imgreturn, _ := json.MarshalIndent(imagestatus, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(imgreturn))
	log.Println("Image successfully returned", imagestatus.Metadata.GetName())

}
func get_image_list(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	images, err := image_list()
	if err != nil {
		fmt.Fprintf(w, "Error returning Image List:"+err.Error())
		return
	}
	if len(images.Items) == 0 {
		fmt.Fprintf(w, "No Images Found")
		log.Println("No Images Found")
		return
	}
	var imagename []string
	for _, element := range images.Items {
		// index is the index where we are
		// element is the element from someSlice for where we are
		imagename = append(imagename, element.GetMetadata().GetName())
	}
	imgreturn, _ := json.MarshalIndent(imagename, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(imgreturn))
	log.Println("Image list returned")

}
func del_image(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	image := ps.ByName("user") + "-" + ps.ByName("code")

	if len(image) < 1 {
		log.Println("Image Name is missing")
		fmt.Fprintf(w, "Image Name is missing")
		return
	}
	if err := delete_image(image); err != nil {
		fmt.Fprintf(w, "Image Failed to Delete with error: "+err.Error())
		return
	}
	fmt.Fprintf(w, "Image successfully deleted:"+image)
	log.Println("Image successfully deleted:", image)
}

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()
		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Println("Unauthorized attempt to access:", r.URL, "by:", r.RemoteAddr)
		}
	}
}

type ImgReq struct {
	Name     string `json:"name"`
	Repo     string `json:"repo"`
	Revision string `json:"revision"`
}
