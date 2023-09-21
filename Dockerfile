FROM golang:1.20.4

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o golang-tg_cloud ./main.go

CMD ["./golang-tg_cloud"]

# docker build -t golang-tg_cloud .
# docker run --name=golang-tg_cloud -p 8080:8080 golang-tg_cloud