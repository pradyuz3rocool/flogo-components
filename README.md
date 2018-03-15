# flogo-components
Collection custom built flogo components

## Components

### Activities
* [randomnumber](activity/randomnumber): Generate a random unique number between the min and max value
* [randomstring](activity/randomstring): Generate a random string consisting with the length you specify
* [addtodate](activity/addtodate): Adds a specified number of days to either the current date or a chosen date
* [dynamodbinsert](activity/dynamodbinsert): Insert a record in an Amazon DynamoDB
* [dynamodbquery](activity/dynamodbquery): Execute a query against an Amazon DynamoDB

### Apps
* [invoiceservice](apps/invoiceservice): A simple service listening to requests on a port exposed as environment variable, sending back random data and leveraging the [paymentservice](apps/paymentservice)
* [paymentservice](apps/paymentservice): A simple service listening to requests on a port exposed as environment variable, sending back a random date
