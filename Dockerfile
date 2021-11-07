FROM alpine:3.14

WORKDIR /app
# copy binary
COPY ./bin/app /app/bin/app
# run binary 
CMD [ "./app/bin/app" ]
