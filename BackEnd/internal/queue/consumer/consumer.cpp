#include <iostream>
#include <string>
#include <sstream>
#include <amqp.h>
#include <amqp_framing.h>
#include <amqp_tcp_socket.h>

int main(){

    // Crear una conexión a RabbitMQ
    amqp_connection_state_t conn = amqp_new_connection();
    amqp_socket_t* socket = amqp_tcp_socket_new(conn);
    if (!socket) {
    std::cerr << "Error creando el socket TCP." << std::endl;
    }

    if (amqp_socket_open(socket, "localhost", 5672)) {
        std::cerr << "Error abriendo la conexión TCP." << std::endl;
        return 1;
    }
    amqp_rpc_reply_t login_reply = amqp_login(conn, "/", 0, 131072, 0, AMQP_SASL_METHOD_PLAIN, "guest", "guest");
    if (login_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        std::cerr << "Error al iniciar sesión en RabbitMQ" << std::endl;
        std::cerr << "Error: " << login_reply.reply_type << std::endl;
        return 1;
    }
    amqp_channel_open(conn, 1);

    // Crear una cola y un intercambio en RabbitMQ
    // amqp_exchange_declare(conn, 1, amqp_cstring_bytes("my_exchange"), amqp_cstring_bytes("fanout"), 0, 0, 0, 0, amqp_empty_table);
    amqp_queue_declare(conn, 1, amqp_cstring_bytes("QueueService1"), 0, 0, 0, 1, amqp_empty_table);
    amqp_queue_bind(conn, 1, amqp_cstring_bytes("QueueService1"), amqp_cstring_bytes("my_exchange"), amqp_cstring_bytes(""), amqp_empty_table);


    amqp_basic_consume_ok_t *consume_ok = amqp_basic_consume(conn, 1, amqp_cstring_bytes("QueueService1"), amqp_empty_bytes, 0, 1, 0, amqp_empty_table);

    amqp_envelope_t envelope;

    while (1) {
        amqp_rpc_reply_t ret = amqp_consume_message(conn, &envelope, NULL, 0);
        if (AMQP_RESPONSE_NORMAL != ret.reply_type) {
            break;
        }

        amqp_message_t message = envelope.message;

        // Procesar el mensaje recibido
        // ...

        // Convertir la cadena de caracteres del mensaje a un parámetro float
        std::string mensaje((char*)envelope.message.body.bytes, (char*)envelope.message.body.bytes + envelope.message.body.len);
        float temperatura = std::stof(mensaje);

        std::cout << "La temperatura es: " << temperatura << std::endl;

        // Cerrar la conexión a RabbitMQ
        amqp_destroy_message(&message);
        amqp_destroy_envelope(&envelope);
    }

    amqp_basic_cancel(conn, 1, consume_ok->consumer_tag);
    free(consume_ok);
    consume_ok = NULL;
    amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
    amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
    amqp_destroy_connection(conn);

    return 0;
}