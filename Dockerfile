FROM alpine:3.14

WORKDIR /app
# copy binary
COPY ./app ./app
# run binary 
CMD [ "./app" ]
