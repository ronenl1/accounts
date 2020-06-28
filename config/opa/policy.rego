package demo 

default allow = false

import data

# Allow costumers to access their own accounts
allow {
    input.method == "GET"
    some id
    input.path = ["accounts",id]
    data.accounts[id].username == input.username
}

# Allow support to access accounts of costumers
allow {
    input.method == "GET"
    some id
    input.path = ["accounts",id]
    input.roles[_] == "customer-service"
    input.region == data.accounts[id].region
}
