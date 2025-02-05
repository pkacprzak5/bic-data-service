# BIC Data Service

A RESTful microservice for managing BIC/SWIFT codes and bank information, built in Go with PostgreSQL.

## Table of Contents
- [About the Project](#about-the-project)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Exposed endpoints](#exposed-endpoints)
- [Technical Details](#technical-details)
- [License](#license)

## About the Project
BIC Data Service is an appliaction useful for storing bank's data. This service provides:
- **Centralized Storage**: Securely store BIC/SWIFT codes, bank details, and branch information
- **Instant Access**: REST API endpoints for real-time data retrieval
- **Regulatory Compliance**: Built-in ISO2 standard support for country codes validation
- **Operational Clarity**: Human-readable error messages and OpenAPI documentation

## Getting Started
### Prerequisites
- **Docker**: Make sure Docker is installed on your system. Follow [this guide](https://docs.docker.com/get-docker/) to install Docker on your platform.

## Usage
### 1. Clone the repo
```
git clone https://github.com/pkacprzak5/bic-data-service.git
```

### 2. Navigate to repository path
```
cd <your-path>/bic-data-service
```

### 3. Download dependencies and run app in Docker
```
docker-compose up --build
```

### 4. Explore application
After everything is set up correctly endpoints will be exposed at ```http://localhost:8080```

### 5. Running tests
- unit tests and integration tests
  ```
  make docker-test
  ```
- unit tests and end2end tests
  ```
  make docker-test-end2end
  ```

## Exposed endpoints:
1. Retrieve details of a single SWIFT code whether for a headquarters or branches.</br>

   #### **GET** `/v1/swift-codes/{swift-code}`</br>
   
   If given SWIFT code is valid and exists in database returns following json:
   - For headquarter SWIFT code:
     ```json
     {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "countryName": "string",
      "isHeadquarter": "bool",
      "swiftCode": "string",
      "branches": [
       {
         "address": "string",
         "bankName": "string",
         "countryISO2": "string",
         "isHeadquarter": "bool",
         "swiftCode": "string"
       },
       {
         "address": "string",
         "bankName": "string",
         "countryISO2": "string",
         "isHeadquarter": "bool",
         "swiftCode": "string"
       },
      ]
     }
     ```
   - For branch SWIFT code:
     ```json
     {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "countryName": "string",
      "isHeadquarter": "bool",
      "swiftCode": "string",
     }
     ```

     
2. Return all SWIFT codes with details for a specific country (both headquarters and branches).</br>

   #### **GET** ` /v1/swift-codes/country/{countryISO2code}`</br>
   
   If given ISO2 code is valid and there are records of banks for this country, it returns following json:
   ```json
     {
      "countryISO2": "string",
      "countryName": "string",
      "swiftCodes": [
       {
         "address": "string",
         "bankName": "string",
         "countryISO2": "string",
         "isHeadquarter": "bool",
         "swiftCode": "string"
       },
       {
         "address": "string",
         "bankName": "string",
         "countryISO2": "string",
         "isHeadquarter": "bool",
         "swiftCode": "string"
       },
      ]
     }
     ```

   
3. Adds new SWIFT code entries to the database for a specific country.</br>

   #### **POST** `/v1/swift-codes`</br>
   
   It requires following request structure:
   ```json
     {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "countryName": "string",
      "isHeadquarter": "bool",
      "swiftCode": "string",
     }
     ```
   In case request structure is valid, bank's data is added to database.


4. Deletes swift-code data if swiftCode matches the one in the database.</br>

   #### **DELETE** `/v1/swift-codes/{swift-code}`</br>
   
   If given SWIFT code is valid and there exist bank with this SWIFT code in database it is removed from storage.


## Technical details
This app ensures all data is valid by checking regex in SWIFT codes, using  `"github.com/mikekonan/go-countries"` package to test whether given country ISO2 code is valid and correctly paired with given country. </br>
*Warning! This package may have different names for some countries. For example required name for US ISO2 code is `United States of America (the)`.*
*You are able to find full list of names [here](https://github.com/mikekonan/go-countries/blob/main/name_gen.go)*.

## License
Distributed under the MIT License. See ```LICENSE``` for more information.
