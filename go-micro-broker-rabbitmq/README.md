## Install RabbitMQ:

Run the RabbitMQ container:

    docker run --name rabbitmq1 -p 5672:5672 -p 15672:15672 -d rabbitmq
    docker exec -it rabbitmq1 /bin/bash

Config RabbitMQ:

    rabbitmq-plugins enable rabbitmq_management
    cd /etc/rabbitmq/conf.d/
    echo management_agent.disable_metrics_collector = false > management_agent.disable_metrics_collector.conf

Restart RabbitMQ:

    docker restart rabbitmq1

## Run the demo:

    go run main.go

or

    go run consumer/consumer.go
    go run publisher/publisher.go
