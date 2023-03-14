package mosquitto

/*
#include <stdio.h>

extern char* mosquitto_connack_string(int);

// The gateway function
void publish_callback_wrapper_cgo(struct mosquitto* inst, void* p, int mid)
{
	printf("Message published callback: message id = %d\n", mid);
	void publish_callback_wrapper(struct mosquitto*, void*, int);
	publish_callback_wrapper(inst, p, mid);
}

void connect_callback_wrapper_cgo(struct mosquitto* inst, void* p, int reason_code)
{
	printf("on_connect: %s\n", mosquitto_connack_string(reason_code));
	if(reason_code != 0){
		// You may wish to set a flag or something (context?) here to indicate to your application that the
		// client is now connected.
	}
	void connect_callback_wrapper(struct mosquitto*, void*, int);
	connect_callback_wrapper(inst, p, reason_code);
}
*/
import "C"
