# Tarea 2
Desarrollo de tarea 2 de sistemas distribuidos.

El inicio de esta tarea comienza con el final de la primera en el cual terminamos con un sistema de proteccion que funciona pero no puede atender todas las peticiones.

Para comenzar a solucionar este problema dividimos nuestra aplicacion en 3 componentes:

* Cliente (VM 1)
* Servicio de proteccion (VM 2)
* Servicio de registros y mensajeria (VM 3)

Nuestra solucion debera ser capaz de:
 
* Realizar la protección de archivos PDF ya desarrollada
* Guardar  registros  de  las  operaciones  y  sus  resultados  una  vez  hecho  el  procesamiento del archivo. 
* Notificar al cliente para el cual el archivo está siendo protegido, a través de un mensaje por correo electrónico.


# Ejecucion

## Correr cliente

Correr en la primera maquina virtual (VM 1)

```shell
go run ./cliente <rut> <correo> <ruta>
```

## Correr servicio proteccion

Correr en la segunda maquina virtual (VM 2)
```shell
go run ./proteccion
```

## Correr servicio registros
Correr en la tercera maquina virtual (VM 3)


```shell
go run ./registros
```

## Correr servicio mensajeria
Correr en la tercera maquina virtual (VM 3)

```shell
go run ./mensajeria
```


# Instalacion y configuracion

Todas las maquinas
```shell
# GOLANG
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz && sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

Instalacion VM 2 (jammy)
```shell
# MongoDB
sudo apt update
sudo apt install gnupg wget apt-transport-https ca-certificates software-properties-common
wget -qO- \
  https://pgp.mongodb.com/server-7.0.asc | \
  gpg --dearmor | \
  sudo tee /usr/share/keyrings/mongodb-server-7.0.gpg >/dev/null

echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] \
  https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/7.0 multiverse" | \
  sudo tee -a /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install mongodb-org

# Comandos Mongo
sudo systemctl start mongod
sudo systemctl daemon-reload
sudo systemctl status mongod
sudo systemctl stop mongod
sudo systemctl restart mongod
netstat -plntu
mongosh

# gRPC
sudo apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+
```


# Contenido

## Interface Description Language (IDL)

Describe una interfaz en un lenguaje neutral, que permite la comunicación entre componentes de software desarrollados en diferentes lenguajes de programación, como por ejemplo entre componentes escritos en C++ y otros escritas en Java.

### Usos

* Son utilizadas con frecuencia en el software de las llamadas a procedimiento remoto (RPC, Remote Procedure Call), lo que permite a los sistemas de computadoras utilizar lenguajes y sistemas operativos diferentes. IDL ofrece un puente entre dos sistemas diferentes.

## gRPC

* [gRPC explicado - Youtube](https://www.youtube.com/watch?v=NHw2cjcMN9g&t=60s)
* [Introduction to gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
* [Tutorial basico gRPC en Go](https://grpc.io/docs/languages/go/basics/)

## Protobuffers 

El código generado por los protocol buffers proporciona métodos de utilidad para recuperar datos de archivos y flujos, extraer valores individuales de los datos, verificar si existen datos, serializar datos de nuevo a un archivo o flujo, y otras funciones útiles.

![alt text](https://protobuf.dev/images/protocol-buffers-concepts.png)


* [Lo basico de protobuffer en Go](https://protobuf.dev/getting-started/gotutorial/)
