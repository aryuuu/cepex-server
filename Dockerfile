FROM alpine:3.14

RUN mkdir /app
WORKDIR /app
# copy binary
COPY ./app ./app
# run binary 
CMD [ "./app" ]
