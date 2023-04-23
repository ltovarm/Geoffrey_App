# First approximation
This is a first approximation, we have a single queue where the sender is listening on port 3000. When it receives a message it publishes it in the QueueService1 queue. The consumer receives all the messages in the queue.

To run this queue it is necessary to run the dockers using the command 
```
make run
```
To stop and delete the containers use the command
```
make stop 
```
# TO-DO
    Añadir exchange para las colas
    Configurar las colas desde un json
    Levantar tantos procesos como colas tenga
    Recibir datos desde el script de data-collector
    Mandar datos al container process-data
