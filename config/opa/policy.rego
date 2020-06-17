package demo 

default allow = false

import data

allow {
    input.method == "GET"
    some id
    input.path = ["accounts",id]
    data.accounts[ID].username == input.userName 
}

allow {
    input.method == "GET"
    some id
    input.path = ["accounts",id]
    input.roles[_] == "customer-service"
    input.region  == data.accounts[id].region
}
