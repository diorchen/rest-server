
# Grocery Item Store API

The Grocery Item Store API is a RESTful API for managing grocery items with expiration dates and other related information.

## Installation

1. Clone the repository: `git clone https://github.com/diorchen/rest-server`
2. Install the required dependencies: `go mod tidy`

## Usage

1. Build the project: `go build`
2. Run the compiled executable: `./rest-server`

## Endpoints

### Create Food Item

- **URL**: `/food/`
- **Method**: `POST`
- **Request Body**:
```json
{
  "name": "Apple",
  "description": "Fresh and juicy apple",
  "ingredients": ["Apple"],
  "expiration": "2023-07-31T00:00:00Z",
  "nutrition": {
    "calories": 52,
    "protein": 0.3,
    "carbohydrates": 14,
    "fat": 0.2,
    "fiber": 2.4
  }
}
````

### Get All Food Items

- **URL**: `/food/`
- **Method**: `GET`

### Get Food Item

- **URL**: `/food/{id}`
- **Method**: `GET`

### Delete Food Item

- **URL**: `/food/{id}`
- **Method**: `DELETE`

### Get Foods by Ingredient

- **URL**: `/ing/{ingredient}`
- **Method**: `GET`

### Get Foods by Expiration Date

- **URL**: `/exp/{year}/{month}/{day}`
- **Method**: `GET`
