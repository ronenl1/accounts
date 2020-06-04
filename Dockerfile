FROM openpolicyagent/opa 

WORKDIR /app

COPY . .

CMD ["run","--server","./policies"]