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
        "/goods/v1/brands": {
            "post": {
                "description": "创建品牌",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Brand"
                ],
                "summary": "创建品牌",
                "operationId": "/goods/v1/brands",
                "parameters": [
                    {
                        "description": "body",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.BrandForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/goods/v1/categorys": {
            "post": {
                "description": "创建分类",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Category"
                ],
                "summary": "创建分类",
                "operationId": "/goods/v1/categorys",
                "parameters": [
                    {
                        "description": "body",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.CategoryForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/goods/v1/goods": {
            "post": {
                "description": "创建商品",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Goods"
                ],
                "summary": "创建商品",
                "operationId": "/goods/v1/goods",
                "parameters": [
                    {
                        "description": "body",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/forms.GoodsForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "forms.BrandForm": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "logo": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 3
                }
            }
        },
        "forms.CategoryForm": {
            "type": "object",
            "required": [
                "is_tab",
                "level",
                "name"
            ],
            "properties": {
                "is_tab": {
                    "type": "boolean"
                },
                "level": {
                    "type": "integer",
                    "enum": [
                        1,
                        2,
                        3
                    ]
                },
                "name": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 3
                },
                "parent": {
                    "type": "integer"
                }
            }
        },
        "forms.GoodsForm": {
            "type": "object",
            "required": [
                "brand",
                "category",
                "desc_images",
                "front_image",
                "goods_brief",
                "goods_sn",
                "images",
                "market_price",
                "name",
                "ship_free",
                "shop_price",
                "stocks"
            ],
            "properties": {
                "brand": {
                    "type": "integer"
                },
                "category": {
                    "type": "integer"
                },
                "desc_images": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                },
                "front_image": {
                    "type": "string"
                },
                "goods_brief": {
                    "type": "string",
                    "minLength": 3
                },
                "goods_sn": {
                    "type": "string",
                    "minLength": 2
                },
                "images": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "string"
                    }
                },
                "market_price": {
                    "type": "number",
                    "minimum": 0
                },
                "name": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 2
                },
                "ship_free": {
                    "type": "boolean"
                },
                "shop_price": {
                    "type": "number",
                    "minimum": 0
                },
                "stocks": {
                    "type": "integer",
                    "minimum": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
