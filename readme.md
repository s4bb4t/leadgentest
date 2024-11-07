# Golang + PostgreSQL + Gin API

## Task Overview

This project implements a simple service with three endpoints for managing building records using Golang, PostgreSQL, and the Gin framework. The service includes the following functionality:

1. **Create Building**  
   Accepts a building object with the following fields:
   - `Title` (string): Name of the building
   - `City` (string): City where the building is located
   - `Year` (int): Year of completion
   - `Floors` (int): Number of floors  
   The service inserts the building into a PostgreSQL database.
   

2. **Get Building**  
   Retrieves a building by `Title` in url:
    - `Title`
    - `City`
    - `Year` 
    - `Floors` 


3. **Get Buildings**  
   Retrieves a list of buildings with optional filtering by:
   - `City`
   - `Year`
   - `Floors`  
   Filters are not mandatory.


4. **OpenAPI Documentation**  
   The API documentation is automatically generated using `swaggo` (https://github.com/swaggo/swag) based on code annotations.

## Database Configuration

The connection to the PostgreSQL database is configured via a `config` file with the following parameters:
- `host`: Host of the PostgreSQL server
- `port`: Port number for the PostgreSQL connection
- `user`: Database user
- `password`: Password for the database user
- `db`: Name of the database

## Contact Information

- **Dmitriy Bratishkin**
- **Phone**: +7 (913) 475-08-87 (preferred method of contact)
- **Email**: s4bb4t@yandex.ru
- **Telegram**: [@s4bb4t](https://t.me/s4bb4t)
```