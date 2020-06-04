package demo 

default allow = false

import data


allow {
    input.method == "GET"
    some id
    input.path = ["account",id]
    data.accounts[id].userName == input.userName 
}

allow {
    input.method == "GET"
    some id
    input.path = ["account",id]
    input.roles[_] == "customer-service"
    input.region  == data.accounts[id].region
}
