FROM ubuntu:18.04

# copy binary
COPY bin/app /
# run binary 
CMD [ "./app" ]
