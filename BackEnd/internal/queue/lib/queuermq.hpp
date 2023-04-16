#include <iostream>
#include <string>
#include <sstream>
#include <amqp.h>
#include <amqp_framing.h>
#include <amqp_tcp_socket.h>

#ifndef QUEUERMQ_HPP
#define QUEUERMQ_HPP

typedef enum _typeQueue{
    Consumer,
    Publisher
}typeQueue_t;

typedef struct QueueCredentials_{
    const char* queueName = "";
    const char* host = "localhost";
    int port = 5672;
    const char* vhost = "/";
    const char* user = "guest";
    const char* pass = "guest";
}QueueCredentials_t;

class QueueRMQ{
    private:
        bool typeQueue = false; //false = Consumer true = Publisher
        int status = 0;
        amqp_connection_state_t conn;
        amqp_socket_t* socket;
        QueueCredentials_t queueCredentials;
        std::string msg;
        int _CreateNewConnection();
        amqp_rpc_reply_t _Loggin();
        int _OpenChannelConn();
        int _StatusConnection();
        amqp_rpc_reply_t _ConsumeMessage(std::string &data);
        amqp_rpc_reply_t _PublishMessage(std::string data);
        void _closeconnection();
    public:
        QueueRMQ(typeQueue_t typeQueue, const char* queueName, const char* host, int port, const char* vhost, const char* user, const char* pass);
        ~QueueRMQ();
        int InitQueue();
        int Dequeue(std::string &data);
        int Queue(std::string data); 
        uint32_t NumberOfData();
        void CloseQueue();
};

#endif // QUEUERMQ_HPP