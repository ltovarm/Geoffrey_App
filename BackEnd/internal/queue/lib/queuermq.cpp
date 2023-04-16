#include "queuermq.hpp"
#include <amqp_framing.h>

// Constructor
QueueRMQ::QueueRMQ(typeQueue_t typeQueue, const char* queueName, const char* host, int port, const char* vhost, const char* user, const char* pass) {
    // Initialize
    if (host!=NULL) queueCredentials.host = host;
    if (port!=0) queueCredentials.port = port;
    if (vhost!=NULL) queueCredentials.vhost = vhost;
    if (user!=NULL) queueCredentials.user = user;
    if (pass!=NULL) queueCredentials.pass = pass;
    queueCredentials.queueName = queueName;
}

QueueRMQ::~QueueRMQ(){
    if (status != 0){
        _closeconnection();
    }
}

int QueueRMQ::_CreateNewConnection(){
    // Create a RabbitMQ's connection
    conn = amqp_new_connection();
    socket = amqp_tcp_socket_new(conn);
    if (!socket) {
        std::cerr << "Error creating TCP socket." << std::endl;
        return -1;
    }
    if (int ret = amqp_socket_open(socket, queueCredentials.host, queueCredentials.port)) {
        std::cerr << "Error opening TCP socket." << std::endl;
        return ret;
    }

    status = 1;

    return EXIT_SUCCESS;
}

amqp_rpc_reply_t QueueRMQ::_Loggin(){
    amqp_rpc_reply_t login_reply = amqp_login(conn, queueCredentials.vhost, 0, 131072, 0, AMQP_SASL_METHOD_PLAIN, queueCredentials.user, queueCredentials.pass);
    if (login_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        std::cerr << "Error initializing RabbitMQ session." << std::endl;
    }
    return login_reply;
}

int QueueRMQ::_OpenChannelConn(){
    // Open a communication channel with RabbitMQ
    amqp_channel_open(conn, 1);

    if (typeQueue == Consumer){
        amqp_rpc_reply_t channel_reply = amqp_get_rpc_reply(conn);
        if (channel_reply.reply_type != AMQP_RESPONSE_NORMAL) {
            std::cerr << "Error opennig a channel of comunication with RabbitMQ" << std::endl;
            return (int)channel_reply.reply_type;
        }
    }

    return EXIT_SUCCESS; 
}
// Need to solved
int QueueRMQ::_StatusConnection(){
    amqp_frame_t frame;

    int res = amqp_simple_wait_frame(conn, &frame);
    if (res < 0) {
        // Error waiting for a frame from RabbitMQ
        std::cerr << "Error waiting for a frame from RabbitMQ" << std::endl;
        return res;
        // Perform some action in case of connection error
    } else if (frame.frame_type == AMQP_FRAME_METHOD && frame.payload.method.id == AMQP_BASIC_CANCEL_OK_METHOD) {
        // The connection with RabbitMQ has been cancelled
        std::cout << "The connection with RabbitMQ has been cancelled" << std::endl;
        // Perform some action in case of cancelled connection
        // To-do
        return -100;
    } else if (frame.frame_type == AMQP_FRAME_METHOD && frame.payload.method.id == AMQP_CONNECTION_CLOSE_METHOD) {
        // The connection with RabbitMQ has been closed
        std::cout << "The connection with RabbitMQ has been closed" << std::endl;
        // Perform some action in case of closed connection
        //To-do
        return -101;
    } else {
        // A frame has been received from RabbitMQ
        std::cout << "A frame has been received from RabbitMQ" << std::endl;
        // Perform some action after receiving the frame
    }

    return EXIT_SUCCESS;
}

amqp_rpc_reply_t QueueRMQ::_ConsumeMessage(std::string &data){

    amqp_basic_consume(conn, 1, amqp_cstring_bytes(queueCredentials.queueName), amqp_empty_bytes, 0, 1, 0, amqp_empty_table);
    amqp_envelope_t envelope;

    amqp_rpc_reply_t consume_replay = amqp_consume_message(conn, &envelope, NULL, 0);
    if (AMQP_RESPONSE_NORMAL != consume_replay.reply_type) {
        std::cerr << "Error consuming data." << std::endl;
        return consume_replay;
    }
    // Processing the received message
    // amqp_message_t message = envelope.message;

    data += std::string((char*)envelope.message.body.bytes, envelope.message.body.len);

    //float temp = std::stof(mensaje);
    // std::cout << "Data received." << std::endl;
    // amqp_destroy_message(&message);
    // close consume
    amqp_destroy_envelope(&envelope);
    //delete consume_ok;
    //consume_ok = NULL;
    return consume_replay;
}

amqp_rpc_reply_t QueueRMQ::_PublishMessage(std::string data){

    // Publish message
    amqp_basic_properties_t props;
    props._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
    props.content_type = amqp_cstring_bytes("application/octet-stream");
    props.delivery_mode = 2;
    amqp_bytes_t message_bytes = amqp_cstring_bytes(data.c_str());
    amqp_basic_publish(conn, 1, amqp_empty_bytes, amqp_cstring_bytes(queueCredentials.queueName), 0, 0, &props, message_bytes);
    amqp_rpc_reply_t publish_reply = amqp_get_rpc_reply(conn);
    if (publish_reply.reply_type != AMQP_RESPONSE_NORMAL) {
        std::cerr << "Error publishing a message in RabbitMQ" << std::endl;
    }
    return publish_reply;
}

void QueueRMQ::_closeconnection(){
    // Close connection
    amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
    amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
    amqp_destroy_connection(conn);
}

int QueueRMQ::InitQueue(){

    if(int ret = _CreateNewConnection()) return ret;
    amqp_rpc_reply_t ret = _Loggin();
    if (ret.reply_type != AMQP_RESPONSE_NORMAL) return ret.reply_type;
    if(int ret = _OpenChannelConn()) return ret;

    return EXIT_SUCCESS;
}

int QueueRMQ::Dequeue(std::string &data){

    // Check status connection To-do
    // if (_StatusConnection() != 0){
    //     std::cerr << "Error with the connection with RabbitMQ" << std::endl;
    //     return -1;
    // }

    // Dequeue message
    amqp_rpc_reply_t consume_replay = _ConsumeMessage(data);
    if (AMQP_RESPONSE_NORMAL != consume_replay.reply_type) {
        return consume_replay.reply_type;
    }

    return EXIT_SUCCESS;
}

int QueueRMQ::Queue(std::string data){

    // Check status connection
    // if (_StatusConnection() != 0){
    //     std::cerr << "Error with the connection with RabbitMQ" << std::endl;
    //     return -1;
    // }
    // Queue message
    amqp_rpc_reply_t publish_replay = _PublishMessage(data);
    if (AMQP_RESPONSE_NORMAL != publish_replay.reply_type) {
        return publish_replay.reply_type;
    }

    return EXIT_SUCCESS;
}

uint32_t QueueRMQ::NumberOfData(){
    // amqp_queue_declare_passive_ok_t *res;

    // res = amqp_queue_declare_passive(conn, 1, amqp_cstring_bytes("my_queue"), 0, 0, 0, 0, amqp_empty_table);

    // if (res) {
    //     printf("La cola tiene %d mensajes\n", res->message_count);
    // }
    // else {
    //     printf("La cola no existe\n");
    // }
    return 1;
}

void QueueRMQ::CloseQueue(){

    _closeconnection();
    status = 0;
}