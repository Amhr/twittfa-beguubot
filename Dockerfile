FROM golang

WORKDIR /app
COPY . ./
RUN go mod download
RUN go get -u github.com/cosmtrek/air
ENTRYPOINT air