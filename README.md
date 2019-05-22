## **Pre-paid card Challenge**

I built this as a REST API as the spec didn't specify otherwise. 
 
 **Requirements**
 Go
 Postgres
  

**Usage**
Edit postgress database connection details as needed in the *.env* file and the run ./go-card

All requests to the API should be made as "application/json"
See bellow for the endpoints and example JSON for their use.
All monetary values are handled as integers for the sake of simplicity. They can then be turned into a decimal when needed on the frontend.


**End Points**
*"/api/card/new"*
Creates a new Card account in the system. 
cardNo: string. must be unique
password: string. Is set and saved in the DB as a hash

    {"cardNo":"123456", "password":"TestPassword"}


*"/api/card/addfunds"*
Adds funds to the pre-paid account.
cardNo: string  must be a valid card Number in the DB
funds: integer 

    {"cardNo":"123456", "funds":12}

*"/api/card/balance"*
Allows the account holder to view the current balance of funds on their account.
cardNo: string
password: string. the same as was set during the account creation 

    {"cardNo":"123456", "password":"testng"}



*"/api/transaction/new"*
Creates a new transaction. Blocking out funds on the user's account
funds: integer 
merchant: string
amount: integer

    {"cardNo":"123456", "merchant":"Steve", "amount":12}


*"/api/transaction/capture"*
Allows merchant to capture the available funds from a transaction.
transactionId: integer. found in the response creating a transaction
amount: integer. amount to capture.

    {"transactionId":1, "amount":12}

"/api/transaction/change"
Allows the merchant to reduce the amount of uncaptured authorised funds for a transaction
transactionId: integer.
amount: integer.

    {"transactionId":1, "amount":12}


"/api/transaction/refund"
Allows merchant to refund all or part of the captured funds of transaction.
transactionId: integer.
amount: integer.

    {"transactionId":1, "amount":12}




**Todo**
Sadly, I ran out of time to get these points done:
 - better authentication
 - Authentication for merchants
 - proper use of tokens to validate transactions
 - better responses back to the user
 - Better validation of inputs
 - Generate account statement (the data is all going into the DB, so it just needs a means of being output)
