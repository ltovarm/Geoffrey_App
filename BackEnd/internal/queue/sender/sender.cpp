#include <amqp.h>
#include <amqp_tcp_socket.h>
#include <amqp_framing.h>
#include <iostream>
#include <string>

const std::string QUEUE_NAME = "QueueService1";

int main() {
    
    // Establecer la conexión
    amqp_connection_state_t conn = amqp_new_connection();
    amqp_socket_t* socket = amqp_tcp_socket_new(conn);
    int sockfd = amqp_socket_open(socket, "localhost", 5672);
    amqp_tcp_socket_set_sockfd(socket, sockfd);

    // Iniciar sesión en RabbitMQ
    amqp_rpc_reply_t login_reply = amqp_login(conn, "/", 0, 131072, 0, AMQP_SASL_METHOD_PLAIN, "guest", "guest");
    if (login_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        std::cerr << "Error al iniciar sesión en RabbitMQ" << std::endl;
        std::cerr << "Error: " << login_reply.reply_type << std::endl;
        return 1;
    }
    // Abrir un canal de comunicación con RabbitMQ
    amqp_channel_open(conn, 1);
    amqp_rpc_reply_t channel_reply = amqp_get_rpc_reply(conn);
    if (channel_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        std::cerr << "Error al abrir un canal de comunicación con RabbitMQ" << std::endl;
        return 1;
    }

    // Declarar la cola
    amqp_queue_declare_ok_t* queue_declare_ok = amqp_queue_declare(conn, 1, amqp_cstring_bytes(QUEUE_NAME.c_str()), 0, 0, 0, 1, amqp_empty_table);
    if (queue_declare_ok == NULL) {
        std::cerr << "Error al crear una cola en RabbitMQ" << std::endl;
        return 1;
    }

    // Publicar un mensaje con un valor float en la cola
    amqp_basic_properties_t props;
    props._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
    props.content_type = amqp_cstring_bytes("application/octet-stream");
    props.delivery_mode = 2;
    float temp = 25.5;
    amqp_bytes_t message_bytes = amqp_cstring_bytes(reinterpret_cast<const char*>(&temp));
    amqp_basic_publish(conn, 1, amqp_empty_bytes, amqp_cstring_bytes("QueueService1"), 0, 0, &props, message_bytes);
    amqp_rpc_reply_t publish_reply = amqp_get_rpc_reply(conn);
    if (publish_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        std::cerr << "Error al publicar un mensaje en RabbitMQ" << std::endl;
        return 1;
    }

    // Cerrar la conexión
    amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
    amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
    amqp_destroy_connection(conn);
    
    return 0;
}
