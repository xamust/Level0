FROM golang
MAINTAINER Stepan K. <xamust@gmail.com>
WORKDIR /client/
VOLUME ["/opt/client"]
COPY . ./
RUN make build
EXPOSE 4223
CMD [ "build/client" ]