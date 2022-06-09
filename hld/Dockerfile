FROM ubuntu:20.04

ENV TZ=Europe/London
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get -y update && apt-get -y upgrade && apt-get -y install python3-pip
RUN apt-get update &&    apt-get install -y cmake gcc make git g++
RUN apt -y install libevent-dev
RUN apt -y install libconfig++-dev

ADD . .

RUN chmod +x build.sh && ./build.sh

EXPOSE 8080

CMD ./ws conf.cfg