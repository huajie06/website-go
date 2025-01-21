FROM arm32v7/golang:1.20

# RUN go get -u github.com/gorilla/mux@1.8.0

# RUN mkdir /app_test
# ADD . /app_test/
WORKDIR /app_test

COPY . .

RUN go mod download

RUN go build -o /test-website

# COPY go.mod go.sum ./
# RUN go mod download

EXPOSE 5000

CMD ["/test-website"]

