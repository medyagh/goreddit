FROM golang:1.6

RUN mkdir /goReddit/
WORKDIR	/goReddit/

COPY . /goReddit/

## AWS Credentials as Env variables
ENV AWS_ACCESS_KEY_ID=AKIAIDIZNSZ5LDOF7UAA
ENV AWS_SECRET_ACCESS_KEY=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

## Facebook Credentials as Env variables
ENV FACEBOOK_KEY=155292901564325
ENV FACEBOOK_SECRET=9cc0744140bfd9726491ec85faca839f

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