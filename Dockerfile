FROM golang:1.6

RUN mkdir /goReddit/
WORKDIR	/goReddit/

COPY . /goReddit/

## Install Gorilla tool Pat
RUN go get github.com/gorilla/pat
RUN go install github.com/gorilla/pat

## Install Gorilla tool Session
RUN go get github.com/gorilla/sessions
RUN go get github.com/gorilla/sessions

## Install Goth package
RUN go get github.com/markbates/goth
RUN go install github.com/markbates/goth

## Install AWS SDK 
RUN go get github.com/aws/aws-sdk-go
RUN go install github.com/aws/aws-sdk-go

# Build Go binary
RUN go build .

EXPOSE 8080
ENTRYPOINT ["./goReddit"]
CMD ["./goReddit"]