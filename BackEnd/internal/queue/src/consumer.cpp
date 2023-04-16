#include "queuermq.hpp"

const std::string QUEUE_NAME = "QueueService1";

int main(){

    std::cout << "Empezamos dequeue..." << std::endl;
    QueueRMQ qrmq(Consumer, QUEUE_NAME.c_str(), NULL, 0, NULL, NULL, NULL);
    if(int ret = qrmq.InitQueue()){
        std::cout << "Error iniciando queue: " << std::to_string(ret) << std::endl;
        return -1;
    }
    std::cout << "Queue inicializada con exito." << std::endl;
    std::string msg;
    for (int i = 0; i < 100000000; ++i){

        // if (qrmq.NumberOfData() < 0){
        //     std::cout << "La cola esta vacia" << std::endl;
        //     break;
        // }
        // std::cout << "i =  " << i << std::endl;
        if (int ret = qrmq.Dequeue(msg)){
            std::cout << "Error desencolando: " << std::to_string(ret) << std::endl;
            return -1;
        }
        std::cout << "Desencolando con exito: " << msg << std::endl;
        msg = "";
    }

    qrmq.CloseQueue();

    return EXIT_SUCCESS;
}

/*
// int main(){

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
    // amqp_queue_declare(conn, 1, amqp_cstring_bytes("QueueService1"), 0, 0, 0, 1, amqp_empty_table);
    // amqp_queue_bind(conn, 1, amqp_cstring_bytes("QueueService1"), amqp_cstring_bytes("my_exchange"), amqp_cstring_bytes(""), amqp_empty_table);


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
        std::string mensaje((char*)envelope.message.body.bytes, envelope.message.body.len);
        //float temperatura = std::stof(mensaje);

        std::cout << "La temperatura es: " << mensaje << std::endl;
        amqp_destroy_message(&message);
    }

    // Cerrar la conexión a RabbitMQ
    amqp_destroy_envelope(&envelope);

    amqp_basic_cancel(conn, 1, consume_ok->consumer_tag);
    free(consume_ok);
    consume_ok = NULL;
    amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
    amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
    amqp_destroy_connection(conn);

    return 0;
}
*/