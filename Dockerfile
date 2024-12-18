FROM golang:1.23.2 AS build

WORKDIR /boombox

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./dist/app 

FROM golang:1.23.2 AS runner

COPY --from=build /boombox/dist/app /usr/bin/app

CMD ["/usr/bin/app"]
