package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	net_url "net/url"

	"github.com/gin-gonic/gin"
)

func (s *Server) handler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var urls []string
		if err := c.BindJSON(&urls); err != nil {
			msg := fmt.Sprintf("Oops, something were wrong: %v", err)
			s.logger.Errorf(msg)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if err := validateUrls(urls); err != nil {
			msg := fmt.Sprintf("Bad request: %v", err)
			s.logger.Errorf(msg)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		schemas, err := s.service.GetSchema(urls)
		if err != nil {
			c.String(http.StatusInternalServerError, "Can't get schemas: %v\n", err)
		}

		b, err := json.Marshal(schemas)
		if err != nil {
			c.String(http.StatusInternalServerError, "Can't marshall json: %v\n", err)
			return
		}

		c.String(http.StatusOK, "%s\n", b)
	}
}

func validateUrls(urls []string) error {
	for _, url := range urls {
		if _, err := net_url.ParseRequestURI(url); err != nil {
			return err
		}
	}
	return nil
}
