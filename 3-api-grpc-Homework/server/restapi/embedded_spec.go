// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "This is the main entry point for the Ikakbolit application, which sets up and runs the application server.",
    "title": "ikakbolit API",
    "contact": {
      "name": "Константин Троицкий",
      "url": "https://t.me/debussy3",
      "email": "varrr7@gmail.com"
    },
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "paths": {
    "/next_takings": {
      "get": {
        "description": "Retrieve upcoming medication schedules for a user",
        "produces": [
          "application/json"
        ],
        "summary": "Get next scheduled takings",
        "parameters": [
          {
            "type": "integer",
            "description": "User ID",
            "name": "user_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/controller.Schedule"
            }
          },
          "204": {
            "description": "No content"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      }
    },
    "/schedule": {
      "get": {
        "description": "Retrieve a schedule by user ID and schedule ID",
        "produces": [
          "application/json"
        ],
        "summary": "Get a specific schedule",
        "parameters": [
          {
            "type": "integer",
            "description": "User ID",
            "name": "user_id",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "description": "Schedule ID",
            "name": "schedule_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/controller.Schedule"
            }
          },
          "204": {
            "description": "No content"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      },
      "post": {
        "description": "Create a new schedule for a user",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Add a new schedule [! Предпочтительнее использовать http-клиент для post-запросов, например Postman]",
        "parameters": [
          {
            "description": "Schedule data",
            "name": "schedule",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/controller.Schedule"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/controller.responseScheduleID"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      }
    },
    "/schedules": {
      "get": {
        "description": "Retrieve schedule IDs for a given user",
        "produces": [
          "application/json"
        ],
        "summary": "Get user schedules",
        "parameters": [
          {
            "type": "integer",
            "description": "User ID",
            "name": "user_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "type": "integer"
              }
            }
          },
          "204": {
            "description": "No content"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "controller.APIError": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "controller.Schedule": {
      "type": "object",
      "properties": {
        "created_at": {
          "type": "string"
        },
        "cure_name": {
          "type": "string"
        },
        "doses_per_day": {
          "type": "integer"
        },
        "duration": {
          "type": "integer"
        },
        "id": {
          "type": "integer"
        },
        "intakes": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "user_id": {
          "type": "integer"
        }
      }
    },
    "controller.responseScheduleID": {
      "type": "object",
      "properties": {
        "schedule_id": {
          "type": "string"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "This is the main entry point for the Ikakbolit application, which sets up and runs the application server.",
    "title": "ikakbolit API",
    "contact": {
      "name": "Константин Троицкий",
      "url": "https://t.me/debussy3",
      "email": "varrr7@gmail.com"
    },
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "paths": {
    "/next_takings": {
      "get": {
        "description": "Retrieve upcoming medication schedules for a user",
        "produces": [
          "application/json"
        ],
        "summary": "Get next scheduled takings",
        "parameters": [
          {
            "type": "integer",
            "description": "User ID",
            "name": "user_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/controller.Schedule"
            }
          },
          "204": {
            "description": "No content"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      }
    },
    "/schedule": {
      "get": {
        "description": "Retrieve a schedule by user ID and schedule ID",
        "produces": [
          "application/json"
        ],
        "summary": "Get a specific schedule",
        "parameters": [
          {
            "type": "integer",
            "description": "User ID",
            "name": "user_id",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "description": "Schedule ID",
            "name": "schedule_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/controller.Schedule"
            }
          },
          "204": {
            "description": "No content"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      },
      "post": {
        "description": "Create a new schedule for a user",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Add a new schedule [! Предпочтительнее использовать http-клиент для post-запросов, например Postman]",
        "parameters": [
          {
            "description": "Schedule data",
            "name": "schedule",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/controller.Schedule"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": {
              "$ref": "#/definitions/controller.responseScheduleID"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      }
    },
    "/schedules": {
      "get": {
        "description": "Retrieve schedule IDs for a given user",
        "produces": [
          "application/json"
        ],
        "summary": "Get user schedules",
        "parameters": [
          {
            "type": "integer",
            "description": "User ID",
            "name": "user_id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "array",
              "items": {
                "type": "integer"
              }
            }
          },
          "204": {
            "description": "No content"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "404": {
            "description": "Resource not found",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/controller.APIError"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "controller.APIError": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "controller.Schedule": {
      "type": "object",
      "properties": {
        "created_at": {
          "type": "string"
        },
        "cure_name": {
          "type": "string"
        },
        "doses_per_day": {
          "type": "integer"
        },
        "duration": {
          "type": "integer"
        },
        "id": {
          "type": "integer"
        },
        "intakes": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "user_id": {
          "type": "integer"
        }
      }
    },
    "controller.responseScheduleID": {
      "type": "object",
      "properties": {
        "schedule_id": {
          "type": "string"
        }
      }
    }
  }
}`))
}
