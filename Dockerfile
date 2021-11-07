FROM alpine:3.14

# copy binary
COPY ./bin/app /
# run binary 
CMD [ "./app" ]
