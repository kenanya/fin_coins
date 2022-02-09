## A. Description
This is a backend service (API and database) for a generic wallet service. This service contains 2 entity:
- Account
- Payment

Tech Stack: golang, go-kit framework, PostgreSQL, docker

## B. How to set up and start the backend server
The config file is located at:
- *meal_be/common/configGlobal.yaml*

You can change the values according to your configuration. 


#### Clone the Project
```bash
git clone https://github.com/kenanya/fin_coins.git
```

#### Build Project
```bash
docker-compose -f docker-compose-local.yml build
```

#### Run Project
```bash
docker-compose -f docker-compose-local.yml up
```

## C. Integration Test
Firstly we have to create database that is defined in fin_coins/common/configGlobal.yaml 

These are the steps to run the test:
```bash
go test -run Test_RaceCondition
APP_ENV=local go test -run Test_RaceCondition -v
```

Those test will do account creation and payment transfer.


## D. The API documentation
Below are the sample requests and expected responses for each API:

### D.1. Create Account
**Purpose:** Creating account
**Endpoint:** *localhost:9595/account/v1/account*
| Param | Value |
| ------ | ------ |
| id | unique id for account  
| balance | float/decimal  
| currency | currency, such as USD, IDR, etc  

#### Positive Scenario 1 - Request
```
curl --location --request POST 'localhost:9595/account/v1/account' \
--header 'Content-Type: application/json' \
--data-raw '{
   "id": "panda101",
   "balance": 2500,
   "currency": "USD"
}'
```

#### Positive Scenario 1 - Response
```json
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
**Endpoint:** *localhost:9595/account/v1/account*
| Param | Value |
| ------ | ------ |
|  |  |

#### Positive Scenario 1 - Request
```json
curl --location --request GET 'localhost:9595/account/v1/account' \
--data-raw ''
```

#### Positive Scenario 1 - Response
```json
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
**Endpoint:** *localhost:9595/account/v1/account/{id}*
| Param | Value |
| ------ | ------ |
| id | account id  |

#### Positive Scenario 1 - Request
```json
curl --location --request GET 'localhost:9595/account/v1/account/panda101' \
--data-raw ''
```

#### Positive Scenario 1 - Response
```json
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
**Endpoint:** *localhost:9595/payment/v1/payment*
| Param | Value |
| ------ | ------ |
| account_id | account id of the sender  |
| amount | amount of money that will be sent  |
| to_account | account id of the receiver  |

#### Positive Scenario 1 - Request
```json
curl --location --request POST 'localhost:9595/payment/v1/payment' \
--header 'Content-Type: application/json' \
--data-raw '{
   "account_id": "bob123",
   "amount": 50,
   "to_account": "alice456"
}'
```

#### Positive Scenario 1 - Response
```json
200 OK
Content-Type: "application/json"
{}
```

### D.5. Get All Payment
**Purpose:** Listing all payments
**Endpoint:** *localhost:9595/account/v1/account*
| Param | Value |
| ------ | ------ |
|  |  |

#### Positive Scenario 1 - Request
```json
curl --location --request GET 'localhost:9595/payment/v1/payment' \
--data-raw ''
```

#### Positive Scenario 1 - Response
```json
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


## E. The URL to the API
https://www.getpostman.com/collections/fe4b851b64be8c8ee1da


