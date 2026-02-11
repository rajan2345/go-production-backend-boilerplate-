# Go Boilerplate Backend

make this backend as go-module , so that dependencies can be installed 
main.go -- initialize all the other services and modules and make them run (starting the server)
handler deals with validation and passes the data the service needs (basically first layer)
second layer -- service , calls the repository methods which is a unit level methods and deals with business logic 
sqlerr -- database error holder 
validation -- logic relate to validation 
static folder -- contains openAPI json files (we will be generating this file not writing )