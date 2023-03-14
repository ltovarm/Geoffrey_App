#ifndef CLIBRARY_H
#define CLIBRARY_H
typedef void (*publish_callback_fcn)(struct mosquitto* inst, void* p, int mid);
typedef void (*connect_callback_fcn)(struct mosquitto* inst, void* p, int mid);
#endif