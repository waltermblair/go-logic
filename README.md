# go-logic
This is the configurable component for dockerized programmable logic device. 

https://hub.docker.com/r/waltermblair/go-logic/

message-passing via RabbitMQ

Data Flow
```mermaid
graph LR
    rmq[incoming message]-->prs((Process Message))
    prs-- Configuration Message -->cfg((Apply Config))
    cfg-->update[Updated Component Configuration]
    prs-- Input Message -->fn((Apply Logic))
    fn-->msg((Build Message))
    msg-->pub((Publish Message))
    pub-->out[outgoing message]
```

Sequence Diagram
```mermaid
sequenceDiagram 
    participant rmq as RabbitMQ
    participant client as RabbitMQ Client
    participant prs as Process()
    participant cfg as ApplyConfig()
    participant fn as ApplyLogic()
    participant msg as BuildMessage()
    
    rmq->>client: Incoming message
    client-->rmq: Acknowledge message
    client->>prs: Process message
    alt Config Message
        prs->>cfg: Apply configuration
    else Input Message
        prs->>fn: Apply logic function to transform input
        fn->>msg: Build output message for each recipient
        msg->>prs: Return each output message to Process()
        prs->>client: Publish output message(s)
        client->>rmq: Outgoing message
    end
```

