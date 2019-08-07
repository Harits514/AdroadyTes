package app

import (
	"github.com/gorilla/mux"
)

//RouterStat : struct isi variabel router
type RouterStat struct {
	Route *mux.Router
}

//Router : u/ inisialisasi ruter
var Router RouterStat

func init() {
	Router.Route = mux.NewRouter()
	Router.Route.HandleFunc("/InsertDriver", InsertData)
	Router.Route.HandleFunc("/UpdateDriverInfo", UpdateData)
	Router.Route.HandleFunc("/DeleteDriverData", DeleteData)
	Router.Route.HandleFunc("/UpdateDriverPhoto", UpdatePhoto)
	Router.Route.HandleFunc("/DriverPhoto", ViewPhoto)
	Router.Route.HandleFunc("/DriverData", LoadData)
}
