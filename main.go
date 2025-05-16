package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type message struct {
	Command   string `json:"command"`
	Timestamp string `json:"timestamp"`
}

type sendmsg struct {
	Command   string `json:"command"`
	Target    string `json:"target"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// get .env file var
	var pwd = os.Getenv("control-password")
	fmt.Println(pwd)

	socketClients := make(map[string]map[string]interface{}) // clients [websocket.Conn]
	router := gin.Default()

	var Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	router.GET("/ws", func(c *gin.Context) {
		client := c.Query("client")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if client != "display" && client != "operation" && client != "control" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client type"})
			return
		} else if client == "control" {
			authHeader := c.Request.Header.Get("Authorization")
			if authHeader != "Bearer "+pwd {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
		}

		ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Error during connection upgrade:", err)
			return
		}

		uuidV4, _ := uuid.NewRandom()
		clientId := uuidV4.String()

		defer func(ws *websocket.Conn) {
			delete(socketClients, clientId)
			err := ws.Close()
			if err != nil {
				log.Println("Error closing connection:", err)
			}
		}(ws)

		socketClients[clientId] = make(map[string]interface{})
		socketClients[clientId]["ws"] = ws
		socketClients[clientId]["type"] = client
		socketClients[clientId]["lastactive"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
		socketClients[clientId]["status"] = "1"

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			var msgData message
			if err = json.Unmarshal(msg, &msgData); err != nil {
				log.Println("Error unmarshalling message:", err)
				break
			}

			if msgData.Command == "ping" {
				socketClients[clientId]["lastactive"] = msgData.Timestamp
				pongMessage := message{
					Command:   "pong",
					Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
				}
				pongMessageBytes, _ := json.Marshal(pongMessage)
				if err = ws.WriteMessage(websocket.TextMessage, pongMessageBytes); err != nil {
					log.Println("Error sending pong message:", err)
					break
				}
			}

			if client == "display" {
				if msgData.Command == "ready" {
					socketClients[clientId]["status"] = "2"
				} else {
					errMsg := message{
						Command:   "unknown command",
						Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
					}
					errMsgBytes, _ := json.Marshal(errMsg)

					if err := ws.WriteMessage(websocket.TextMessage, errMsgBytes); err != nil {
						return
					}
				}
			} else if client == "operation" {
				if msgData.Command == "start" {

					if client == "display" {
						return
					}

					for id, clientData := range socketClients {
						if clientData["type"] == "display" {
							fmt.Println("b")
							clientWs := clientData["ws"].(*websocket.Conn)
							playMessage := sendmsg{
								Command:   "start",
								Target:    id,
								Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
							}
							playMessageBytes, _ := json.Marshal(playMessage)
							if err = clientWs.WriteMessage(websocket.TextMessage, playMessageBytes); err != nil {
								log.Println("Error sending play message:", err)
								break
							}
						}
					}

				}
			}

		}

	})

	router.GET("/files/:file", func(c *gin.Context) {
		c.File("public/files/" + c.Param("file"))
	})

	router.GET("/operation/:file", func(c *gin.Context) {
		path := c.Param("file")
		c.File("public/operation/" + path)
	})

	router.GET("/operation", func(c *gin.Context) {
		c.File("./public/operation/index.html")
	})

	router.GET("/display/:file", func(c *gin.Context) {
		c.File("./public/display/" + c.Param("file"))
	})

	router.GET("/display", func(c *gin.Context) {
		c.File("./public/display/index.html")
	})

	router.GET("/:file", func(c *gin.Context) {
		path := c.Param("file")
		c.File("public/" + path)
	})

	router.GET("/", func(c *gin.Context) {
		c.File("public/index.html")
	})

	err := router.Run(":7380")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
