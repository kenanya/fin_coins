## A. Description
This is a backend service (API and database) for a generic wallet service. This service contains 2 entities:
- Account
- Payment

Tech Stack: golang, go-kit framework, PostgreSQL, docker

## B. How to set up and start the backend server
#### Clone the Project
```bash
git clone https://github.com/kenanya/fin_coins.git
```

#### Assign Config Variable
These are config variables that has been registered. Each of the variable already has a default value, but each value can be replaced according to your configuration. 
- PORT
- DB_HOST
- DB_PORT
- DB_USER
- DB_PASSWORD
- DB_NAME
- DB_SCHEMA_NAME

#### Build Project
Change your current directory to the project directory.
```bash
docker-compose -f docker-compose-local.yml build
```

#### Run Project
```bash
docker-compose -f docker-compose-local.yml up
```

## C. Unit Test
```bash
go test .\repository\ -v
```


## D. The API documentation
Below are the sample requests and expected responses for each API:

### D.1. Create Account
**Purpose:** Creating account
**Endpoint:** *{BASE_URL}/account/v1/account*
| Param | Value |
| ------ | ------ |
| id | unique id for account  
| balance | float/decimal  
| currency | currency, such as USD, IDR, etc  

#### Positive Scenario 1 - Request
```
curl --location --request POST '{BASE_URL}/account/v1/account' \
--header 'Content-Type: application/json' \
--data-raw '{
   "id": "panda101",
   "balance": 2500,
   "currency": "USD"
}'
```

#### Positive Scenario 1 - Response
```
200 Ok
Content-Type: "application/json"

{
    "account": {
        "id": "panda101",
        "balance": 2500,
        "currency": "USD",
        "created_at": "2022-02-09T13:11:02.604785+07:00",
        "updated_at": "2022-02-09T13:11:02.604785+07:00"
    }
}
```

### D.2. Get All Account
**Purpose:** Listing all available accounts
**Endpoint:** *{BASE_URL}/account/v1/account*
| Param | Value |
| ------ | ------ |
|  |  |

#### Positive Scenario 1 - Request
```
curl --location --request GET '{BASE_URL}/account/v1/account' \
--data-raw ''
```

#### Positive Scenario 1 - Response
```
200 Ok
Content-Type: "application/json"

{
    "accounts": [
        {
            "id": "alice456",
            "balance": 2100,
            "currency": "USD",
            "created_at": "2022-02-09T11:19:40.786259Z",
            "updated_at": "2022-02-09T11:49:52.361321Z"
        },
        {
            "id": "panda101",
            "balance": 2500,
            "currency": "USD",
            "created_at": "2022-02-09T13:11:02.604785Z",
            "updated_at": "2022-02-09T13:11:02.604785Z"
        }
    ]
}
```

### D.3. Get Account by ID
**Purpose:** Getting account by account id
**Endpoint:** *{BASE_URL}/account/v1/account/{id}*
| Param | Value |
| ------ | ------ |
| id | account id  |

#### Positive Scenario 1 - Request
```
curl --location --request GET '{BASE_URL}/account/v1/account/panda101' \
--data-raw ''
```

#### Positive Scenario 1 - Response
```
200 OK
Content-Type: "application/json"

{
    "account": {
        "id": "panda101",
        "balance": 2500,
        "currency": "USD",
        "created_at": "2022-02-09T13:11:02.604785Z",
        "updated_at": "2022-02-09T13:11:02.604785Z"
    }
}
```

### D.4. Send Payment
**Purpose:** Sending payment from one account to another account that registered in this wallet system
**Endpoint:** *{BASE_URL}/payment/v1/payment*
| Param | Value |
| ------ | ------ |
| account_id | account id of the sender  |
| amount | amount of money that will be sent  |
| to_account | account id of the receiver  |

#### Positive Scenario 1 - Request
```
curl --location --request POST '{BASE_URL}/payment/v1/payment' \
--header 'Content-Type: application/json' \
--data-raw '{
   "account_id": "bob123",
   "amount": 50,
   "to_account": "alice456"
}'
```

#### Positive Scenario 1 - Response
```
200 OK
Content-Type: "application/json"
{}
```

### D.5. Get All Payment
**Purpose:** Listing all payments
**Endpoint:** *{BASE_URL}/account/v1/account*
| Param | Value |
| ------ | ------ |
|  |  |

#### Positive Scenario 1 - Request
```
curl --location --request GET '{BASE_URL}/payment/v1/payment' \
--data-raw ''
```

#### Positive Scenario 1 - Response
```
200 Ok
Content-Type: "application/json"

{
    "payments": [
        {
            "id": "e850688a-ef15-4da6-8486-b7656752c87b",
            "account_id": "bob123",
            "transaction_id": "21a3ec76-c017-452d-9eb2-c5d2c0ebbca4",
            "amount": 50,
            "to_account": "alice456",
            "from_account": "",
            "direction": "outgoing",
            "created_at": "2022-02-09T14:15:00.421456Z"
        },
        {
            "id": "9479d132-4553-4ad7-b3b1-978412c88b41",
            "account_id": "alice456",
            "transaction_id": "21a3ec76-c017-452d-9eb2-c5d2c0ebbca4",
            "amount": 50,
            "to_account": "",
            "from_account": "bob123",
            "direction": "incoming",
            "created_at": "2022-02-09T14:15:00.457359Z"
        }
    ]
}
```


## E. The URL to the API
https://www.getpostman.com/collections/fe4b851b64be8c8ee1da


