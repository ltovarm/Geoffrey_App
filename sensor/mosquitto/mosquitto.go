// / mosquitto V3.1.1 of MQTT protocol
package mosquitto // https://mosquitto.org/api/files/mosquitto-h.html

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -lmosquitto
#include <mosquitto.h>
#include <stdlib.h>

#include "clibrary.h"

void publish_callback_wrapper_cgo(struct mosquitto* inst, void* p, int mid);
void connect_callback_wrapper_cgo(struct mosquitto* inst, void* p, int mid);
*/
import "C"

//export publish_callback_wrapper
func publish_callback_wrapper(mosq *C.struct_mosquitto, p unsafe.Pointer, mid int) {
	fmt.Printf("Go.callOnMeGo(): called with arg = %d\n", mid)
}

//export connect_callback_wrapper
func connect_callback_wrapper(mosq *C.struct_mosquitto, p unsafe.Pointer, mid int) {
	fmt.Printf("Mosquitto client connected. mid = %d\n", mid)
}

type mosqErrT int8

/* Enum: mosq_err_t_([a-zA-Z])
 * Integer values returned from many libmosquitto functions.
 * TODO Move to go style
 */
const (
	mosqErrAuthContinue         mosqErrT = -4
	mosqErrNoSubscribers        mosqErrT = -3
	mosqErrSubExists            mosqErrT = -2
	mosqErrConnPending          mosqErrT = -1
	mosqErrSuccess              mosqErrT = 0
	mosqErrNomem                mosqErrT = 1
	mosqErrProtocol             mosqErrT = 2
	mosqErrInval                mosqErrT = 3
	mosqErrNoConn               mosqErrT = 4
	mosqErrConnRefused          mosqErrT = 5
	mosqErrNotFound             mosqErrT = 6
	mosqErrConnLost             mosqErrT = 7
	mosqErrTls                  mosqErrT = 8
	mosqErrPayloadSize          mosqErrT = 9
	mosqErrNotSupported         mosqErrT = 10
	mosqErrAuth                 mosqErrT = 11
	mosqErrAclDenied            mosqErrT = 12
	mosqErrUnknown              mosqErrT = 13
	mosqErrErrno                mosqErrT = 14
	mosqErrEai                  mosqErrT = 15
	mosqErrProxy                mosqErrT = 16
	mosqErrPluginDefer          mosqErrT = 17
	mosqErrMalformedUtf8        mosqErrT = 18
	mosqErrKeepalive            mosqErrT = 19
	mosqErrLookup               mosqErrT = 20
	mosqErrMalformedPacket      mosqErrT = 21
	mosqErrDuplicateProperty    mosqErrT = 22
	mosqErrTlsHandshake         mosqErrT = 23
	mosqErrQosNotSupported      mosqErrT = 24
	mosqErrOversizePacket       mosqErrT = 25
	mosqErrOcsp                 mosqErrT = 26
	mosqErrTimeout              mosqErrT = 27
	mosqErrRetainNotSupported   mosqErrT = 28
	mosqErrTopicAliasInvalid    mosqErrT = 29
	mosqErrAdministrativeAction mosqErrT = 30
	mosqErrAlreadyExists        mosqErrT = 31
)

type Mosquitto struct {
	instance *C.struct_mosquitto
}

func ConnackString(conn_code int) string {
	return C.GoString(C.mosquitto_connack_string(C.int(conn_code)))
}

func Strerror(mqtt_errno int) string {
	return C.GoString(C.mosquitto_strerror(C.int(mqtt_errno)))
}

// / MosquittoLibInit Required before calling other mosquitto functions
func LibInit() error {
	_, err := C.mosquitto_lib_init()
	return err
}

func LibCleanup() error {
	_, err := C.mosquitto_lib_cleanup()
	return err
}

// TODO add callback function
func New(id string, clean_session bool) (Mosquitto, error) {
	var id_conn *C.char = nil

	if len(id) != 0 {
		id_conn = C.CString(id)
	}

	mqtt_client, err := C.mosquitto_new(id_conn, C.bool(clean_session), nil)
	C.free(unsafe.Pointer(id_conn))
	if err != nil {
		return Mosquitto{}, err
	}

	return Mosquitto{instance: mqtt_client}, nil
}

func (mqtt *Mosquitto) Destroy() {
	C.mosquitto_destroy(mqtt.instance)
}

func (mqtt *Mosquitto) UsernamePwSet(user string, password string) error {

	_user := C.CString(user)
	defer C.free(unsafe.Pointer(_user))
	_password := C.CString(password)
	defer C.free(unsafe.Pointer(_password))

	raw_err, err := C.mosquitto_username_pw_set(mqtt.instance, _user, _password)
	if err != nil {
		return err
	}
	if raw_err != C.int(mosqErrSuccess) {
		return errors.New(fmt.Sprint("Mosquitto error: ", Strerror(int(raw_err))))
	}

	return nil
}

// TODO context?
func (mqtt *Mosquitto) Connect(host string, port int, keepalive int) error {
	_host := C.CString(host)
	defer C.free(unsafe.Pointer(_host))

	raw_err, err := C.mosquitto_connect(mqtt.instance, _host, C.int(port), C.int(keepalive))
	if err != nil {
		return err
	}
	if raw_err != C.int(mosqErrSuccess) {
		// TODO translate
		return errors.New(fmt.Sprint("Mosquitto error: ", Strerror(int(raw_err))))
	}

	return nil
}

func (mqtt *Mosquitto) Disconnect() error {
	raw_err := int(C.mosquitto_disconnect(mqtt.instance))
	if raw_err != int(mosqErrSuccess) {
		// TODO translate
		return errors.New(fmt.Sprint("Mosquitto error: ", Strerror(int(raw_err))))
	}
	return nil
}

func (mqtt *Mosquitto) LoopStart() error {
	raw_err := int(C.mosquitto_loop_start(mqtt.instance))
	if raw_err != int(mosqErrSuccess) {
		// TODO translate
		return errors.New(fmt.Sprint("Mosquitto error: ", Strerror(int(raw_err))))
	}
	return nil
}

func (mqtt *Mosquitto) LoopStop(force bool) error {
	raw_err := int(C.mosquitto_loop_stop(mqtt.instance, C.bool(force)))
	if raw_err != int(mosqErrSuccess) {
		return errors.New(fmt.Sprintf("Mosquitto error: ", Strerror(int(raw_err))))
	}
	return nil
}

func (mqtt *Mosquitto) Publish(topic string, payloadlen int, payload unsafe.Pointer, qos int, retain bool) (int, error) {
	var mid C.int = 0

	_topic := C.CString(topic)
	defer C.free(unsafe.Pointer(_topic))

	raw_err, err := C.mosquitto_publish(mqtt.instance, &mid, _topic, C.int(payloadlen), payload, C.int(qos), C.bool(retain))
	if err != nil {
		return 0, err
	}
	if raw_err != C.int(mosqErrSuccess) {
		// TODO translate
		return 0, errors.New(fmt.Sprint("Mosquitto error: ", Strerror(int(raw_err))))
	}

	return int(mid), nil
}

// PublishCallbackSet
func (mqtt *Mosquitto) PublishCallbackSet(callback func(Mosquitto, unsafe.Pointer, int)) {
	//_publishCallbackFunc = callback
	C.mosquitto_publish_callback_set(mqtt.instance, (C.publish_callback_fcn)(unsafe.Pointer(C.publish_callback_wrapper_cgo)))
}

// PublishCallbackSet
func (mqtt *Mosquitto) ConnectCallbackSet(callback func(Mosquitto, unsafe.Pointer, int)) {
	//_publishCallbackFunc = callback
	C.mosquitto_connect_callback_set(mqtt.instance, (C.connect_callback_fcn)(unsafe.Pointer(C.connect_callback_wrapper_cgo)))
}
