// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT License",
            "url": "https://opensource.org/license/mit"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/activate": {
            "patch": {
                "security": [],
                "description": "Activate a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Activate a user",
                "parameters": [
                    {
                        "description": "activate user payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.activateUserPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user successfully activated",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "security": [],
                "description": "Log in a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Log in a user",
                "parameters": [
                    {
                        "description": "log in payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.loginUserPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user successfully logged in",
                        "schema": {
                            "$ref": "#/definitions/response.DocsSuccessResponseLoginUser"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/auth/password-reset-request": {
            "post": {
                "security": [],
                "description": "Request to reset a users password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Request to reset a users password",
                "parameters": [
                    {
                        "description": "password reset request payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.passwordResetRequestPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "security": [],
                "description": "Register a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a user",
                "parameters": [
                    {
                        "description": "register user payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.registerUserPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.DocsSuccessResponseRegisterUser"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/auth/resend-verification-email": {
            "post": {
                "security": [],
                "description": "Resend verification email to user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Resend verification email to user",
                "parameters": [
                    {
                        "description": "resend verification email payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.resendVerificationEmailPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "email sent if account exists",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/auth/reset-password": {
            "post": {
                "security": [],
                "description": "Reset a users password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reset a users password",
                "parameters": [
                    {
                        "description": "reset user password payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.resetUserPasswordPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/auth/users/{userID}/deactivate": {
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Deactivate a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Deactivate a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "userID to deactivate",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user successfully deactivated",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/organizations": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create an organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Create an organization",
                "parameters": [
                    {
                        "description": "create organization payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/organizations.createOrganizationPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "organization successfully created",
                        "schema": {
                            "$ref": "#/definitions/organizations.orgResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/organizations/{orgID}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update an organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Update an organization",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "orgID to update",
                        "name": "orgID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update organization payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/organizations.updateOrganizationPyload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "organization successfully updated",
                        "schema": {
                            "$ref": "#/definitions/organizations.orgResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseForbidden"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete an organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Delete an organization",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "orgID to delete",
                        "name": "orgID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "organization successfully deleted",
                        "schema": {
                            "$ref": "#/definitions/response.DocsSuccessResponseDoneMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseForbidden"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Deactivate an organization(staff or admin)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Deactivate an organization(staff or admin)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "orgID to deactivate",
                        "name": "orgID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "organization successfully deactivated",
                        "schema": {
                            "$ref": "#/definitions/response.DocsSuccessResponseDoneMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseForbidden"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        },
        "/profiles/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a users details based on token provided",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profiles"
                ],
                "summary": "Get a users details based on token provided",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/profiles.userProfile"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update a users profile details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profiles"
                ],
                "summary": "Update a users profile details",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/profiles.userProfile"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete a users profile details(soft delete)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profiles"
                ],
                "summary": "Delete a users account details(soft delete)",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.DocsResponseMessageOnly"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.DocsErrorResponseInternalServerErr"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.activateUserPayload": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "auth.loginUserPayload": {
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
        "auth.passwordResetRequestPayload": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.registerUserPayload": {
            "type": "object",
            "required": [
                "dateOfBirth",
                "email",
                "password",
                "profilePic",
                "username"
            ],
            "properties": {
                "dateOfBirth": {
                    "type": "string",
                    "example": "mm/dd/yyyy"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 10
                },
                "profilePic": {
                    "type": "string",
                    "example": "https://fake.link/img.png"
                },
                "username": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 3
                }
            }
        },
        "auth.resendVerificationEmailPayload": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.resetUserPasswordPayload": {
            "type": "object",
            "required": [
                "password",
                "token"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 10
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "organizations.createOrganizationPayload": {
            "type": "object",
            "required": [
                "description",
                "name",
                "profilePic"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 100
                },
                "profilePic": {
                    "type": "string"
                }
            }
        },
        "organizations.orgResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "profilePic": {
                    "type": "string"
                }
            }
        },
        "organizations.updateOrganizationPyload": {
            "type": "object",
            "required": [
                "description",
                "name",
                "profilePic"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 100
                },
                "profilePic": {
                    "type": "string"
                }
            }
        },
        "profiles.userProfile": {
            "type": "object",
            "properties": {
                "dateOfBirth": {
                    "type": "string",
                    "example": "mm/dd/yyyy"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "profilePic": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "example": "client"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "response.DocsErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "object",
                    "properties": {
                        "fieldName": {
                            "type": "string",
                            "example": "error message"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "Invalid request body"
                }
            }
        },
        "response.DocsErrorResponseForbidden": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "forbidden"
                }
            }
        },
        "response.DocsErrorResponseInternalServerErr": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "internal server error"
                }
            }
        },
        "response.DocsErrorResponseUnauthorized": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "unauthorized"
                }
            }
        },
        "response.DocsResponseMessageOnly": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "response.DocsSuccessResponseDoneMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Done"
                }
            }
        },
        "response.DocsSuccessResponseLoginUser": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "properties": {
                        "token": {
                            "type": "string",
                            "example": "jwt.access.token"
                        }
                    }
                }
            }
        },
        "response.DocsSuccessResponseRegisterUser": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Done."
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

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "MeetUp Clone API",
	Description:      "API for MeetUp Clone",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
