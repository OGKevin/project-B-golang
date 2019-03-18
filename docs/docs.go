// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-03-18 18:21:55.902425 +0100 CET m=+0.057255536

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Well you know, nothing important. Just making sure people can capture memories",
        "title": "Project b",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "project-b.ogkevin.nl",
    "basePath": "/api/v1",
    "paths": {
        "/user": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Register a new user",
                "operationId": "register-new-user",
                "parameters": [
                    {
                        "description": "The expected request body. Username must be length(5|255) and Password length(10|255).",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/user.userRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The response will include the id of the newly created user.",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/responses.Created"
                        }
                    },
                    "400": {
                        "description": "The error object will explain why the request failed.",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "on success, you will get a JWT token to put in the auth header",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "logs a user in",
                "operationId": "user-login",
                "parameters": [
                    {
                        "description": "The expected request body.",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/user.userRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The user",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/user.jwtToken"
                        }
                    },
                    "400": {
                        "description": "The error object will explain why the request failed.",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    }
                }
            }
        },
        "/user/{userId}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "gets user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "gets user by id",
                "operationId": "get-user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The id to get the user",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The BEARER token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The user",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "The error object will explain why the request failed.",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    },
                    "404": {
                        "description": "The error object will explain why the entity was not found.",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/responses.NotFound"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.Ack": {
            "type": "object",
            "properties": {
                "ack": {
                    "description": "Ack Defines if the server could acknowledge the request.",
                    "type": "boolean"
                }
            }
        },
        "responses.BadRequest": {
            "type": "object",
            "properties": {
                "ack": {
                    "description": "Ack Defines if the request was successful or not.",
                    "type": "object",
                    "$ref": "#/definitions/responses.Ack"
                },
                "error": {
                    "description": "Error Explains why the server is responding with a bad request.",
                    "type": "object",
                    "$ref": "#/definitions/responses.Error"
                }
            }
        },
        "responses.Created": {
            "type": "object",
            "properties": {
                "ack": {
                    "description": "Ack Defines if the request was successful or not.",
                    "type": "object",
                    "$ref": "#/definitions/responses.Ack"
                },
                "id": {
                    "description": "ID The id of the newly created entity.",
                    "type": "string"
                }
            }
        },
        "responses.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Coode The http status code that belongs to this error.",
                    "type": "integer"
                },
                "message": {
                    "description": "Message The message explaining the error.",
                    "type": "string"
                }
            }
        },
        "responses.NotFound": {
            "type": "object",
            "properties": {
                "ack": {
                    "description": "Ack Defines if the request was successful or not.",
                    "type": "object",
                    "$ref": "#/definitions/responses.Ack"
                },
                "error": {
                    "description": "Error Explains why the server is responding with a bad request.",
                    "type": "object",
                    "$ref": "#/definitions/responses.Error"
                }
            }
        },
        "user.User": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "the userId",
                    "type": "string"
                },
                "username": {
                    "description": "the username",
                    "type": "string"
                }
            }
        },
        "user.jwtToken": {
            "type": "object",
            "properties": {
                "ack": {
                    "type": "object",
                    "$ref": "#/definitions/responses.Ack"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "user.userRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "description": "Password The user's password, must be length(5|255)",
                    "type": "string"
                },
                "username": {
                    "description": "Username The user's username, must be unique and length(5|255)",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
