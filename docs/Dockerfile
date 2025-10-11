FROM python:3.14-alpine
RUN apk add git
WORKDIR /opt
COPY . /opt
RUN pip3 install --no-cache-dir  -r requirements.txt
EXPOSE 8000
CMD ["mkdocs", "serve", "--dev-addr", "0.0.0.0:8000"]