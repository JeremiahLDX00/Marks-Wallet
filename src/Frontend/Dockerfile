FROM nginx:1.17-alpine

COPY ./module.html /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/

EXPOSE 9070

# (Optional) environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
#ENV HTTP_PORT=8081

# Run
CMD [ "/Frontend" ]