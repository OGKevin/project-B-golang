// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-03-17 11:18:41.78546 +0100 CET m=+0.057865037

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
                            "$ref": "#/definitions/user.createUserRequest"
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
        }
    },
    "definitions": {
        "responses.Ack": {
            "type": "object",
            "properties": {
                "ack": {
                    "type": "boolean"
                }
            }
        },
        "responses.BadRequest": {
            "type": "object",
            "properties": {
                "ack": {
                    "type": "object",
                    "$ref": "#/definitions/responses.Ack"
                },
                "error": {
                    "type": "object",
                    "$ref": "#/definitions/responses.Error"
                }
            }
        },
        "responses.Created": {
            "type": "object",
            "properties": {
                "ack": {
                    "type": "object",
                    "$ref": "#/definitions/responses.Ack"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "responses.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "user.createUserRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
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
