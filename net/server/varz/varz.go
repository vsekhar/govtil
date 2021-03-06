// Package varz provides a simple varz implementation
package varz

import (
	"io"
	"net/http"

	"github.com/vsekhar/govtil/log"
	"github.com/vsekhar/govtil/net/server/multihandler"
)

// A function that writes varz data and returns an error
type VarzFunc func(io.Writer) error

type subVarzHandler struct {
	VarzFunc
	string
}

func (svh *subVarzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := svh.VarzFunc(w)
	if err != nil {
		log.Println("Varz failed:", svh.string, ",", err)
	}
}

// VarzHandler is an http.Handler that aggregates varz values from any number
// of registered VarzFunc's
type varzHandler struct {
	multihandler.MultiHandler
}

// Create a new varz handler, an http.Handler that aggregates varz responses
// from a number of registered VarzFunc's
func NewHandler() *varzHandler {
	return &varzHandler{}
}

// Register a VarzFunc to be included in varz output
func (vh *varzHandler) Register(vf VarzFunc, name string) {
	svh := subVarzHandler{vf, name}
	vh.MultiHandler.Register(&svh)
}

// Helper function to write a single key-value pair
func Write(k string, v string, w io.Writer) (err error) {
	_, err = w.Write([]byte(k + "=" + v + "\n"))
	return err
}

// Helper function to write a map of values (often the way varz's are stored)
func WriteMap(m map[string]string, w io.Writer) (err error) {
	for k, v := range m {
		err = Write(k, v, w)
		if err != nil {
			return
		}
	}
	return
}
