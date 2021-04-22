package web

import (
	"encoding/json"
	"fmt"
	"log"
	"my-app/cors"
	"my-app/db"
	"my-app/model"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	d        db.DB
	handlers map[string]http.HandlerFunc
}

const productsPath = "Plant"

func NewApp(d db.DB, corsbool bool) App {
	app := App{
		d:        d,
		handlers: make(map[string]http.HandlerFunc),
	}
	techHandler := app.GetPlants
	plantsHandler := app.handlePlants
	plantHandler := app.handlePlant
	if !corsbool {
		//抓cors內的middleware,在拋給ＡＰＩ之前先做一步處理
		techHandler = cors.Middleware(techHandler)
		plantsHandler = cors.Middleware(plantsHandler)
		plantHandler = cors.Middleware(plantHandler)
		log.Println("add middleware")
	}
	app.handlers["/api/technologies"] = techHandler
	app.handlers["/api/Plant"] = plantsHandler
	app.handlers["/api/Plant/"] = plantHandler
	app.handlers["/"] = http.FileServer(http.Dir("/webapp")).ServeHTTP
	return app
}

//啟動時執行
func (a *App) Serve() error {
	for path, handler := range a.handlers {
		http.Handle(path, handler)
	}
	log.Println("Web server is available on port 8080")
	return http.ListenAndServe(":8080", nil)
}

func (a *App) GetPlants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	technologies, err := a.d.GetPlants()
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(technologies)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}
}
func (a *App) handlePlants(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		product, err := a.d.GetPlants()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if product == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

		/*case http.MethodPut:
		var product model.SysPlant
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if *product.Id != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateProduct(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}*/

	case http.MethodPost:
		var plant model.SysPlant

		/*err := json.NewDecoder(r.Body)*/
		err := json.NewDecoder(r.Body).Decode(&plant)
		//product.Plantname = r.PostFormValue("plantName")
		//product.Plantname = json.NewDecoder(r.Body)
		//product.Plantcode = r.PostFormValue("plantCode")
		//product.Plantdesc = r.PostFormValue("plantDesc")

		//fmt.Println(plant.Plantdesc)

		if err != nil {
			log.Print(err)
			fmt.Println("test")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		plantID, err := a.d.InsertPlant(plant)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"productId":%d}`, plantID)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (a *App) handlePlant(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", productsPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		product, err := a.d.GetPlant(ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if product == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

		/*case http.MethodPut:
		var product model.SysPlant
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if *product.Id != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateProduct(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}*/
	case http.MethodDelete:
		err = a.d.RemovePlant(ID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var plant model.SysPlant

		/*err := json.NewDecoder(r.Body)*/
		err := json.NewDecoder(r.Body).Decode(&plant)
		//product.Plantname = r.PostFormValue("plantName")
		//product.Plantname = json.NewDecoder(r.Body)
		//product.Plantcode = r.PostFormValue("plantCode")
		//product.Plantdesc = r.PostFormValue("plantDesc")

		//fmt.Println(plant.Plantdesc)

		if err != nil {
			log.Print(err)
			fmt.Println("test")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		plantID, err := a.d.InsertPlant(plant)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"productId":%d}`, plantID)))
	case http.MethodPut:
		//update
		var plant model.SysPlant
		log.Println("update test")
		err := json.NewDecoder(r.Body).Decode(&plant)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(plant.PlantID)
		if plant.PlantID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.d.UpdatePlant(ID, plant)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}

/*// Needed in order to disable CORS for local development
func disableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h(w, r)
	}
}*/
