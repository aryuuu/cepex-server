FROM alpine:3.14

WORKDIR /app
# copy binary
COPY bin/app ./app
# run binary 
CMD [ "./app" ]
