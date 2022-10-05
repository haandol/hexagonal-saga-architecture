// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/efs": {
            "get": {
                "description": "list all files in the given path",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "efs"
                ],
                "summary": "list files",
                "parameters": [
                    {
                        "type": "string",
                        "description": "path",
                        "name": "path",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/trips": {
            "get": {
                "description": "list all trips",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trips"
                ],
                "summary": "list all trips",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Trip"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "create new trip",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trips"
                ],
                "summary": "create new trip",
                "parameters": [
                    {
                        "description": "trip id is required",
                        "name": "trip",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Trip"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Trip"
                        }
                    }
                }
            }
        },
        "/trips/{trip_id}/recover/backward": {
            "put": {
                "description": "recover backward",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trips"
                ],
                "summary": "recover backward",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "trip id",
                        "name": "trip_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Trip"
                        }
                    }
                }
            }
        },
        "/trips/{trip_id}/recover/forward": {
            "put": {
                "description": "recover forward",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trips"
                ],
                "summary": "recover forward",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "trip id",
                        "name": "trip_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Trip"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.Trip": {
            "type": "object",
            "required": [
                "carId",
                "flightId",
                "hotelId",
                "id",
                "userId"
            ],
            "properties": {
                "carBookingId": {
                    "type": "integer"
                },
                "carId": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "flightBookingId": {
                    "type": "integer"
                },
                "flightId": {
                    "type": "integer"
                },
                "hotelBookingId": {
                    "type": "integer"
                },
                "hotelId": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Hexagonal Saga API",
	Description:      "hexagonal saga orchestration example api server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
