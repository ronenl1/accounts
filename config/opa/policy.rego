package demo 

default allow = false

import data

# Allow customers to access their own accounts
allow {
    input.method == "GET"
    some id
    input.path = ["accounts",id]
    input.username == data.accounts[id].username 
}

# Allow customer service to access accounts of customers
allow {
    input.method == "GET"
    some id
    input.path = ["accounts",id]
    input.roles[_] == "customer-service"
    input.region == data.accounts[id].region
}
