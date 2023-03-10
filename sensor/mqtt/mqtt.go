package mqtt

// #cgo LDFLAGS: -lmosquitto
// #include <mosquitto.h>
import "C"

/// MosquittoLibInit Required before calling other mosquitto functions
func LibInit() error {
	_, err := C.mosquitto_lib_init()
	return err
}

func LibCleanup() error {
	_, err := C.mosquitto_lib_cleanup()
	return err
}

/*func New() error {
	_, err := C.mosquitto_new()
	return err
}
*/
