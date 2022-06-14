FROM golang
MAINTAINER Stepan K. <xamust@gmail.com>
WORKDIR /service/
VOLUME ["/opt/service"]
COPY . ./
RUN make build
EXPOSE 4223
EXPOSE 8223
EXPOSE 8080
EXPOSE 5432
CMD [ "build/service" ]