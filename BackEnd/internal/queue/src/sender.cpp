#include <queuermq.hpp>
#include <unistd.h> 

const std::string QUEUE_NAME = "QueueService1";


int main(){

    std::cout << "Empezamos..." << std::endl;
    QueueRMQ qrmq(Publisher, QUEUE_NAME.c_str(), NULL, 0, NULL, NULL, NULL);
    if(int ret = qrmq.InitQueue()){
        std::cout << "Error iniciando queue: " << std::to_string(ret) << std::endl;
        return -1;
    }

    for (int i = 0; i < 100000000; ++i){
        if (int ret = qrmq.Queue(std::to_string(i))){
            std::cout << "Error encolando: " << std::to_string(ret) << std::endl;
            return -1;
        }
        std::cout << "Encolando con exito: " << std::to_string(i) << std::endl;
        sleep(3);
    }

    qrmq.CloseQueue();
    
    return EXIT_SUCCESS;
}

/*
int main() {
    
    // Establecer la conexión
    amqp_connection_state_t conn = amqp_new_connection();
    amqp_socket_t* socket = amqp_tcp_socket_new(conn);
    if (amqp_socket_open(socket, "localhost", 5672)) {
        std::cerr << "Error abriendo la conexión TCP." << std::endl;
        return 1;
    }

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
    // amqp_queue_declare_ok_t* queue_declare_ok = amqp_queue_declare(conn, 1, amqp_cstring_bytes(QUEUE_NAME.c_str()), 0, 0, 0, 1, amqp_empty_table);
    // amqp_queue_declare_ok_t* queue_declare_ok = amqp_queue_declare(conn, 1, amqp_cstring_bytes("QueueService1"), 0, 0, 0, 1, amqp_empty_table);
    // if (queue_declare_ok == NULL) {
    //     std::cerr << "Error al crear una cola en RabbitMQ" << std::endl;
    //     return 1;
    // }

    // Publicar un mensaje con un valor float en la cola
    amqp_basic_properties_t props;
    props._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
    props.content_type = amqp_cstring_bytes("application/octet-stream");
    props.delivery_mode = 2;
    float temp = 25.5;
    // amqp_bytes_t message_bytes = amqp_cstring_bytes(reinterpret_cast<const char*>(&temp));
    amqp_bytes_t message_bytes = amqp_cstring_bytes("esto es una prueba");
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
*/