// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "rahmat.putra@spesolution.com"
        },
        "license": {
            "name": "Apache 2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/api/auth/check-token": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Check Token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Check Token",
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/entity.GeneralResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid Access Token",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Payload Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/api/auth/login": {
            "post": {
                "description": "Login by using registered account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Payload Request Body",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Invalid Access Token",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Payload Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/api/auth/register": {
            "post": {
                "description": "Create User for Guest",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Create User as Guest",
                "parameters": [
                    {
                        "description": "Payload Request Body",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.CreateUserReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.CreateUserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Invalid Access Token",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Payload Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/api/todo-list": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Retrieve a list of Todo Lists belonging to a user by their User ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo List"
                ],
                "summary": "Retrieve Todo Lists by User ID",
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.TodoListResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Update an existing Todo List",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo List"
                ],
                "summary": "Update an existing Todo List by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the todo list",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Payload Request Body",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.TodoListReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/entity.GeneralResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create a new Todo List",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo List"
                ],
                "summary": "Create a new Todo List",
                "parameters": [
                    {
                        "description": "Payload Request Body",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.TodoListReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.TodoListReq"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/api/todo-lists/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get a Todo List by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo List"
                ],
                "summary": "Get Todo List by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the Todo List",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.GeneralResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.TodoListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Delete an existing Todo List by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo List"
                ],
                "summary": "Delete Todo List by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the todo list",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/entity.GeneralResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Invalid Request Body",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.CustomErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.CreateUserReq": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "phone",
                "reenter_password",
                "role_access"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "reenter_password": {
                    "type": "string"
                },
                "role_access": {
                    "type": "integer"
                }
            }
        },
        "entity.CreateUserResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "role_access": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "entity.CustomErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "http_code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.GeneralResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.LoginReq": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.LoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role_access": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "entity.TodoListReq": {
            "type": "object",
            "required": [
                "description",
                "doing_at",
                "title",
                "user_id"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "doing_at": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "entity.TodoListResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "doing_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Skeleton!",
	Description:      "This is a sample swagger for Go Skeleton",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
